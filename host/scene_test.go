package host

import (
	"bytes"
	"fmt"
	"image"
	"strings"
	"testing"

	"linkz/host/window"
)

// The canvas keeps the drawing as tiles and repaints only what an edit
// touched (lib/ui/view_scene.liz). These tests drive the editor through edits
// and check, after each one, that the canvas is pixel for pixel what a full
// repaint of every tile gives, and close to what raster.paint draws in one go.

// sceneDrawings are the drawings the golden test edits: one of literal rows
// only (edits patch one shape), one with a let, a gradient and blend layers
// (every edit runs the whole drawing), and one that is zoomed out.
var sceneDrawings = map[string]string{
	"literal": perfDrawing(300),
	"computed": `-~> lzdraw 1
let y = 300
layer glow
  opacity(0.7)
  circle((400, y), 20, #e0a526)
  circle((430, y + 10), 60, #c0392b80)
  blend("multiply")
layer back
  rect((100, 100), (300, 200), radial((250, 200), 150, [#fff3c4, #2d6a8f]))
  stroke([(50, 500), (300, 400), (500, 550), (700, 420)], #1d1d1f, 6)
  brush([(60, 60), (200, 120), (320, 70)], [0.2, 0.9, 0.4], #3b2a6b, pens.ink)
  polygon([(600, 100), (700, 150), (650, 250)], #3fa34d)
canvas = {"size": (800, 600), "background": #fbf8ef}
`,
	"zoomed": strings.Replace(perfDrawing(200), `"size": (800, 600)`, `"size": (1300, 1100)`, 1),
}

// sceneStep is one edit: lizard code to run, or events, then a tick.
type sceneStep struct {
	name   string
	lizard string
	events func(r *editorRig, t *testing.T) []window.Event
}

// dragFrom presses on drawing point (x, y), moves by (dx, dy) in steps, releases.
func dragFrom(x, y, dx, dy int) func(r *editorRig, t *testing.T) []window.Event {
	return func(r *editorRig, t *testing.T) []window.Event {
		wx, wy := r.window(t, x, y)
		evs := []window.Event{{Type: "down", X: wx, Y: wy}}
		for i := 1; i <= 4; i++ {
			evs = append(evs, window.Event{Type: "move", X: wx + dx*i/4, Y: wy + dy*i/4, Injected: true})
		}
		return append(evs, window.Event{Type: "up", X: wx + dx, Y: wy + dy})
	}
}

var sceneSteps = []sceneStep{
	{name: "select", lizard: "ed_select(5)\nview_after()"},
	{name: "prop", lizard: "ed_set_prop(\"r\", \"35\")\nview_after()"},
	{name: "drag shape", events: dragFrom(400, 300, 60, 40)},
	{name: "drag far", events: dragFrom(460, 340, -300, 200)},
	{name: "type", lizard: "ed_set_caret(5, 10)\nview_after()\ncmp_ed_type(\"7\")\nview_prepare()\nview_after()"},
	{name: "hide row", lizard: "ed_set_hidden(6, true)\nview_prepare()\nview_after()"},
	{name: "show row", lizard: "ed_set_hidden(6, false)\nview_prepare()\nview_after()"},
	{name: "new row", lizard: "ed_set_caret(5, 0)\ned_type(\"  rect((380, 280), (90, 70), #3b2a6b)\\n\")\nview_prepare()\nview_after()"},
	{name: "move rows", lizard: "ed_set_caret(5, 0)\ned_key(\"ArrowDown\", {\"alt\": true})\nview_prepare()\nview_after()"},
	{name: "delete row", lizard: "ed_select(6)\ned_key(\"Backspace\", null)\nview_prepare()\nview_after()"},
	{name: "undo", lizard: "ed_undo()\nview_prepare()\nview_after()"},
	{name: "undo again", lizard: "ed_undo()\nview_prepare()\nview_after()"},
	{name: "caret away", lizard: "ed_set_caret(1, 0)\nview_after()"},
}

