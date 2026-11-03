package host

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"linkz/host/window"
	"lizardscript/object"
)

// homeHarness stands in for the app's screen router, which app.liz owns: a
// screens table, show_screen, and an editor screen that only remembers it was
// shown. It routes the app's events and paints to the current screen, and
// puts the home screen's backdoor commands where the real app will.
const homeHarness = `
let screens = {}
let harness = {"screen": "home", "editor_shown": 0}
screens["editor"] = {
  "event": function(e) { null },
  "paint": function() { fb.target("ui"); fb.fill((0, 0), fb.size(), #223322); null },
  "enter": function(args) { harness.editor_shown += 1; null },
  "needs_paint": function() { false },
  "wants_hover": function() { false }
}
home_register()
function show_screen(name, args) {
  harness.screen = name
  screens[name].enter(args)
  app["needs_paint"] = true
}
app["event"] = function(e) {
  screens[harness.screen].event(e)
  app["needs_paint"] = screens[harness.screen].needs_paint()
  null
}
app["paint"] = function() {
  screens[harness.screen].paint()
  app["needs_paint"] = screens[harness.screen].needs_paint()
  null
}
app["debug_commands"] = home_debug_commands()
app["debug_state"] = function() { home_debug_state() }
`

// drawings by name with how long ago they were edited
var homeDrawings = []struct {
	name string
	ago  time.Duration
}{
	{"layers", 2 * time.Hour}, {"strokes", 26 * time.Hour}, {"night-city", 50 * time.Hour},
	{"hello", 8 * 24 * time.Hour}, {"mountains", 15 * 24 * time.Hour}, {"lily", 4*24*time.Hour + time.Hour},
	{"badge", 5*24*time.Hour + time.Hour}, {"waves", 22 * 24 * time.Hour}, {"poster", 40 * 24 * time.Hour},
	{"grid", 45 * 24 * time.Hour}, {"bird", 70 * 24 * time.Hour}, {"tree", 75 * 24 * time.Hour},
}

// newHomeApp makes a headless app in a temporary home folder with the
// drawings of testdata/home (or none), shows the home screen and returns it.
func newHomeApp(t *testing.T, withDrawings bool, lang string) (*App, string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := filepath.Join(home, "Documents", "Linkz")
	if withDrawings {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		for _, d := range homeDrawings {
			body, err := os.ReadFile(filepath.Join("testdata", "home", d.name+".lzdraw"))
			if err != nil {
				t.Fatal(err)
			}
			p := filepath.Join(dir, d.name+".lzdraw")
			writeFile(t, p, string(body))
			when := time.Now().Add(-d.ago)
			os.Chtimes(p, when, when)
		}
	}
	app := NewApp(1024, 720)
	app.headless = true
	app.Settings = fromGo(map[string]any{"language": lang}).(*object.Hash)
	if err := app.Load(""); err != nil {
		t.Fatal(err)
	}
	if err := app.call(homeHarness, nil); err != nil {
		t.Fatal(err)
	}
	if err := app.call(`show_screen("home", null)`, nil); err != nil {
		t.Fatal(err)
	}
	settle(t, app)
	return app, dir
}

// settle ticks until nothing wants a paint, so every thumbnail is made.
func settle(t *testing.T, app *App) {
	t.Helper()
	for i := 0; i < 80; i++ {
		app.Tick(window.Input{})
		if !app.needsPaint() {
			return
		}
	}
	t.Fatal("the screen still wants to paint after 80 ticks")
}

func evalInt(t *testing.T, app *App, src string) int {
	t.Helper()
	v, err := app.eval(src, nil)
	if err != nil {
		t.Fatalf("%s: %v", src, err)
	}
	n, ok := toGo(v, 0).(int64)
	if !ok {
		t.Fatalf("%s gave %s, want an integer", src, v.Inspect())
	}
	return int(n)
}

func evalStr(t *testing.T, app *App, src string) string {
	t.Helper()
	v, err := app.eval(src, nil)
	if err != nil {
		t.Fatalf("%s: %v", src, err)
	}
	return v.Inspect()
}

