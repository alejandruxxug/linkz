package host

import (
	"fmt"
	"image"
	"os"
	"time"

	"linkz/host/window"
	"linkz/lib"
	"lizardscript/ast"
	"lizardscript/evaluator"
	"lizardscript/lexer"
	"lizardscript/object"
	"lizardscript/parser"
)

// Stats are the timings the backdoor's Stats command reports.
type Stats struct {
	PaintMS   float64
	EventMS   float64
	RunMS     float64
	RunErrors int
	Frames    int
	Events    int
	Coalesced int
}

// App runs the lizard-script UI: it loads lib/ui, feeds it events and asks it
// to paint. All lizard code runs on the goroutine that calls Tick; other
// goroutines (the backdoor) hand work over with Do.
type App struct {
	FB    *FB
	Src   *Src
	Stats Stats

	env     *object.Environment
	app     *object.Hash
	progs   map[string]*ast.Program
	lastErr string
	started bool
	path    string
	quit    bool

	jobs     chan job
	injected []window.Event
	waiters  []*waiter
	resizeTo [2]int
	forcePnt bool
	headless bool // no window: resize requests apply directly
}

// job is work for the UI goroutine. run returns a function that makes the
// reply; it is called after the tick's paint, so replies show the result.
type job struct {
	run   func() func() any
	reply chan any
}

// waiter is a Wait command: it replies at the first idle tick.
type waiter struct {
	reply chan any
	ticks int
}

// NewApp makes an app with a window of w by h pixels.
func NewApp(w, h int) *App {
	a := &App{FB: NewFB(w, h), progs: map[string]*ast.Program{}, jobs: make(chan job, 64)}
	a.Src = NewSrc(&a.Stats)
	a.forcePnt = true
	return a
}

// Load builds the UI environment from the embedded libraries and starts the
// app on path.
func (a *App) Load(path string) error {
	evaluator.Output = os.Stderr
	a.env = object.NewEnvironment()
	a.env.Define("fb", a.FB.Builtins())
	a.env.Define("src", a.Src.Builtins())
	a.env.Define("fs", fsHash())
	a.env.Define("time", timeHash())
	a.env.Define("clip", clipHash())
	a.env.Define("os", osHash())

	for _, name := range lib.UIFiles() {
		code, err := lib.Files.ReadFile(name)
		if err != nil {
			return err
		}
		if err := a.exec(name, string(code)); err != nil {
			return err
		}
	}
	v, ok := a.env.Get("app")
	h, isHash := v.(*object.Hash)
	if !ok || !isHash {
		return fmt.Errorf("the UI libraries must define a hash called app")
	}
	a.app, a.path = h, path
	return a.call("app.start(__path)", map[string]object.Object{"__path": newStr(path)})
}

// exec runs library code in the UI environment.
func (a *App) exec(name, code string) error {
	p := parser.New(lexer.New(code))
	prog := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		return fmt.Errorf("%s: %s", name, errs[0])
	}
	if e, ok := evaluator.Eval(prog, a.env).(*object.Error); ok {
		return fmt.Errorf("%s: %s", name, e.Message)
	}
	return nil
}

// eval runs a snippet in the UI environment with some names bound.
func (a *App) eval(src string, bind map[string]object.Object) (object.Object, error) {
	prog, ok := a.progs[src]
	if !ok {
		p := parser.New(lexer.New(src))
		prog = p.ParseProgram()
		if len(p.Errors()) > 0 {
			return nil, fmt.Errorf("%v", p.Errors())
		}
		a.progs[src] = prog
	}
	for k, v := range bind {
		a.env.Define(k, v)
	}
	evaluator.Output = os.Stderr
	res := evaluator.Eval(prog, a.env)
	if e, ok := res.(*object.Error); ok {
		return nil, fmt.Errorf("%s", e.Message)
	}
	return res, nil
}

func (a *App) call(src string, bind map[string]object.Object) error {
	_, err := a.eval(src, bind)
	return err
}

// fail remembers a UI error and prints it; the last good frame stays.
func (a *App) fail(err error) {
	a.lastErr = err.Error()
	fmt.Fprintln(os.Stderr, "ui error:", a.lastErr)
}

// Do runs f on the UI goroutine and waits for the reply f's result makes.
func (a *App) Do(f func() func() any) any {
	j := job{run: f, reply: make(chan any, 1)}
	a.jobs <- j
	return <-j.reply
}

// Inject queues events as if the user made them.
func (a *App) Inject(evs ...window.Event) {
	for i := range evs {
		evs[i].Injected = true
	}
	a.injected = append(a.injected, evs...)
}

// WantsHover is read by the window: should it send moves with no button?
func (a *App) WantsHover() bool {
	if a.app == nil {
		return false
	}
	return truthy(hashGet(a.app, "wants_hover"))
}

func (a *App) needsPaint() bool {
	return a.forcePnt || (a.app != nil && truthy(hashGet(a.app, "needs_paint")))
}

// coalesce collapses consecutive moves from real input to the latest one.
// Injected events are never dropped.
func (a *App) coalesce(evs []window.Event) []window.Event {
	var out []window.Event
	for _, e := range evs {
		if n := len(out); n > 0 && e.Type == "move" && !e.Injected && out[n-1].Type == "move" && !out[n-1].Injected {
			out[n-1] = e
			a.Stats.Coalesced++
			continue
		}
		out = append(out, e)
	}
	return out
}