// canvasPixels copies the drawing's part of the canvas buffer.
func canvasPixels(t *testing.T, r *editorRig) *image.RGBA {
	t.Helper()
	v, err := r.app.eval("let __cl = view_canvas_layout()\nreturn [__cl.offset, (__cl.sw, __cl.sh)]", nil)
	if err != nil {
		t.Fatal(err)
	}
	var x, y, w, h int
	fmt.Sscanf(v.Inspect(), "[(%d, %d), (%d, %d)]", &x, &y, &w, &h)
	src := r.app.FB.canvas.Img.SubImage(image.Rect(x, y, x+w, y+h)).(*image.RGBA)
	out := image.NewRGBA(image.Rect(0, 0, w, h))
	for row := 0; row < h; row++ {
		copy(out.Pix[row*out.Stride:row*out.Stride+w*4], src.Pix[row*src.Stride:row*src.Stride+w*4])
	}
	return out
}

// fullRepaint paints every tile again and returns the canvas.
func fullRepaint(t *testing.T, r *editorRig) *image.RGBA {
	t.Helper()
	r.lizard(t, "tiles_mark_all(V.scene.tiles)\nview_force_all()\napp.needs_paint = true")
	r.tick()
	return canvasPixels(t, r)
}

// tilesPixels is the drawing as the tiles hold it, without the selection
// drawn over it on the canvas.
func tilesPixels(t *testing.T, r *editorRig) *image.RGBA {
	t.Helper()
	r.lizard(t, "fb.offscreen(V.scene.size)\ntiles_blit(V.scene.tiles, (0, 0))\nfb.keep(\"tiled\")")
	return r.app.FB.Kept("tiled").Img
}

// oneGo draws the drawing with raster.paint into a buffer of its own, the way
// export and the thumbnails do.
func oneGo(t *testing.T, r *editorRig) *image.RGBA {
	t.Helper()
	r.lizard(t, `
let __cl = view_canvas_layout()
let __shapes = []
let __list = shapes_to_paint(ed.shapes, view_rows())
for (let i = 0; i < len(__list); i += 1) {
  push(__shapes, view_scale_shape(__list[i], __cl.zoom))
}
fb.offscreen((__cl.sw, __cl.sh))
fb.fill((0, 0), (__cl.sw, __cl.sh), ed.canvas.background)
raster.paint(__shapes, (0, 0))
fb.keep("golden")`)
	return r.app.FB.Kept("golden").Img
}

// maxDiff is the largest channel difference between two pictures, and how
// many pixels differ at all.
func maxDiff(a, b *image.RGBA) (int, int) {
	worst, n := 0, 0
	for i := 0; i < len(a.Pix); i += 4 {
		d := 0
		for c := 0; c < 4; c++ {
			d = max(d, abs(int(a.Pix[i+c])-int(b.Pix[i+c])))
		}
		if d > 0 {
			n++
		}
		worst = max(worst, d)
	}
	return worst, n
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func TestSceneMatchesFullRepaint(t *testing.T) {
	for name, code := range sceneDrawings {
		t.Run(name, func(t *testing.T) {
			r := newEditorRig(t, code)
			for _, st := range sceneSteps {
				var evs []window.Event
				if st.events != nil {
					evs = st.events(r, t)
				} else {
					r.lizard(t, st.lizard)
				}
				r.tick(evs...)
				got := canvasPixels(t, r)
				want := fullRepaint(t, r)
				if !bytes.Equal(got.Pix, want.Pix) {
					worst, n := maxDiff(got, want)
					t.Fatalf("%s: the repainted tiles differ from a full repaint in %d pixels (by up to %d)", st.name, n, worst)
				}
				if worst, n := maxDiff(tilesPixels(t, r), oneGo(t, r)); worst > 1 {
					t.Fatalf("%s: the tiles differ from raster.paint in %d pixels (by up to %d)", st.name, n, worst)
				}
			}
		})
	}
}

// After a drag of one shape only the tiles around it are painted again, and
// a literal drawing is not run again as a whole.
func TestSceneRepaintsOnlyWhatChanged(t *testing.T) {
	r := newEditorRig(t, perfDrawing(300))
	r.lizard(t, "ed_select(5)\nview_after()")
	r.tick()
	evs := dragFrom(400, 300, 6, 4)(r, t)
	for _, ev := range evs[:len(evs)-1] {
		r.tick(ev)
	}
	v, err := r.app.eval(`[V.scene.tiles.painted, V.scene.tiles.area < 80 * 80, ed.patch != null]`, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := v.Inspect(); got != "[1, true, true]" {
		t.Errorf("after a small drag [tiles painted, small area, patched] = %s, want [1, true, true]", got)
	}
}
