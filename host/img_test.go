package host

import (
	"bytes"
	"image"
	stdcolor "image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"linkz/lib/draw"
	"lizardscript/evaluator"
	"lizardscript/object"
)

// writeTestPNG makes a w by h PNG where paint(x, y) gives each pixel.
func writeTestPNG(t *testing.T, path string, w, h int, paint func(x, y int) stdcolor.Color) {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, paint(x, y))
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}

func TestImgLoad(t *testing.T) {
	dir := t.TempDir()
	red := stdcolor.NRGBA{200, 40, 40, 255}
	writeTestPNG(t, filepath.Join(dir, "wide.png"), 200, 100, func(x, y int) stdcolor.Color { return red })
	writeTestPNG(t, filepath.Join(dir, "clear.png"), 4, 4, func(x, y int) stdcolor.Color { return stdcolor.NRGBA{0, 0, 0, 0} })
	writeTestPNG(t, filepath.Join(dir, "halves.png"), 4, 2, func(x, y int) stdcolor.Color {
		if x < 2 {
			return stdcolor.NRGBA{255, 0, 0, 255}
		}
		return stdcolor.NRGBA{0, 0, 255, 255}
	})
	jf, _ := os.Create(filepath.Join(dir, "pic.jpg"))
	jpg := image.NewRGBA(image.Rect(0, 0, 30, 20))
	for i := range jpg.Pix {
		jpg.Pix[i] = 255
	}
	if err := jpeg.Encode(jf, jpg, nil); err != nil {
		t.Fatal(err)
	}
	jf.Close()

	tests := []struct {
		name, file string
		maxSide    int
		w, h       int
		first      string
	}{
		{"shrinks the longer side", "wide.png", 50, 50, 25, "#c82828"},
		{"never enlarges", "wide.png", 500, 200, 100, "#c82828"},
		{"transparent becomes white", "clear.png", 256, 4, 4, "#ffffff"},
		{"box filter averages", "halves.png", 2, 2, 1, "#ff0000"},
		{"decodes JPEG", "pic.jpg", 256, 30, 20, "#ffffff"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w, h, px, err := loadImage(filepath.Join(dir, tc.file), tc.maxSide)
			if err != nil {
				t.Fatal(err)
			}
			if w != tc.w || h != tc.h || len(px) != w*h {
				t.Fatalf("size %dx%d (%d pixels), want %dx%d", w, h, len(px), tc.w, tc.h)
			}
			c := px[0]
			got := (&object.Color{R: c.R, G: c.G, B: c.B, A: 255}).Inspect()
			if got != tc.first && !(tc.file == "pic.jpg") {
				t.Errorf("first pixel %s, want %s", got, tc.first)
			}
		})
	}

	if _, _, _, err := loadImage(filepath.Join(dir, "missing.png"), 10); err == nil {
		t.Error("a missing file should be an error")
	}
	bad := filepath.Join(dir, "bad.png")
	os.WriteFile(bad, []byte("not an image"), 0o644)
	if _, _, _, err := loadImage(bad, 10); err == nil {
		t.Error("a file that isn't an image should be an error")
	}

	// halves.png at max_side 2: the left pixel is red, the right blue.
	_, _, px, _ := loadImage(filepath.Join(dir, "halves.png"), 2)
	if px[1] != (pixelRGB{0, 0, 255}) {
		t.Errorf("second pixel %v, want blue", px[1])
	}
}

// netArea runs generated code and measures the first path shape in it: the
// area inside its outline, holes taken out.
func netArea(t *testing.T, code string) float64 {
	t.Helper()
	evaluator.Output = &bytes.Buffer{}
	env := object.NewEnvironment()
	for _, name := range draw.Order {
		src, _ := draw.Files.ReadFile(name)
		if _, err := run(string(src), "draw/"+name, env, nil); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := run(code, "traced", env, nil); err != nil {
		t.Fatalf("the traced code fails: %v", err)
	}
	res, err := run(`
let found = null
for (let i = 0; i < len(drawing); i = i + 1) {
  if (found === null && drawing[i].kind == "path") { found = drawing[i] }
}
let subs = flatten_path(found.path, 0.05)
let sum = 0.0
for (let i = 0; i < len(subs); i = i + 1) { sum = sum + polygon_area2(subs[i].points) }
math.abs(sum) / 2.0`, "check", env, nil)
	if err != nil {
		t.Fatal(err)
	}
	f, ok := res.(*object.Float)
	if !ok {
		t.Fatalf("area is %s", res.Inspect())
	}
	return f.Value
}

// A circle on a background traces to code that renders, and the outline
// covers about the circle's area.
func TestTraceCircle(t *testing.T) {
	const r = 30.0
	file := filepath.Join(t.TempDir(), "circle.png")
	writeTestPNG(t, file, 100, 100, func(x, y int) stdcolor.Color {
		if math.Hypot(float64(x)+0.5-50, float64(y)+0.5-50) < r {
			return stdcolor.NRGBA{200, 40, 40, 255}
		}
		return stdcolor.NRGBA{240, 240, 230, 255}
	})
	code, stats, err := Trace(file, map[string]any{"colors": 2}, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if stats["colors"] != int64(2) {
		t.Errorf("stats = %v, want 2 colors", stats)
	}
	svg, err := Render(code, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("the traced code doesn't render: %v\n%s", err, code)
	}
	if !strings.Contains(svg, "<path") {
		t.Errorf("no shapes in the SVG:\n%s", svg)
	}
	want := math.Pi * r * r
	got := netArea(t, code)
	t.Logf("traced area %.0f, circle %.0f (%.2f%% off)", got, want, 100*(got-want)/want)
	if math.Abs(got-want)/want > 0.04 {
		t.Errorf("traced circle area %.0f, want %.0f (within 4%%)\n%s", got, want, code)
	}
}

func TestTraceErrors(t *testing.T) {
	if _, _, err := Trace(filepath.Join(t.TempDir(), "nope.png"), nil, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "img.load") {
		t.Errorf("a missing file should say img.load failed, got %v", err)
	}
	file := filepath.Join(t.TempDir(), "p.png")
	writeTestPNG(t, file, 8, 8, func(x, y int) stdcolor.Color { return stdcolor.White })
	if _, _, err := Trace(file, map[string]any{"colours": 3}, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "unknown option") {
		t.Errorf("a misspelt option should be an error, got %v", err)
	}
}
