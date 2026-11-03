package gfx

import (
	"os"
	"strings"
	"testing"
)

// TestRecommend runs the fixed-input programs for calibrate.liz.
func TestRecommend(t *testing.T) {
	_, out := scene(t, "recommend", 16, 16)
	want, err := os.ReadFile("tests/recommend.expected")
	if err != nil {
		t.Fatal(err)
	}
	if out != string(want) {
		t.Errorf("output:\n%s\nwant:\n%s", out, want)
	}
}

// TestCalibrateShape runs the real benchmark once and checks the result.
func TestCalibrateShape(t *testing.T) {
	_, out := paint(t, "calibrate", `
let m = calibrate()
for (let i = 0; i < len(calib_keys); i = i + 1) {
  puts(calib_keys[i], type(m[calib_keys[i]]))
}
puts(m.max_objects > 0, m.ms_per_object > 0, m.bytes_per_object > 0, m.limited_by)
`, 512, 512)
	for _, want := range []string{"max_objects INTEGER", "ms_per_object FLOAT", "warm_ms_per_object FLOAT", "cold_ms_per_object FLOAT", "run_ms_per_object FLOAT", "bytes_per_object", "memory_cap INTEGER",
		"speed_cap INTEGER", "limited_by STRING", "measured_at INTEGER", "app_version STRING", "true true true"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}
