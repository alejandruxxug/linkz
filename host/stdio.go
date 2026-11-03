package host

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"io"
	"sync"

	"linkz/lib"
	"linkz/lib/draw"
	"linkz/lib/mcp"
	"linkz/lib/trace"
	"lizardscript/evaluator"
	"lizardscript/lexer"
	"lizardscript/object"
	"lizardscript/parser"
)

// The stdio primitives of `linkz mcp`: reading lines from stdin, writing to
// stdout and stderr, and encoding the current buffer as a PNG. The protocol
// itself (JSON, JSON-RPC, the tools) is written in lizard, in lib/mcp.

// MaxLine is the longest request line the server reads.
const MaxLine = 32 << 20

// Stdio is the io hash: the process's three standard streams.
type Stdio struct {
	in  *bufio.Reader
	out *bufio.Writer
	err io.Writer
	mu  sync.Mutex
}

// NewStdio wraps the three streams.
func NewStdio(in io.Reader, out, errw io.Writer) *Stdio {
	return &Stdio{in: bufio.NewReaderSize(in, 1<<16), out: bufio.NewWriter(out), err: errw}
}

// Builtins is the io hash.
func (s *Stdio) Builtins() *object.Hash {
	return builtinHash("io", map[string]object.BuiltinFunction{
		// read_line() is the next line of stdin without its line break, or
		// null at the end of the input.
		"read_line": func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return errorf("wrong number of arguments: want 0, got %d", len(args))
			}
			line, err := s.readLine()
			if err != nil {
				return errorf("%v", err)
			}
			if line == nil {
				return null
			}
			return newStr(string(line))
		},
		// write(text) writes to stdout and flushes.
		"write": func(args ...object.Object) object.Object {
			t, err := oneString(args, "text")
			if err != nil {
				return err
			}
			s.mu.Lock()
			defer s.mu.Unlock()
			if _, werr := s.out.WriteString(t); werr != nil {
				return errorf("%v", werr)
			}
			if werr := s.out.Flush(); werr != nil {
				return errorf("%v", werr)
			}
			return null
		},
		// log(text) writes a line to stderr.
		"log": func(args ...object.Object) object.Object {
			t, err := oneString(args, "text")
			if err != nil {
				return err
			}
			fmt.Fprintln(s.err, t)
			return null
		},
		// step_limit(n) starts counting steps from zero, with n as the most
		// the program may take before it is stopped. The server calls it at
		// the start of every request.
		"step_limit": func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return errorf("wrong number of arguments: want 1 (steps), got %d", len(args))
			}
			n, ok := args[0].(*object.Integer)
			if !ok || n.Value < 0 {
				return errorf("steps must be a whole number, 0 or more")
			}
			evaluator.MaxSteps = int(n.Value)
			evaluator.ResetSteps()
			return null
		},
	})
}

// readLine reads up to the next "\n" (a "\r" before it is dropped). nil
// means the input ended. A last line without a line break still counts.
func (s *Stdio) readLine() ([]byte, error) {
	var line []byte
	for {
		chunk, err := s.in.ReadSlice('\n')
		line = append(line, chunk...)
		if len(line) > MaxLine {
			return nil, fmt.Errorf("a request line is longer than %d bytes", MaxLine)
		}
		if err == bufio.ErrBufferFull {
			continue
		}
		if err == io.EOF {
			if len(line) == 0 {
				return nil, nil
			}
			break
		}
		if err != nil {
			return nil, err
		}
		break
	}
	line = bytes.TrimSuffix(line, []byte("\n"))
	line = bytes.TrimSuffix(line, []byte("\r"))
	return line, nil
}

// pngBase64 is fb.png_base64(): the current target as a PNG, base64 text.
func (f *FB) pngBase64(args ...object.Object) object.Object {
	if len(args) != 0 {
		return errorf("wrong number of arguments: want 0, got %d", len(args))
	}
	f.mu.Lock()
	b := f.current()
	f.mu.Unlock()
	var buf bytes.Buffer
	if err := png.Encode(&buf, b.Img); err != nil {
		return errorf("%v", err)
	}
	return newStr(base64.StdEncoding.EncodeToString(buf.Bytes()))
}

// ServeMCP runs the MCP server: it loads the drawing, graphics and MCP
// libraries, with the same environment as `linkz paint` plus src and io, and
// calls mcp_serve() from lib/mcp/server.liz. It returns when stdin ends.
func ServeMCP(in io.Reader, out, errw io.Writer) error {
	env, err := NewMCPEnv(in, out, errw)
	if err != nil {
		return err
	}
	return EvalIn(env, "mcp_serve", "mcp_serve()")
}

// EvalIn parses and runs code in env.
func EvalIn(env *object.Environment, name, code string) error {
	p := parser.New(lexer.New(code))
	prog := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		return fmt.Errorf("%s: %s", name, errs[0])
	}
	if e, ok := Eval(prog, env).(*object.Error); ok {
		return fmt.Errorf("%s: %s", name, e.Message)
	}
	return nil
}

// NewMCPEnv makes the environment of the MCP server with every library
// loaded. What the programs print goes to errw, never to the protocol's out.
func NewMCPEnv(in io.Reader, out, errw io.Writer) (*object.Environment, error) {
	evaluator.Output = errw
	evaluator.ResetMath()
	fb := NewFBUI(1, 1)
	fbh := fb.Builtins()
	hashSet(fbh, "png_base64", &object.Builtin{Name: "fb.png_base64", Fn: fb.pngBase64})

	settings, err := LoadSettings()
	if err != nil {
		fmt.Fprintln(errw, "linkz: settings ignored, using defaults:", err)
	}

	env := object.NewEnvironment()
	env.Define("fb", fbh)
	env.Define("src", NewSrc(nil).Builtins())
	env.Define("io", NewStdio(in, out, errw).Builtins())
	env.Define("fs", fsHash())
	env.Define("time", timeHash())
	env.Define("os", osHash())
	env.Define("img", imgHash())
	env.Define("settings", settings)

	load := func(name, code string) error { return EvalIn(env, name, code) }
	for _, name := range draw.Order {
		code, err := draw.Files.ReadFile(name)
		if err != nil {
			return nil, err
		}
		if err := load(name, string(code)); err != nil {
			return nil, err
		}
	}
	for _, name := range lib.Dir("gfx") {
		code, err := lib.Files.ReadFile(name)
		if err != nil {
			return nil, err
		}
		if err := load(name, string(code)); err != nil {
			return nil, err
		}
	}
	for _, name := range trace.Order {
		code, err := trace.Files.ReadFile(name)
		if err != nil {
			return nil, err
		}
		if err := load(name, string(code)); err != nil {
			return nil, err
		}
	}
	for _, name := range mcp.Order {
		code, err := mcp.Files.ReadFile(name)
		if err != nil {
			return nil, err
		}
		if err := load(name, string(code)); err != nil {
			return nil, err
		}
	}
	return env, nil
}
