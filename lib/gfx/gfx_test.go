// Package gfx has no Go code: the graphics libraries (raster.liz, font.liz)
// are lizard-script. This file runs the programs in tests/ through the
// headless painter and checks the pixels and the printed output.
package gfx

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"linkz/host"
)

// paint runs source on a w by h buffer and returns the picture and what the
// program printed. Set GFX_PNG_DIR to keep the pictures.
func paint(t *testing.T, name, source string, w, h int) (image.Image, string) {
	t.Helper()
	dir := t.TempDir()
	if keep := os.Getenv("GFX_PNG_DIR"); keep != "" {
		dir = keep
	}
	out := filepath.Join(dir, name+".png")
	var stdout bytes.Buffer
	if err := host.Paint(source, out, w, h, &stdout); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	f, err := os.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	return img, stdout.String()
}

func scene(t *testing.T, file string, w, h int) (image.Image, string) {
	t.Helper()
	src, err := os.ReadFile(filepath.Join("tests", file+".liz"))
	if err != nil {
		t.Fatal(err)
	}
	return paint(t, file, string(src), w, h)
}

type px struct {
	x, y int
	want [3]uint8 // the color
	tol  uint8    // how far each channel may be off (partial coverage)
	note string
}

var (
	white = [3]uint8{255, 255, 255}
	red   = [3]uint8{255, 0, 0}
	blue  = [3]uint8{0, 0, 255}
)

func check(t *testing.T, img image.Image, pixels []px) {
	t.Helper()
	for _, p := range pixels {
		r, g, b, _ := img.At(p.x, p.y).RGBA()
		got := [3]uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
		for i := range got {
			d := int(got[i]) - int(p.want[i])
			if d < 0 {
				d = -d
			}
			if d > int(p.tol) {
				t.Errorf("(%d,%d) %s: got %v, want %v (tolerance %d)", p.x, p.y, p.note, got, p.want, p.tol)
				break
			}
		}
	}
}

func TestRect(t *testing.T) {
	img, out := scene(t, "rect", 120, 80)
	// 30 rows (10 to 39), each from 10.0 to 50.0.
	if want := "30 [10, 10.0, 50.0] [39, 10.0, 50.0]\n"; out != want {
		t.Errorf("output %q, want %q", out, want)
	}
	check(t, img, []px{
		{10, 10, red, 0, "first pixel"},
		{49, 39, red, 0, "last pixel"},
		{30, 25, red, 0, "inside"},
		{9, 25, white, 0, "left of it"},
		{50, 25, white, 0, "right of it"},
		{30, 9, white, 0, "above"},
		{30, 40, white, 0, "below: no off-by-one on the last row"},
		// 60.5 to 100.25: half of pixel 60 and a quarter of pixel 100.
		{59, 25, white, 0, "before the soft edge"},
		{60, 25, [3]uint8{128, 128, 255}, 4, "left edge covers half"},
		{61, 25, blue, 0, "full"},
		{99, 25, blue, 0, "full"},
		{100, 25, [3]uint8{191, 191, 255}, 4, "right edge covers a quarter"},
		{101, 25, white, 0, "after the soft edge"},
		// Cut by the buffer: x from -20 to 30, y from 50 to 120.
		{0, 50, [3]uint8{0, 170, 0}, 0, "clipped shape starts at its row"},
		{29, 79, [3]uint8{0, 170, 0}, 0, "clipped shape, last buffer row"},
		{30, 60, white, 0, "clipped shape, right edge"},
	})
}

func TestCircle(t *testing.T) {
	img, out := scene(t, "circle", 120, 120)
	// A radius 40 circle at (60, 60) takes rows 20 to 99, and the second
	// shape() call returns the very same cached result.
	if want := "80 null\ntrue\n"; out != want {
		t.Errorf("output %q, want %q", out, want)
	}
	green := [3]uint8{0x3f, 0xa3, 0x4d}
	check(t, img, []px{
		{60, 60, green, 0, "center"},
		{60, 21, green, 0, "just inside the top"},
		{60, 19, white, 0, "just above the top"},
		{60, 99, green, 0, "just inside the bottom"},
		{60, 100, white, 0, "just below the bottom"},
		{21, 60, green, 24, "just inside the left (soft edge)"},
		{18, 60, white, 0, "just outside the left"},
		{99, 60, green, 24, "just inside the right (soft edge)"},
		{101, 60, white, 0, "just outside the right"},
		{30, 30, white, 0, "corner of the box is outside the circle"},
	})
}

