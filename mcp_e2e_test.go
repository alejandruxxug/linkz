package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// session is a running `linkz mcp` with its pipes.
type session struct {
	t     *testing.T
	cmd   *exec.Cmd
	in    *bufio.Writer
	inRaw interface{ Close() error }
	out   *bufio.Reader
	err   bytes.Buffer
}

// start builds linkz once per test binary and runs `linkz mcp` with a home
// folder of its own, so nothing it saves reaches the real one.
func start(t *testing.T, bin string) *session {
	t.Helper()
	return startIn(t, bin, t.TempDir())
}

// startIn is start with the given home folder.
func startIn(t *testing.T, bin, home string) *session {
	t.Helper()
	cmd := exec.Command(bin, "mcp")
	cmd.Env = append(os.Environ(), "HOME="+home, "USERPROFILE="+home)
	in, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	s := &session{t: t, cmd: cmd, in: bufio.NewWriter(in), inRaw: in, out: bufio.NewReaderSize(out, 1<<20)}
	cmd.Stderr = &s.err
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	return s
}

// send writes one line.
func (s *session) send(line string) {
	s.t.Helper()
	s.in.WriteString(line + "\n")
	if err := s.in.Flush(); err != nil {
		s.t.Fatal(err)
	}
}

// reply reads one line and decodes it, failing if it is not one JSON object.
func (s *session) reply() map[string]any {
	s.t.Helper()
	done := make(chan struct{})
	var line string
	var err error
	go func() {
		line, err = s.out.ReadString('\n')
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(60 * time.Second):
		s.cmd.Process.Kill()
		s.t.Fatalf("no reply in 60 s; stderr:\n%s", s.err.String())
	}
	if err != nil {
		s.t.Fatalf("reading a reply: %v; stderr:\n%s", err, s.err.String())
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		s.t.Fatalf("a reply that is not JSON: %q (%v)", line, err)
	}
	if m["jsonrpc"] != "2.0" {
		s.t.Fatalf("reply without jsonrpc 2.0: %s", line)
	}
	return m
}

// call sends a tools/call and returns its result object.
func (s *session) call(id int, name string, args map[string]any) map[string]any {
	s.t.Helper()
	req, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": id, "method": "tools/call",
		"params": map[string]any{"name": name, "arguments": args}})
	s.send(string(req))
	m := s.reply()
	if int(m["id"].(float64)) != id {
		s.t.Fatalf("reply for id %v, want %d", m["id"], id)
	}
	res, ok := m["result"].(map[string]any)
	if !ok {
		s.t.Fatalf("no result: %v", m)
	}
	return res
}

func TestMCPSession(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "linkz")
	if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	s := start(t, bin)

	// initialize
	s.send(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test","version":"0"}}}`)
	m := s.reply()
	res := m["result"].(map[string]any)
	if res["protocolVersion"] != "2025-06-18" || res["serverInfo"].(map[string]any)["name"] != "linkz" {
		t.Fatalf("initialize: %v", m)
	}
	s.send(`{"jsonrpc":"2.0","method":"notifications/initialized"}`) // no reply

	// tools/list
	s.send(`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`)
	m = s.reply()
	var names []string
	for _, tool := range m["result"].(map[string]any)["tools"].([]any) {
		names = append(names, tool.(map[string]any)["name"].(string))
	}
	if got, want := len(names), 7; got != want {
		t.Fatalf("tools = %v, want %d", names, want)
	}

	// run
	res = s.call(3, "run", map[string]any{"code": "circle((10, 10), 5, #ff0000)\nputs(\"hi\")"})
	var report struct {
		OK      bool
		Printed string
		Shapes  int
		Steps   int
		Errors  []any
	}
	text := res["content"].([]any)[0].(map[string]any)["text"].(string)
	if err := json.Unmarshal([]byte(text), &report); err != nil {
		t.Fatal(err)
	}
	if !report.OK || report.Printed != "hi\n" || report.Shapes != 1 || report.Steps == 0 {
		t.Errorf("run report = %+v", report)
	}

	// render: a 64 by 48 canvas with a red circle in the middle.
	code := `canvas = {"size": (64, 48), "background": #336699}` + "\n" + `circle((32, 24), 10, #ff0000)`
	res = s.call(4, "render", map[string]any{"code": code})
	item := res["content"].([]any)[0].(map[string]any)
	if item["type"] != "image" || item["mimeType"] != "image/png" {
		t.Fatalf("render item = %v", item)
	}
	raw, err := base64.StdEncoding.DecodeString(item["data"].(string))
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	if b := img.Bounds(); b.Dx() != 64 || b.Dy() != 48 {
		t.Errorf("image size = %v, want 64x48", b)
	}
	if r, g, b, _ := img.At(32, 24).RGBA(); r>>8 != 255 || g>>8 != 0 || b>>8 != 0 {
		t.Errorf("center pixel = %d %d %d, want red", r>>8, g>>8, b>>8)
	}
	if r, g, b, _ := img.At(2, 2).RGBA(); r>>8 != 0x33 || g>>8 != 0x66 || b>>8 != 0x99 {
		t.Errorf("corner pixel = %d %d %d, want #336699", r>>8, g>>8, b>>8)
	}

	// render scaled to fit a width: the picture is twice as wide.
	res = s.call(5, "render", map[string]any{"code": code, "width": 128})
	raw, _ = base64.StdEncoding.DecodeString(res["content"].([]any)[0].(map[string]any)["data"].(string))
	img, err = png.Decode(bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	if b := img.Bounds(); b.Dx() != 128 || b.Dy() != 96 {
		t.Errorf("scaled size = %v, want 128x96", b)
	}
	if r, g, b, _ := img.At(64, 48).RGBA(); r>>8 != 255 || g>>8 != 0 || b>>8 != 0 {
		t.Errorf("scaled center pixel = %d %d %d, want red", r>>8, g>>8, b>>8)
	}

	// A malformed line gets a parse error and the server goes on.
	s.send(`{"jsonrpc": "2.0", "id": 6, "method": "ping"`)
	m = s.reply()
	if e, ok := m["error"].(map[string]any); !ok || int(e["code"].(float64)) != -32700 {
		t.Fatalf("malformed line: %v", m)
	}
	s.send(`{"jsonrpc":"2.0","id":7,"method":"ping"}`)
	if m = s.reply(); m["result"] == nil || int(m["id"].(float64)) != 7 {
		t.Fatalf("ping after the malformed line: %v", m)
	}

	// An unknown method and a bad tool call do not stop it either.
	s.send(`{"jsonrpc":"2.0","id":8,"method":"nope"}`)
	if e, ok := s.reply()["error"].(map[string]any); !ok || int(e["code"].(float64)) != -32601 {
		t.Fatal("unknown method should be -32601")
	}
	res = s.call(9, "run", map[string]any{"code": "while (true) { }"})
	text = res["content"].([]any)[0].(map[string]any)["text"].(string)
	if !bytes.Contains([]byte(text), []byte("ran too long")) {
		t.Errorf("an endless drawing should be stopped, got %s", text)
	}
	res = s.call(10, "save", map[string]any{"path": "/etc/linkz-test.lzdraw", "code": "puts(1)"})
	if res["isError"] != true {
		t.Errorf("save outside home should be refused: %v", res)
	}
	if _, err := os.Stat("/etc/linkz-test.lzdraw"); err == nil {
		t.Error("the file outside home was written")
	}

	// EOF: the server exits 0.
	s.inRaw.Close()
	done := make(chan error, 1)
	go func() { done <- s.cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("exit after EOF: %v; stderr:\n%s", err, s.err.String())
		}
	case <-time.After(10 * time.Second):
		s.cmd.Process.Kill()
		t.Fatal("the server did not exit at EOF")
	}
}

