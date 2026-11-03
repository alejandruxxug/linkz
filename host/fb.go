package host

import (
	"image"
	"image/png"
	"lizardscript/object"
	"math"
	"os"
	"sync"
)

// FB is the framebuffer: a half-resolution UI buffer shown at 2x, plus a
// full-resolution canvas buffer composited over a rectangle of the UI. Go
// only fills pixels here; what to draw is decided by lizard-script.
type FB struct {
	mu sync.Mutex

	winW, winH int // window size in pixels
	ui         *Buf
	canvas     *Buf
	canvasRect image.Rectangle // in UI pixels
	off        *Buf            // set while an offscreen buffer is the target
	target     string          // "ui" or "canvas" (the named target)

	// Per target: the region changed since the last present, and whether
	// anything was drawn at all. A draw with no declared region counts as
	// the whole buffer.
	dirty   map[string]image.Rectangle
	touched map[string]bool
	final   *image.RGBA     // the composited window image
	all     bool            // present everything next time
	kept    map[string]*Buf // offscreen buffers lizard asked to keep
}

// NewFB makes a framebuffer for a window of w by h pixels. The UI buffer is
// half that size.
func NewFB(w, h int) *FB {
	f := &FB{target: "ui", dirty: map[string]image.Rectangle{}, touched: map[string]bool{}}
	f.Resize(w, h)
	return f
}

// NewFBUI makes a framebuffer whose UI buffer is exactly w by h (headless
// painting). It has no canvas.
func NewFBUI(w, h int) *FB {
	f := NewFB(w*2, h*2)
	return f
}

// Resize reallocates the buffers for a new window size. Contents are lost;
// lizard repaints after the resize event.
func (f *FB) Resize(w, h int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.winW, f.winH = w, h
	f.ui = NewBuf(w/2, h/2)
	f.final = image.NewRGBA(image.Rect(0, 0, w, h))
	f.canvas, f.canvasRect, f.off, f.target = nil, image.Rectangle{}, nil, "ui"
	f.all = true
	f.touched["ui"] = true
}

// WindowSize is the size of the window in pixels.
func (f *FB) WindowSize() (int, int) { return f.winW, f.winH }

// UISize is the size of the UI buffer.
func (f *FB) UISize() (int, int) { return f.ui.W(), f.ui.H() }

// current is the buffer drawing goes to now.
func (f *FB) current() *Buf {
	if f.off != nil {
		return f.off
	}
	if f.target == "canvas" && f.canvas != nil {
		return f.canvas
	}
	return f.ui
}

// mark records that the current target was drawn on.
func (f *FB) mark() {
	if f.off == nil {
		f.touched[f.target] = true
	}
}

// Screenshot is a copy of the window as presented (UI scaled, canvas over).
func (f *FB) Screenshot() *image.RGBA {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.compose(f.final.Rect)
	out := image.NewRGBA(f.final.Rect)
	copy(out.Pix, f.final.Pix)
	return out
}

// Compose builds the window image for what changed. It returns the image and
// the changed rectangle (empty if nothing changed).
func (f *FB) Compose() (*image.RGBA, image.Rectangle) {
	f.mu.Lock()
	defer f.mu.Unlock()
	full := f.final.Rect
	var region image.Rectangle
	if f.all {
		region = full
	} else {
		region = f.regionOf("ui", 2, image.Point{}).Union(f.regionOf("canvas", 1, f.canvasRect.Min.Mul(2)))
	}
	f.all = false
	f.touched = map[string]bool{}
	f.dirty = map[string]image.Rectangle{}
	region = region.Intersect(full)
	if !region.Empty() {
		f.compose(region)
	}
	return f.final, region
}

// regionOf is the changed part of a target in window pixels: its dirty rect
// scaled by scale and moved by origin.
func (f *FB) regionOf(name string, scale int, origin image.Point) image.Rectangle {
	if !f.touched[name] {
		return image.Rectangle{}
	}
	buf := f.ui
	if name == "canvas" {
		buf = f.canvas
	}
	if buf == nil {
		return image.Rectangle{}
	}
	r, declared := f.dirty[name]
	if !declared {
		r = buf.Img.Rect
	}
	r = r.Intersect(buf.Img.Rect)
	return image.Rect(r.Min.X*scale, r.Min.Y*scale, r.Max.X*scale, r.Max.Y*scale).Add(origin)
}

