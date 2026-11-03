package host

import (
	"image"
	"lizardscript/object"
	"testing"
)

var white = rgba{255, 255, 255, 255}

// px reads a pixel as r, g, b, a.
func px(b *Buf, x, y int) [4]uint8 {
	o := y*b.Img.Stride + x*4
	return [4]uint8{b.Img.Pix[o], b.Img.Pix[o+1], b.Img.Pix[o+2], b.Img.Pix[o+3]}
}

func TestSpanCoverage(t *testing.T) {
	tests := []struct {
		name   string
		x0, x1 float64
		want   [6]uint8 // red channel of pixels 0..5 on black
	}{
		{"whole pixels", 1, 4, [6]uint8{0, 255, 255, 255, 0, 0}},
		{"float ends", 1.5, 4.25, [6]uint8{0, 128, 255, 255, 64, 0}},
		{"inside one pixel", 2.25, 2.75, [6]uint8{0, 0, 128, 0, 0, 0}},
		{"empty", 3, 3, [6]uint8{}},
		{"backwards", 4, 1, [6]uint8{}},
		{"cut at the left", -5, 2, [6]uint8{255, 255, 0, 0, 0, 0}},
		{"cut at the right", 4, 99, [6]uint8{0, 0, 0, 0, 255, 255}},
		{"far outside", 100, 200, [6]uint8{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuf(6, 1)
			for x := 0; x < 6; x++ {
				b.Img.Pix[x*4+3] = 255 // opaque black
			}
			b.Span(0, tt.x0, tt.x1, white)
			for x := 0; x < 6; x++ {
				if got := px(b, x, 0)[0]; got != tt.want[x] {
					t.Errorf("pixel %d: red %d, want %d", x, got, tt.want[x])
				}
			}
		})
	}
}

func TestSpanOutOfRangeRowsNeverPanic(t *testing.T) {
	b := NewBuf(4, 4)
	for _, y := range []int{-1, 4, 1 << 30, -(1 << 30)} {
		b.Span(y, 0, 4, white)
	}
	NewBuf(0, 0).Span(0, 0, 4, white)
}

func TestBlendSourceOver(t *testing.T) {
	b := NewBuf(2, 1)
	b.Span(0, 0, 1, rgba{0, 0, 0, 255})
	b.Span(0, 0, 1, rgba{255, 255, 255, 128}) // half white over black
	got := px(b, 0, 0)
	if got[0] < 127 || got[0] > 129 || got[3] != 255 {
		t.Errorf("half white over black = %v, want gray 128 and opaque", got)
	}
	// Over a transparent pixel the result is premultiplied.
	b.Span(0, 1, 2, rgba{255, 0, 0, 128})
	got = px(b, 1, 0)
	if got[0] != 128 || got[3] != 128 {
		t.Errorf("half red over nothing = %v, want premultiplied (128,0,0,128)", got)
	}
}

func TestClip(t *testing.T) {
	b := NewBuf(10, 10)
	b.SetClip(image.Rect(2, 2, 5, 5))
	b.Span(3, 0, 10, white)
	b.Span(7, 0, 10, white)
	for x := 0; x < 10; x++ {
		want := uint8(0)
		if x >= 2 && x < 5 {
			want = 255
		}
		if got := px(b, x, 3)[3]; got != want {
			t.Errorf("row 3 pixel %d alpha %d, want %d", x, got, want)
		}
		if px(b, x, 7)[3] != 0 {
			t.Errorf("row 7 is outside the clip but pixel %d was drawn", x)
		}
	}
	b.Unclip()
	b.Span(7, 0, 10, white)
	if px(b, 9, 7)[3] != 255 {
		t.Error("unclip did not lift the clip")
	}
}

// call runs a builtin of the fb hash.
func call(t *testing.T, h *object.Hash, name string, args ...object.Object) object.Object {
	t.Helper()
	b, ok := hashGet(h, name).(*object.Builtin)
	if !ok {
		t.Fatalf("fb.%s is not a builtin", name)
	}
	return b.Fn(args...)
}

func color(r, g, b, a uint8) object.Object { return &object.Color{R: r, G: g, B: b, A: a} }