func TestStar(t *testing.T) {
	img, _ := scene(t, "star", 200, 100)
	check(t, img, []px{
		{50, 50, [3]uint8{0xcc, 0x22, 0}, 0, "nonzero fills the middle"},
		{150, 50, white, 0, "evenodd leaves the middle empty"},
		{50, 12, [3]uint8{0xcc, 0x22, 0}, 0, "nonzero arm"},
		{150, 12, [3]uint8{0xcc, 0x22, 0}, 0, "evenodd arm"},
		{150, 20, [3]uint8{0xcc, 0x22, 0}, 0, "evenodd arm, lower"},
		{5, 5, white, 0, "outside"},
	})
}

func TestStroke(t *testing.T) {
	img, _ := scene(t, "stroke", 200, 100)
	ink := [3]uint8{0x22, 0x33, 0x44}
	check(t, img, []px{
		// The middle of the first slope is on the line (20,70)-(60,20).
		{40, 45, ink, 0, "on the line"},
		{43, 47, ink, 0, "4 px from the line, inside half the width"},
		{48, 51, white, 0, "10 px from the line, outside"},
		{15, 70, ink, 0, "inside the round cap"},
		{12, 70, white, 0, "outside the cap"},
		{60, 20, ink, 0, "the sharp corner is round, not a gap"},
		{100, 66, ink, 0, "valley corner"},
		{100, 50, [3]uint8{0x22, 0x88, 0x33}, 0, "the dot"},
		{100, 57, white, 0, "outside the dot (radius 5)"},
		{5, 5, white, 0, "outside"},
	})
}

func TestBrush(t *testing.T) {
	img, _ := scene(t, "brush", 200, 100)
	ink := [3]uint8{0x1d, 0x1d, 0x1f}
	check(t, img, []px{
		{70, 30, ink, 0, "widest point (width 14)"},
		{70, 21, white, 0, "outside the widest point"},
		{5, 5, white, 0, "outside"},
	})
}

func TestText(t *testing.T) {
	img, out := scene(t, "text", 200, 100)
	if want := "(35, 13) (70, 26) (63, 78) 26\n"; out != want {
		t.Errorf("output %q, want %q", out, want)
	}
	// Some ink in the first line, none past the end of the text.
	dark := func(x0, y0, x1, y1 int) int {
		n := 0
		for y := y0; y < y1; y++ {
			for x := x0; x < x1; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				if r>>8 < 128 && g>>8 < 128 && b>>8 < 128 {
					n++
				}
			}
		}
		return n
	}
	if n := dark(4, 4, 4+13*7, 17); n < 40 {
		t.Errorf("Hello, Linkz! has only %d dark pixels", n)
	}
	if n := dark(4+13*7, 4, 100, 17); n != 0 {
		t.Errorf("%d dark pixels after the end of the text", n)
	}
	if n := dark(100, 4, 100+7*7, 17); n < 30 {
		t.Errorf("non-ASCII line has only %d dark pixels", n)
	}
}

// '!' is a bar of rows 2 to 8 and a dot on row 10, in columns 3 to 4, so its
// pixels are known without reading the font table.
func TestGlyphPixels(t *testing.T) {
	source := `
fb.fill((0, 0), (60, 60), #ffffff)
font.draw("!", (10, 10), #000000, 1)
font.draw("!", (30, 10), #000000, 2)
font.draw("!", (44, 10), #000000, 9)
`
	img, _ := paint(t, "glyph", source, 60, 60)
	black := [3]uint8{}
	check(t, img, []px{
		// scale 1: column 13, rows 12 to 18 and row 20
		{13, 12, black, 0, "top of the bar"},
		{13, 18, black, 0, "bottom of the bar"},
		{13, 19, white, 0, "gap"},
		{13, 20, black, 0, "the dot"},
		{12, 15, white, 0, "left of the bar"},
		{14, 15, white, 0, "right of the bar"},
		// scale 2: columns 36 and 37, rows 14 to 27 (10 + 2 * 2 to 10 + 2 * 8 + 1), gap 28 to 29, dot 30 to 31
		{36, 14, black, 0, "scaled bar top"},
		{37, 14, black, 0, "scaled bar is two wide"},
		{36, 27, black, 0, "scaled bar bottom"},
		{36, 28, white, 0, "scaled gap"},
		{36, 29, white, 0, "scaled gap"},
		{36, 30, black, 0, "scaled dot"},
		{36, 31, black, 0, "scaled dot is two tall"},
		// a scale of 9 counts as 3: the bar is 3 wide at columns 53 to 55
		{53, 20, black, 0, "clamped scale"},
		{55, 20, black, 0, "clamped scale"},
		{56, 20, white, 0, "clamped scale: not wider"},
		{38, 20, white, 0, "right of the scaled bar"},
	})
}
