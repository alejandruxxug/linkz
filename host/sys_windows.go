//go:build windows

package host

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Windows side of the OS primitives. Nothing here could be run by the
// developer who wrote it: it compiles and vets, and is listed in the Build
// Plan as needing a test on a real Windows PC.

// hideWindow keeps PowerShell from flashing a console window.
func hideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: windows.CREATE_NO_WINDOW}
}

// powershell runs a script with no profile and returns its output.
func powershell(stdin, script string) ([]byte, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	hideWindow(cmd)
	return cmd.Output()
}

const psUTF8 = "[Console]::OutputEncoding=[Text.Encoding]::UTF8;[Console]::InputEncoding=[Text.Encoding]::UTF8;"

func platformClipGet() (string, bool) {
	out, err := powershell("", psUTF8+"Get-Clipboard -Raw")
	return strings.TrimPrefix(string(out), "\xef\xbb\xbf"), err == nil
}

func platformClipSet(text string) {
	if text == "" {
		powershell("", "Set-Clipboard -Value $null")
		return
	}
	powershell(text, psUTF8+"Set-Clipboard -Value ([Console]::In.ReadToEnd())")
}

// platformLocale asks Windows for the user's locale name ("es-CO").
func platformLocale() string {
	proc := windows.NewLazySystemDLL("kernel32.dll").NewProc("GetUserDefaultLocaleName")
	var buf [85]uint16 // LOCALE_NAME_MAX_LENGTH
	n, _, _ := proc.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if n == 0 {
		return ""
	}
	return twoLetters(windows.UTF16ToString(buf[:]))
}

// platformReveal opens Explorer with the file selected. explorer.exe often
// exits with status 1 even when it worked, so the exit code is ignored.
func platformReveal(path string) error {
	err := runCommand("explorer", "/select,"+filepath.FromSlash(path))
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return nil
	}
	return err
}

// There is no folder to move into: the Recycle Bin is only reachable
// through the shell.
func trashDir() (string, error) { return "", errNoTrashDir }

func trashNote(src, dest string) error { return nil }

// platformTrash sends the file to the Recycle Bin.
func platformTrash(path string) error {
	script := "Add-Type -AssemblyName Microsoft.VisualBasic;" +
		"[Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFile($env:LINKZ_TRASH,'OnlyErrorDialogs','SendToRecycleBin')"
	return runCommandEnv([]string{"LINKZ_TRASH=" + filepath.FromSlash(path)}, "powershell", "-NoProfile", "-NonInteractive", "-Command", script)
}
