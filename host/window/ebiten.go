package window

import (
	"fmt"
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"linkz/host/pressure"
)

// Ebiten is the Ebitengine window. Only WritePixels is used to show pixels;
// none of Ebitengine's drawing, text or vector APIs are called.
type Ebiten struct {
	// Frameless opens the window without a title bar, transparent, floating
	// above other windows and not resizable: the splash screen. A Chrome
	// output turns it into a normal window. Transparency can only be chosen
	// before the window opens, so the window stays "transparent" for its
	// whole life; that only matters for pixels the app leaves clear.
	Frameless bool
}

// Run implements Window.
func (e Ebiten) Run(title string, w, h int, handler Handler) error {
	g := &game{handler: handler, w: w, h: h}
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if e.Frameless {
		ebiten.SetWindowDecorated(false)
		ebiten.SetWindowFloating(true)
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
		center(w, h)
	}
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetTPS(30)
	opts := &ebiten.RunGameOptions{ScreenTransparent: e.Frameless}
	if err := ebiten.RunGameWithOptions(g, opts); err != nil && err != ebiten.Termination {
		return err
	}
	return nil
}

type game struct {
	handler Handler
	w, h    int // size from Layout
	started bool
	frame   *Frame // waiting to be shown
	shownW  int    // width of the screen the last full frame went to
	shownH  int

	chromeLog int

	lastX, lastY int
	held         map[int]bool
}

// Layout makes the buffer the size of the window: logical pixels equal
// buffer pixels.
func (g *game) Layout(outW, outH int) (int, int) {
	g.w, g.h = outW, outH
	return outW, outH
}

var mouseButtons = []ebiten.MouseButton{ebiten.MouseButtonLeft, ebiten.MouseButtonRight, ebiten.MouseButtonMiddle}

func (g *game) Update() error {
	if !g.started {
		g.started = true
		pressure.Start()
		g.held = map[int]bool{}
		g.lastX, g.lastY = -1, -1
	}
	in := Input{W: g.w, H: g.h}
	in.Hover = false
	if hv, ok := g.handler.(interface{ WantsHover() bool }); ok {
		in.Hover = hv.WantsHover()
	}
	in.Events = g.poll(in.Hover)

	out := g.handler.Tick(in)
	if out.Frame != nil && !out.Frame.Dirty.Empty() {
		g.frame = out.Frame
	}
	if out.ResizeTo[0] > 0 && out.ResizeTo[1] > 0 {
		ebiten.SetWindowSize(out.ResizeTo[0], out.ResizeTo[1])
	}
	if out.Chrome != nil {
		g.applyChrome(out.Chrome)
	} else if g.chromeLog > 0 {
		g.chromeLog--
		if g.chromeLog == 0 {
			fmt.Fprintf(os.Stderr, "window: decorated=%v floating=%v\n",
				ebiten.IsWindowDecorated(), ebiten.IsWindowFloating())
		}
	}
	if out.Quit {
		return ebiten.Termination
	}
	return nil
}

// center puts a window of w by h in the middle of the screen.
func center(w, h int) {
	if m := ebiten.Monitor(); m != nil {
		mw, mh := m.Size()
		ebiten.SetWindowPosition(max(0, (mw-w)/2), max(0, (mh-h)/2))
	}
}

// applyChrome changes the window frame while the window runs.
func (g *game) applyChrome(c *Chrome) {
	ebiten.SetWindowDecorated(c.Decorated)
	ebiten.SetWindowFloating(c.Floating)
	if c.Decorated {
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	}
	if c.Size[0] > 0 && c.Size[1] > 0 {
		ebiten.SetWindowSize(c.Size[0], c.Size[1])
		if c.Center {
			center(c.Size[0], c.Size[1])
		}
	}
	g.chromeLog = 3 // report what the system made of it, a few ticks later
}

// Draw writes the pixels lizard produced. When the screen size changed the
// whole image is written, otherwise only the changed part.
func (g *game) Draw(screen *ebiten.Image) {
	f := g.frame
	if f == nil {
		return
	}
	b := screen.Bounds()
	if b.Dx() != f.Img.Rect.Dx() || b.Dy() != f.Img.Rect.Dy() {
		return
	}
	if b.Dx() != g.shownW || b.Dy() != g.shownH || f.Dirty == f.Img.Rect {
		screen.WritePixels(f.Img.Pix)
		g.shownW, g.shownH = b.Dx(), b.Dy()
	} else {
		r := f.Dirty.Intersect(f.Img.Rect)
		sub := make([]byte, 0, r.Dx()*r.Dy()*4)
		for y := r.Min.Y; y < r.Max.Y; y++ {
			o := y*f.Img.Stride + r.Min.X*4
			sub = append(sub, f.Img.Pix[o:o+r.Dx()*4]...)
		}
		screen.SubImage(r).(*ebiten.Image).WritePixels(sub)
	}
	g.frame = nil
}

