// Package host is the thin native layer under Linkz. It runs lizard-script
// and gives it the few things only Go can do. The drawing itself, from the
// shapes to the SVG text, is written in lizard-script in lib/.
package host

import (
	"errors"
	"fmt"
	"io"

	"linkz/lib/draw"
	"lizardscript/evaluator"
	"lizardscript/object"
)

// Render runs a .lzdraw program and returns the SVG of what it drew. What the
// program prints with puts goes to out. The program runs under the same step
// limit as the editor (RunSteps), so a loop that never ends is an error. An
// error raised inside a library is reported at the line of the user's
// statement that called it.
func Render(source string, out io.Writer) (string, error) {
	evaluator.Output = out
	env := object.NewEnvironment()

	for _, name := range draw.Order {
		code, err := draw.Files.ReadFile(name)
		if err != nil {
			return "", err
		}
		if _, err := run(string(code), "draw/"+name, env, nil); err != nil {
			return "", fmt.Errorf("library %s: %w", name, err)
		}
	}

	defer limitSteps()()
	defer trackLayers(env)()
	if _, err := run(source, "", env, &attribution{}); err != nil {
		return "", err
	}

	// The closing call is not on any user line: a library error in it (a bad
	// canvas) is reported with its own message.
	result, err := run("svg_document(canvas, drawing)", "draw/svg.liz", env, noLines())
	if err != nil {
		return "", err
	}
	svg, ok := result.(*object.String)
	if !ok {
		return "", fmt.Errorf("svg_document gave %s, not a STRING", result.Type())
	}
	return svg.Value, nil
}

// run parses and evaluates source (a source called name) in env. Syntax
// errors and runtime errors come back as a Go error with lizard-script's own
// message. With a non-nil attr, a runtime error from library code is moved to
// the user's line, and attr watches the user's statements while it runs.
func run(source, name string, env *object.Environment, attr *attribution) (object.Object, error) {
	program, err := parseNamed(source, name)
	if err != nil {
		return nil, err
	}
	if attr != nil {
		if attr.user == nil {
			attr.user = userStatements(program)
		}
		evaluator.OnStatement = attr.hook
	}
	result := evaluator.Eval(program, env)
	if attr != nil {
		evaluator.OnStatement = nil
	}
	if e, ok := result.(*object.Error); ok {
		if attr != nil {
			return nil, errors.New(attr.message(e))
		}
		return nil, errors.New(e.Message)
	}
	return result, nil
}
