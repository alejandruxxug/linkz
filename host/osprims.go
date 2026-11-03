package host

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"lizardscript/object"
)

// The OS primitives: files, time, clipboard, locale. Nothing here decides
// anything about drawing.

var startTime = time.Now()

// strArgs checks that every argument is a STRING.
func strArgs(args []object.Object, n int, want string) ([]string, *object.Error) {
	if len(args) != n {
		return nil, errorf("wrong number of arguments: want %s, got %d", want, len(args))
	}
	out := make([]string, n)
	for i, a := range args {
		s, ok := a.(*object.String)
		if !ok {
			return nil, errorf("argument %d must be a STRING, got %s", i+1, a.Type())
		}
		out[i] = s.Value
	}
	return out, nil
}

// WriteFileAtomic writes to a temp file next to path and renames it over, so
// a crash never leaves half a file.
func WriteFileAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".linkz-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(name)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}
	mode := os.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	}
	if err := os.Chmod(name, mode); err != nil {
		os.Remove(name)
		return err
	}
	if err := os.Rename(name, path); err != nil {
		os.Remove(name)
		return err
	}
	return nil
}

func fsHash() *object.Hash {
	h := builtinHash("fs", map[string]object.BuiltinFunction{
		"read": func(args ...object.Object) object.Object {
			s, err := strArgs(args, 1, "1 (path)")
			if err != nil {
				return err
			}
			data, rerr := os.ReadFile(s[0])
			if rerr != nil {
				return errorf("%v", rerr)
			}
			return newStr(string(data))
		},
		"write": func(args ...object.Object) object.Object {
			s, err := strArgs(args, 2, "2 (path, text)")
			if err != nil {
				return err
			}
			if werr := WriteFileAtomic(s[0], []byte(s[1])); werr != nil {
				return errorf("%v", werr)
			}
			return null
		},
		"exists": func(args ...object.Object) object.Object {
			s, err := strArgs(args, 1, "1 (path)")
			if err != nil {
				return err
			}
			_, serr := os.Stat(s[0])
			return newBool(serr == nil)
		},
	})
	addFSPrims(h)
	return h
}

func timeHash() *object.Hash {
	h := builtinHash("time", map[string]object.BuiltinFunction{
		// ms is a FLOAT: milliseconds since the program started.
		"ms": func(args ...object.Object) object.Object {
			return newFloat(float64(time.Since(startTime).Microseconds()) / 1000)
		},
	})
	addTimePrims(h)
	return h
}

// The clipboard goes through the system's own tools (pbcopy on macOS,
// PowerShell on Windows, wl-copy or xclip on Linux; no cgo). The last text set
// stays in the app, so copy and paste still work when no tool is there.
var (
	clipMu    sync.Mutex
	clipLocal string
)

func clipGet() string {
	clipMu.Lock()
	defer clipMu.Unlock()
	if text, ok := platformClipGet(); ok {
		return text
	}
	return clipLocal
}

func clipSet(text string) {
	clipMu.Lock()
	defer clipMu.Unlock()
	clipLocal = text
	platformClipSet(text)
}

func clipHash() *object.Hash {
	return builtinHash("clip", map[string]object.BuiltinFunction{
		"get": func(args ...object.Object) object.Object { return newStr(clipGet()) },
		"set": func(args ...object.Object) object.Object {
			s, err := strArgs(args, 1, "1 (text)")
			if err != nil {
				return err
			}
			clipSet(s[0])
			return null
		},
	})
}

// Locale is the user's UI language as a lowercase two-letter code. It asks
// the system (macOS preferences, Windows locale), then LANG, and falls back
// to "en".
func Locale() string {
	if code := platformLocale(); code != "" {
		return code
	}
	for _, name := range []string{"LC_ALL", "LANG"} {
		if code := twoLetters(os.Getenv(name)); code != "" {
			return code
		}
	}
	return "en"
}

// twoLetters returns the leading language code, like "es" from "es_CO.UTF-8".
func twoLetters(s string) string {
	if len(s) < 2 {
		return ""
	}
	a, b := s[0]|0x20, s[1]|0x20
	if a < 'a' || a > 'z' || b < 'a' || b > 'z' {
		return ""
	}
	if len(s) > 2 && s[2] >= 'a' && s[2] <= 'z' { // "eng": not a two-letter code
		return ""
	}
	return string([]byte{a, b})
}

// ConfigDir is ~/.linkz, created if missing.
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".linkz")
	return dir, os.MkdirAll(dir, 0o755)
}

func osHash() *object.Hash {
	h := builtinHash("os", map[string]object.BuiltinFunction{
		"locale": func(args ...object.Object) object.Object { return newStr(Locale()) },
		"config_dir": func(args ...object.Object) object.Object {
			dir, err := ConfigDir()
			if err != nil {
				return errorf("%v", err)
			}
			return newStr(filepath.ToSlash(dir))
		},
	})
	addOSPrims(h)
	addMemPrims(h)
	return h
}
