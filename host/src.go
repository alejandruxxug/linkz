package host

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"linkz/lib/draw"
	"lizardscript/ast"
	"lizardscript/evaluator"
	"lizardscript/lexer"
	"lizardscript/object"
	"lizardscript/parser"
	"lizardscript/token"
)

// RunSteps is the step limit for running a drawing.
const RunSteps = 2_000_000

// libLineBase is added to every line of the drawing libraries' source, so an
// error's line tells whether it came from library code (>= libLineBase) or
// from the user's drawing, which keeps its own numbers.
const libLineBase = 100_000

// Src is the "src" hash: the language's own parser and evaluator as data, so
// the editor (written in lizard) can read and run the user's drawing.
type Src struct {
	stats *Stats

	mu    sync.Mutex
	cache map[string]*ast.Program // parsed libraries
}

// NewSrc makes the src primitives. stats may be nil.
func NewSrc(stats *Stats) *Src {
	if stats == nil {
		stats = &Stats{}
	}
	return &Src{stats: stats, cache: map[string]*ast.Program{}}
}

// Builtins is the src hash for the UI environment.
func (s *Src) Builtins() *object.Hash {
	return builtinHash("src", map[string]object.BuiltinFunction{
		"parse": func(args ...object.Object) object.Object {
			code, err := oneString(args, "code")
			if err != nil {
				return err
			}
			return fromGo(parseResult(code))
		},
		"run": func(args ...object.Object) object.Object {
			code, err := oneString(args, "code")
			if err != nil {
				return err
			}
			return s.Run(code)
		},
		"library_docs": func(args ...object.Object) object.Object {
			return fromGo(s.libraryDocs())
		},
	})
}

func oneString(args []object.Object, name string) (string, *object.Error) {
	s, err := strArgs(args, 1, "1 ("+name+")")
	if err != nil {
		return "", err
	}
	return s[0], nil
}

// ---- src.parse ----

var posRe = regexp.MustCompile(`^\[line (\d+), column (\d+)\] (.*)$`)

// splitPosition turns "[line L, column C] message" into its parts.
func splitPosition(msg string) SrcError {
	if m := posRe.FindStringSubmatch(msg); m != nil {
		l, _ := strconv.Atoi(m[1])
		c, _ := strconv.Atoi(m[2])
		return SrcError{l, c, m[3]}
	}
	return SrcError{0, 0, msg}
}

// syntaxErrors parses code with the language's parser and returns its errors.
func syntaxErrors(code string) (*ast.Program, []SrcError) {
	p := parser.New(lexer.New(code))
	prog := p.ParseProgram()
	var errs []SrcError
	for _, e := range p.Errors() {
		errs = append(errs, splitPosition(e))
	}
	return prog, errs
}

func parseResult(code string) map[string]any {
	sc := ScanSource(code)
	_, errs := syntaxErrors(code)
	errList := []any{}
	for _, e := range errs {
		errList = append(errList, errorMap(e))
	}
	lines := []any{}
	for _, li := range sc.Lines {
		flags := []any{}
		for _, f := range li.Flags {
			flags = append(flags, f)
		}
		lines = append(lines, map[string]any{
			"line": li.Line, "indent": li.Indent, "text": li.Text,
			"comment": li.Comment, "name": li.Name, "flags": flags,
		})
	}
	return map[string]any{"lines": lines, "statements": stmtList(sc.Stmts), "errors": errList}
}

func errorMap(e SrcError) map[string]any {
	return map[string]any{"line": e.Line, "column": e.Column, "message": e.Message}
}

func stmtList(stmts []*Stmt) []any {
	out := []any{}
	for _, s := range stmts {
		out = append(out, stmtMap(s))
	}
	return out
}

func stmtMap(s *Stmt) map[string]any {
	flags := []any{}
	for _, f := range s.Flags {
		flags = append(flags, f)
	}
	m := map[string]any{
		"line": s.Line, "end_line": s.EndLine, "column": s.Column, "depth": s.Depth,
		"kind": s.Kind, "name": s.Name, "callee": s.Callee, "label": s.Label, "flags": flags,
		"args": argList(s.Args), "children": stmtList(s.Children),
	}
	if s.Kind == "function" {
		m["params"] = paramList(s.Params)
		m["doc"] = docMap(s.Doc)
	}
	return m
}

