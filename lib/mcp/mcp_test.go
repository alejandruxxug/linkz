package mcp_test

import (
	"bytes"
	"flag"
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
