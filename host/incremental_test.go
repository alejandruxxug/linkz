package host

import (
	"bytes"
	"strings"
	"testing"

	"linkz/host/window"
)

// The editor works out what it shows from the last text when it can: the
// parse and the rows by splicing (lib/ui/doc_splice.liz), the panel's items
// (view_items.liz), the shapes by running only the statements an edit
// touched (editor_rerun.liz). These tests drive it through edits that add
// and take away lines, break and mend the code, and after each one check
// that all of it is what working it out from scratch gives.

// incrementalSteps are the edits, on perfDrawing(300) (one layer line per
// 50 shapes: line 30 is a shape row inside layer g0).
var incrementalSteps = []sceneStep{
	{name: "caret", lizard: "ed_set_caret(30, 1000)\nview_after()"},
	{name: "enter", events: keyStep("Enter")},
	{name: "type in the new row", events: textStep("circle((10, 10), 5, #000)")},
	{name: "enter again", events: keyStep("Enter")},
	{name: "empty row with the caret", lizard: "view_after()"},
	{name: "join the rows back", events: keyStep("Backspace", "Backspace", "Backspace")},
	{name: "select a row", lizard: "ed_select_range(40, 40)\nview_after()"},
	{name: "delete it", events: keyStep("Backspace")},
	{name: "paste rows", lizard: "__clip.text = edit_rows_text(ed_model(), 60, 69)\ned_set_caret(45, 0)\nview_after()"},
	{name: "cmd-v", events: func(r *editorRig, t *testing.T) []window.Event {
		return []window.Event{{Type: "key", Key: "V", Mods: window.Mods{Cmd: true}, Injected: true}}
	}},
	{name: "break a row", lizard: "ed_set_caret(50, 1000)\nview_after()"},
	{name: "take its bracket off", events: keyStep("Backspace")},
	{name: "type while broken", events: textStep(" 1")},
	{name: "take that back", events: keyStep("Backspace", "Backspace")},
	{name: "mend it", events: textStep(")")},
	{name: "a row with an error", lizard: "ed_set_caret(52, 1000)\nview_after()"},
	{name: "an unknown name", events: textStep(" + nope")},
	{name: "a new row under it", events: keyStep("Enter")},
	{name: "mend that", events: keyStep("Backspace", "Backspace", "Backspace", "Backspace", "Backspace", "Backspace", "Backspace", "Backspace", "Backspace")},
	{name: "a layer line", lizard: "ed_set_caret(54, 0)\nview_after()"},
	{name: "new layer", events: textStep("layer extra\n")},
	{name: "undo", lizard: "ed_undo()\nview_after()"},
	{name: "undo again", lizard: "ed_undo()\nview_after()"},
	{name: "redo", lizard: "ed_redo()\nview_after()"},
	{name: "fold a layer", lizard: "ed_toggle_fold(4)\nview_after()"},
	{name: "enter while folded", lizard: "ed_set_caret(60, 1000)\nview_after()\ned_key(\"Enter\", null)\nview_after()"},
	{name: "unfold", lizard: "ed_toggle_fold(4)\nview_after()"},
	{name: "fold a layer further down", lizard: foldLayer("g2")},
	{name: "enter above it", lizard: "ed_set_caret(30, 1000)\nview_after()\ned_key(\"Enter\", null)\nview_after()"},
	{name: "and take it back", events: keyStep("Backspace", "Backspace", "Backspace")},
	{name: "an error further down", lizard: "ed_set_caret(200, 1000)\nview_after()\ned_type(\" + nope\")\nview_after()"},
	{name: "enter above the error", lizard: "ed_set_caret(30, 1000)\nview_after()\ned_key(\"Enter\", null)\nview_after()"},
}

// foldLayer is lizard that folds (or unfolds) the layer called name.
func foldLayer(name string) string {
	return `let __ls = doc_layers(ed_model())
for (let i = 0; i < len(__ls); i += 1) {
  if (__ls[i].name == "` + name + `") {
    ed_toggle_fold(__ls[i].line)
  }
}
view_after()`
}

func keyStep(keys ...string) func(r *editorRig, t *testing.T) []window.Event {
	return func(r *editorRig, t *testing.T) []window.Event {
		var evs []window.Event
		for _, k := range keys {
			evs = append(evs, window.Event{Type: "key", Key: k, Injected: true})
		}
		return evs
	}
}

func textStep(s string) func(r *editorRig, t *testing.T) []window.Event {
	return func(r *editorRig, t *testing.T) []window.Event {
		return []window.Event{{Type: "text", Text: s, Injected: true}}
	}
}

// incrementalCheck is lizard that compares what the editor holds with what a
// fresh start gives: [parse, rows, items, shapes] (each true when the same).
const incrementalCheck = `
let __m = ed_model()
let __fresh = doc_new(doc_text(__m))
let __p = doc_parsed(__m)
let __fp = doc_parsed(__fresh)
let __same_parse = len(__fp.errors) > 0 || (parse_flat(__p).statements === __fp.statements && len(__p.errors) == 0)
let __rows_plain = function(rows) {
  let out = []
  for (let i = 0; i < len(rows); i += 1) {
    let r = hash_copy(rows[i])
    row_stmt(r)
    push(out, r)
  }
  out
}
let __same_rows = len(__fp.errors) > 0 ||
  __rows_plain(doc_rows(__m, ed.errors, ed.folded)) === __rows_plain(doc_rows(__fresh, ed.errors, ed.folded))
let __items = view_items()
let __at = V.items_at
let __pinned = V.pinned
view_items_build(view_rows(), __m)
let __same_items = __items === V.items && __at === V.items_at && __pinned === V.pinned
let __run = src.run(doc_text(__m))
let __same_shapes = len(__run.errors) > 0 || ed.shapes === __run.shapes
[__same_parse, __same_rows, __same_items, __same_shapes]
`

// incrementalDrawings: literal rows only, and the same zoomed out with a
// hidden row (the canvas keeps the scaled shapes and leaves the hidden out).
var incrementalDrawings = map[string]string{
	"literal":        perfDrawing(300),
	"zoomed, hidden": strings.Replace(perfDrawing(300), `"size": (800, 600)`, `"size": (1300, 1100)`, 1),
}

func TestIncrementalMatchesFresh(t *testing.T) {
	for name, code := range incrementalDrawings {
		t.Run(name, func(t *testing.T) {
			if name != "literal" {
				// line 20, a shape row of layer g0
				lines := strings.Split(code, "\n")
				lines[19] += " -~> [hidden]"
				code = strings.Join(lines, "\n")
			}
			r := newEditorRig(t, code)
			r.fakeClipboard(t)
			for _, st := range incrementalSteps {
				var evs []window.Event
				if st.events != nil {
					evs = st.events(r, t)
				} else {
					r.lizard(t, st.lizard)
				}
				r.tick(evs...)
				v, err := r.app.eval(incrementalCheck, nil)
				if err != nil {
					t.Fatalf("%s: %v", st.name, err)
				}
				if got := v.Inspect(); got != "[true, true, true, true]" {
					t.Errorf("%s: [parse, rows, items, shapes] the same as from scratch = %s", st.name, got)
				}
				got := canvasPixels(t, r)
				want := fullRepaint(t, r)
				if !bytes.Equal(got.Pix, want.Pix) {
					worst, n := maxDiff(got, want)
					t.Fatalf("%s: the canvas differs from a full repaint in %d pixels (by up to %d)", st.name, n, worst)
				}
			}
		})
	}
}