func paramList(names []string) []any {
	out := []any{}
	for _, n := range names {
		out = append(out, n)
	}
	return out
}

func docMap(d *Doc) map[string]any {
	params := map[string]any{}
	for k, v := range d.Params {
		params[k] = v
	}
	var ex any
	if d.HasEx {
		ex = d.Example
	}
	return map[string]any{"summary": d.Summary, "params": params, "example": ex}
}

func argList(args []Arg) []any {
	out := []any{}
	for _, a := range args {
		out = append(out, map[string]any{
			"text": a.Text, "line": a.Line, "end_line": a.EndLine, "start": a.Start, "end": a.End,
			"kind": a.Kind, "items": argList(a.Items),
		})
	}
	return out
}

// libraryDocs collects the docs of every function the drawing libraries define.
func (s *Src) libraryDocs() map[string]any {
	out := map[string]any{}
	for _, name := range draw.Order {
		code, err := draw.Files.ReadFile(name)
		if err != nil {
			continue
		}
		var walk func(stmts []*Stmt)
		walk = func(stmts []*Stmt) {
			for _, st := range stmts {
				if st.Kind == "function" {
					out[st.Name] = map[string]any{"params": paramList(st.Params), "doc": docMap(st.Doc)}
				}
			}
		}
		walk(ScanSource(string(code)).Stmts)
	}
	return out
}

// ---- src.run ----

// libProgram parses a drawing library once and keeps the tree.
func (s *Src) libProgram(name string) (*ast.Program, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p, ok := s.cache[name]; ok {
		return p, nil
	}
	code, err := draw.Files.ReadFile(name)
	if err != nil {
		return nil, err
	}
	prog, errs := syntaxErrors(strings.Repeat("\n", libLineBase) + string(code))
	if len(errs) > 0 {
		return nil, fmt.Errorf("library %s: [line %d, column %d] %s", name, errs[0].Line-libLineBase, errs[0].Column, errs[0].Message)
	}
	s.cache[name] = prog
	return prog, nil
}

// userStatements collects the statements of a program that count for tagging:
// the top level, and inside layers, loops and ifs, but not function bodies.
func userStatements(prog *ast.Program) map[ast.Statement]bool {
	set := map[ast.Statement]bool{}
	var block func(b *ast.BlockStatement)
	var stmt func(st ast.Statement)
	block = func(b *ast.BlockStatement) {
		if b == nil {
			return
		}
		for _, st := range b.Statements {
			stmt(st)
		}
	}
	stmt = func(st ast.Statement) {
		set[st] = true
		switch n := st.(type) {
		case *ast.LayerStatement:
			block(n.Body)
		case *ast.ExpressionStatement:
			switch e := n.Expression.(type) {
			case *ast.IfExpression:
				block(e.Consequence)
				block(e.Alternative)
			case *ast.WhileExpression:
				block(e.Body)
			case *ast.ForExpression:
				block(e.Body)
			}
		}
	}
	for _, st := range prog.Statements {
		stmt(st)
	}
	return set
}