// TestMCPRenderTime renders examples/layers.lzdraw and reports the time of
// each call (run with -v to see it).
func TestMCPRenderTime(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "linkz")
	if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	code, err := os.ReadFile("examples/layers.lzdraw")
	if err != nil {
		t.Fatal(err)
	}
	s := start(t, bin)
	for i := 1; i <= 3; i++ {
		start := time.Now()
		res := s.call(i, "render", map[string]any{"code": string(code)})
		took := time.Since(start)
		if res["isError"] == true || res["content"].([]any)[0].(map[string]any)["type"] != "image" {
			t.Fatalf("render failed: %v", res)
		}
		t.Logf("render %d of examples/layers.lzdraw: %v", i, took)
		if took > 10*time.Second {
			t.Errorf("render took %v", took)
		}
	}
	s.inRaw.Close()
	s.cmd.Wait()
}

// TestMCPTrace writes a PNG of a circle in the home folder, traces it, and
// checks the three items: code that `run` accepts, the stats, and an image.
func TestMCPTrace(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "linkz")
	if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	home := t.TempDir()
	pic := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			c := color.RGBA{250, 250, 240, 255}
			if (x-32)*(x-32)+(y-32)*(y-32) <= 20*20 {
				c = color.RGBA{30, 90, 200, 255}
			}
			pic.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, pic); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, "circle.png"), buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	s := startIn(t, bin, home)

	res := s.call(1, "trace", map[string]any{"path": "~/circle.png", "colors": 2})
	items, _ := res["content"].([]any)
	if res["isError"] == true || len(items) != 3 {
		t.Fatalf("trace result = %v", res)
	}
	code := items[0].(map[string]any)["text"].(string)
	var stats map[string]any
	if err := json.Unmarshal([]byte(items[1].(map[string]any)["text"].(string)), &stats); err != nil {
		t.Fatalf("stats are not JSON: %v", err)
	}
	if stats["colors"] == nil {
		t.Errorf("stats = %v", stats)
	}
	img := items[2].(map[string]any)
	if img["type"] != "image" || img["mimeType"] != "image/png" {
		t.Fatalf("third item = %v", img["type"])
	}
	raw, err := base64.StdEncoding.DecodeString(img["data"].(string))
	if err != nil {
		t.Fatal(err)
	}
	traced, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	r, g, b, _ := traced.At(traced.Bounds().Dx()/2, traced.Bounds().Dy()/2).RGBA()
	if b>>8 < 150 || r>>8 > 100 || g>>8 > 140 {
		t.Errorf("center of the traced picture = %d %d %d, want blue", r>>8, g>>8, b>>8)
	}

	// The code is a drawing that run accepts.
	res = s.call(2, "run", map[string]any{"code": code})
	var report struct {
		OK     bool
		Shapes int
	}
	if err := json.Unmarshal([]byte(res["content"].([]any)[0].(map[string]any)["text"].(string)), &report); err != nil {
		t.Fatal(err)
	}
	if !report.OK || report.Shapes == 0 {
		t.Errorf("run of the traced code = %+v", report)
	}

	// Refused paths.
	for i, p := range []string{"/etc/hosts.png", "~/circle.txt", "~/.hidden/x.png", "~/gone.png"} {
		if res = s.call(3+i, "trace", map[string]any{"path": p}); res["isError"] != true {
			t.Errorf("trace %s should be refused: %v", p, res)
		}
	}
	s.inRaw.Close()
	s.cmd.Wait()
}
