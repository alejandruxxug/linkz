package host

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"linkz/host/pressure"
	"linkz/host/window"
	"lizardscript/object"
)

// The backdoor: an HTTP server on 127.0.0.1 that lets an agent drive and
// inspect the running editor. Input commands go through the same event path
// as real input and reply after the resulting paint.

// TokenHeader carries the secret from ~/.linkz/debug.json.
const TokenHeader = "X-Linkz-Token"

// Debug serves the backdoor for an App.
type Debug struct {
	App    *App
	Token  string
	srv    *http.Server
	ln     net.Listener
	file   string
	shots  int
	cmds   map[string]func(args []string) (any, error)
	onExit func()
}

// NewDebug makes a backdoor with a random token. It doesn't listen yet.
func NewDebug(app *App) *Debug {
	buf := make([]byte, 16)
	rand.Read(buf)
	d := &Debug{App: app, Token: hex.EncodeToString(buf)}
	d.cmds = map[string]func([]string) (any, error){
		"List": d.list, "Dump": d.dump, "Screenshot": d.screenshot,
		"Click": d.click, "DoubleClick": d.doubleClick, "Drag": d.drag, "Scroll": d.scroll,
		"Key": d.key, "Type": d.typeText,
		"Resize": d.resize, "Stats": d.stats, "Wait": d.wait, "Quit": d.quit,
		"Call": d.callCmd,
		// The file-state commands are the editor's to define in lizard.
		"SetCode": d.viaLizard("SetCode"), "Open": d.viaLizard("Open"),
		"Save": d.viaLizard("Save"), "Export": d.viaLizard("Export"),
	}
	return d
}

// Start listens on 127.0.0.1 on a free port and writes ~/.linkz/debug.json.
func (d *Debug) Start() error {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	d.ln = ln
	mux := http.NewServeMux()
	mux.HandleFunc("/cmd", d.serve)
	d.srv = &http.Server{Handler: mux}
	go d.srv.Serve(ln)

	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	d.file = filepath.Join(dir, "debug.json")
	info, _ := json.Marshal(map[string]any{"port": ln.Addr().(*net.TCPAddr).Port, "token": d.Token, "pid": os.Getpid()})
	return os.WriteFile(d.file, info, 0o600)
}

// Stop closes the server and removes debug.json.
func (d *Debug) Stop() {
	if d.srv != nil {
		d.srv.Close()
	}
	if d.file != "" {
		os.Remove(d.file)
	}
}

// Addr is where the server listens.
func (d *Debug) Addr() string { return d.ln.Addr().String() }

type request struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

