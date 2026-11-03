package host

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"lizardscript/object"
)

func TestPidAlive(t *testing.T) {
	if !pidAlive(os.Getpid()) {
		t.Error("our own pid should be alive")
	}
	for _, pid := range []int{0, -1, 1 << 30} {
		if pidAlive(pid) {
			t.Errorf("pid %d should not be alive", pid)
		}
	}
}

func TestPercentEncodePath(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"/home/a/art.lzdraw", "/home/a/art.lzdraw"},
		{"/home/a b/ñ.lzdraw", "/home/a%20b/%C3%B1.lzdraw"},
		{"/x/100%.liz", "/x/100%25.liz"},
	} {
		if got := percentEncodePath(c.in); got != c.want {
			t.Errorf("percentEncodePath(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestTrashInfoText(t *testing.T) {
	when := time.Date(2026, 11, 3, 9, 5, 7, 0, time.UTC)
	want := "[Trash Info]\nPath=/a%20b.lzdraw\nDeletionDate=2026-11-03T09:05:07\n"
	if got := trashInfoText("/a b.lzdraw", when); got != want {
		t.Errorf("trashInfoText = %q, want %q", got, want)
	}
}

func TestTwoLetters(t *testing.T) {
	for in, want := range map[string]string{"es_CO.UTF-8": "es", "en-US": "en", "ES": "es", "C": "", "eng": "", "": "", "C.UTF-8": ""} {
		if got := twoLetters(in); got != want {
			t.Errorf("twoLetters(%q) = %q, want %q", in, got, want)
		}
	}
}

// Lizard sees "/" in every path, whatever the system.
func TestFSPathPrims(t *testing.T) {
	h := fsHash()
	if got := call(t, h, "sep").(*object.String).Value; got != "/" {
		t.Errorf("fs.sep = %q", got)
	}
	if got := call(t, h, "join", newStr("a"), newStr("b.lzdraw")).(*object.String).Value; got != "a/b.lzdraw" {
		t.Errorf("fs.join = %q", got)
	}
	if got := call(t, h, "join", newStr("a/"), newStr("../b")).(*object.String).Value; got != "b" {
		t.Errorf("fs.join with .. = %q", got)
	}
	in := filepath.Join("a", "b", "c.lzdraw")
	if got := call(t, h, "slash", newStr(in)).(*object.String).Value; got != "a/b/c.lzdraw" {
		t.Errorf("fs.slash(%q) = %q", in, got)
	}
	for _, p := range []object.Object{call(t, osHash(), "home"), call(t, osHash(), "config_dir")} {
		if strings.Contains(p.(*object.String).Value, `\`) {
			t.Errorf("path %q has a backslash", p.Inspect())
		}
	}
}
