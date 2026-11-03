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

	if !strings.HasPrefix(svg, "<svg xmlns='http://www.w3.org/2000/svg' width='400' height='200'") {
		t.Errorf("wrong start of the SVG:\n%s", svg)
	}
	if got := strings.Count(svg, "<circle"); got != 6 {
		t.Errorf("expected 6 circles, got %d:\n%s", got, svg)
	}
	if !strings.Contains(svg, "fill='#e0a526'") {
		t.Errorf("changing a shape after drawing it didn't reach the SVG:\n%s", svg)
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