func (d *Debug) serve(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(TokenHeader) != d.Token {
		http.Error(w, `{"error":"missing or wrong token"}`, http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST only"}`, http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(io.LimitReader(r.Body, 8<<20))
	var req request
	if err := json.Unmarshal(body, &req); err != nil {
		reply(w, nil, fmt.Errorf("bad request: %v", err))
		return
	}
	run, ok := d.cmds[req.Cmd]
	if !ok {
		reply(w, nil, fmt.Errorf("unknown command %q (try List)", req.Cmd))
		return
	}
	out, err := run(req.Args)
	reply(w, out, err)
}

func reply(w http.ResponseWriter, out any, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(out)
}

// ---- argument helpers ----

// split separates --flags from positional words.
func split(args []string) (pos []string, flags map[string]string) {
	flags = map[string]string{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if strings.HasPrefix(a, "--") {
			name, val, has := strings.Cut(a[2:], "=")
			if !has && i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") && name != "full" &&
				name != "shift" && name != "alt" && name != "cmd" && name != "ctrl" {
				i++
				val = args[i]
			}
			flags[name] = val
			continue
		}
		pos = append(pos, a)
	}
	return
}

func floats(pos []string, n int, usage string) ([]float64, error) {
	if len(pos) < n {
		return nil, fmt.Errorf("usage: %s", usage)
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		f, err := strconv.ParseFloat(pos[i], 64)
		if err != nil {
			return nil, fmt.Errorf("%q is not a number", pos[i])
		}
		out[i] = f
	}
	return out, nil
}

// modsFrom reads shift/alt/cmd/ctrl from words or flags.
func modsFrom(words []string, flags map[string]string) window.Mods {
	var m window.Mods
	set := func(w string) {
		switch strings.ToLower(w) {
		case "shift":
			m.Shift = true
		case "alt", "option":
			m.Alt = true
		case "cmd", "meta", "command":
			m.Cmd = true
		case "ctrl", "control":
			m.Ctrl = true
		}
	}
	for _, w := range words {
		for _, part := range strings.Split(w, "+") {
			set(part)
		}
	}
	for k := range flags {
		set(k)
	}
	return m
}

// coords converts to window pixels. Coordinates are UI pixels (half the
// screenshot's pixels, what events call pos); --full means window pixels.
func coords(x, y float64, flags map[string]string) (int, int) {
	if _, full := flags["full"]; full {
		return int(math.Round(x)), int(math.Round(y))
	}
	return int(math.Round(x * 2)), int(math.Round(y * 2))
}

// afterPaint runs inject on the UI goroutine and replies once the tick that
// handled the events has painted.
func (d *Debug) afterPaint(inject func(), result func() any) (any, error) {
	return d.App.Do(func() func() any {
		inject()
		return result
	}), nil
}

func (d *Debug) state() any {
	return map[string]any{"frames": d.App.Stats.Frames, "error": d.App.lastErr}
}

// ---- commands ----

func (d *Debug) list(args []string) (any, error) {
	host := make([]string, 0, len(d.cmds))
	for k := range d.cmds {
		host = append(host, k)
	}
	sort.Strings(host)
	return d.App.Do(func() func() any {
		var app []string
		if dc, ok := hashGet(d.App.app, "debug_commands").(*object.Hash); ok {
			dc.Each(func(p object.HashPair) { app = append(app, p.Key.Inspect()) })
		}
		sort.Strings(app)
		return func() any { return map[string]any{"host": host, "app": app} }
	}), nil
}

func (d *Debug) dump(args []string) (any, error) {
	return d.App.Do(func() func() any {
		return func() any {
			w, h := d.App.FB.WindowSize()
			uw, uh := d.App.FB.UISize()
			out := map[string]any{
				"window": map[string]any{"width": w, "height": h, "ui_width": uw, "ui_height": uh},
				"path":   d.App.path,
				"error":  d.App.lastErr,
				"stats":  d.statsMap(),
			}
			if r, ok := d.App.FB.CanvasRect(); ok {
				out["canvas"] = map[string]any{"x": r.Min.X, "y": r.Min.Y, "width": r.Dx(), "height": r.Dy()}
			}
			if res, err := d.App.eval("app.debug_state()", nil); err != nil {
				out["app_error"] = err.Error()
			} else {
				out["app"] = toGo(res, 0)
			}
			return out
		}
	}), nil
}

func (d *Debug) statsMap() map[string]any {
	s := d.App.Stats
	return map[string]any{
		"paint_ms": s.PaintMS, "event_ms": s.EventMS, "run_ms": s.RunMS, "run_errors": s.RunErrors,
		"frames_painted": s.Frames, "events": s.Events, "events_coalesced": s.Coalesced,
	}
}

func (d *Debug) stats(args []string) (any, error) {
	return d.App.Do(func() func() any { return func() any { return d.statsMap() } }), nil
}

func (d *Debug) screenshot(args []string) (any, error) {
	pos, _ := split(args)
	path := ""
	if len(pos) > 0 {
		path = pos[0]
	}
	if path == "" {
		d.shots++
		path = filepath.Join(os.TempDir(), fmt.Sprintf("linkz-shot-%d-%d.png", os.Getpid(), d.shots))
	}
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}
	var werr error
	var w, h int
	d.App.Do(func() func() any {
		return func() any {
			img := d.App.Screenshot()
			w, h = img.Rect.Dx(), img.Rect.Dy()
			werr = writePNG(path, img)
			return nil
		}
	})
	if werr != nil {
		return nil, werr
	}
	return map[string]any{"path": path, "width": w, "height": h}, nil
}

func (d *Debug) click(args []string) (any, error) { return d.clickN(args, 1) }

func (d *Debug) doubleClick(args []string) (any, error) { return d.clickN(args, 2) }

func (d *Debug) clickN(args []string, n int) (any, error) {
	pos, flags := split(args)
	f, err := floats(pos, 2, "Click x y [shift|alt|cmd|ctrl] [--full]")
	if err != nil {
		return nil, err
	}
	x, y := coords(f[0], f[1], flags)
	mods := modsFrom(pos[2:], flags)
	return d.afterPaint(func() {
		for i := 0; i < n; i++ {
			d.App.Inject(
				window.Event{Type: "down", X: x, Y: y, Mods: mods, Pen: window.Pen{Pressure: 1, Source: "mouse"}},
				window.Event{Type: "up", X: x, Y: y, Mods: mods, Pen: window.Pen{Pressure: 1, Source: "mouse"}})
		}
	}, d.state)
}

func (d *Debug) drag(args []string) (any, error) {
	pos, flags := split(args)
	f, err := floats(pos, 4, "Drag x1 y1 x2 y2 [steps] [--pressure 0.2,1.0] [--full]")
	if err != nil {
		return nil, err
	}
	steps := 10
	if len(pos) > 4 {
		if n, err := strconv.Atoi(pos[4]); err == nil && n > 0 {
			steps = n
		}
	}
	var ramp []float64
	if p := flags["pressure"]; p != "" {
		for _, part := range strings.Split(p, ",") {
			v, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return nil, fmt.Errorf("bad pressure %q", part)
			}
			ramp = append(ramp, v)
		}
	}
	mods := modsFrom(pos[min(5, len(pos)):], flags)
	penAt := func(t float64) window.Pen {
		if len(ramp) == 0 {
			return window.Pen{Pressure: 1, Source: "mouse"}
		}
		// piecewise linear through the ramp values
		p := ramp[0]
		if len(ramp) > 1 {
			seg := t * float64(len(ramp)-1)
			i := min(int(seg), len(ramp)-2)
			p = ramp[i] + (ramp[i+1]-ramp[i])*(seg-float64(i))
		}
		return window.Pen{Pressure: p, Source: "trackpad"}
	}
	return d.afterPaint(func() {
		x0, y0 := coords(f[0], f[1], flags)
		d.App.Inject(window.Event{Type: "down", X: x0, Y: y0, Mods: mods, Pen: penAt(0)})
		for i := 1; i <= steps; i++ {
			t := float64(i) / float64(steps)
			x, y := coords(f[0]+(f[2]-f[0])*t, f[1]+(f[3]-f[1])*t, flags)
			d.App.Inject(window.Event{Type: "move", X: x, Y: y, Mods: mods, Pen: penAt(t),
				Samples: []window.Sample{{X: x, Y: y, Pressure: penAt(t).Pressure, T: pressure.Now()}}})
		}
		x1, y1 := coords(f[2], f[3], flags)
		d.App.Inject(window.Event{Type: "up", X: x1, Y: y1, Mods: mods, Pen: penAt(1)})
	}, d.state)
}

