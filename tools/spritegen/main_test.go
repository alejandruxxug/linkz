package main

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestRuns(t *testing.T) {
	leaf := color.NRGBA{0x3f, 0xa3, 0x4d, 0xff}
	ink := color.NRGBA{0x1b, 0x1a, 0x17, 0xff}
	tests := []struct {
		name string
		row  []color.NRGBA // one row; zero value is transparent
		want string
		err  string
	}{
		{"empty", []color.NRGBA{{}, {}}, "", ""},
		{"one run", []color.NRGBA{{}, leaf, leaf, {}}, `[0, 1, 3, #3fa34d]`, ""},
		{"color change splits", []color.NRGBA{ink, leaf, leaf}, `[0, 0, 1, #1b1a17] [0, 1, 3, #3fa34d]`, ""},
		{"any opaque color", []color.NRGBA{{1, 2, 0xab, 0xff}}, `[0, 0, 1, #0102ab]`, ""},
		{"half transparent", []color.NRGBA{{0x3f, 0xa3, 0x4d, 0x80}}, "", "half-transparent"},
	}
	for _, tt := range tests {
		img := image.NewNRGBA(image.Rect(0, 0, len(tt.row), 1))
		for x, c := range tt.row {
			img.SetNRGBA(x, 0, c)
		}
		rs, err := runs(img)
		if tt.err != "" {
			if err == nil || !strings.Contains(err.Error(), tt.err) {
				t.Errorf("%s: want error %q, got %v", tt.name, tt.err, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: %v", tt.name, err)
			continue
		}
		if got := strings.Join(rs, " "); got != tt.want {
			t.Errorf("%s: got %s, want %s", tt.name, got, tt.want)
		}
	}
}
