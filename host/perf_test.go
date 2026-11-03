package host

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"linkz/host/window"
)

// The editor's speed on big drawings. The benchmarks open a generated
// drawing headless and time the four things a user feels: a full repaint, a
// small edit, one frame of dragging a shape and one keystroke.
//
//	go test ./host -run '^$' -bench Editor -benchtime 20x -cpuprofile cpu.out
//
// go test ./host -run TestPerfFixtures -update writes the drawings to
// docs/Batches/fixtures, so they can be opened in the editor too.

var updateFixtures = flag.Bool("update", false, "rewrite the perf fixtures in docs/Batches/fixtures")

// perfCanvas is the drawing's size: it fits the perf window at 100%.
const perfCanvas = "(800, 600)"

// perfDrawing is a drawing of n literal shapes in layers of 50: circles,
// rects, strokes and polygons, spread over the canvas by a fixed
// pseudo-random sequence, so every run gets the same file. The first shape
// row (the frontmost) is circle((400, 300), 20, ...), the one the benchmarks
// drag and edit.
func perfDrawing(n int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "-~> lzdraw 1\n-~> %d shapes for the perf benchmarks (host/perf_test.go).\n\n", n)
	colors := []string{"#3fa34d", "#2d6a8f", "#e0a526", "#c0392b80", "#3b2a6b", "#f2a65a66"}
	seed := uint32(12345)
	rnd := func(k int) int {
		seed = seed*1664525 + 1013904223
		return int(seed>>8) % k
	}
	for i := 0; i < n; i++ {
		if i%50 == 0 {
			fmt.Fprintf(&b, "layer g%d\n", i/50)
		}
		x, y := 20+rnd(760), 20+rnd(560)
		c := colors[rnd(len(colors))]
		r := 6 + rnd(24)
		if i == 0 {
			x, y, r = 400, 300, 20
		}
		switch i % 4 {
		case 0:
			fmt.Fprintf(&b, "  circle((%d, %d), %d, %s)\n", x, y, r, c)
		case 1:
			fmt.Fprintf(&b, "  rect((%d, %d), (%d, %d), %s)\n", x, y, 8+rnd(40), 8+rnd(40), c)
		case 2:
			fmt.Fprintf(&b, "  stroke([(%d, %d), (%d, %d), (%d, %d), (%d, %d)], %s, %d)\n",
				x, y, x+10+rnd(20), y-rnd(30), x+30+rnd(20), y+rnd(30), x+50+rnd(20), y-rnd(20), c, 1+rnd(5))
		default:
			fmt.Fprintf(&b, "  polygon([(%d, %d), (%d, %d), (%d, %d)], %s)\n", x, y, x+10+rnd(30), y+rnd(30), x-rnd(20), y+10+rnd(30), c)
		}
	}
	fmt.Fprintf(&b, "\ncanvas = {\"size\": %s, \"background\": #fbf8ef}\n", perfCanvas)
	return b.String()
}