// Tick is one round of the loop: jobs, events, paint, replies. It implements
// window.Handler.
func (a *App) Tick(in window.Input) window.Output {
	busy := false

	if in.W > 0 && in.H > 0 {
		if w, h := a.FB.WindowSize(); w != in.W || h != in.H {
			a.FB.Resize(in.W, in.H)
			a.forcePnt = true
			uw, uh := a.FB.UISize()
			a.deliver(window.Event{Type: "resize", X: uw, Y: uh})
			busy = true
		}
	}

	var replies []func()
	for drained := false; !drained; {
		select {
		case j := <-a.jobs:
			thunk := j.run()
			reply := j.reply
			replies = append(replies, func() { reply <- thunk() })
			busy = true
		default:
			drained = true
		}
	}

	if a.headless && a.resizeTo[0] > 0 {
		a.FB.Resize(a.resizeTo[0], a.resizeTo[1])
		a.resizeTo = [2]int{}
		a.forcePnt = true
		uw, uh := a.FB.UISize()
		a.deliver(window.Event{Type: "resize", X: uw, Y: uh})
		busy = true
	}

	evs := a.coalesce(append(append([]window.Event{}, in.Events...), a.injected...))
	a.injected = nil
	for _, e := range evs {
		a.deliver(e)
		busy = true
	}

	painted := true
	if a.needsPaint() {
		busy = true
		start := time.Now()
		if err := a.call("app.paint()", nil); err != nil {
			a.fail(err)
			painted = false
		} else {
			a.lastErr = ""
			hashSet(a.app, "needs_paint", newBool(false))
			a.Stats.Frames++
		}
		a.Stats.PaintMS = float64(time.Since(start).Microseconds()) / 1000
		a.forcePnt = !painted
	}

	var out window.Output
	if painted && busy {
		img, dirty := a.FB.Compose()
		if !dirty.Empty() {
			out.Frame = &window.Frame{Img: img, Dirty: dirty}
		}
	}
	if a.resizeTo[0] > 0 {
		out.ResizeTo, a.resizeTo = a.resizeTo, [2]int{}
	}
	out.Quit = a.quit

	for _, r := range replies {
		r()
	}
	a.settleWaiters(busy)
	return out
}

// settleWaiters answers Wait commands at the first tick with nothing to do.
func (a *App) settleWaiters(busy bool) {
	keep := a.waiters[:0]
	for _, w := range a.waiters {
		w.ticks++
		if (!busy && w.ticks > 1) || w.ticks > 90 {
			w.reply <- map[string]any{"idle": !busy, "ticks": w.ticks}
		} else {
			keep = append(keep, w)
		}
	}
	a.waiters = keep
}

// deliver sends one event to app.event.
func (a *App) deliver(e window.Event) {
	if a.app == nil {
		return
	}
	start := time.Now()
	a.Stats.Events++
	if err := a.call("app.event(__e)", map[string]object.Object{"__e": a.eventHash(e)}); err != nil {
		a.fail(err)
	}
	a.Stats.EventMS = float64(time.Since(start).Microseconds()) / 1000
}

// eventHash turns an event into the hash lizard sees. pos is in UI pixels;
// canvas_pos is full resolution, relative to the canvas, when the pointer is
// inside it.
func (a *App) eventHash(e window.Event) *object.Hash {
	h := object.NewHash()
	hashSet(h, "type", newStr(e.Type))
	switch e.Type {
	case "down", "up", "move", "wheel":
		hashSet(h, "pos", newPoint(e.X/2, e.Y/2))
		if r, ok := a.FB.CanvasRect(); ok {
			cx, cy := e.X-r.Min.X*2, e.Y-r.Min.Y*2
			if cx >= 0 && cy >= 0 && cx < r.Dx()*2 && cy < r.Dy()*2 {
				hashSet(h, "canvas_pos", newPoint(cx, cy))
			}
		}
		hashSet(h, "button", newInt(e.Button))
		if e.Type == "wheel" {
			hashSet(h, "dy", newFloat(e.DY))
		}
	case "key":
		hashSet(h, "key", newStr(e.Key))
	case "text":
		hashSet(h, "text", newStr(e.Text))
	case "resize":
		hashSet(h, "size", newPoint(e.X, e.Y))
	}
	mods := object.NewHash()
	hashSet(mods, "shift", newBool(e.Mods.Shift))
	hashSet(mods, "alt", newBool(e.Mods.Alt))
	hashSet(mods, "cmd", newBool(e.Mods.Cmd))
	hashSet(mods, "ctrl", newBool(e.Mods.Ctrl))
	hashSet(h, "mods", mods)

	pen := e.Pen
	if pen.Source == "" {
		pen = window.Pen{Pressure: 1, Source: "mouse"}
	}
	ph := object.NewHash()
	hashSet(ph, "pressure", newFloat(pen.Pressure))
	hashSet(ph, "tilt", &object.Point{Coords: []object.Object{newFloat(pen.Tilt[0]), newFloat(pen.Tilt[1])}})
	hashSet(ph, "rotation", newFloat(pen.Rotation))
	hashSet(ph, "source", newStr(pen.Source))
	hashSet(ph, "eraser", newBool(pen.Eraser))
	hashSet(h, "pen", ph)
	return h
}

// Screenshot is the window as shown.
func (a *App) Screenshot() *image.RGBA { return a.FB.Screenshot() }
