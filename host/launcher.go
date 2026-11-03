package host

import (
	"fmt"
	"os"
	"sync"
	"time"

	"linkz/host/window"
)

// Splash timing. The editor never appears before MinSplash, and waits a
// moment on the finished picture so the bloom can be seen.
const (
	MinSplash  = 1200 * time.Millisecond
	HoldVines  = 300 * time.Millisecond // last vine stage, when there is no bloom
	HoldBloom  = 400 * time.Millisecond // last bloom frame
	editorW    = 1024
	editorH    = 720
	switchWait = 30 // ticks to wait for the window to take its new size
)

type launchPhase int

const (
	phaseSplash launchPhase = iota
	phaseSwitch
	phaseEditor
)

// Launcher is the window handler for `linkz edit` with a splash screen. While
// the editor libraries load on a goroutine, it cycles the banner frames lizard
// drew; when loading is done and the picture is finished, it turns the window
// into the editor's and hands every tick to the App.
//
// The interpreter runs one program at a time: the loading goroutine is the
// only thing evaluating until it reports back, and Tick touches the App only
// after that.
type Launcher struct {
	Splash  *Splash
	App     *App
	Path    string
	OnReady func()                             // called once the editor is up (starts the backdoor)
	Now     func() time.Time                   // the clock; tests replace it
	Err     error                              // why loading failed, if it did
	Load    func(progress func(float64)) error // loads the editor; tests replace it

	phase     launchPhase
	started   bool
	start     time.Time
	loaded    chan error
	mu        sync.Mutex
	progress  float64
	done      bool
	n         int // animation frame
	stage     int // vine stage shown
	bloom     int // bloom frame shown, -1 before it starts
	finalAt   time.Time
	loadErr   string
	switching int
}

// NewLauncher makes a launcher; loading starts at the first Tick.
func NewLauncher(sp *Splash, app *App, path string) *Launcher {
	l := &Launcher{Splash: sp, App: app, Path: path, Now: time.Now, bloom: -1, n: -1}
	l.Load = func(progress func(float64)) error { return app.LoadWith(path, progress) }
	return l
}

// WantsHover forwards the app's answer once the editor is up.
func (l *Launcher) WantsHover() bool {
	return l.phase == phaseEditor && l.App.WantsHover()
}

// Progress is how much of the loading is done, from 0 to 1; it never goes back.
func (l *Launcher) Progress() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.progress
}

func (l *Launcher) setProgress(p float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if p > l.progress {
		l.progress = p
	}
}

// Tick implements window.Handler.
func (l *Launcher) Tick(in window.Input) window.Output {
	switch l.phase {
	case phaseSplash:
		return l.tickSplash(in)
	case phaseSwitch:
		l.switching++
		if (in.W == editorW && in.H == editorH) || l.switching >= switchWait {
			l.phase = phaseEditor
			return l.App.Tick(in)
		}
		return window.Output{}
	}
	return l.App.Tick(in)
}

func (l *Launcher) tickSplash(in window.Input) window.Output {
	now := l.Now()
	if !l.started {
		l.started, l.start = true, now
		l.loaded = make(chan error, 1)
		go func() { l.loaded <- l.Load(l.setProgress) }()
	}
	for _, e := range in.Events {
		if e.Type == "key" && e.Key == "Escape" {
			return window.Output{Quit: true}
		}
		if e.Type == "down" && l.loadErr != "" {
			return window.Output{Quit: true}
		}
	}

	changed := false
	if !l.done {
		select {
		case err := <-l.loaded:
			l.done = true
			if err != nil {
				l.Err = err
				l.loadErr = err.Error()
				fmt.Fprintln(os.Stderr, "linkz:", err)
				l.Splash.ShowError(l.loadErr)
				changed = true
			}
		default:
		}
	}

	elapsed := now.Sub(l.start)
	if n := int(elapsed.Seconds() * float64(l.Splash.FPS)); n != l.n {
		l.n = n
		changed = true
		if l.loadErr == "" {
			l.advance(now)
		}
	}

	if l.loadErr == "" && l.ready(now, elapsed) {
		l.phase = phaseSwitch
		if l.OnReady != nil {
			l.OnReady()
		}
		return window.Output{Chrome: &window.Chrome{Decorated: true, Size: [2]int{editorW, editorH}, Center: true}}
	}

	var out window.Output
	if changed {
		img := l.Splash.Frame(l.n, l.stage, l.bloom, l.loadErr != "")
		out.Frame = &window.Frame{Img: img, Dirty: img.Rect}
	}
	return out
}

// advance moves the picture one animation frame: the vines grow one stage
// toward the stage the progress asks for, and once loading is done and the
// vines are full, the bloom plays once.
func (l *Launcher) advance(now time.Time) {
	sp := l.Splash
	if sp.Vines > 0 {
		target := int(l.Progress() * float64(sp.Vines-1))
		if l.done {
			target = sp.Vines - 1
		}
		if l.stage < target {
			l.stage++
		}
	}
	if !l.done || (sp.Vines > 0 && l.stage < sp.Vines-1) {
		return
	}
	atEnd := true
	if sp.Bloom > 0 {
		if l.bloom < sp.Bloom-1 {
			l.bloom++
		}
		atEnd = l.bloom == sp.Bloom-1
	}
	if atEnd && l.finalAt.IsZero() {
		l.finalAt = now
	}
}

// ready: loading is done, the finished picture has been shown long enough,
// and the splash has been up for the minimum time.
func (l *Launcher) ready(now time.Time, elapsed time.Duration) bool {
	if !l.done || l.finalAt.IsZero() || elapsed < MinSplash {
		return false
	}
	hold := time.Duration(0)
	switch {
	case l.Splash.Bloom > 0:
		hold = HoldBloom
	case l.Splash.Vines > 0:
		hold = HoldVines
	}
	return now.Sub(l.finalAt) >= hold
}
