package host

import (
	"encoding/json"
	"fmt"
	"lizardscript/evaluator"
	"lizardscript/object"
	"strings"
	"testing"
)

// asJSON makes a lizard result easy to compare.
func asJSON(t *testing.T, o object.Object) map[string]any {
	t.Helper()
	b, err := json.Marshal(toGo(o, 0))
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	json.Unmarshal(b, &m)
	return m
}

func list(v any) []any          { l, _ := v.([]any); return l }
func hash(v any) map[string]any { m, _ := v.(map[string]any); return m }
func num(v any) int             { f, _ := v.(float64); return int(f) }

func parse(t *testing.T, code string) map[string]any {
	t.Helper()
	return fromGoRoundTrip(t, parseResult(code))
}

func fromGoRoundTrip(t *testing.T, m map[string]any) map[string]any {
	t.Helper()
	b, _ := json.Marshal(m)
	var out map[string]any
	json.Unmarshal(b, &out)
	return out
}

func TestParseCallsAndPoints(t *testing.T) {
	code := `circle((380, 90), 40, #e0a526)
stroke([(40, 300), (120, 210), (260, 240)], #1d1d1f, 4)
`
	r := parse(t, code)
	stmts := list(r["statements"])
	if len(stmts) != 2 {
		t.Fatalf("got %d statements, want 2", len(stmts))
	}
	c := hash(stmts[0])
	if c["kind"] != "call" || c["callee"] != "circle" || num(c["line"]) != 1 || num(c["column"]) != 1 {
		t.Errorf("first statement = %v", c)
	}
	args := list(c["args"])
	want := []struct{ text, kind string }{{"(380, 90)", "point"}, {"40", "number"}, {"#e0a526", "color"}}
	for i, w := range want {
		a := hash(args[i])
		if a["text"] != w.text || a["kind"] != w.kind {
			t.Errorf("arg %d = %v, want %s %s", i, a, w.kind, w.text)
		}
	}
	// the point's span: "(380, 90)" starts at column 8 and ends before the comma
	p := hash(args[0])
	if num(p["start"]) != 8 || num(p["end"]) != 17 {
		t.Errorf("point span %v-%v, want 8-17", p["start"], p["end"])
	}
	if items := list(p["items"]); len(items) != 2 || hash(items[0])["text"] != "380" {
		t.Errorf("point items = %v", p["items"])
	}

	// one point inside an array literal can be rewritten
	arr := hash(list(hash(stmts[1])["args"])[0])
	if arr["kind"] != "array" {
		t.Fatalf("array arg kind = %v", arr["kind"])
	}
	items := list(arr["items"])
	if len(items) != 3 {
		t.Fatalf("array has %d items, want 3", len(items))
	}
	second := hash(items[1])
	line := strings.Split(code, "\n")[1]
	s, e := num(second["start"]), num(second["end"])
	if got := line[s-1 : e-1]; got != "(120, 210)" || second["text"] != "(120, 210)" {
		t.Errorf("second point span covers %q (text %q)", got, second["text"])
	}
}

