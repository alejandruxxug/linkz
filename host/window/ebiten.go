package window

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"linkz/host/pressure"
)

// Ebiten is the Ebitengine window. Only WritePixels is used to show pixels;
// none of Ebitengine's drawing, text or vector APIs are called.
type Ebiten struct{}

// Run implements Window.
func (Ebiten) Run(title string, w, h int, handler Handler) error {
	g := &game{handler: handler, w: w, h: h}
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetTPS(30)
	if err := ebiten.RunGame(g); err != nil && err != ebiten.Termination {
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
	if out.Quit {
		return ebiten.Termination
	}
	return nil
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

func (g *game) poll(hover bool) []Event {
	var evs []Event
	x, y := ebiten.CursorPosition()
	mods := g.mods()

	for i, b := range mouseButtons {
		if inpututil.IsMouseButtonJustPressed(b) {
			g.held[i] = true
			evs = append(evs, Event{Type: "down", X: x, Y: y, Button: i, Mods: mods, Pen: g.pen(true)})
		}
	}
	anyDown := false
	for _, b := range mouseButtons {
		if ebiten.IsMouseButtonPressed(b) {
			anyDown = true
		}
	}
	if (x != g.lastX || y != g.lastY) && (anyDown || hover) {
		evs = append(evs, Event{Type: "move", X: x, Y: y, Mods: mods, Pen: g.pen(anyDown)})
	}
	g.lastX, g.lastY = x, y
	for i, b := range mouseButtons {
		if inpututil.IsMouseButtonJustReleased(b) {
			delete(g.held, i)
			evs = append(evs, Event{Type: "up", X: x, Y: y, Button: i, Mods: mods, Pen: g.pen(false)})
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
