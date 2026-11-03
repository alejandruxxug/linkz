//go:build !darwin && !windows

package host

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Linux (and other Unix) side of the OS primitives. It uses the tools a
// desktop normally has and falls back to doing nothing clever.

func hideWindow(cmd *exec.Cmd) {}

// clipTools lists the clipboard programs to try, in order: the name and
// the arguments that read (clipWriters: that write).
var clipTools = []struct {
	name string
	read []string
}{
	{"wl-paste", []string{"--no-newline"}},
	{"xclip", []string{"-selection", "clipboard", "-o"}},
	{"xsel", []string{"--clipboard", "--output"}},
}

var clipWriters = []struct {
	name string
	args []string
}{
	{"wl-copy", nil},
	{"xclip", []string{"-selection", "clipboard", "-i"}},
	{"xsel", []string{"--clipboard", "--input"}},
}

func platformClipGet() (string, bool) {
	for _, t := range clipTools {
		if _, err := exec.LookPath(t.name); err != nil {
			continue
		}
		if out, err := exec.Command(t.name, t.read...).Output(); err == nil {
			return string(out), true
		}
	}
	return "", false
}

func platformClipSet(text string) {
	for _, w := range clipWriters {
		if _, err := exec.LookPath(w.name); err != nil {
			continue
		}
		cmd := exec.Command(w.name, w.args...)
		cmd.Stdin = strings.NewReader(text)
		if cmd.Run() == nil {
			return
		}
	}
}

// platformLocale has nothing beyond LANG on Linux.
func platformLocale() string { return "" }

func platformReveal(path string) error { return runCommand("xdg-open", filepath.Dir(path)) }

// trashDir is the freedesktop Trash, $XDG_DATA_HOME/Trash/files.
func trashDir() (string, error) {
	if x := os.Getenv("XDG_DATA_HOME"); x != "" && filepath.IsAbs(x) {
		return filepath.Join(x, "Trash", "files"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "Trash", "files"), nil
}

// trashNote writes the .trashinfo for a file about to be moved to dest, so
// file managers can restore it.
func trashNote(src, dest string) error {
	abs, err := filepath.Abs(src)
	if err != nil {
		return err
	}
	dir, err := trashDir()
	if err != nil {
		return err
	}
	info := filepath.Join(filepath.Dir(dir), "info")
	if err := os.MkdirAll(info, 0o700); err != nil {
		return err
	}
	text := trashInfoText(abs, time.Now())
	return os.WriteFile(filepath.Join(info, filepath.Base(dest)+".trashinfo"), []byte(text), 0o600)
}

// platformTrash asks gio when the plain move was not allowed (another disk).
func platformTrash(path string) error {
	if _, err := exec.LookPath("gio"); err != nil {
		return fmt.Errorf("gio not found")
	}
	return runCommand("gio", "trash", path)
}