func TestFBPrimitives(t *testing.T) {
	f := NewFBUI(8, 6)
	h := f.Builtins()

	if p, ok := call(t, h, "size").(*object.Point); !ok || p.Inspect() != "(8, 6)" {
		t.Errorf("size = %v, want (8, 6)", call(t, h, "size").Inspect())
	}

	call(t, h, "fill", newPoint(1, 1), newPoint(3, 2), color(255, 0, 0, 255))
	if px(f.ui, 1, 1) != [4]uint8{255, 0, 0, 255} || px(f.ui, 3, 2) != [4]uint8{255, 0, 0, 255} {
		t.Error("fill missed its rectangle")
	}
	if px(f.ui, 4, 1)[3] != 0 || px(f.ui, 1, 3)[3] != 0 {
		t.Error("fill went past its rectangle")
	}

	// spans: [y, x0, x1] runs moved by an offset
	runs := &object.Array{Elements: []object.Object{
		&object.Array{Elements: []object.Object{newInt(0), newInt(0), newInt(2)}},
		&object.Array{Elements: []object.Object{newInt(1), newInt(1), newInt(2)}},
	}}
	call(t, h, "spans", runs, newPoint(4, 3), color(0, 255, 0, 255))
	if px(f.ui, 4, 3)[1] != 255 || px(f.ui, 5, 3)[1] != 255 || px(f.ui, 5, 4)[1] != 255 || px(f.ui, 4, 4)[1] != 0 {
		t.Error("spans did not land at the offset")
	}

	// errors, not panics
	for _, args := range [][]object.Object{
		{newInt(1)},
		{newInt(0), newInt(0), newInt(1), newInt(5)},
	} {
		if _, ok := call(t, h, "span", args...).(*object.Error); !ok {
			t.Errorf("span with bad arguments %v should be an error", args)
		}
	}
	call(t, h, "span", newInt(-4), newInt(-9), newInt(99), white2())
	call(t, h, "span", newInt(400), newInt(0), newInt(3), white2())

	// clear makes pixels transparent again
	call(t, h, "clear", newPoint(0, 0), newPoint(8, 6))
	if px(f.ui, 1, 1)[3] != 0 {
		t.Error("clear left a pixel")
	}
}

func white2() object.Object { return color(255, 255, 255, 255) }

func TestOffscreenAndTargets(t *testing.T) {
	f := NewFB(40, 40) // UI 20x20
	h := f.Builtins()

	call(t, h, "offscreen", newPoint(3, 3))
	call(t, h, "fill", newPoint(0, 0), newPoint(3, 3), white2())
	if px(f.ui, 0, 0)[3] != 0 {
		t.Error("offscreen drawing reached the UI")
	}
	if p := call(t, h, "size").Inspect(); p != "(3, 3)" {
		t.Errorf("offscreen size = %s", p)
	}
	call(t, h, "screen")
	if p := call(t, h, "size").Inspect(); p != "(20, 20)" {
		t.Errorf("size after screen() = %s", p)
	}

	if _, ok := call(t, h, "target", newStr("canvas")).(*object.Error); !ok {
		t.Error("target canvas without a canvas should fail")
	}
	call(t, h, "canvas", newPoint(2, 2), newPoint(10, 8))
	call(t, h, "target", newStr("canvas"))
	if p := call(t, h, "size").Inspect(); p != "(20, 16)" {
		t.Errorf("canvas buffer size = %s, want full resolution (20, 16)", p)
	}
	call(t, h, "target", newStr("ui"))
	if p := call(t, h, "size").Inspect(); p != "(20, 20)" {
		t.Errorf("ui size = %s", p)
	}
	if _, ok := call(t, h, "target", newStr("nope")).(*object.Error); !ok {
		t.Error("unknown target should fail")
	}
}