func (d *Debug) scroll(args []string) (any, error) {
	pos, flags := split(args)
	f, err := floats(pos, 1, "Scroll dy [x y]")
	if err != nil {
		return nil, err
	}
	x, y := 0, 0
	if len(pos) >= 3 {
		if p, err := floats(pos[1:], 2, ""); err == nil {
			x, y = coords(p[0], p[1], flags)
		}
	}
	return d.afterPaint(func() {
		d.App.Inject(window.Event{Type: "wheel", X: x, Y: y, DY: f[0]})
	}, d.state)
}

func (d *Debug) key(args []string) (any, error) {
	pos, flags := split(args)
	if len(pos) < 1 {
		return nil, fmt.Errorf("usage: Key name [shift|alt|cmd|ctrl]")
	}
	mods := modsFrom(pos[1:], flags)
	return d.afterPaint(func() {
		d.App.Inject(window.Event{Type: "key", Key: pos[0], Mods: mods})
	}, d.state)
}

func (d *Debug) typeText(args []string) (any, error) {
	text := strings.Join(args, " ")
	return d.afterPaint(func() {
		for _, r := range text {
			d.App.Inject(window.Event{Type: "text", Text: string(r)})
		}
	}, d.state)
}

func (d *Debug) resize(args []string) (any, error) {
	pos, _ := split(args)
	f, err := floats(pos, 2, "Resize w h")
	if err != nil {
		return nil, err
	}
	w, h := int(f[0]), int(f[1])
	if w < 100 || h < 100 || w > 8000 || h > 8000 {
		return nil, fmt.Errorf("size must be between 100 and 8000")
	}
	d.App.Do(func() func() any {
		d.App.resizeTo = [2]int{w, h}
		return func() any { return nil }
	})
	return d.wait(nil)
}

func (d *Debug) wait(args []string) (any, error) {
	ch := make(chan any, 1)
	d.App.Do(func() func() any {
		d.App.waiters = append(d.App.waiters, &waiter{reply: ch})
		return func() any { return nil }
	})
	select {
	case v := <-ch:
		return v, nil
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("timed out waiting for the app to go idle")
	}
}

func (d *Debug) quit(args []string) (any, error) {
	return d.App.Do(func() func() any {
		d.App.quit = true
		return func() any { return map[string]any{"quit": true} }
	}), nil
}

// callCmd runs app.debug_commands[Name](args).
func (d *Debug) callCmd(args []string) (any, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("usage: Call Name args...")
	}
	return d.callLizard(args[0], args[1:])
}

func (d *Debug) viaLizard(name string) func([]string) (any, error) {
	return func(args []string) (any, error) { return d.callLizard(name, args) }
}

// parseArg reads a number if the word is one, else keeps the text.
func parseArg(s string) object.Object {
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return &object.Integer{Value: n}
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil && !math.IsInf(f, 0) && !math.IsNaN(f) {
		return newFloat(f)
	}
	return newStr(s)
}

func (d *Debug) callLizard(name string, args []string) (any, error) {
	var cerr error
	var result any
	d.App.Do(func() func() any {
		list := &object.Array{}
		for _, a := range args {
			list.Elements = append(list.Elements, parseArg(a))
		}
		cmds, _ := hashGet(d.App.app, "debug_commands").(*object.Hash)
		var fn object.Object
		if cmds != nil {
			fn = hashGet(cmds, name)
		}
		if fn == nil {
			cerr = fmt.Errorf("the app has no command %q (try List)", name)
			return func() any { return nil }
		}
		res, err := d.App.eval("app.debug_commands[__name](__args)", map[string]object.Object{
			"__name": newStr(name), "__args": list})
		if err != nil {
			d.App.fail(err)
			cerr = err
			return func() any { return nil }
		}
		return func() any { result = toGo(res, 0); return nil }
	})
	return result, cerr
}
