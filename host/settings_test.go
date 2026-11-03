package host

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"lizardscript/evaluator"
	"lizardscript/object"
)

func writeSettings(t *testing.T, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "settings.liz")
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadSettingsFile(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		skip       bool
		language   string
		wantErr    string
	}{
		{"on", "-~> mine\nlet skip_loading_anim = true\nlet language = \"es\"\n", true, "es", ""},
		{"off", "let skip_loading_anim = false\n", false, "", ""},
		{"not set", "let language = \"en\"\n", false, "en", ""},
		{"not a boolean", "let skip_loading_anim = \"yes\"\n", false, "", ""},
		{"computed", "let a = 2\nlet skip_loading_anim = a > 1\n", true, "", ""},
		{"syntax error", "let skip_loading_anim = \n", false, "", "syntax error"},
		{"runtime error", "let language = nope\n", false, "", "nope"},
		{"runs too long", "while (true) { let x = 1 }\nlet skip_loading_anim = true\n", false, "", "ran too long"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h, err := LoadSettingsFile(writeSettings(t, tc.body))
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("err = %v, want one with %q", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if h == nil {
				t.Fatal("no hash: broken settings must still give an empty one")
			}
			if got := SkipLoadingAnim(h); got != tc.skip {
				t.Errorf("skip_loading_anim = %v, want %v", got, tc.skip)
			}
			lang := ""
			if s, ok := hashGet(h, "language").(*object.String); ok {
				lang = s.Value
			}
			if lang != tc.language {
				t.Errorf("language = %q, want %q", lang, tc.language)
			}
		})
	}
}

func TestLoadSettingsMissingFile(t *testing.T) {
	h, err := LoadSettingsFile(filepath.Join(t.TempDir(), "nope.liz"))
	if err != nil || h == nil || len(h.Pairs) != 0 {
		t.Errorf("missing file: %v, %v", h, err)
	}
}

func TestSettingsReachTheUI(t *testing.T) {
	app := NewApp(100, 100)
	h, _ := LoadSettingsFile(writeSettings(t, "let language = \"es\"\n"))
	app.Settings = h
	if err := app.Load(""); err != nil {
		t.Fatal(err)
	}
	v, err := app.eval("lang_get()", nil)
	if err != nil || v.Inspect() != "es" {
		t.Errorf("the UI language = %v, %v; want es from settings", v, err)
	}
}

func TestSettingsRestoreEvaluatorState(t *testing.T) {
	prev := evaluator.MaxSteps
	if _, err := LoadSettingsFile(writeSettings(t, "let a = 1\n")); err != nil {
		t.Fatal(err)
	}
	if evaluator.MaxSteps != prev {
		t.Error("evaluating settings left the step limit changed")
	}
}
