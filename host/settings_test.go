package host

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"lizardscript/evaluator"
	"lizardscript/object"

	"linkz/host/pressure"
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

func TestApplyPressureCurves(t *testing.T) {
	trackpad, tablet := pressure.DefaultTrackpad, pressure.Identity
	for _, tc := range []struct {
		name, body string
		wantTrack  pressure.Curve
		wantTablet pressure.Curve
		wantErr    string
	}{
		{"none", "let language = \"en\"\n", trackpad, tablet, ""},
		{"full trackpad", "let trackpad_pressure = {\"low\": 0.7, \"high\": 0.95, \"gamma\": 2}\n", pressure.Curve{Low: 0.7, High: 0.95, Gamma: 2}, tablet, ""},
		{"partial keeps defaults", "let trackpad_pressure = {\"gamma\": 1}\n", pressure.Curve{Low: 0.8, High: 1, Gamma: 1}, tablet, ""},
		{"tablet", "let tablet_pressure = {\"low\": 0.1}\n", trackpad, pressure.Curve{Low: 0.1, High: 1, Gamma: 1}, ""},
		{"integers", "let tablet_pressure = {\"low\": 0, \"high\": 1, \"gamma\": 2}\n", trackpad, pressure.Curve{Low: 0, High: 1, Gamma: 2}, ""},
		{"not a hash", "let trackpad_pressure = 3\n", trackpad, tablet, "trackpad_pressure: want a hash"},
		{"not a number", "let trackpad_pressure = {\"low\": \"a\"}\n", trackpad, tablet, "low must be a number"},
		{"unknown key", "let trackpad_pressure = {\"lo\": 0.1}\n", trackpad, tablet, "unknown key \"lo\""},
		{"out of range", "let trackpad_pressure = {\"gamma\": 50}\n", trackpad, tablet, "gamma must be"},
		{"a bad curve applies none", "let tablet_pressure = {\"low\": 0.1}\nlet trackpad_pressure = {\"high\": 0}\n", trackpad, tablet, "high must be above low"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				pressure.SetCurve("trackpad", trackpad)
				pressure.SetCurve("tablet", tablet)
			})
			h, err := LoadSettingsFile(writeSettings(t, tc.body))
			if err != nil {
				t.Fatal(err)
			}
			err = ApplyPressureCurves(h)
			if (err == nil) != (tc.wantErr == "") || err != nil && !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error = %v, want %q", err, tc.wantErr)
			}
			// Probe the curves through their effect on two raw values.
			for _, c := range []struct {
				source string
				want   pressure.Curve
			}{{"trackpad", tc.wantTrack}, {"tablet", tc.wantTablet}} {
				for _, raw := range []float64{0.3, 0.9} {
					if got, want := pressure.Shape(c.source, raw), c.want.Apply(raw); math.Abs(got-want) > 1e-12 {
						t.Errorf("%s Shape(%v) = %v, want %v", c.source, raw, got, want)
					}
				}
			}
		})
	}
	if err := ApplyPressureCurves(nil); err != nil {
		t.Errorf("no settings gave %v", err)
	}
}

// With a temp HOME, ~/.linkz/settings.liz is read the way the editor does.
func TestPressureCurvesFromHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Cleanup(func() { pressure.SetCurve("trackpad", pressure.DefaultTrackpad) })
	dir := filepath.Join(home, ".linkz")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "let trackpad_pressure = {\"low\": 0.5, \"high\": 1.0, \"gamma\": 1.0}\n"
	if err := os.WriteFile(filepath.Join(dir, "settings.liz"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	h, err := LoadSettings()
	if err != nil {
		t.Fatal(err)
	}
	if err := ApplyPressureCurves(h); err != nil {
		t.Fatal(err)
	}
	if got := pressure.Shape("trackpad", 0.75); math.Abs(got-0.5) > 1e-12 {
		t.Errorf("Shape(0.75) = %v, want 0.5", got)
	}
}
