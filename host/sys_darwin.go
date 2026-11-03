//go:build darwin

package host

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// macOS side of the OS primitives: the clipboard, the language, Finder and
// the Trash. Windows and Linux have their own files.

func hideWindow(cmd *exec.Cmd) {}

func platformClipGet() (string, bool) {
	out, err := exec.Command("pbpaste").Output()
	return string(out), err == nil
}

func platformClipSet(text string) {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	cmd.Run()
}

func platformLocale() string {
	out, err := exec.Command("defaults", "read", "-g", "AppleLanguages").Output()
	if err != nil {
		return ""
	}
	// ( "es-CO", "en-US" ) : the first quoted entry
	text := string(out)
	if i := strings.IndexByte(text, '"'); i >= 0 {
		return twoLetters(text[i+1:])
	}
	return ""
}

func platformReveal(path string) error { return runCommand("open", "-R", path) }

func trashDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".Trash"), nil
}

// trashNote has nothing to write on macOS.
func trashNote(src, dest string) error { return nil }

// platformTrash asks Finder when the plain move was not allowed.
func platformTrash(path string) error {
	script := fmt.Sprintf(`tell application "Finder" to delete POSIX file %q`, path)
	return runCommand("osascript", "-e", script)
}