// compose rebuilds the window pixels in region: the UI scaled by two with
// nearest neighbor, then the canvas over its rectangle.
func (f *FB) compose(region image.Rectangle) {
	uw, uh := f.ui.W(), f.ui.H()
	if uw == 0 || uh == 0 {
		return
	}
	for y := region.Min.Y; y < region.Max.Y; y++ {
		uy := min(y/2, uh-1)
		for x := region.Min.X; x < region.Max.X; x++ {
			ux := min(x/2, uw-1)
			s := f.ui.Img.Pix[uy*f.ui.Img.Stride+ux*4:]
			d := f.final.Pix[y*f.final.Stride+x*4:]
			d[0], d[1], d[2], d[3] = s[0], s[1], s[2], s[3]
		}
	}
	if f.canvas == nil {
		return
	}
	origin := f.canvasRect.Min.Mul(2)
	over := region.Intersect(f.canvas.Img.Rect.Add(origin))
	for y := over.Min.Y; y < over.Max.Y; y++ {
		for x := over.Min.X; x < over.Max.X; x++ {
			s := f.canvas.Img.Pix[(y-origin.Y)*f.canvas.Img.Stride+(x-origin.X)*4:]
			if s[3] == 0 {
				continue
			}
			d := f.final.Pix[y*f.final.Stride+x*4:]
			inv := 255 - int(s[3])
			for i := 0; i < 4; i++ {
				d[i] = uint8(min(255, int(s[i])+(int(d[i])*inv+127)/255))
			}
		}
	}
}

// CanvasRect is where the canvas shows, in UI pixels. ok is false when no
// canvas was declared.
func (f *FB) CanvasRect() (image.Rectangle, bool) {
	return f.canvasRect, f.canvas != nil
}

// SavePNG writes the current target (or part of it) to path.
func (f *FB) SavePNG(path string, r *image.Rectangle) error {
	f.mu.Lock()
	b := f.current()
	f.mu.Unlock()
	var img image.Image = b.Img
	if r != nil {
		img = b.Img.SubImage(*r)
	}
	return writePNG(path, img)
}

