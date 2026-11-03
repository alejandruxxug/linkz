package draw

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"lizardscript/evaluator"
	"lizardscript/lexer"
	"lizardscript/object"
	"lizardscript/parser"
	"lizardscript/token"
)

var update = flag.Bool("update", false, "rewrite the .expected files with the current output")

// run parses and evaluates source in env, like the host does.
func run(source string, env *object.Environment) error {
	p := parser.New(lexer.New(source))
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return errors.New(strings.Join(p.Errors(), "\n"))
	}
	if e, ok := evaluator.Eval(program, env).(*object.Error); ok {
		return errors.New(e.Message)
	}
	return nil
}

// load makes an environment with every library loaded.
func load(t *testing.T, out *bytes.Buffer) *object.Environment {
	t.Helper()
	evaluator.Output = out
	evaluator.ResetMath()
	env := object.NewEnvironment()
	for _, name := range Order {
		code, err := Files.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if err := run(string(code), env); err != nil {
			t.Fatalf("library %s: %v", name, err)
		}
	}
	return env
}

// check compares got with the contents of the expected file, or rewrites the
// file with -update.
func check(t *testing.T, expectedFile, got string) {
	t.Helper()
	if *update {
		if err := os.WriteFile(expectedFile, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("missing %s: %v", expectedFile, err)
	}
	if got != string(want) {
		t.Errorf("wrong output\n--- expected:\n%s\n--- got:\n%s", want, got)
	}
}

// Every tests/*.liz program must run, and tests/errors/*.liz must fail. What
// they print (and the error) is compared with the .expected file next to them.
func TestPrograms(t *testing.T) {
	for dir, shouldSucceed := range map[string]bool{"tests": true, "tests/errors": false} {
		files, err := filepath.Glob(filepath.Join(dir, "*.liz"))
		if err != nil || len(files) == 0 {
			t.Fatalf("no programs in %s: %v", dir, err)
		}
		for _, file := range files {
			t.Run(strings.TrimSuffix(file, ".liz"), func(t *testing.T) {
				source, err := os.ReadFile(file)
				if err != nil {
					t.Fatal(err)
				}
				var out bytes.Buffer
				env := load(t, &out)
				runErr := run(string(source), env)
				if runErr != nil {
					fmt.Fprintln(&out, runErr)
				}
				if (runErr == nil) != shouldSucceed {
					t.Errorf("success = %t, want %t. Output:\n%s", runErr == nil, shouldSucceed, out.String())
				}
				check(t, strings.TrimSuffix(file, ".liz")+".expected", out.String())
			})
		}
	}
}

// Every example renders, and its SVG is the golden one.
func TestExamples(t *testing.T) {
	files, err := filepath.Glob("../../examples/*.lzdraw")
	if err != nil || len(files) == 0 {
		t.Fatalf("no examples: %v", err)
	}
	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".lzdraw")
		t.Run(name, func(t *testing.T) {
			source, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			var out bytes.Buffer
			env := load(t, &out)
			defer trackLayers(env)()
			if err := run(string(source), env); err != nil {
				t.Fatal(err)
			}
			result := evaluator.Eval(parser.New(lexer.New("svg_document(canvas, drawing)")).ParseProgram(), env)
			svg, ok := result.(*object.String)
			if !ok {
				t.Fatalf("svg_document gave %v", result.Inspect())
			}
			check(t, filepath.Join("tests", "golden", name+".svg.expected"), svg.Value)
		})
	}
}

// trackLayers makes here() report the layer path, like the host does, so the
// examples' layers reach the shapes. It returns a function that undoes it.
func trackLayers(env *object.Environment) func() {
	prev := evaluator.OnLayer
	var layers []string
	evaluator.OnLayer = func(name string, enter bool, tok token.Token) {
		if enter {
			layers = append(layers, name)
		} else if len(layers) > 0 {
			layers = layers[:len(layers)-1]
		}
	}
	env.Define("here", &object.Builtin{Name: "here", Fn: func(args ...object.Object) object.Object {
		path := &object.Array{}
		for _, l := range layers {
			path.Elements = append(path.Elements, &object.String{Value: l})
		}
		h := object.NewHash()
		h.Set(&object.String{Value: "line"}, &object.Null{})
		h.Set(&object.String{Value: "layer"}, path)
		return h
	}})
	return func() { evaluator.OnLayer = prev }
}
