package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTraceCmd(t *testing.T) {
	dir := t.TempDir()
	img := image.NewNRGBA(image.Rect(0, 0, 20, 20))
	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			c := color.NRGBA{255, 255, 255, 255}
			if x >= 5 && x < 15 && y >= 5 && y < 15 {
				c = color.NRGBA{0, 0, 200, 255}
			}
			img.Set(x, y, c)
		}
	}
	in, out := filepath.Join(dir, "in.png"), filepath.Join(dir, "out.lzdraw")
	f, _ := os.Create(in)
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	f.Close()

	if err := traceCmd([]string{in, out, "--colors", "2", "--detail", "1"}); err != nil {
		t.Fatal(err)
	}
	code, err := os.ReadFile(out)
	if err != nil || !strings.Contains(string(code), "layer color_2") {
		t.Fatalf("no traced code written: %v\n%s", err, code)
	}
	if err := render([]string{out, filepath.Join(dir, "out.svg")}); err != nil {
		t.Fatalf("the traced file doesn't render: %v", err)
	}

	for name, args := range map[string][]string{
		"no files":           {},
		"unknown flag":       {in, out, "--bogus", "1"},
		"flag no number":     {in, out, "--colors", "many"},
		"missing flag value": {in, out, "--colors"},
	} {
		if err := traceCmd(args); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
}
