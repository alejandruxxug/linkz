package host

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"lizardscript/object"
)

// fakeMeasure counts how often it is asked to measure.
func fakeMeasure(n *int, version string) func() (string, error) {
	return func() (string, error) {
		*n++
		return `let machine = {"max_objects": 120, "app_version": "` + version + `"}` + "\n", nil
	}
}

func TestLoadMachineCache(t *testing.T) {
	old := AppVersion
	defer func() { AppVersion = old }()
	path := filepath.Join(t.TempDir(), "machine.liz")
	n := 0

	AppVersion = "1.0"
	m, measured, err := LoadMachine(path, false, fakeMeasure(&n, "1.0"))
	if err != nil || !measured || n != 1 {
		t.Fatalf("first load: measured=%v n=%d err=%v", measured, n, err)
	}
	if got := hashGet(m, "max_objects").(*object.Integer).Value; got != 120 {
		t.Errorf("max_objects = %d", got)
	}

	// Same version: read from the file, no new measuring.
	if _, measured, err = LoadMachine(path, false, fakeMeasure(&n, "1.0")); err != nil || measured || n != 1 {
		t.Errorf("second load: measured=%v n=%d err=%v", measured, n, err)
	}

	// A new app version measures again.
	AppVersion = "1.1"
	if _, measured, err = LoadMachine(path, false, fakeMeasure(&n, "1.1")); err != nil || !measured || n != 2 {
		t.Errorf("new version: measured=%v n=%d err=%v", measured, n, err)
	}
	if _, measured, _ = LoadMachine(path, false, fakeMeasure(&n, "1.1")); measured {
		t.Error("the new measurement should be kept")
	}

	// force measures again.
	if _, measured, _ = LoadMachine(path, true, fakeMeasure(&n, "1.1")); !measured || n != 3 {
		t.Errorf("forced: measured=%v n=%d", measured, n)
	}

	// A broken file measures again.
	os.WriteFile(path, []byte("let machine = ("), 0o644)
	if _, measured, _ = LoadMachine(path, false, fakeMeasure(&n, "1.1")); !measured || n != 4 {
		t.Errorf("broken file: measured=%v n=%d", measured, n)
	}
}

func TestLoadMachineMeasureFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "machine.liz")
	_, _, err := LoadMachine(path, false, func() (string, error) { return "", io.ErrUnexpectedEOF })
	if err == nil {
		t.Error("a failed measurement should be an error")
	}
	if _, serr := os.Stat(path); serr == nil {
		t.Error("nothing should be written after a failure")
	}
}

// The real benchmark: the JSON that `linkz calibrate` prints has every key,
// and the numbers are plausible.
func TestCalibrateOutputShape(t *testing.T) {
	old := AppVersion
	defer func() { AppVersion = old }()
	AppVersion = "test"
	text, err := runCalibration("machine_json(calibrate())", io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"max_objects", "ms_per_object", "bytes_per_object", "memory_cap", "speed_cap", "limited_by", "measured_at", "app_version"} {
		if !strings.Contains(text, `"`+key+`": `) {
			t.Errorf("missing %s in %s", key, text)
		}
	}
	if !strings.HasPrefix(text, "{") || !strings.HasSuffix(text, "}") || !strings.Contains(text, `"app_version": "test"`) {
		t.Errorf("not the expected JSON: %s", text)
	}
	h, err := evalSettings("let machine = " + text)
	if err != nil {
		t.Fatalf("the JSON is not lizard data either: %v", err)
	}
	m := hashGet(h, "machine").(*object.Hash)
	if v := hashGet(m, "max_objects").(*object.Integer).Value; v < 10 || v > 1_000_000 {
		t.Errorf("max_objects = %d", v)
	}
}
