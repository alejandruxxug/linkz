package host

import (
	"errors"
	"image"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"linkz/host/window"
	"linkz/lib"
)

func newTestSplash(t *testing.T) *Splash {
	t.Helper()
	sp, err := NewSplash()
	if err != nil {
		t.Fatal(err)
	}
	return sp
}

func opaque(img *image.RGBA, x, y int) bool { return img.Pix[y*img.Stride+x*4+3] == 255 }

func TestSplashFrameLayers(t *testing.T) {
	sp := newTestSplash(t)
	if sp.W != 200 || sp.H != 80 || sp.Lizards != 6 {
		t.Fatalf("banner %dx%d with %d lizard frames", sp.W, sp.H, sp.Lizards)
	}
	img := sp.Frame(0, -1, -1, false)
	if img.Rect.Dx() != 600 || img.Rect.Dy() != 240 {
		t.Fatalf("window image is %v, want 600x240", img.Rect)
	}
	// The banner is transparent where nothing is drawn: all four corners.
	for _, p := range [][2]int{{0, 0}, {599, 0}, {0, 239}, {599, 239}} {
		if img.Pix[p[1]*img.Stride+p[0]*4+3] != 0 {
			t.Errorf("corner %v is not transparent", p)
		}
	}
	// The logo sits at (48, 14) and is opaque (outline on its first row).
	if !opaque(img, (48+3)*SplashScale, 14*SplashScale) {
		t.Error("no logo pixel at the top of the logo's first letter")
	}
	// The tagline is drawn under the logo: some opaque pixel in its band.
	found := false
	for y := 56 * SplashScale; y < 66*SplashScale && !found; y++ {
		for x := 0; x < img.Rect.Dx(); x++ {
			if opaque(img, x, y) {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("no tagline pixels under the logo")
	}
}

func TestSplashLizardWalksAndWraps(t *testing.T) {
	sp := newTestSplash(t)
	same := func(a, b *image.RGBA) bool { return string(a.Pix) == string(b.Pix) }
	f0, f1 := sp.Frame(0, -1, -1, false), sp.Frame(1, -1, -1, false)
	if same(f0, f1) {
		t.Error("frame 0 and frame 1 are identical: the lizard doesn't move")
	}
	// After one full trip across the banner the lizard is back at the start
	// (6 frames of walk cycle divide the period).
	n := sp.period / sp.Step
	if n%sp.Lizards == 0 && !same(sp.Frame(0, -1, -1, false), sp.Frame(n, -1, -1, false)) {
		t.Error("the lizard does not wrap around")
	}
}

func TestSplashError(t *testing.T) {
	sp := newTestSplash(t)
	before := sp.Frame(0, -1, -1, false)
	sp.ShowError("boom: something went wrong with a very long message indeed")
	after := sp.Frame(0, -1, -1, true)
	if string(before.Pix) == string(after.Pix) {
		t.Fatal("the error line was not drawn")
	}
	// Ember is #d94a3a.
	red := false
	for i := 0; i < len(after.Pix); i += 4 {
		if after.Pix[i] == 0xd9 && after.Pix[i+1] == 0x4a && after.Pix[i+2] == 0x3a && after.Pix[i+3] == 255 {
			red = true
			break
		}
	}
	if !red {
		t.Error("no ember pixels in the error line")
	}
}

// fakeLauncher runs the launcher against a clock the test moves, with a
// loader that finishes when the test says so.
type fakeLauncher struct {
	*Launcher
	now     time.Time
	release chan error
}

func newFakeLauncher(t *testing.T) *fakeLauncher {
	sp := newTestSplash(t)
	f := &fakeLauncher{now: time.Unix(1000, 0), release: make(chan error, 1)}
	f.Launcher = NewLauncher(sp, NewApp(100, 100), "")
	f.Now = func() time.Time { return f.now }
	f.Load = func(progress func(float64)) error {
		progress(0.5)
		return <-f.release
	}
	return f
}

// tick advances the clock by d and runs one tick.
func (f *fakeLauncher) tick(d time.Duration, evs ...window.Event) window.Output {
	f.now = f.now.Add(d)
	return f.Tick(window.Input{W: 600, H: 240, Events: evs})
}

// settle waits for the loader goroutine's result to be seen by the launcher.
func (f *fakeLauncher) settle() {
	for i := 0; i < 200 && !f.done; i++ {
		f.tick(0)
		time.Sleep(time.Millisecond)
	}
}

func TestLauncherWaitsForLoadAndMinimum(t *testing.T) {
	f := newFakeLauncher(t)
	if out := f.tick(0); out.Frame == nil {
		t.Fatal("no first frame")
	}
	for i := 0; i < 30; i++ { // 2.5 s of 12 fps ticks while loading
		if out := f.tick(80 * time.Millisecond); out.Chrome != nil {
			t.Fatal("the editor appeared before loading finished")
		}
	}
	if f.Progress() != 0.5 {
		t.Errorf("progress = %v, want 0.5", f.Progress())
	}
	f.release <- nil
	f.settle()
	var out window.Output
	for i := 0; i < 100 && out.Chrome == nil; i++ {
		out = f.tick(40 * time.Millisecond)
	}
	if out.Chrome == nil || !out.Chrome.Decorated || out.Chrome.Size != [2]int{1024, 720} || !out.Chrome.Center {
		t.Fatalf("chrome = %+v, want the decorated 1024x720 editor window", out.Chrome)
	}
	if f.stage != max(f.Splash.Vines-1, 0) {
		t.Errorf("editor appeared at vine stage %d of %d", f.stage, f.Splash.Vines)
	}
	if f.Splash.Bloom > 0 && f.bloom != f.Splash.Bloom-1 {
		t.Errorf("editor appeared at bloom frame %d of %d", f.bloom, f.Splash.Bloom)
	}
}

func TestLauncherMinimumTime(t *testing.T) {
	f := newFakeLauncher(t)
	f.release <- nil // loading is instant
	f.tick(0)
	f.settle()
	for i := 0; i < 20; i++ { // 20 x 40 ms = 0.8 s
		if out := f.tick(40 * time.Millisecond); out.Chrome != nil {
			t.Fatalf("the editor appeared after %d ms, before the 1.2 s minimum", (i+1)*40)
		}
	}
	for i := 0; i < 100; i++ {
		if out := f.tick(40 * time.Millisecond); out.Chrome != nil {
			return
		}
	}
	t.Error("the editor never appeared")
}

func TestLauncherEscapeQuits(t *testing.T) {
	f := newFakeLauncher(t)
	defer func() { f.release <- nil }()
	if out := f.tick(0, window.Event{Type: "down"}); out.Quit {
		t.Error("a click on the splash quit")
	}
	if out := f.tick(0, window.Event{Type: "key", Key: "Escape"}); !out.Quit {
		t.Error("Escape did not quit")
	}
}

func TestLauncherLoadError(t *testing.T) {
	f := newFakeLauncher(t)
	f.release <- errors.New("ui/x.liz: boom")
	f.tick(0)
	f.settle()
	if f.Err == nil {
		t.Fatal("the load error was not kept")
	}
	for i := 0; i < 50; i++ {
		if out := f.tick(100 * time.Millisecond); out.Chrome != nil {
			t.Fatal("the editor appeared after a load error")
		}
	}
	if out := f.tick(0, window.Event{Type: "down"}); !out.Quit {
		t.Error("a click did not close the splash that shows an error")
	}
}

func TestLauncherProgressNeverGoesBack(t *testing.T) {
	l := NewLauncher(newTestSplash(t), NewApp(100, 100), "")
	l.setProgress(0.6)
	l.setProgress(0.3)
	if l.Progress() != 0.6 {
		t.Errorf("progress = %v, want 0.6", l.Progress())
	}
}

// withOverlays gives the splash 6 vine stages and 4 bloom frames of empty
// layers, whatever the art has, so the sequencing can be tested.
func withOverlays(sp *Splash) {
	sp.Vines, sp.Bloom = 6, 4
	sp.vine, sp.bloom = nil, nil
	for i := 0; i < 6; i++ {
		sp.vine = append(sp.vine, NewBuf(sp.W, sp.H))
	}
	for i := 0; i < 4; i++ {
		sp.bloom = append(sp.bloom, NewBuf(sp.W, sp.H))
	}
}

func TestLauncherVinesGrowThenBloomThenEditor(t *testing.T) {
	f := newFakeLauncher(t)
	withOverlays(f.Splash)
	frame := time.Second / time.Duration(f.Splash.FPS)
	f.release <- nil // done at once: progress asks for the last stage
	f.tick(0)
	f.settle()
	prev, bloomStarted := f.stage, -1
	var out window.Output
	for i := 0; i < 400 && out.Chrome == nil; i++ {
		out = f.tick(frame)
		if f.stage-prev > 1 {
			t.Fatalf("stage jumped from %d to %d", prev, f.stage)
		}
		if f.bloom >= 0 && bloomStarted < 0 {
			bloomStarted = f.stage
		}
		prev = f.stage
	}
	if bloomStarted != 5 {
		t.Errorf("bloom started at vine stage %d, want the last (5)", bloomStarted)
	}
	if out.Chrome == nil {
		t.Fatal("the editor never appeared")
	}
	if f.bloom != 3 {
		t.Errorf("editor appeared at bloom frame %d, want the last (3)", f.bloom)
	}
	if held := f.now.Sub(f.finalAt); held < HoldBloom {
		t.Errorf("final bloom frame held %v, want at least %v", held, HoldBloom)
	}
}

func TestUIFilesOrder(t *testing.T) {
	files := lib.UIFiles()
	idx := map[string]int{}
	for i, f := range files {
		idx[f] = i
	}
	for _, f := range []string{"draw/geometry.liz", "draw/shapes.liz", "gfx/raster.liz", "ui/app.liz"} {
		if _, ok := idx[f]; !ok {
			t.Fatalf("%s missing from UIFiles %v", f, files)
		}
	}
	if !(idx["draw/shapes.liz"] < idx["gfx/raster.liz"]) || !(idx["draw/geometry.liz"] < idx["gfx/raster.liz"]) {
		t.Error("lib/draw must load before lib/gfx")
	}
	if idx["ui/app.liz"] != len(files)-1 {
		t.Error("app.liz must load last")
	}
	for _, f := range files {
		if f == "ui/splash.liz" || f == "ui/splash_art.liz" {
			t.Errorf("%s belongs to the splash's own environment", f)
		}
		if _, err := lib.ReadUI(f); err != nil {
			t.Errorf("ReadUI(%s): %v", f, err)
		}
	}
	if _, err := lib.Files.ReadFile("gfx/tests/text.liz"); err == nil {
		t.Error("gfx/tests must not be embedded")
	}
}

func TestBackdoorFilePerPid(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dir, _ := ConfigDir()
	stale := filepath.Join(dir, "debug-999999.json") // no such process
	if err := os.WriteFile(stale, []byte(`{"port":1,"token":"x","pid":999999}`), 0o600); err != nil {
		t.Fatal(err)
	}
	d := NewDebug(NewApp(100, 100))
	if err := d.Start(); err != nil {
		t.Fatal(err)
	}
	defer d.Stop()
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Error("the debug file of a dead pid was not cleaned")
	}
	if filepath.Base(d.file) != "debug-"+strconv.Itoa(os.Getpid())+".json" {
		t.Errorf("debug file is %s", d.file)
	}
	info, err := FindBackdoor(0)
	if err != nil || info.PID != os.Getpid() || info.Token != d.Token {
		t.Errorf("FindBackdoor(0) = %+v, %v", info, err)
	}
	if _, err := FindBackdoor(os.Getpid() + 1); err == nil {
		t.Error("FindBackdoor for another pid found something")
	}
	d.Stop()
	if _, err := FindBackdoor(0); err == nil {
		t.Error("a stopped window is still found")
	}
}