func TestParseLayersCommentsFlags(t *testing.T) {
	code := `layer sun -~> [hidden]
  circle((1, 2), 3, #fff) -~> core [locked]
  layer inner
    rect((0, 0), (2, 2), #000)
circle((5, 5), 1, #000)
`
	r := parse(t, code)
	stmts := list(r["statements"])
	if len(stmts) != 2 {
		t.Fatalf("got %d top statements, want 2 (the layer and the loose circle)", len(stmts))
	}
	sun := hash(stmts[0])
	if sun["kind"] != "layer" || sun["name"] != "sun" || num(sun["end_line"]) != 4 {
		t.Errorf("layer = %v", sun)
	}
	if f := list(sun["flags"]); len(f) != 1 || f[0] != "hidden" {
		t.Errorf("layer flags = %v", sun["flags"])
	}
	kids := list(sun["children"])
	if len(kids) != 2 {
		t.Fatalf("layer has %d children, want 2", len(kids))
	}
	core := hash(kids[0])
	if core["label"] != "core" || num(core["depth"]) != 1 {
		t.Errorf("core = %v", core)
	}
	if f := list(core["flags"]); len(f) != 1 || f[0] != "locked" {
		t.Errorf("core flags = %v", core["flags"])
	}
	inner := hash(kids[1])
	if inner["kind"] != "layer" || len(list(inner["children"])) != 1 || num(hash(list(inner["children"])[0])["depth"]) != 2 {
		t.Errorf("nested layer = %v", inner)
	}
	lines := list(r["lines"])
	l2 := hash(lines[1])
	if num(l2["indent"]) != 2 || l2["name"] != "core" || l2["comment"] != "core [locked]" {
		t.Errorf("line 2 = %v", l2)
	}
	if hash(stmts[1])["kind"] != "call" || num(hash(stmts[1])["depth"]) != 0 {
		t.Errorf("loose circle after the layer = %v", stmts[1])
	}
}

func TestParseMultiLineAndKinds(t *testing.T) {
	code := `let horizon = 300
canvas = {"size": (512, 512)}
stroke([
  (1, 2),
  (3, 4)
], #000, 2)
for (let i = 0; i < 3; i = i + 1) {
  circle((i, i), 1, #000)
}
if (true) { rect((0, 0), (1, 1), #000) } else { puts("no") }
`
	r := parse(t, code)
	stmts := list(r["statements"])
	kinds := []string{"let", "assign", "call", "block", "block"}
	if len(stmts) != len(kinds) {
		t.Fatalf("got %d statements, want %d: %v", len(stmts), len(kinds), stmts)
	}
	for i, k := range kinds {
		if hash(stmts[i])["kind"] != k {
			t.Errorf("statement %d kind %v, want %s", i, hash(stmts[i])["kind"], k)
		}
	}
	multi := hash(stmts[2])
	if num(multi["line"]) != 3 || num(multi["end_line"]) != 6 {
		t.Errorf("multi-line call lines %v-%v, want 3-6", multi["line"], multi["end_line"])
	}
	first := hash(list(multi["args"])[0])
	if first["kind"] != "array" || len(list(first["items"])) != 2 || num(first["end_line"]) != 6 {
		t.Errorf("multi-line array arg = %v", first)
	}
	loop := hash(stmts[3])
	if kids := list(loop["children"]); len(kids) != 1 || hash(kids[0])["callee"] != "circle" || num(hash(kids[0])["depth"]) != 1 {
		t.Errorf("loop children = %v", loop["children"])
	}
	if hash(stmts[4])["name"] != "if" || len(list(hash(stmts[4])["children"])) != 2 {
		t.Errorf("if statement = %v", stmts[4])
	}
}

func TestParseDocsAndParams(t *testing.T) {
	code := `-~> Draws a star.
-~> Second sentence.
-~> center: where it sits
-~> r: the radius
-~>   in pixels
-~> example: star((1, 2), 5)
function star(center, r) {
  circle(center, r, #000)
}
`
	r := parse(t, code)
	fn := hash(list(r["statements"])[0])
	if fn["kind"] != "function" || fn["name"] != "star" {
		t.Fatalf("function = %v", fn)
	}
	if p := list(fn["params"]); len(p) != 2 || p[0] != "center" || p[1] != "r" {
		t.Errorf("params = %v", fn["params"])
	}
	doc := hash(fn["doc"])
	if doc["summary"] != "Draws a star. Second sentence." {
		t.Errorf("summary = %q", doc["summary"])
	}
	params := hash(doc["params"])
	if params["center"] != "where it sits" || params["r"] != "the radius in pixels" {
		t.Errorf("param docs = %v", params)
	}
	if doc["example"] != "star((1, 2), 5)" {
		t.Errorf("example = %v", doc["example"])
	}
	if kids := list(fn["children"]); len(kids) != 1 {
		t.Errorf("function body statements = %v", fn["children"])
	}
}

