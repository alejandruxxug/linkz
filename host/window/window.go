// Package window is the native window under Linkz. It only opens a window,
// polls input and shows an RGBA buffer. It never draws: every pixel comes
// from the host's framebuffer, decided by lizard-script.
//
// It sits behind a small interface so a hand-written Cocoa layer can replace
// Ebitengine later. Ebitengine owns the main loop, so the window calls the
// handler once per tick (Run replaces Open/Events/Present).
package window

import "image"

// Mods are the modifier keys held during an event.
type Mods struct{ Shift, Alt, Cmd, Ctrl bool }

// Pen is what the pointer knows besides its position.
type Pen struct {
	Pressure float64    // 0 to 1; 1 for a mouse
	Tilt     [2]float64 // tablet tilt, -1 to 1
	Rotation float64    // degrees
	Source   string     // "mouse", "trackpad" or "tablet"
	Eraser   bool
}

// Sample is one point of a drag: every native pointer event the system sent,
// not just the one the tick polled. X and Y are window pixels, T is seconds.
type Sample struct {
	X, Y     int
	Pressure float64
	T        float64
}

// Event is one input event. X and Y are window pixels.
type Event struct {
	Type     string // down, up, move, key, text, wheel, focus, leave, enter
	X, Y     int
	Button   int
	Key      string
	Text     string
	Focused  bool // focus: the window gained (true) or lost (false) focus
	DY       float64
	Mods     Mods
	Pen      Pen
	Samples  []Sample // move: every point since the last tick, last one at X, Y
	Injected bool     // made by the backdoor: never coalesced
}

// Smallest window, in logical pixels: below it the toolbar and the panels
// no longer fit. The window system enforces it; the backdoor clamps to it.
const (
	MinW = 640
	MinH = 420
)

// Input is what happened since the last tick.
type Input struct {
	W, H   int // window size in pixels
	Events []Event
	Hover  bool // lizard wants move events with no button down
}

// Frame is the window image to show. Pix is RGBA, premultiplied, W*H*4.
// Dirty is the part that changed; empty means nothing changed.
type Frame struct {
	Img   *image.RGBA
	Dirty image.Rectangle
}

// Chrome is a change of the window's frame, made while it runs: the splash
// screen becoming the editor window.
type Chrome struct {
	Decorated bool   // title bar and border
	Floating  bool   // above all other windows
	Size      [2]int // new size, if both are positive
	Center    bool   // move to the middle of the screen after resizing
}

// Output is what the handler wants back.
type Output struct {
	Frame    *Frame
	ResizeTo [2]int  // new window size, if both are positive
	Chrome   *Chrome // change the window frame, if set
	Quit     bool
}

// Handler runs the app once per tick.
type Handler interface {
	Tick(in Input) Output
}

// Window shows the app.
type Window interface {
	// Run opens the window and blocks until it closes or the handler quits.
	Run(title string, w, h int, handler Handler) error
}