// Run runs a drawing in a fresh environment and returns the result hash.
// It restores the caller's output, step limit and hooks, so the UI can call
// it in the middle of its own evaluation.
func (s *Src) Run(code string) object.Object {
	start := time.Now()
	result := map[string]any{"shapes": []any{}, "errors": []any{}, "printed": "", "steps": 0, "ms": 0.0}
	var errs []SrcError

	// Save what the caller had.
	prevOut, prevMax := evaluator.Output, evaluator.MaxSteps
	prevStmt, prevLayer := evaluator.OnStatement, evaluator.OnLayer
	prevSteps := evaluator.Steps()
	var printed bytes.Buffer
	defer func() {
		evaluator.Output, evaluator.MaxSteps = prevOut, prevMax
		evaluator.OnStatement, evaluator.OnLayer = prevStmt, prevLayer
		evaluator.SetSteps(prevSteps)
	}()
	evaluator.Output = &printed
	evaluator.MaxSteps = RunSteps
	evaluator.OnStatement, evaluator.OnLayer = nil, nil
	evaluator.ResetSteps()
	evaluator.ResetMath()

	env := object.NewEnvironment()
	var shapes object.Object = &object.Array{}
	var canvas object.Object = null

	prog, perrs := syntaxErrors(code)
	if len(perrs) > 0 {
		errs = perrs
	} else {
		// Tag shapes: the hooks track the innermost user statement and the
		// layer path; the libraries' draw() reads them through here().
		user := userStatements(prog)
		line := 0
		var layers []string
		libsLoaded := false
		evaluator.OnStatement = func(n ast.Statement) {
			if libsLoaded && user[n] {
				line = stmtLine(n)
			}
		}
		evaluator.OnLayer = func(name string, enter bool, tok token.Token) {
			if !libsLoaded {
				return
			}
			if enter {
				layers = append(layers, name)
			} else if len(layers) > 0 {
				layers = layers[:len(layers)-1]
			}
		}
		here := &object.Builtin{Name: "here", Fn: func(args ...object.Object) object.Object {
			path := &object.Array{}
			for _, l := range layers {
				path.Elements = append(path.Elements, newStr(l))
			}
			h := object.NewHash()
			hashSet(h, "line", newInt(line))
			hashSet(h, "layer", path)
			return h
		}}
		if lerr := s.loadLibs(env, func() {
			env.Define("here", here)
			env.Define("src", builtinHash("src", map[string]object.BuiltinFunction{"here": here.Fn}))
		}); lerr != nil {
			errs = []SrcError{{0, 0, lerr.Error()}}
		} else {
			libsLoaded = true
			res := evaluator.Eval(prog, env)
			if e, ok := res.(*object.Error); ok {
				errs = []SrcError{runError(e, line)}
			}
		}
	}

	if v, ok := env.Get("drawing"); ok {
		shapes = v
	}
	if v, ok := env.Get("canvas"); ok {
		canvas = v
	}
	ms := float64(time.Since(start).Microseconds()) / 1000
	errList := []any{}
	for _, e := range errs {
		errList = append(errList, errorMap(e))
	}
	result["errors"], result["printed"], result["ms"] = errList, printed.String(), ms
	result["steps"] = evaluator.Steps()
	out := fromGo(result).(*object.Hash)
	hashSet(out, "shapes", shapes)
	hashSet(out, "canvas", canvas)
	s.stats.RunMS, s.stats.RunErrors = ms, len(errs)
	return out
}

// runError turns a runtime error into a SrcError. An error raised inside a
// drawing library has a position in the library's source, which means
// nothing to the user: it is reported at userLine (the innermost user
// statement that was running), column 0, with the library's message alone.
// An error in the user's own code keeps its position.
func runError(e *object.Error, userLine int) SrcError {
	se := splitPosition(e.Message)
	if e.Line >= libLineBase {
		se.Line, se.Column = userLine, 0
		se.Message = strings.TrimPrefix(se.Message, "runtime error: ")
	}
	return se
}

// loadLibs evaluates the drawing libraries into env. After here.liz it
// calls afterHere, so the host can replace here() with the real one.
func (s *Src) loadLibs(env *object.Environment, afterHere func()) error {
	for _, name := range draw.Order {
		prog, err := s.libProgram(name)
		if err != nil {
			return err
		}
		if res := evaluator.Eval(prog, env); res != nil {
			if e, ok := res.(*object.Error); ok {
				return fmt.Errorf("library %s: %s", name, e.Message)
			}
		}
		if name == "here.liz" && afterHere != nil {
			afterHere()
		}
	}
	return nil
}

// stmtLine is the line a statement starts on.
func stmtLine(n ast.Statement) int {
	switch n := n.(type) {
	case *ast.LetStatement:
		return n.Token.Line
	case *ast.ReturnStatement:
		return n.Token.Line
	case *ast.BreakStatement:
		return n.Token.Line
	case *ast.ExpressionStatement:
		return n.Token.Line
	case *ast.BlockStatement:
		return n.Token.Line
	case *ast.LayerStatement:
		return n.Token.Line
	}
	return 0
}