func (g *game) mods() Mods {
	return Mods{
		Shift: ebiten.IsKeyPressed(ebiten.KeyShift),
		Alt:   ebiten.IsKeyPressed(ebiten.KeyAlt),
		Cmd:   ebiten.IsKeyPressed(ebiten.KeyMeta),
		Ctrl:  ebiten.IsKeyPressed(ebiten.KeyControl),
	}
}

// pen reads the pressure the host/pressure monitor saw. A button must be
// down for it to count; otherwise the pointer is a plain mouse.
func (g *game) pen(down bool) Pen {
	p := Pen{Pressure: 1, Source: "mouse"}
	s := pressure.Latest()
	if down && s.OK && s.Source != "" && s.Source != "mouse" {
		p = Pen{Pressure: s.Pressure, Tilt: s.Tilt, Rotation: s.Rotation, Source: s.Source, Eraser: s.Eraser}
	}
	return p
}

// repeating: a held key repeats after 15 ticks, then every 2 ticks (30 TPS).
func repeats(k ebiten.Key) bool {
	d := inpututil.KeyPressDuration(k)
	return d == 1 || (d >= 15 && (d-15)%2 == 0)
}

// penOf turns a native sample into the pen state of an event.
func penOf(s pressure.Sample, eraser bool) Pen {
	return Pen{Pressure: s.Pressure, Tilt: s.Tilt, Rotation: s.Rotation, Source: s.Source, Eraser: eraser}
}

// lastPhase finds the newest native sample of a phase.
func lastPhase(ss []pressure.Sample, phase string) (pressure.Sample, bool) {
	for i := len(ss) - 1; i >= 0; i-- {
		if ss[i].Phase == phase {
			return ss[i], true
		}
	}
	return pressure.Sample{}, false
}

// dragSamples are the native drag points of this tick as window samples.
// The last one is moved onto (x, y), the position the event itself carries.
func dragSamples(ss []pressure.Sample, x, y int) ([]Sample, pressure.Sample, bool) {
	var out []Sample
	var last pressure.Sample
	for _, s := range ss {
		if s.Phase != pressure.PhaseDrag {
			continue
		}
		out = append(out, Sample{X: int(s.X), Y: int(s.Y), Pressure: s.Pressure, T: s.T})
		last = s
	}
	if len(out) > 0 {
		out[len(out)-1].X, out[len(out)-1].Y = x, y
	}
	return out, last, len(out) > 0
}

func (g *game) poll(hover bool) []Event {
	var evs []Event
	x, y := ebiten.CursorPosition()
	mods := g.mods()
	native := pressure.Drain()

	for i, b := range mouseButtons {
		if inpututil.IsMouseButtonJustPressed(b) {
			g.held[i] = true
			pen := g.pen(true)
			if s, ok := lastPhase(native, pressure.PhaseDown); ok {
				pen = penOf(s, pressure.Latest().Eraser)
			}
			evs = append(evs, Event{Type: "down", X: x, Y: y, Button: i, Mods: mods, Pen: pen})
		}
	}
	anyDown := false
	for _, b := range mouseButtons {
		if ebiten.IsMouseButtonPressed(b) {
			anyDown = true
		}
	}
	drags, lastDrag, haveDrags := dragSamples(native, x, y)
	if (x != g.lastX || y != g.lastY || (anyDown && haveDrags)) && (anyDown || hover) {
		pen := g.pen(anyDown)
		if anyDown && haveDrags {
			pen = penOf(lastDrag, pressure.Latest().Eraser)
		}
		evs = append(evs, Event{Type: "move", X: x, Y: y, Mods: mods, Pen: pen, Samples: drags})
	}
	g.lastX, g.lastY = x, y
	for i, b := range mouseButtons {
		if inpututil.IsMouseButtonJustReleased(b) {
			delete(g.held, i)
			pen := g.pen(false)
			if s, ok := lastPhase(native, pressure.PhaseUp); ok {
				pen = penOf(s, pressure.Latest().Eraser)
			}
			evs = append(evs, Event{Type: "up", X: x, Y: y, Button: i, Mods: mods, Pen: pen})
		}
	}
	if _, dy := ebiten.Wheel(); dy != 0 {
		evs = append(evs, Event{Type: "wheel", X: x, Y: y, DY: dy, Mods: mods, Pen: g.pen(false)})
	}

	for _, k := range inpututil.AppendPressedKeys(nil) {
		if repeats(k) {
			evs = append(evs, Event{Type: "key", X: x, Y: y, Key: k.String(), Mods: mods})
		}
	}
	// Full Unicode text, as the system composed it (dead keys, option combos).
	for _, r := range ebiten.AppendInputChars(nil) {
		if r >= 32 && r != 127 {
			evs = append(evs, Event{Type: "text", Text: string(r), Mods: mods})
		}
	}
	return evs
}

var _ image.Image = (*image.RGBA)(nil)
