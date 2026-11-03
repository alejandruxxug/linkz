package trace

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"linkz/lib/draw"
	"lizardscript/evaluator"
	"lizardscript/lexer"
	"lizardscript/object"
	"lizardscript/parser"
)

var update = flag.Bool("update", false, "rewrite the .expected files with the current output")

// run parses and evaluates source in env, like the host does, and gives the
// value of the last statement.
func run(source string, env *object.Environment) (object.Object, error) {
	p := parser.New(lexer.New(source))
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return nil, errors.New(strings.Join(p.Errors(), "\n"))
	}
	res := evaluator.Eval(program, env)
	if e, ok := res.(*object.Error); ok {
		return nil, errors.New(e.Message)
	}
	return res, nil
}

// load makes an environment with the drawing libraries, the tracing library
// and a time.ms that always says 0 (so the output is the same every run).
func load(t *testing.T) *object.Environment {
	t.Helper()
	evaluator.Output = &bytes.Buffer{}
	evaluator.ResetMath()
	env := object.NewEnvironment()
	clock := object.NewHash()
	clock.Set(&object.String{Value: "ms"}, &object.Builtin{Name: "time.ms", Fn: func(args ...object.Object) object.Object {
		return &object.Float{Value: 0}
	}})
	env.Define("time", clock)
	for _, name := range draw.Order {
		code, err := draw.Files.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := run(string(code), env); err != nil {
			t.Fatalf("library %s: %v", name, err)
		}
	}
	code, err := Files.ReadFile("trace.liz")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := run(string(code), env); err != nil {
		t.Fatalf("library trace.liz: %v", err)
	}
	return env
}

// pictures builds tiny pictures in lizard: grid(w, h, f) calls f(x, y) for
// every pixel, row by row.
const pictures = `
let red = #c82828
let white = #f0f0e6
let green = #2d8a3e
let blue = #2850c8
function grid(w, h, f) {
  let px = []
  for (let y = 0; y < h; y = y + 1) {
    for (let x = 0; x < w; x = x + 1) {
      push(px, f(x, y))
    }
  }
  {"size": (w, h), "pixels": px}
}
`

// area is the net area of the first path shape in the drawing, holes
// subtracted, from the flattened outlines.
const area = `
function net_area(shape) {
  let subs = flatten_path(shape.path, 0.05)
  let sum = 0.0
  for (let i = 0; i < len(subs); i = i + 1) {
    sum = sum + polygon_area2(subs[i].points)
  }
  math.abs(sum) / 2.0
}
`