// click presses and releases the left button at a point in UI pixels.
func click(t *testing.T, app *App, x, y int) {
	t.Helper()
	app.Inject(window.Event{Type: "move", X: x * 2, Y: y * 2},
		window.Event{Type: "down", X: x * 2, Y: y * 2},
		window.Event{Type: "up", X: x * 2, Y: y * 2})
	settle(t, app)
}

func hoverAt(t *testing.T, app *App, x, y int) {
	t.Helper()
	app.Inject(window.Event{Type: "move", X: x * 2, Y: y * 2})
	settle(t, app)
}

func typeText(t *testing.T, app *App, s string) {
	t.Helper()
	for _, r := range s {
		app.Inject(window.Event{Type: "text", Text: string(r)})
	}
	settle(t, app)
}

func pressKey(t *testing.T, app *App, key string, mods window.Mods) {
	t.Helper()
	app.Inject(window.Event{Type: "key", Key: key, Mods: mods})
	settle(t, app)
}

// saveShot writes a screenshot when LINKZ_SHOTS names a folder.
func saveShot(t *testing.T, app *App, name string) {
	t.Helper()
	dir := os.Getenv("LINKZ_SHOTS")
	if dir == "" {
		return
	}
	os.MkdirAll(dir, 0o755)
	f, err := os.Create(filepath.Join(dir, name+".png"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, app.Screenshot()); err != nil {
		t.Fatal(err)
	}
}

func thumbFiles(t *testing.T) map[string]time.Time {
	t.Helper()
	cfg, _ := ConfigDir()
	entries, _ := os.ReadDir(filepath.Join(cfg, "thumbs"))
	out := map[string]time.Time{}
	for _, e := range entries {
		info, _ := e.Info()
		out[e.Name()] = info.ModTime()
	}
	return out
}

func TestHomeGalleryAndThumbnailCache(t *testing.T) {
	app, dir := newHomeApp(t, true, "en")
	saveShot(t, app, "home-en")

	// 3 "working on" cards and the first two rows of the grid are visible.
	if n := evalInt(t, app, `len(home.cards)`); n != 3+12 {
		t.Fatalf("cards = %d, want 3 working + 12 all", n)
	}
	if n := evalInt(t, app, `len(home_thumbs.made)`); n != 3+10 {
		t.Errorf("thumbnails made = %d, want 13 (3 big + 10 visible small)", n)
	}
	renders := evalInt(t, app, `home_thumbs.renders`)
	if renders != 13 {
		t.Errorf("renders = %d, want one per thumbnail", renders)
	}
	files := thumbFiles(t)
	if len(files) != 13 {
		t.Errorf("%d PNG files in ~/.linkz/thumbs, want 13", len(files))
	}
	for name := range files {
		if !strings.HasSuffix(name, ".png") || !strings.Contains(name, "-") {
			t.Errorf("odd thumbnail name %q", name)
		}
	}

	// The first big card shows the sunset: its sky color is on screen.
	if !hasColor(app.Screenshot(), image.Rect(2*(12+80+8), 2*46, 2*(12+80+8+128), 2*(46+108)), 0xf2, 0xa6, 0x5a) {
		t.Error("the sunset's sky is not on the first working-on card")
	}

	// Going home again is free: nothing is painted again.
	app.call(`show_screen("home", null)`, nil)
	settle(t, app)
	if n := evalInt(t, app, `home_thumbs.renders`); n != renders {
		t.Errorf("renders went from %d to %d after coming back", renders, n)
	}

	// A restart: the memory is gone, the PNGs are read, not painted.
	app.call(`home_thumbs.made = {}`, nil)
	app.call(`show_screen("home", null)`, nil)
	settle(t, app)
	if n := evalInt(t, app, `home_thumbs.renders`); n != renders {
		t.Errorf("renders = %d after a restart, want %d (cached PNGs)", n, renders)
	}
	if n := evalInt(t, app, `len(home_thumbs.made)`); n != 13 {
		t.Errorf("thumbnails after a restart = %d, want 13", n)
	}
	for name, mt := range thumbFiles(t) {
		if !mt.Equal(files[name]) {
			t.Errorf("%s was written again", name)
		}
	}

	// A drawing that changes gets a new thumbnail; the others keep theirs.
	p := filepath.Join(dir, "layers.lzdraw")
	body, _ := os.ReadFile(p)
	writeFile(t, p, string(body)+"\ncircle((50, 50), 20, #ff0000)\n")
	app.call(`show_screen("home", null)`, nil)
	settle(t, app)
	if n := evalInt(t, app, `home_thumbs.renders`); n != renders+2 { // big and small
		t.Errorf("renders = %d after one drawing changed, want %d", n, renders+2)
	}
	if n := len(thumbFiles(t)); n != 15 {
		t.Errorf("%d thumbnails on disk, want 15 (13 + the 2 new sizes of the changed file)", n)
	}
}

func hasColor(img *image.RGBA, r image.Rectangle, cr, cg, cb uint8) bool {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			p := img.Pix[y*img.Stride+x*4:]
			if p[0] == cr && p[1] == cg && p[2] == cb {
				return true
			}
		}
	}
	return false
}

