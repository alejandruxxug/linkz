package host

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"linkz/host/window"
)

func TestWriteFileAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.lzdraw")
	if err := os.WriteFile(path, []byte("old"), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := WriteFileAtomic(path, []byte("new contents")); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != "new contents" {
		t.Errorf("content = %q", got)
	}
	if info, _ := os.Stat(path); info.Mode().Perm() != 0o640 {
		t.Errorf("mode = %v, want the old 0640 kept", info.Mode().Perm())
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Errorf("temp files left behind: %v", entries)
	}
	// A failing write leaves the old file alone and no temp behind.
	if err := WriteFileAtomic(filepath.Join(dir, "missing", "b"), []byte("x")); err == nil {
		t.Error("writing into a missing folder should fail")
	}
	if got, _ := os.ReadFile(path); string(got) != "new contents" {
		t.Error("the file changed after an unrelated failure")
	}
}

func TestFSHash(t *testing.T) {
	h := fsHash()
	path := filepath.Join(t.TempDir(), "x.txt")
	call(t, h, "write", newStr(path), newStr("héllo"))
	if got := call(t, h, "read", newStr(path)).Inspect(); got != "héllo" {
		t.Errorf("read = %q", got)
	}
	if call(t, h, "exists", newStr(path)).Inspect() != "true" || call(t, h, "exists", newStr(path+"no")).Inspect() != "false" {
		t.Error("exists is wrong")
	}
	if call(t, h, "read", newStr(path+"no")).Type() != "ERROR" {
		t.Error("reading a missing file should be an error")
	}
}

func TestLocaleAndConfig(t *testing.T) {
	for in, want := range map[string]string{"es_CO.UTF-8": "es", "en-US": "en", "EN": "en", "C": "", "": "", "eng": "", "pt_BR": "pt"} {
		if got := twoLetters(in); got != want {
			t.Errorf("twoLetters(%q) = %q, want %q", in, got, want)
		}
	}
	if l := Locale(); len(l) != 2 || strings.ToLower(l) != l {
		t.Errorf("Locale() = %q, want two lowercase letters", l)
	}
	t.Setenv("HOME", t.TempDir())
	dir, err := ConfigDir()
	if err != nil || filepath.Base(dir) != ".linkz" {
		t.Fatalf("ConfigDir = %q, %v", dir, err)
	}
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		t.Error("ConfigDir did not create the folder")
	}
}

// newTestApp starts the UI headless and pumps its loop on a goroutine, like
// the window would.
func newTestApp(t *testing.T) (*App, func()) {
	t.Helper()
	app := NewApp(400, 300)
	app.headless = true
	if err := app.Load("test.lzdraw"); err != nil {
		t.Fatal(err)
	}
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-stop:
				return
			default:
			}
			app.Tick(window.Input{})
			time.Sleep(2 * time.Millisecond)
		}
	}()
	return app, func() { close(stop); <-done }
}

func TestAppPlaceholderEvents(t *testing.T) {
	app, stop := newTestApp(t)
	defer stop()
	dbg := NewDebug(app)

	mustCall := func(cmd string, args ...string) map[string]any {
		t.Helper()
		out, err := dbg.cmds[cmd](args)
		if err != nil {
			t.Fatalf("%s: %v", cmd, err)
		}
		b, _ := json.Marshal(out)
		var m map[string]any
		json.Unmarshal(b, &m)
		return m
	}

	mustCall("Wait")
	// A press inside the UI: pos is in UI pixels (window / 2), x*2 in the window.
	mustCall("Click", "50", "60")
	d := mustCall("Dump")
	state := d["app"].(map[string]any)
	if mouse := state["mouse"].([]any); mouse[0].(float64) != 50 || mouse[1].(float64) != 60 {
		t.Errorf("mouse = %v, want [50 60] (UI pixels)", mouse)
	}
	if state["held"] != false {
		t.Errorf("held = %v after a click", state["held"])
	}

	// A click inside the canvas rect carries canvas_pos at full resolution.
	mustCall("Click", "100.5", "20")
	state = mustCall("Dump")["app"].(map[string]any)
	cp, _ := state["canvas_pos"].([]any)
	if len(cp) != 2 || cp[0].(float64) != 21 || cp[1].(float64) != 24 { // canvas at UI (90,8): (201-180, 40-16)
		t.Errorf("canvas_pos = %v, want [21 24]", state["canvas_pos"])
	}
	// ...and not outside it.
	mustCall("Click", "10", "10")
	if state = mustCall("Dump")["app"].(map[string]any); state["canvas_pos"] != nil {
		t.Errorf("canvas_pos outside the canvas = %v", state["canvas_pos"])
	}

	// Drag injects every step (no coalescing), with the pressure ramp.
	before := app.Stats.Events
	mustCall("Drag", "10", "10", "50", "50", "8", "--pressure", "0.2,1.0")
	if n := app.Stats.Events - before; n != 10 {
		t.Errorf("drag delivered %d events, want 10 (down + 8 moves + up)", n)
	}
	if app.Stats.Coalesced != 0 {
		t.Errorf("injected moves were coalesced: %d", app.Stats.Coalesced)
	}

	// Lizard commands and the command list.
	out := mustCall("Call", "Ping", "7")
	if out["pong"] != true {
		t.Errorf("Ping = %v", out)
	}
	l := mustCall("List")
	if names := l["app"].([]any); len(names) != 1 || names[0] != "Ping" {
		t.Errorf("app commands = %v", l["app"])
	}
	if _, err := dbg.cmds["Call"]([]string{"Nope"}); err == nil {
		t.Error("calling an unknown command should fail")
	}
	if _, err := dbg.cmds["Save"](nil); err == nil {
		t.Error("the placeholder has no Save: expected an error")
	}

	// Screenshot is the composited window (full resolution).
	shot := filepath.Join(t.TempDir(), "s.png")
	r := mustCall("Screenshot", shot)
	if r["width"].(float64) != 400 || r["height"].(float64) != 300 {
		t.Errorf("screenshot size = %v", r)
	}
	if info, err := os.Stat(shot); err != nil || info.Size() == 0 {
		t.Error("screenshot file missing")
	}

	// Resize reallocates the buffers, and lizard repaints.
	frames := app.Stats.Frames
	mustCall("Resize", "600", "400")
	if w, h := app.FB.WindowSize(); w != 600 || h != 400 {
		t.Errorf("window size %dx%d after Resize", w, h)
	}
	if app.Stats.Frames == frames {
		t.Error("no repaint after the resize")
	}
}