func TestComposeScalesUIAndOverlaysCanvas(t *testing.T) {
	f := NewFB(40, 40) // UI 20x20
	h := f.Builtins()
	call(t, h, "fill", newPoint(0, 0), newPoint(20, 20), color(10, 20, 30, 255))
	call(t, h, "canvas", newPoint(5, 5), newPoint(4, 4)) // canvas buffer 8x8
	call(t, h, "target", newStr("canvas"))
	call(t, h, "fill", newPoint(0, 0), newPoint(8, 8), color(200, 0, 0, 255))
	call(t, h, "clear", newPoint(0, 0), newPoint(2, 2)) // transparent corner shows the UI
	call(t, h, "target", newStr("ui"))

	img, region := f.Compose()
	if region != img.Rect {
		t.Errorf("first compose region %v, want the whole window", region)
	}
	at := func(x, y int) [4]uint8 {
		o := y*img.Stride + x*4
		return [4]uint8{img.Pix[o], img.Pix[o+1], img.Pix[o+2], img.Pix[o+3]}
	}
	if at(1, 1) != [4]uint8{10, 20, 30, 255} || at(39, 39) != [4]uint8{10, 20, 30, 255} {
		t.Error("the UI is not scaled to the window")
	}
	if at(10, 10) != [4]uint8{10, 20, 30, 255} {
		t.Errorf("transparent canvas corner should show the UI, got %v", at(10, 10))
	}
	if at(12, 12) != [4]uint8{200, 0, 0, 255} || at(17, 17) != [4]uint8{200, 0, 0, 255} {
		t.Errorf("canvas not over its rectangle: %v", at(12, 12))
	}
	if at(18, 18) != [4]uint8{10, 20, 30, 255} {
		t.Error("canvas drawn past its rectangle")
	}

	// A dirty rectangle on the UI composes only that part (scaled by 2).
	call(t, h, "fill", newPoint(0, 0), newPoint(1, 1), color(1, 1, 1, 255))
	call(t, h, "dirty", newPoint(0, 0), newPoint(1, 1))
	_, region = f.Compose()
	if region != image.Rect(0, 0, 2, 2) {
		t.Errorf("dirty region = %v, want (0,0)-(2,2)", region)
	}
	// Nothing drawn: nothing to present.
	if _, region = f.Compose(); !region.Empty() {
		t.Errorf("idle compose region = %v, want empty", region)
	}
}

func TestCanvasDirtyIsPerTarget(t *testing.T) {
	f := NewFB(40, 40)
	h := f.Builtins()
	call(t, h, "canvas", newPoint(5, 5), newPoint(4, 4))
	f.Compose()
	call(t, h, "target", newStr("canvas"))
	call(t, h, "dirty", newPoint(2, 2), newPoint(3, 3))
	_, region := f.Compose()
	// canvas origin (5,5)*2 = (10,10) in the window; dirty (2,2)-(5,5)
	if region != image.Rect(12, 12, 15, 15) {
		t.Errorf("canvas dirty region = %v, want (12,12)-(15,15)", region)
	}
}

// Two dirty rectangles far apart are composed each on its own: what lies
// between them is left as it was.
func TestComposeEachDirtyRect(t *testing.T) {
	f := NewFB(40, 40)
	h := f.Builtins()
	f.Compose()
	call(t, h, "fill", newPoint(0, 0), newPoint(20, 20), color(9, 9, 9, 255))
	call(t, h, "dirty", newPoint(0, 0), newPoint(1, 1))
	call(t, h, "dirty", newPoint(18, 18), newPoint(1, 1))
	img, region := f.Compose()
	if region != image.Rect(0, 0, 38, 38) {
		t.Errorf("region = %v, want the box around both (0,0)-(38,38)", region)
	}
	px := func(x, y int) uint8 { return img.Pix[y*img.Stride+x*4] }
	if px(1, 1) != 9 || px(37, 37) != 9 {
		t.Error("the dirty rectangles were not composed")
	}
	if px(20, 20) != 0 {
		t.Error("a pixel between the dirty rectangles was composed")
	}
}

func TestDisjoint(t *testing.T) {
	r := image.Rect
	tests := []struct {
		name string
		in   []image.Rectangle
		want int
	}{
		{"apart", []image.Rectangle{r(0, 0, 2, 2), r(5, 5, 6, 6)}, 2},
		{"overlapping", []image.Rectangle{r(0, 0, 4, 4), r(2, 2, 6, 6)}, 1},
		{"a chain", []image.Rectangle{r(0, 0, 3, 3), r(10, 10, 12, 12), r(2, 2, 11, 11)}, 1},
		{"empty ones", []image.Rectangle{{}, r(1, 1, 2, 2)}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := disjoint(tt.in); len(got) != tt.want {
				t.Errorf("disjoint(%v) = %v, want %d rectangles", tt.in, got, tt.want)
			}
		})
	}
}
