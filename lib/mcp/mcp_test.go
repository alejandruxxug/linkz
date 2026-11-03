package mcp_test

import (
	"bytes"
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"linkz/host"
	"linkz/lib/mcp"
	"lizardscript/evaluator"
	"lizardscript/object"
)

var update = flag.Bool("update", false, "rewrite the .expected files with the current output")

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

// jsonEnv is an environment with only json.liz loaded.
func jsonEnv(t *testing.T) *object.Environment {
	t.Helper()
	env := object.NewEnvironment()
	code, err := mcp.Files.ReadFile("json.liz")
	if err != nil {
		t.Fatal(err)
	}
	if err := host.EvalIn(env, "json.liz", string(code)); err != nil {
		t.Fatal(err)
	}
	return env
}

// The json_*.liz programs run with json.liz alone; the server_*.liz programs
// run in the whole MCP environment, with HOME in a temporary folder so the
// files they save never touch the real one.
func TestPrograms(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("tests", "*.liz"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no programs: %v", err)
	}
	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".liz")
		t.Run(name, func(t *testing.T) {
			source, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			var out bytes.Buffer
			var env *object.Environment
			if strings.HasPrefix(name, "json_") {
				evaluator.Output = &out
				env = jsonEnv(t)
			} else {
				home := t.TempDir()
				t.Setenv("HOME", home)
				t.Setenv("USERPROFILE", home)
				env, err = host.NewMCPEnv(strings.NewReader(""), &bytes.Buffer{}, &out)
				if err != nil {
					t.Fatal(err)
				}
				evaluator.Output = &out
			}
			if err := host.EvalIn(env, name, string(source)); err != nil {
				out.WriteString(err.Error() + "\n")
				t.Errorf("the program failed: %v", err)
			}
			got := strings.ReplaceAll(out.String(), os.Getenv("HOME"), "$HOME")
			check(t, strings.TrimSuffix(file, ".liz")+".expected", got)
		})
	}
}

// writeCircle writes a 48 by 48 PNG: a red circle on white.
func writeCircle(t *testing.T, path string) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 48, 48))
	for y := 0; y < 48; y++ {
		for x := 0; x < 48; x++ {
			c := color.RGBA{255, 255, 255, 255}
			if (x-24)*(x-24)+(y-24)*(y-24) <= 15*15 {
				c = color.RGBA{220, 30, 30, 255}
			}
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}

// The trace tool on a generated picture, with HOME in a temporary folder.
func TestTraceTool(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	writeCircle(t, filepath.Join(home, "circle.png"))
	os.WriteFile(filepath.Join(home, "notes.txt"), []byte("x"), 0o644)
	var out bytes.Buffer
	env, err := host.NewMCPEnv(strings.NewReader(""), &bytes.Buffer{}, &out)
	if err != nil {
		t.Fatal(err)
	}
	evaluator.Output = &out
	program := `
function call(args) {
  let line = json.stringify({"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "trace", "arguments": args}})
  json.parse(mcp_handle_line(line)).result
}
let r = call({"path": os.home() + "/circle.png", "colors": 2})
puts(r.isError, len(r.content), r.content[0].type, r.content[1].type, r.content[2].type, r.content[2].mimeType)
puts(src.run(r.content[0].text).errors)
puts(json.parse(r.content[1].text).colors)
puts(len(r.content[2].data) > 100)
let bad = [
  {"path": os.home() + "/notes.txt"},
  {"path": "/etc/hosts.png"},
  {"path": os.home() + "/missing.png"},
  {"path": os.home() + "/circle.png", "colors": 0},
  {"path": os.home() + "/circle.png", "background": "yes"},
  {}
]
for (let i = 0; i < len(bad); i += 1) {
  let e = call(bad[i])
  puts(e.isError, text.replace(e.content[0].text, os.home(), "$HOME"))
}
`
	if err := host.EvalIn(env, "trace_test", program); err != nil {
		t.Fatalf("%v\n%s", err, out.String())
	}
	got := out.String()
	for _, want := range []string{
		"false 3 text text image image/png",
		"[]\n",
		"true only .png, .jpg, .jpeg files, not notes.txt",
		"true an absolute path must be inside your home folder",
		"true there is no file $HOME/missing.png",
		"true the argument colors must be a whole number from 1 to 64",
		"true the argument background must be true or false",
		"true the argument path is required",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output lacks %q:\n%s", want, got)
		}
	}
}