func writePNG(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := png.Encode(file, img); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

// Builtins is the fb hash that lizard-script sees.
func (f *FB) Builtins() *object.Hash {
	return builtinHash("fb", map[string]object.BuiltinFunction{
		"size":      f.size,
		"span":      f.span,
		"spans":     f.spans,
		"fill":      f.fill,
		"clear":     f.clear,
		"clip":      f.clip,
		"unclip":    f.unclip,
		"offscreen": f.offscreen,
		"keep":      f.keep,
		"screen":    f.screen,
		"target":    f.setTarget,
		"canvas":    f.setCanvas,
		"dirty":     f.markDirty,
		"save_png":  f.savePNG,
		"load_png":  f.loadPNG,
		"blit":      f.blit,
	})
}

// size gives the size of the current target as a point.
func (f *FB) size(args ...object.Object) object.Object {
	if len(args) != 0 {
		return errorf("wrong number of arguments: want 0, got %d", len(args))
	}
	b := f.current()
	return newPoint(b.W(), b.H())
}

func (f *FB) span(args ...object.Object) object.Object {
	if len(args) != 4 {
		return errorf("wrong number of arguments: want 4 (y, x0, x1, color), got %d", len(args))
	}
	y, ok1 := number(args[0])
	x0, ok2 := number(args[1])
	x1, ok3 := number(args[2])
	c, ok4 := colorOf(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return errorf("want (number y, number x0, number x1, COLOR)")
	}
	f.mark()
	f.current().Span(int(math.Floor(y)), x0, x1, c)
	return null
}

// spans draws many [y, x0, x1] runs at an offset in one color.
func (f *FB) spans(args ...object.Object) object.Object {
	if len(args) != 3 {
		return errorf("wrong number of arguments: want 3 (list, offset, color), got %d", len(args))
	}
	list, ok1 := args[0].(*object.Array)
	dx, dy, ok2 := pointXY(args[1])
	c, ok3 := colorOf(args[2])
	if !ok1 || !ok2 || !ok3 {
		return errorf("want (ARRAY of [y, x0, x1], POINT offset, COLOR)")
	}
	f.mark()
	b := f.current()
	for _, item := range list.Elements {
		run, ok := item.(*object.Array)
		if !ok || len(run.Elements) != 3 {
			return errorf("each span must be [y, x0, x1]")
		}
		y, ok1 := number(run.Elements[0])
		x0, ok2 := number(run.Elements[1])
		x1, ok3 := number(run.Elements[2])
		if !ok1 || !ok2 || !ok3 {
			return errorf("each span must be [y, x0, x1] of numbers")
		}
		b.Span(int(math.Floor(y+dy)), x0+dx, x1+dx, c)
	}
	return null
}

// rectArgs reads (pos, size) into float bounds.
func rectArgs(pos, size object.Object) (x, y, w, h float64, ok bool) {
	x, y, ok1 := pointXY(pos)
	w, h, ok2 := pointXY(size)
	return x, y, w, h, ok1 && ok2
}

func (f *FB) fill(args ...object.Object) object.Object {
	if len(args) != 3 {
		return errorf("wrong number of arguments: want 3 (pos, size, color), got %d", len(args))
	}
	x, y, w, h, ok := rectArgs(args[0], args[1])
	c, okc := colorOf(args[2])
	if !ok || !okc {
		return errorf("want (POINT pos, POINT size, COLOR)")
	}
	f.mark()
	b := f.current()
	for row := int(math.Round(y)); row < int(math.Round(y+h)); row++ {
		b.Span(row, x, x+w, c)
	}
	return null
}

// clear makes a rectangle transparent (the canvas has no background of its own).
func (f *FB) clear(args ...object.Object) object.Object {
	if len(args) != 2 {
		return errorf("wrong number of arguments: want 2 (pos, size), got %d", len(args))
	}
	x, y, w, h, ok := rectArgs(args[0], args[1])
	if !ok {
		return errorf("want (POINT pos, POINT size)")
	}
	f.mark()
	b := f.current()
	r := image.Rect(int(x), int(y), int(x+w), int(y+h)).Intersect(b.clip)
	for row := r.Min.Y; row < r.Max.Y; row++ {
		p := b.Img.Pix[row*b.Img.Stride+r.Min.X*4 : row*b.Img.Stride+r.Max.X*4]
		for i := range p {
			p[i] = 0
		}
	}
	return null
}

func (f *FB) clip(args ...object.Object) object.Object {
	if len(args) != 2 {
		return errorf("wrong number of arguments: want 2 (pos, size), got %d", len(args))
	}
	x, y, w, h, ok := rectArgs(args[0], args[1])
	if !ok {
		return errorf("want (POINT pos, POINT size)")
	}
	f.current().SetClip(image.Rect(int(x), int(y), int(x+w), int(y+h)))
	return null
}

func (f *FB) unclip(args ...object.Object) object.Object {
	f.current().Unclip()
	return null
}

// offscreen switches drawing to a new transparent buffer.
func (f *FB) offscreen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return errorf("wrong number of arguments: want 1 (size), got %d", len(args))
	}
	w, h, ok := pointXY(args[0])
	if !ok || w < 0 || h < 0 || w*h > 1<<26 {
		return errorf("want a POINT size, at most 64M pixels")
	}
	f.off = NewBuf(int(w), int(h))
	return null
}

