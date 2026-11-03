// Package host is the thin native layer under Linkz. It runs lizard-script
// and gives it the few things only Go can do. The drawing itself, from the
// shapes to the SVG text, is written in lizard-script in lib/.
package host

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"linkz/lib/draw"
	"lizardscript/evaluator"
	"lizardscript/lexer"
	"lizardscript/object"
	"lizardscript/parser"
)

// Render runs a .lzdraw program and returns the SVG of what it drew. What the
// program prints with puts goes to out.
func Render(source string, out io.Writer) (string, error) {
	evaluator.Output = out
	env := object.NewEnvironment()

	for _, name := range draw.Order {
		code, err := draw.Files.ReadFile(name)
		if err != nil {
			return "", err
		}
		if _, err := run(string(code), env); err != nil {
			return "", fmt.Errorf("library %s: %w", name, err)
		}
	}

	if _, err := run(source, env); err != nil {
		return "", err
	}

	result, err := run("svg_document(canvas, drawing)", env)
	if err != nil {
		return "", fmt.Errorf("library svg.liz: %w", err)
	}
	svg, ok := result.(*object.String)
	if !ok {
		return "", fmt.Errorf("svg_document gave %s, not a STRING", result.Type())
	}
	return svg.Value, nil
}

// run parses and evaluates source in env. Syntax errors and runtime errors
// come back as a Go error with lizard-script's own message.
func run(source string, env *object.Environment) (object.Object, error) {
	p := parser.New(lexer.New(source))
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return nil, errors.New(strings.Join(p.Errors(), "\n"))
	}

	result := evaluator.Eval(program, env)
	if err, ok := result.(*object.Error); ok {
		return nil, errors.New(err.Message)
	}
	return result, nil
}
