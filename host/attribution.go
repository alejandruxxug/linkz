package host

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"lizardscript/ast"
	"lizardscript/evaluator"
	"lizardscript/lexer"
	"lizardscript/object"
	"lizardscript/parser"
	"lizardscript/token"
)

// Attribution makes a runtime error point at the user's program. The
// libraries are loaded with names (lexer.NewNamed) and the user's program
// with the name "", so an error knows which source it came from. An error
// raised in library code has a line in the library, which means nothing to
// the user: it is reported at the innermost user statement that was running.
type attribution struct {
	user map[ast.Statement]bool // the statements that count (see userStatements)
	tok  token.Token            // the innermost user statement running now
	seen bool
}

// newAttribution tracks the statements of the user's program. Call hook as
// evaluator.OnStatement while it runs.
func newAttribution(user *ast.Program) *attribution {
	return &attribution{user: userStatements(user)}
}

// noLines is an attribution for code that is not on any user line: library
// errors lose their library position and get none.
func noLines() *attribution {
	return &attribution{user: map[ast.Statement]bool{}}
}

// hook is an evaluator.OnStatement hook.
func (a *attribution) hook(n ast.Statement) {
	if a.user[n] {
		a.tok, a.seen = stmtToken(n), true
	}
}

var libPrefixRe = regexp.MustCompile(`^\[line \d+, column \d+\] (runtime error: )?`)

// message is the text of a runtime error. An error from the user's own
// source keeps its message. An error from a library gets the position of the
// user statement instead, with the library's position cut off.
func (a *attribution) message(e *object.Error) string {
	if e.Source == "" {
		return e.Message
	}
	msg := libPrefixRe.ReplaceAllString(e.Message, "")
	if !a.seen {
		return msg
	}
	return "[line " + strconv.Itoa(a.tok.Line) + ", column " + strconv.Itoa(a.tok.Column) + "] runtime error: " + msg
}

// parseNamed parses source whose tokens carry name. Syntax errors come back
// as one Go error.
func parseNamed(source, name string) (*ast.Program, error) {
	p := parser.New(lexer.NewNamed(source, name))
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return nil, errors.New(strings.Join(p.Errors(), "\n"))
	}
	return prog, nil
}

// stmtToken is the first token of a statement.
func stmtToken(n ast.Statement) token.Token {
	switch n := n.(type) {
	case *ast.LetStatement:
		return n.Token
	case *ast.ReturnStatement:
		return n.Token
	case *ast.BreakStatement:
		return n.Token
	case *ast.ExpressionStatement:
		return n.Token
	case *ast.BlockStatement:
		return n.Token
	case *ast.LayerStatement:
		return n.Token
	}
	return token.Token{}
}

// limitSteps sets the step limit for one run (RunSteps) with fresh counters,
// and returns the function that puts back what the caller had.
func limitSteps() (restore func()) {
	prevMax, prevSteps := evaluator.MaxSteps, evaluator.Steps()
	evaluator.MaxSteps = RunSteps
	evaluator.ResetSteps()
	evaluator.ResetMath()
	return func() {
		evaluator.MaxSteps = prevMax
		evaluator.SetSteps(prevSteps)
	}
}

// MaxCanvasSide is the biggest width or height of a canvas or a paint buffer.
const MaxCanvasSide = 16384
