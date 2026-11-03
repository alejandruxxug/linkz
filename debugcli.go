package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"linkz/host"
)

// debug sends one command to the newest running editor (or --pid N) and prints the JSON reply.
func debug(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: linkz debug [--pid N] Command [args...] (try List)")
	}
	pid := 0
	if len(args) >= 2 && args[0] == "--pid" {
		n, err := strconv.Atoi(args[1])
		if err != nil || n <= 0 {
			return fmt.Errorf("--pid wants a process id, got %q", args[1])
		}
		pid, args = n, args[2:]
		if len(args) < 1 {
			return fmt.Errorf("usage: linkz debug [--pid N] Command [args...]")
		}
	}
	info, err := host.FindBackdoor(pid)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]any{"cmd": args[0], "args": args[1:]})
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://127.0.0.1:%d/cmd", info.Port), bytes.NewReader(body))
	req.Header.Set(host.TokenHeader, info.Token)
	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, _ := io.ReadAll(resp.Body)
	fmt.Print(string(out))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("command failed")
	}
	return nil
}
