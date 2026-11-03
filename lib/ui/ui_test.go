// Package ui has no Go code: the editor (lib/ui/*.liz) is lizard-script. This
// file runs the programs in tests/ against the real src primitives and checks
// what they print, and checks the language tables.
package ui

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"linkz/host"
	"linkz/lib"
	"linkz/lib/draw"
	"lizardscript/evaluator"
	"lizardscript/lexer"
	"lizardscript/object"
	"lizardscript/parser"
)

var update = flag.Bool("update", false, "rewrite the .expected files with the current output")

func run(source string, env *object.Environment) error {
	p := parser.New(lexer.New(source))
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return errors.New(strings.Join(p.Errors(), "\n"))
	}
	if e, ok := host.Eval(program, env).(*object.Error); ok {
		return errors.New(e.Message)
	}
	return nil
}

// uiFiles are the editor's libraries in the order the host loads them: by name,
// without the app and the splash.
func uiFiles(t *testing.T) []string {
	t.Helper()
	files, err := filepath.Glob("*.liz")
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, f := range files {
		if f == "app.liz" || strings.HasPrefix(f, "splash") {
			continue
		}
		out = append(out, f)
	}
	sort.Strings(out)
	return out
}

// load makes an environment: drawing libraries, graphics, the real src, the
// in-memory fs/clip/os, then the editor.
func load(t *testing.T, out *bytes.Buffer) *object.Environment {
	t.Helper()
	evaluator.Output = out
	evaluator.ResetMath()
	env := object.NewEnvironment()
	env.Define("src", host.NewSrc(nil).Builtins())
	for _, name := range draw.Order {
		code, err := draw.Files.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if err := run(string(code), env); err != nil {
			t.Fatalf("draw/%s: %v", name, err)
		}
	}
	for _, name := range lib.Dir("gfx") {
		code, err := lib.Files.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if err := run(string(code), env); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
	prelude, err := os.ReadFile("tests/prelude.liz")
	if err != nil {
		t.Fatal(err)
	}
	if err := run(string(prelude), env); err != nil {
		t.Fatalf("prelude: %v", err)
	}
	for _, name := range uiFiles(t) {
		code, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if err := run(string(code), env); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
	return env
}

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

// Every tests/*.liz program (but the prelude) must run; what it prints is
// compared with the .expected file next to it.
func TestPrograms(t *testing.T) {
	files, err := filepath.Glob("tests/*.liz")
	if err != nil || len(files) == 0 {
		t.Fatalf("no programs: %v", err)
	}
	for _, file := range files {
		// presence.liz uses the real app (app.liz), which this runner
		// doesn't load; tests/ui/programs.liz runs it with modules.
		if filepath.Base(file) == "prelude.liz" || filepath.Base(file) == "presence.liz" {
			continue
		}
		t.Run(strings.TrimSuffix(filepath.Base(file), ".liz"), func(t *testing.T) {
			source, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			var out bytes.Buffer
			env := load(t, &out)
			if err := run(string(source), env); err != nil {
				out.WriteString("ERROR: " + err.Error() + "\n")
				t.Errorf("program failed: %v\n%s", err, out.String())
			}
			check(t, strings.TrimSuffix(file, ".liz")+".expected", out.String())
		})
	}
}

var (
	keyRe    = regexp.MustCompile(`\btf?\("([^"]+)"`)
	prefixRe = regexp.MustCompile(`\btf?\("([^"]*)"\s*\+`)
	tableRe  = regexp.MustCompile(`(?m)^\s*"([^"]+)":`)
)

func tableKeys(t *testing.T, file string) map[string]bool {
	t.Helper()
	src, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	keys := map[string]bool{}
	for _, m := range tableRe.FindAllStringSubmatch(string(src), -1) {
		keys[m[1]] = true
	}
	return keys
}

// Every t("key") the UI code uses must exist in the English table; keys that
// Spanish lacks are listed.
func TestLanguageKeys(t *testing.T) {
	en, es := tableKeys(t, "lang_en.liz"), tableKeys(t, "lang_es.liz")
	files, _ := filepath.Glob("*.liz")
	used := map[string]bool{}
	for _, f := range files {
		if strings.HasPrefix(f, "lang_") {
			continue
		}
		src, _ := os.ReadFile(f)
		var code []string
		for _, line := range strings.Split(string(src), "\n") {
			if !strings.HasPrefix(strings.TrimSpace(line), "-~>") {
				code = append(code, line)
			}
		}
		body := strings.Join(code, "\n")
		for _, m := range prefixRe.FindAllStringSubmatch(body, -1) {
			found := false
			for k := range en {
				found = found || strings.HasPrefix(k, m[1])
			}
			if !found {
				t.Errorf("%s: no key in lang_en.liz starts with %q", f, m[1])
			}
		}
		for _, m := range keyRe.FindAllStringSubmatch(body, -1) {
			used[m[1]] = true
			if !en[m[1]] {
				t.Errorf("%s: key %q is missing from lang_en.liz", f, m[1])
			}
		}
	}
	var missing []string
	for k := range en {
		if !es[k] {
			missing = append(missing, k)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Logf("keys missing from lang_es.liz: %v", missing)
	}
	if len(used) == 0 {
		t.Error("found no t(\"...\") calls at all: the pattern is wrong")
	}
}