func TestPerfFixtures(t *testing.T) {
	if !*updateFixtures {
		t.Skip("run with -update to write the fixtures")
	}
	dir := filepath.Join("..", "docs", "Batches", "fixtures")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, n := range []int{2000, 10000} {
		path := filepath.Join(dir, fmt.Sprintf("shapes-%dk.lzdraw", n/1000))
		if err := os.WriteFile(path, []byte(perfDrawing(n)), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// editorRig is the editor, headless, on one drawing. It ticks the app by
// hand, so every measurement is one whole tick: events, run, paint. It adds
// up the lizard side's own phases of every tick (see lib/ui/perf.liz).
type editorRig struct {
	app    *App
	phases map[string]float64
	ticks  int
}

func newEditorRig(tb testing.TB, code string) *editorRig {
	tb.Helper()
	tb.Setenv("HOME", tb.TempDir())
	path := filepath.Join(tb.TempDir(), "drawing.lzdraw")
	if err := os.WriteFile(path, []byte(code), 0o644); err != nil {
		tb.Fatal(err)
	}
	a := NewApp(1600, 1000)
	a.headless = true
	if err := a.Load(path); err != nil {
		tb.Fatal(err)
	}
	r := &editorRig{app: a}
	r.tick()
	return r
}

// window is the window pixel of a point of the drawing.
func (r *editorRig) window(tb testing.TB, x, y int) (int, int) {
	tb.Helper()
	at, err := r.app.eval(fmt.Sprintf("view_to_window(view_canvas_layout(), (%d, %d))", x, y), nil)
	if err != nil {
		tb.Fatal(err)
	}
	var wx, wy float64
	fmt.Sscanf(at.Inspect(), "(%g, %g)", &wx, &wy)
	return int(wx), int(wy)
}

// tick runs one round of the loop and returns how long it took.
func (r *editorRig) tick(evs ...window.Event) time.Duration {
	start := time.Now()
	r.app.Tick(window.Input{Events: evs})
	if r.app.lastErr != "" {
		panic(r.app.lastErr)
	}
	took := time.Since(start)
	if r.phases != nil {
		for k, v := range r.app.phases() {
			if ms, ok := v.(float64); ok {
				r.phases[k] += ms
			}
		}
		r.ticks++
	}
	return took
}

// measure starts adding up the phases of the ticks that follow.
func (r *editorRig) measure() {
	r.phases = map[string]float64{}
	r.ticks = 0
}

// report gives the benchmark each phase's milliseconds per tick.
func (r *editorRig) report(b *testing.B) {
	b.StopTimer()
	for k, v := range r.phases {
		if r.ticks > 0 {
			b.ReportMetric(v/float64(r.ticks), k+"-ms")
		}
	}
	r.phases = nil
}

func (r *editorRig) lizard(tb testing.TB, code string) {
	tb.Helper()
	if err := r.app.call(code, nil); err != nil {
		tb.Fatal(err)
	}
}

// heapMB is the live heap after a collection, in megabytes. The benchmark
// reports what the editor holds once the drawing is open and painted.
func heapMB() float64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.HeapAlloc) / (1 << 20)
}

// The line of the dragged circle: the first shape row, after the header and
// its layer line.
const perfFirstShape = 5

func BenchmarkEditor(b *testing.B) {
	for _, n := range []int{2000, 10000} {
		before := heapMB()
		r := newEditorRig(b, perfDrawing(n))
		heap := heapMB() - before
		b.Run(fmt.Sprintf("%dk/full_repaint", n/1000), func(b *testing.B) {
			b.ReportMetric(heap, "heap-MB")
			r.measure()
			for i := 0; i < b.N; i++ {
				// Every pixel of the drawing again (the span cache stays warm).
				r.lizard(b, "if (V.scene != null) { tiles_mark_all(V.scene.tiles) }\nview_force_all()\napp.needs_paint = true")
				r.tick()
			}
			r.report(b)
		})
		b.Run(fmt.Sprintf("%dk/small_edit", n/1000), func(b *testing.B) {
			r.lizard(b, fmt.Sprintf("ed_select(%d)\nview_after()", perfFirstShape))
			r.tick()
			r.measure()
			for i := 0; i < b.N; i++ {
				r.lizard(b, fmt.Sprintf("ed_set_prop(\"r\", \"%d\")\nview_after()", 10+i%2))
				r.tick()
			}
			r.report(b)
		})
		b.Run(fmt.Sprintf("%dk/drag_frame", n/1000), func(b *testing.B) {
			x, y := r.window(b, 400, 300)
			r.tick(window.Event{Type: "down", X: x, Y: y})
			r.measure()
			for i := 0; i < b.N; i++ {
				dx := 4 + 2*(i%8)
				r.tick(window.Event{Type: "move", X: x + dx, Y: y + dx/2, Injected: true})
			}
			r.report(b)
			r.tick(window.Event{Type: "up", X: x, Y: y})
			r.lizard(b, "ed_undo()\nview_after()")
			r.tick()
		})
		b.Run(fmt.Sprintf("%dk/keystroke", n/1000), func(b *testing.B) {
			// The caret goes right after the circle's radius; typing a digit
			// and taking it back alternate, so the drawing stays the same size.
			r.lizard(b, fmt.Sprintf("ed_set_caret(%d, 23)\nview_after()", perfFirstShape))
			r.tick()
			r.measure()
			for i := 0; i < b.N; i++ {
				if i%2 == 0 {
					r.tick(window.Event{Type: "text", Text: "1", Injected: true})
				} else {
					r.tick(window.Event{Type: "key", Key: "Backspace", Injected: true})
				}
			}
			r.report(b)
		})
	}
}