// keep stores the current offscreen buffer under a name, so the host can show
// it later (the splash frames). Keeping again under a name replaces it.
func (f *FB) keep(args ...object.Object) object.Object {
	name, ok := object.Object(nil), false
	if len(args) == 1 {
		name = args[0]
		_, ok = name.(*object.String)
	}
	if !ok || f.off == nil {
		return errorf("want (STRING name) while an offscreen buffer is the target")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.kept == nil {
		f.kept = map[string]*Buf{}
	}
	f.kept[name.(*object.String).Value] = f.off
	f.off = nil // the kept buffer is no longer drawn on
	return null
}

// Kept is the buffer stored by fb.keep under name, or nil.
func (f *FB) Kept(name string) *Buf {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.kept[name]
}

// screen goes back to the named target (ui unless fb.target changed it).
func (f *FB) screen(args ...object.Object) object.Object {
	f.off = nil
	return null
}

func (f *FB) setTarget(args ...object.Object) object.Object {
	name, ok := object.Object(nil), false
	if len(args) == 1 {
		name = args[0]
		_, ok = name.(*object.String)
	}
	if !ok {
		return errorf(`want ("ui" or "canvas")`)
	}
	switch n := name.(*object.String).Value; n {
	case "ui":
		f.target = "ui"
	case "canvas":
		if f.canvas == nil {
			return errorf("no canvas: call fb.canvas(pos, size) first")
		}
		f.target = "canvas"
	default:
		return errorf("unknown target %q, want \"ui\" or \"canvas\"", n)
	}
	f.off = nil
	return null
}

// setCanvas declares where the canvas shows (UI pixels). The buffer behind
// it is full resolution, so twice the size. A size of (0, 0) removes it.
func (f *FB) setCanvas(args ...object.Object) object.Object {
	if len(args) != 2 {
		return errorf("wrong number of arguments: want 2 (pos, size), got %d", len(args))
	}
	x, y, w, h, ok := rectArgs(args[0], args[1])
	if !ok || w < 0 || h < 0 {
		return errorf("want (POINT pos, POINT size)")
	}
	r := image.Rect(int(x), int(y), int(x+w), int(y+h))
	f.mu.Lock()
	defer f.mu.Unlock()
	if r.Empty() {
		if f.canvas != nil {
			f.all = true
		}
		f.canvas, f.canvasRect, f.target = nil, image.Rectangle{}, "ui"
		return null
	}
	if f.canvas == nil || f.canvas.W() != r.Dx()*2 || f.canvas.H() != r.Dy()*2 {
		f.canvas = NewBuf(r.Dx()*2, r.Dy()*2)
		f.touched["canvas"] = true
	}
	if r != f.canvasRect {
		f.all = true // the old place must show the UI again
		f.canvasRect = r
	}
	return null
}

// markDirty tells which part of the current target changed.
func (f *FB) markDirty(args ...object.Object) object.Object {
	if len(args) != 2 {
		return errorf("wrong number of arguments: want 2 (pos, size), got %d", len(args))
	}
	x, y, w, h, ok := rectArgs(args[0], args[1])
	if !ok {
		return errorf("want (POINT pos, POINT size)")
	}
	if f.off != nil {
		return null
	}
	r := image.Rect(int(math.Floor(x)), int(math.Floor(y)), int(math.Ceil(x+w)), int(math.Ceil(y+h)))
	f.touched[f.target] = true
	if old, ok := f.dirty[f.target]; ok {
		r = old.Union(r)
	}
	f.dirty[f.target] = r
	return null
}

func (f *FB) savePNG(args ...object.Object) object.Object {
	if len(args) != 1 && len(args) != 3 {
		return errorf("wrong number of arguments: want 1 (path) or 3 (path, pos, size), got %d", len(args))
	}
	path, ok := args[0].(*object.String)
	if !ok {
		return errorf("path must be a STRING")
	}
	var rect *image.Rectangle
	if len(args) == 3 {
		x, y, w, h, ok := rectArgs(args[1], args[2])
		if !ok {
			return errorf("want (path, POINT pos, POINT size)")
		}
		r := image.Rect(int(x), int(y), int(x+w), int(y+h))
		rect = &r
	}
	if err := f.SavePNG(path.Value, rect); err != nil {
		return errorf("%v", err)
	}
	return null
}