func TestCoalesceMoves(t *testing.T) {
	a := NewApp(100, 100)
	mv := func(x int, inj bool) window.Event { return window.Event{Type: "move", X: x, Injected: inj} }
	got := a.coalesce([]window.Event{
		{Type: "down"}, mv(1, false), mv(2, false), mv(3, false), {Type: "up"},
		mv(4, true), mv(5, true), mv(6, false), mv(7, false),
	})
	var xs []int
	for _, e := range got {
		xs = append(xs, e.X)
	}
	want := []int{0, 3, 0, 4, 5, 7}
	if len(xs) != len(want) {
		t.Fatalf("got %v, want %v", xs, want)
	}
	for i := range want {
		if xs[i] != want[i] {
			t.Fatalf("got %v, want %v", xs, want)
		}
	}
	if a.Stats.Coalesced != 3 {
		t.Errorf("coalesced = %d, want 3", a.Stats.Coalesced)
	}
}

func TestBackdoorRejectsBadToken(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	app, stop := newTestApp(t)
	defer stop()
	d := NewDebug(app)
	if err := d.Start(); err != nil {
		t.Fatal(err)
	}
	defer d.Stop()

	if info, err := os.Stat(d.file); err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("debug.json: %v, mode %v, want 0600", err, info)
	}
	if !strings.HasPrefix(d.Addr(), "127.0.0.1:") {
		t.Errorf("listening on %s, want 127.0.0.1", d.Addr())
	}

	post := func(token string) (int, string) {
		req, _ := http.NewRequest("POST", "http://"+d.Addr()+"/cmd", bytes.NewBufferString(`{"cmd":"Stats"}`))
		if token != "" {
			req.Header.Set(TokenHeader, token)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, string(b)
	}
	if code, _ := post(""); code != http.StatusForbidden {
		t.Errorf("no token: status %d, want 403", code)
	}
	if code, _ := post("wrong"); code != http.StatusForbidden {
		t.Errorf("wrong token: status %d, want 403", code)
	}
	if code, body := post(d.Token); code != http.StatusOK || !strings.Contains(body, "frames_painted") {
		t.Errorf("right token: status %d body %s", code, body)
	}
	d.Stop()
	if _, err := os.Stat(d.file); !os.IsNotExist(err) {
		t.Error("debug.json was not removed on stop")
	}
}

func TestPaintHeadless(t *testing.T) {
	out := filepath.Join(t.TempDir(), "o.png")
	code := `fb.fill((0, 0), (8, 8), #ff0000)
fb.span(2, 1.5, 3.5, #0000ff)
puts("done")
`
	var printed bytes.Buffer
	if err := Paint(code, out, 8, 8, &printed); err != nil {
		t.Fatal(err)
	}
	if printed.String() != "done\n" {
		t.Errorf("printed %q", printed.String())
	}
	if info, err := os.Stat(out); err != nil || info.Size() == 0 {
		t.Fatal("no PNG written")
	}
	if err := Paint("nope", out, 8, 8, &printed); err == nil {
		t.Error("a failing program should report an error")
	}
}