func TestParseErrorsAndUnicode(t *testing.T) {
	r := parse(t, "let = 5\n")
	errs := list(r["errors"])
	if len(errs) == 0 || num(hash(errs[0])["line"]) != 1 {
		t.Errorf("errors = %v, want one at line 1", errs)
	}
	// columns count characters, so spans stay right after a non-ASCII string
	r = parse(t, `text((1, 2), "ñandú", #000)`+"\n")
	args := list(hash(list(r["statements"])[0])["args"])
	if len(args) != 3 || hash(args[1])["text"] != `"ñandú"` || hash(args[1])["kind"] != "string" {
		t.Errorf("args after unicode = %v", args)
	}
	if hash(args[2])["text"] != "#000" || num(hash(args[2])["start"]) != 23 {
		t.Errorf("color arg after unicode = %v", args[2])
	}
}

// run runs a drawing through src.run and returns the result as JSON-like data.
func srcRun(t *testing.T, code string) map[string]any {
	t.Helper()
	return asJSON(t, NewSrc(nil).Run(code))
}

func shapeTags(r map[string]any) []string {
	var out []string
	for _, s := range list(r["shapes"]) {
		h := hash(s)
		out = append(out, fmt.Sprintf("%v@%v:%v", h["kind"], h["line"], h["layer"]))
	}
	return out
}

func TestRunTagsShapes(t *testing.T) {
	code := `function dot(p) {
  circle(p, 2, #000)
}
circle((1, 1), 1, #000)
layer sea
  rect((0, 0), (4, 4), #123456)
  layer waves
    for (let i = 0; i < 2; i = i + 1) {
      line((i, 0), (i, 9), #000, 1)
    }
dot((5, 5))
`
	r := srcRun(t, code)
	if errs := list(r["errors"]); len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	got := strings.Join(shapeTags(r), " ")
	want := "circle@4:[] rect@6:[sea] line@9:[sea waves] line@9:[sea waves] circle@11:[]"
	if got != want {
		t.Errorf("tags\n got %s\nwant %s", got, want)
	}
	if r["printed"] != "" {
		t.Errorf("printed = %q", r["printed"])
	}
}

func TestRunCapturesPutsAndErrors(t *testing.T) {
	r := srcRun(t, "puts(\"hi\", 1 + 2)\ncircle((1, 1), 1, #000)\nnope\n")
	if r["printed"] != "hi 3\n" {
		t.Errorf("printed = %q", r["printed"])
	}
	errs := list(r["errors"])
	if len(errs) != 1 || num(hash(errs[0])["line"]) != 3 || !strings.Contains(hash(errs[0])["message"].(string), "nope") {
		t.Errorf("errors = %v", errs)
	}
	if len(list(r["shapes"])) != 1 {
		t.Errorf("shapes before the error should still come back: %v", r["shapes"])
	}

	r = srcRun(t, "let = 5")
	if errs := list(r["errors"]); len(errs) == 0 || num(hash(errs[0])["line"]) != 1 {
		t.Errorf("syntax error = %v", errs)
	}
}

func TestRunStepLimit(t *testing.T) {
	r := srcRun(t, "while (true) { let x = 1 }\n")
	errs := list(r["errors"])
	if len(errs) != 1 || !strings.Contains(hash(errs[0])["message"].(string), "ran too long") {
		t.Errorf("errors = %v, want the step limit", errs)
	}
}

func TestRunRestoresCallerState(t *testing.T) {
	var out strings.Builder
	evaluator.Output = &out
	evaluator.MaxSteps = 0
	evaluator.OnStatement = nil
	srcRun(t, `puts("inside")`)
	if evaluator.MaxSteps != 0 || evaluator.OnStatement != nil || evaluator.OnLayer != nil {
		t.Error("src.run left its limit or hooks behind")
	}
	if evaluator.Output != &out || out.Len() != 0 {
		t.Error("src.run did not restore Output or let puts leak to it")
	}
}

