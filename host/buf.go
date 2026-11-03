package host

import (
	"image"
	"math"
)

// Buf is one pixel buffer. Pixels are RGBA, premultiplied, so a fresh buffer
// is transparent and drawing blends source-over.
type Buf struct {
	Img  *image.RGBA
	clip image.Rectangle
}

// NewBuf returns a transparent buffer. A negative size counts as empty.
func NewBuf(w, h int) *Buf {
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	return &Buf{Img: img, clip: img.Rect}
}

// W and H are the size in pixels.
func (b *Buf) W() int { return b.Img.Rect.Dx() }
func (b *Buf) H() int { return b.Img.Rect.Dy() }

// SetClip limits drawing to r, cut to the buffer.
func (b *Buf) SetClip(r image.Rectangle) { b.clip = r.Intersect(b.Img.Rect) }

// Unclip lifts the clip.
func (b *Buf) Unclip() { b.clip = b.Img.Rect }

// rgba is a straight (not premultiplied) color.
type rgba struct{ r, g, b, a uint8 }

// Span fills the pixels [x0, x1) of row y. A fractional end covers its pixel
// partly, so the edge is soft. Anything outside the buffer or clip is cut.
func (b *Buf) Span(y int, x0, x1 float64, c rgba) {
	if y < b.clip.Min.Y || y >= b.clip.Max.Y || c.a == 0 {
		return
	}
	x0 = math.Max(x0, float64(b.clip.Min.X))
	x1 = math.Min(x1, float64(b.clip.Max.X))
	if !(x1 > x0) {
		return
	}
	first := int(math.Floor(x0))
	last := int(math.Ceil(x1)) - 1 // last pixel touched
	row := b.Img.Pix[y*b.Img.Stride:]
	cover := func(x int) float64 {
		return math.Min(float64(x+1), x1) - math.Max(float64(x), x0)
	}

	blend(row[first*4:first*4+4], c, cover(first))
	if last == first {
		return
	}
	// The pixels between the two ends are fully covered.
	for x := first + 1; x < last; x++ {
		blend(row[x*4:x*4+4], c, 1)
	}
	blend(row[last*4:last*4+4], c, cover(last))
}

// blend mixes color c over the premultiplied pixel p with coverage cov.
func blend(p []byte, c rgba, cov float64) {
	a := int(math.Round(float64(c.a) * cov))
	if a <= 0 {
		return
	}
	if a >= 255 {
		p[0], p[1], p[2], p[3] = c.r, c.g, c.b, 255
		return
	}
	inv := 255 - a
	p[0] = uint8((int(c.r)*a + int(p[0])*inv + 127) / 255)
	p[1] = uint8((int(c.g)*a + int(p[1])*inv + 127) / 255)
	p[2] = uint8((int(c.b)*a + int(p[2])*inv + 127) / 255)
	p[3] = uint8((255*a + int(p[3])*inv + 127) / 255)
}
