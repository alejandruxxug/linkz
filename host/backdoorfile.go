package host

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// Each editor window writes its own ~/.linkz/debug-<pid>.json, so a second
// window (a smoke run, a tester) never takes over the first one's backdoor.

// BackdoorInfo is what a debug file holds.
type BackdoorInfo struct {
	Port  int    `json:"port"`
	Token string `json:"token"`
	PID   int    `json:"pid"`
}

// pidAlive reports whether a process with that id exists.
func pidAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	err := syscall.Kill(pid, 0)
	return err == nil || err == syscall.EPERM
}

// backdoorFiles lists the debug files of dir with their pids.
func backdoorFiles(dir string) map[string]int {
	out := map[string]int{}
	names, _ := filepath.Glob(filepath.Join(dir, "debug-*.json"))
	for _, name := range names {
		base := strings.TrimSuffix(strings.TrimPrefix(filepath.Base(name), "debug-"), ".json")
		if pid, err := strconv.Atoi(base); err == nil {
			out[name] = pid
		}
	}
	return out
}

// CleanStaleBackdoors removes the debug files of windows that no longer run,
// and the old single debug.json if its process is gone.
func CleanStaleBackdoors() {
	dir, err := ConfigDir()
	if err != nil {
		return
	}
	for name, pid := range backdoorFiles(dir) {
		if !pidAlive(pid) {
			os.Remove(name)
		}
	}
	legacy := filepath.Join(dir, "debug.json")
	if data, err := os.ReadFile(legacy); err == nil {
		var info BackdoorInfo
		if json.Unmarshal(data, &info) != nil || !pidAlive(info.PID) {
			os.Remove(legacy)
		}
	}
}

// FindBackdoor reads the debug file of window pid, or of the newest running
// window when pid is 0.
func FindBackdoor(pid int) (BackdoorInfo, error) {
	dir, err := ConfigDir()
	if err != nil {
		return BackdoorInfo{}, err
	}
	best, bestPID := "", 0
	var bestTime int64
	for name, p := range backdoorFiles(dir) {
		if (pid != 0 && p != pid) || !pidAlive(p) {
			continue
		}
		st, err := os.Stat(name)
		if err != nil {
			continue
		}
		if t := st.ModTime().UnixNano(); best == "" || t > bestTime {
			best, bestPID, bestTime = name, p, t
		}
	}
	if best == "" {
		if pid != 0 {
			return BackdoorInfo{}, fmt.Errorf("no running editor with pid %d", pid)
		}
		return BackdoorInfo{}, fmt.Errorf("no running editor")
	}
	data, err := os.ReadFile(best)
	if err != nil {
		return BackdoorInfo{}, err
	}
	var info BackdoorInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return BackdoorInfo{}, err
	}
	info.PID = bestPID
	return info, nil
}