func TestHomeBrokenDrawingGetsErrorThumbnail(t *testing.T) {
	app, dir := newHomeApp(t, true, "en")
	writeFile(t, filepath.Join(dir, "broken.lzdraw"), "circle((1, 2), \n")
	app.call(`show_screen("home", null)`, nil)
	settle(t, app)
	if got := evalStr(t, app, `home_thumb_state(home.cards[0])`); got != "error" {
		t.Errorf("broken drawing's thumbnail = %s, want error", got)
	}
	before := evalInt(t, app, `home_thumbs.renders`)
	app.call(`show_screen("home", null)`, nil)
	settle(t, app)
	if n := evalInt(t, app, `home_thumbs.renders`); n != before {
		t.Error("a drawing that fails is run again on every paint")
	}
}

func TestHomeSpanishAndEmpty(t *testing.T) {
	app, _ := newHomeApp(t, true, "es")
	saveShot(t, app, "home-es")
	if got := evalStr(t, app, `home_when(time.now() - 7200, time.now(), 0)`); got != "hace 2 h" {
		t.Errorf("Spanish time = %s", got)
	}
	app, _ = newHomeApp(t, false, "en")
	saveShot(t, app, "home-empty")
	if got := evalStr(t, app, `home_debug_state().empty`); got != "true" {
		t.Errorf("empty = %s", got)
	}
	app, _ = newHomeApp(t, false, "es")
	saveShot(t, app, "home-empty-es")
}

