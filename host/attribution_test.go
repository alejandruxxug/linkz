package host

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// A loop that never ends stops at the step limit instead of hanging.
func TestRenderStepLimit(t *testing.T) {
	done := make(chan error, 1)
	go func() {
		_, err := Render("while (true) {}", &bytes.Buffer{})
		done <- err
	}()
	select {
	case err := <-done:
		if err == nil || !strings.Contains(err.Error(), "the program ran too long (more than 2000000 steps)") {
			t.Errorf("got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("render did not stop")
	}
}

func TestPaintStepLimit(t *testing.T) {
	out := filepath.Join(t.TempDir(), "o.png")
	done := make(chan error, 1)
	go func() { done <- Paint("while (true) {}", out, 8, 8, &bytes.Buffer{}) }()
	select {
	case err := <-done:
		if err == nil || !strings.Contains(err.Error(), "ran too long") {
			t.Errorf("got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("paint did not stop")
	}
}

// The step counter starts fresh on every run.
func TestRenderResetsSteps(t *testing.T) {
	src := "let n = 0\nfor (let i = 0; i < 200000; i = i + 1) { n = n + 1 }\ncircle((5, 5), 2, #ff0000)"
	for i := 0; i < 20; i++ {
		if _, err := Render(src, &bytes.Buffer{}); err != nil {
			t.Fatalf("run %d: %v", i, err)
		}
	}
}

func TestRenderCanvasErrors(t *testing.T) {
	tests := []struct{ name, src, want string }{
		{"negative", `canvas = {"size": (-5, 0), "background": #ffffff}`, "canvas: size must be between (1, 1) and (16384, 16384), got (-5, 0)"},
		{"huge", `canvas = {"size": (20000, 10), "background": #ffffff}`, "got (20000, 10)"},
		{"point", `canvas = {"size": 5}`, "canvas: size must be a point"},
		{"background", `canvas = {"size": (10, 10), "background": 3}`, "canvas: background must be a color or a gradient, got INTEGER"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Render(tc.src, &bytes.Buffer{})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("got %v, want %q", err, tc.want)
			}
			if err != nil && strings.Contains(err.Error(), "[line 1") {
				t.Errorf("a library position leaked: %v", err)
			}
		})
	}
}

// An error inside library code is reported at the user's line.
func TestRenderAttributesLibraryErrors(t *testing.T) {
	tests := []struct{ name, src, want string }{
		{"top level", "let a = 1\ncircle(1, 3, #ff0000)", "[line 2, column 1] runtime error: circle:"},
		{"in a layer", "layer a\n  let x = 1\n  circle(5, 5, 5)\n", "[line 3, column 3] runtime error: circle:"},
		{"user error keeps its place", "let a = 1\nnope", "[line 2, column 1] runtime error: identifier not found: nope"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Render(tc.src, &bytes.Buffer{})
			if err == nil || !strings.HasPrefix(err.Error(), tc.want) {
				t.Errorf("got %v, want prefix %q", err, tc.want)
			}
		})
	}
}

func TestPaintAttributesAndSize(t *testing.T) {
	out := filepath.Join(t.TempDir(), "o.png")
	err := Paint("let a = 1\ncircle(5, 5, 5)", out, 8, 8, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "[line 2, column 1] runtime error: circle:") {
		t.Errorf("got %v", err)
	}
	for _, sz := range [][2]int{{0, 8}, {8, -1}, {20000, 8}} {
		if err := Paint("", out, sz[0], sz[1], &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "size must be between") {
			t.Errorf("size %v: got %v", sz, err)
		}
	}
	_ = os.Remove(out)
}
