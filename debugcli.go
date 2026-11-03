package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"linkz/host"
)

// debug sends one command to the running editor and prints the JSON reply.
func debug(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: linkz debug Command [args...] (try List)")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(home, ".linkz", "debug.json"))
	if err != nil {
		return fmt.Errorf("no running editor (%v)", err)
	}
	var info struct {
		Port  int    `json:"port"`
		Token string `json:"token"`
	}
	if err := json.Unmarshal(data, &info); err != nil {
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