func TestHomeFlows(t *testing.T) {
	app, dir := newHomeApp(t, true, "en")

	// Opening a card opens the file in the editor and records it as recent.
	click(t, app, 12+80+8+60, 46+30) // the first working-on card
	if got := evalStr(t, app, `harness.screen`); got != "editor" {
		t.Fatalf("screen = %s after clicking a card, want editor", got)
	}
	if got := evalStr(t, app, `ed.file`); got != filepath.Join(dir, "layers.lzdraw") {
		t.Errorf("ed.file = %s", got)
	}
	cfg, _ := ConfigDir()
	recent, _ := os.ReadFile(filepath.Join(cfg, "recent.liz"))
	if !strings.Contains(string(recent), "layers.lzdraw") {
		t.Errorf("recent.liz = %q", recent)
	}

	// Back home: layers is first of the recent ones. Search filters.
	app.call(`show_screen("home", null)`, nil)
	settle(t, app)
	typeText(t, app, "ni")
	if n := evalInt(t, app, `len(home.cards)`); n != 1 {
		t.Errorf("search \"ni\" shows %d cards, want 1 (night-city)", n)
	}
	saveShot(t, app, "home-search")
	pressKey(t, app, "Escape", window.Mods{})
	if n := evalInt(t, app, `len(home.cards)`); n != 3+12 {
		t.Errorf("clearing the search shows %d cards, want 15", n)
	}

	// The menu: hover the 5th grid card, press its menu button, see the menu.
	cardX, cardY := 12+4*99, 180
	hoverAt(t, app, cardX+30, cardY+20)
	click(t, app, cardX+91-18+7, cardY+4+6)
	if got := evalStr(t, app, `home.menu != null`); got != "true" {
		t.Fatalf("the menu did not open")
	}
	saveShot(t, app, "home-menu")

	// Rename through the dialog: click "Rename" in the menu, edit the name.
	click(t, app, evalInt(t, app, `home_menu_row(1).x + 5`), evalInt(t, app, `home_menu_row(1).y + 5`))
	if got := evalStr(t, app, `home.dialog.kind`); got != "rename" {
		t.Fatalf("dialog = %s after clicking Rename", got)
	}
	saveShot(t, app, "home-rename")
	pressKey(t, app, "Backspace", window.Mods{})
	typeText(t, app, "x")
	pressKey(t, app, "Enter", window.Mods{})
	if _, err := os.Stat(filepath.Join(dir, "badgx.lzdraw")); err != nil {
		t.Errorf("renamed file missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "badge.lzdraw")); err == nil {
		t.Error("the old name is still there")
	}

	// Duplicate makes "<name> copy".
	app.call(`home_duplicate(home.files[0])`, nil)
	if got := evalInt(t, app, `len(home.files)`); got != 13 {
		t.Errorf("files after duplicate = %d, want 13", got)
	}

	// Delete asks, and moves to the Trash only when confirmed.
	victim := evalStr(t, app, `home.files[len(home.files) - 1].path`)
	victim = strings.Trim(victim, `"`)
	app.call(`home.dialog = {"kind": "delete", "file": home.files[len(home.files) - 1]}`, nil)
	app.call(`home_dirty()`, nil)
	settle(t, app)
	saveShot(t, app, "home-delete")
	pressKey(t, app, "Escape", window.Mods{})
	if _, err := os.Stat(victim); err != nil {
		t.Fatal("cancel deleted the file")
	}
	app.call(`home.dialog = {"kind": "delete", "file": home.files[len(home.files) - 1]}`, nil)
	pressKey(t, app, "Enter", window.Mods{})
	if _, err := os.Stat(victim); err == nil {
		t.Error("the file is still there after confirming")
	}
	trash, _ := trashDir()
	if _, err := os.Stat(filepath.Join(trash, filepath.Base(victim))); err != nil {
		t.Errorf("the file is not in the Trash: %v", err)
	}

	// New drawing creates Untitled.lzdraw and opens it; the second is Untitled 2.
	app.call(`show_screen("home", null)`, nil)
	settle(t, app)
	click(t, app, 12+30, 46+50)
	if got := evalStr(t, app, `harness.screen`); got != "editor" {
		t.Errorf("screen = %s after + New drawing", got)
	}
	body, err := os.ReadFile(filepath.Join(dir, "Untitled.lzdraw"))
	if err != nil || !strings.HasPrefix(string(body), "-~> lzdraw 1\ncanvas = {") {
		t.Errorf("Untitled.lzdraw = %q, %v", body, err)
	}
	app.call(`show_screen("home", null)`, nil)
	pressKey(t, app, "KeyN", window.Mods{Cmd: true})
	if _, err := os.Stat(filepath.Join(dir, "Untitled 2.lzdraw")); err != nil {
		t.Errorf("Cmd+N did not make Untitled 2: %v", err)
	}
}

func TestHomeEmptyStateButtons(t *testing.T) {
	app, dir := newHomeApp(t, false, "en")
	// "Open an example" copies layers.lzdraw into the folder and opens it.
	x := evalInt(t, app, `home_empty_layout().example.x + 5`)
	y := evalInt(t, app, `home_empty_layout().example.y + 5`)
	click(t, app, x, y)
	want, _ := os.ReadFile(filepath.Join("..", "examples", "layers.lzdraw"))
	got, err := os.ReadFile(filepath.Join(dir, "layers.lzdraw"))
	if err != nil || string(got) != string(want) {
		t.Errorf("the example was not copied as is: %v", err)
	}
	if s := evalStr(t, app, `harness.screen`); s != "editor" {
		t.Errorf("screen = %s", s)
	}
	// The backdoor commands.
	app.call(`show_screen("home", null)`, nil)
	settle(t, app)
	if got := evalStr(t, app, `home_cmd_open_card([0]).opened`); got != "true" {
		t.Errorf("OpenCard 0 = %s", got)
	}
	if got := evalStr(t, app, `home_cmd_open_card([9]).opened`); got != "false" {
		t.Errorf("OpenCard 9 = %s, want false", got)
	}
}
