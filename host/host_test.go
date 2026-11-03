package host

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// The smoke test: the libraries load, a real example runs, and the SVG has
// what it drew.
func TestRenderHello(t *testing.T) {
	source, err := os.ReadFile("../examples/hello.lzdraw")
	if err != nil {
		t.Fatal(err)
	}

	var printed bytes.Buffer
	svg, err := Render(string(source), &printed)
	if err != nil {
		t.Fatal(err)
	}

	golden, err := os.ReadFile("../lib/draw/tests/golden/hello.svg.expected")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(svg) != strings.TrimSpace(string(golden)) {
		t.Errorf("the SVG differs from the golden file:\n%s", svg)
	}
	if got, want := printed.String(), "drew 7 shapes\n"; got != want {
		t.Errorf("puts printed %q, want %q", got, want)
	}
}

func TestRenderReportsErrors(t *testing.T) {
	for name, source := range map[string]string{
		"syntax":  "let = 5",
		"runtime": "circle(1, 2, nope)",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := Render(source, &bytes.Buffer{})
			if err == nil || !strings.Contains(err.Error(), "line 1") {
				t.Errorf("expected an error at line 1, got %v", err)
			}
		})
	}
}