func TestRunIsFresh(t *testing.T) {
	src := NewSrc(nil)
	src.Run("circle((1, 1), 1, #000)\nlet leak = 1")
	r := asJSON(t, src.Run("puts(leak)"))
	if errs := list(r["errors"]); len(errs) != 1 {
		t.Errorf("a variable from the first run is visible in the second: %v", r)
	}
	if n := len(list(r["shapes"])); n != 0 {
		t.Errorf("shapes leaked between runs: %d", n)
	}
}

func TestLibraryDocs(t *testing.T) {
	docs := NewSrc(nil).libraryDocs()
	c, ok := docs["circle"].(map[string]any)
	if !ok {
		t.Fatalf("circle missing from the library docs (%d entries)", len(docs))
	}
	if p, _ := c["params"].([]any); len(p) == 0 {
		t.Errorf("circle params = %v", c["params"])
	}
}

func TestRunReportsAndRestoresSteps(t *testing.T) {
	evaluator.MaxSteps = 1000
	evaluator.SetSteps(7)
	defer func() { evaluator.MaxSteps = 0; evaluator.SetSteps(0) }()
	r := srcRun(t, "let i = 0\nwhile (i < 50) { i = i + 1 }\n")
	if n, _ := r["steps"].(float64); n < 50 {
		t.Errorf("steps = %v, want at least 50", r["steps"])
	}
	if got := evaluator.Steps(); got != 7 {
		t.Errorf("caller steps = %d, want 7 restored", got)
	}
	if evaluator.MaxSteps != 1000 {
		t.Errorf("MaxSteps = %d, want 1000 restored", evaluator.MaxSteps)
	}
}

func TestRunLibraryErrorReportsUserLine(t *testing.T) {
	for _, tc := range []struct {
		name, code string
		line       int
		message    string
	}{
		{"library check", "let a = 1\n\n\ncircle(5, 3, #fff)\n", 4,
			"circle: center must be a point like (10, 20), not INTEGER"},
		{"inside a layer", "layer sky\n  let x = 1\n  circle(5, 3, #fff)\n", 3,
			"circle: center must be a point like (10, 20), not INTEGER"},
		{"inside a loop", "for (let i = 0; i < 2; i = i + 1) {\n  circle(i, 3, #fff)\n}\n", 2, "circle: center must be a point"},
		{"user code keeps its own position", "let a = 1\nlet b = a + nope\n", 2, "identifier not found: nope"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			errs := list(srcRun(t, tc.code)["errors"])
			if len(errs) != 1 {
				t.Fatalf("errors = %v", errs)
			}
			e := hash(errs[0])
			msg := e["message"].(string)
			if num(e["line"]) != tc.line || !strings.Contains(msg, tc.message) {
				t.Errorf("error = %v, want line %d with %q", e, tc.line, tc.message)
			}
			if strings.Contains(msg, "[line") || strings.Contains(msg, "runtime error") && tc.name != "user code keeps its own position" {
				t.Errorf("message still has a library position or prefix: %q", msg)
			}
		})
	}
	// Library errors have column 0; user errors have a real one.
	lib := hash(list(srcRun(t, "circle(5, 3, #fff)\n")["errors"])[0])
	if num(lib["column"]) != 0 {
		t.Errorf("library error column = %v, want 0", lib["column"])
	}
	usr := hash(list(srcRun(t, "let b = 1 + nope\n")["errors"])[0])
	if num(usr["column"]) == 0 {
		t.Errorf("user error lost its column: %v", usr)
	}
}

func TestRunUserFunctionErrorKeepsPosition(t *testing.T) {
	// A function the user wrote runs outside the statements that count for
	// tagging, but its errors are the user's own.
	errs := list(srcRun(t, "function f() {\n  nope\n}\nf()\n")["errors"])
	if len(errs) != 1 || num(hash(errs[0])["line"]) != 2 {
		t.Errorf("errors = %v, want line 2", errs)
	}
}