func TestTrace(t *testing.T) {
	tests := []struct {
		name    string
		picture string
		options string
		colors  int
		regions int
		holes   int
		shapes  int     // shapes the code draws
		area    float64 // net area of the first path of the code, 0 = skip
	}{
		{"two colors, a square", `grid(12, 12, function(x, y) { if (x >= 4 && x < 8 && y >= 4 && y < 8) { red } else { white } })`,
			`{"colors": 2, "min_area": 0}`, 2, 2, 1, 2, 16},
		{"a square with a hole", `grid(14, 14, function(x, y) {
			if (x >= 5 && x < 9 && y >= 5 && y < 9) { white }
			elif (x >= 2 && x < 12 && y >= 2 && y < 12) { red }
			else { white } })`,
			`{"colors": 2, "min_area": 0}`, 2, 3, 2, 2, 84},
		{"three colors in bands", `grid(12, 6, function(x, y) { if (x < 4) { red } elif (x < 8) { green } else { blue } })`,
			`{"colors": 3, "min_area": 0}`, 3, 3, 0, 3, 0},
		{"three colors, nested", `grid(15, 15, function(x, y) {
			if (x >= 5 && x < 10 && y >= 5 && y < 10) { blue }
			elif (x >= 2 && x < 13 && y >= 2 && y < 13) { green }
			else { white } })`,
			`{"colors": 3, "min_area": 0}`, 3, 3, 2, 3, 0},
		{"corners only touching are two regions", `grid(2, 2, function(x, y) { if (x == y) { red } else { white } })`,
			`{"colors": 2, "min_area": 0}`, 2, 4, 0, 3, 0},
		{"a speck is merged", `grid(12, 12, function(x, y) { if (x == 6 && y == 6) { red } else { white } })`,
			`{"colors": 2, "min_area": 4}`, 1, 1, 0, 1, 0},
		{"no background rectangle", `grid(12, 12, function(x, y) { if (x >= 4 && x < 8 && y >= 4 && y < 8) { red } else { white } })`,
			`{"colors": 2, "min_area": 0, "background": false}`, 2, 2, 1, 2, 16},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			env := load(t)
			src := pictures + area + "let pic = " + tc.picture + "\n" +
				"let result = trace_pixels(pic.size, pic.pixels, " + tc.options + ", \"test\")\n" +
				"let stats = result.stats\n" +
				"let summary = [stats.colors, stats.regions, stats.holes]\nsummary"
			res, err := run(src, env)
			if err != nil {
				t.Fatal(err)
			}
			if got, want := res.Inspect(), fmt.Sprintf("[%d, %d, %d]", tc.colors, tc.regions, tc.holes); got != want {
				t.Errorf("colors, regions, holes = %s, want %s", got, want)
			}

			// The code parses and draws.
			code, err := run("result.code", env)
			if err != nil {
				t.Fatal(err)
			}
			text := code.(*object.String).Value
			if _, err := run("drawing = []\n"+text, env); err != nil {
				t.Fatalf("the code fails: %v\n%s", err, text)
			}
			n, _ := run("len(drawing)", env)
			if got, want := n.Inspect(), fmt.Sprint(tc.shapes); got != want {
				t.Errorf("the code draws %s shapes, want %s\n%s", got, want, text)
			}
			if svg, err := run("svg_document(canvas, drawing)", env); err != nil || !strings.Contains(svg.(*object.String).Value, "<path") {
				t.Errorf("no SVG from the code: %v", err)
			}
			if tc.area > 0 {
				a, err := run("let found = null\nfor (let i = 0; i < len(drawing); i = i + 1) { if (found === null && drawing[i].kind == \"path\") { found = drawing[i] } }\nnet_area(found)", env)
				if err != nil {
					t.Fatal(err)
				}
				var got float64
				fmt.Sscan(a.Inspect(), &got)
				if got < tc.area-0.01 || got > tc.area+0.01 {
					t.Errorf("traced area = %v, want %v\n%s", got, tc.area, text)
				}
			}
		})
	}
}

// The code of a small picture is checked against a golden file, so any change
// to what trace writes shows up as a diff.
func TestCodeGolden(t *testing.T) {
	env := load(t)
	src := pictures + `
let pic = grid(14, 14, function(x, y) {
  if (x >= 5 && x < 9 && y >= 5 && y < 9) { white }
  elif (x >= 2 && x < 12 && y >= 2 && y < 12) { red }
  else { white } })
trace_pixels(pic.size, pic.pixels, {"colors": 2, "min_area": 0}, "ring").code
`
	res, err := run(src, env)
	if err != nil {
		t.Fatal(err)
	}
	got := res.(*object.String).Value
	file := "tests/ring.lzdraw.expected"
	if *update {
		if err := os.WriteFile(file, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if got != string(want) {
		t.Errorf("code differs from %s:\n%s", file, got)
	}
}

func TestTraceOptionsAreChecked(t *testing.T) {
	for name, options := range map[string]string{
		"unknown option": `{"colours": 4}`,
		"too few colors": `{"colors": 1}`,
		"not a hash":     `5`,
	} {
		t.Run(name, func(t *testing.T) {
			env := load(t)
			_, err := run(pictures+"let pic = grid(2, 2, function(x, y) { red })\ntrace_pixels(pic.size, pic.pixels, "+options+")", env)
			if err == nil || !strings.Contains(err.Error(), "trace:") {
				t.Errorf("want a trace: error, got %v", err)
			}
		})
	}
}
