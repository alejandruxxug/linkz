package host

import (
	"fmt"
	"io"

	"linkz/lib"
	"linkz/lib/draw"
	"lizardscript/evaluator"
	"lizardscript/lexer"
	"lizardscript/object"
	"lizardscript/parser"
)

// Paint runs a lizard program headless with the fb primitives on a UI buffer
// of w by h pixels, then saves that buffer to out. The drawing libraries and
// the graphics libraries are loaded first. It is how the lizard graphics
// libraries get tested without a window.
func Paint(source, out string, w, h int, stdout io.Writer) error {
	fb := NewFBUI(w, h)
	env := object.NewEnvironment()
	evaluator.Output = stdout
	env.Define("fb", fb.Builtins())
	env.Define("fs", fsHash())
	env.Define("time", timeHash())
	env.Define("clip", clipHash())
	env.Define("os", osHash())

	run := func(name, code string) error {
		p := parser.New(lexer.New(code))
		prog := p.ParseProgram()
		if errs := p.Errors(); len(errs) > 0 {
			return fmt.Errorf("%s: %s", name, errs[0])
		}
		if e, ok := evaluator.Eval(prog, env).(*object.Error); ok {
			return fmt.Errorf("%s: %s", name, e.Message)
		}
		return nil
	}
	for _, name := range draw.Order {
		code, err := draw.Files.ReadFile(name)
		if err != nil {
			return err
		}
		if err := run(name, string(code)); err != nil {
			return err
		}
	}
	for _, name := range lib.Dir("gfx") {
		code, err := lib.Files.ReadFile(name)
		if err != nil {
			return err
		}
		if err := run(name, string(code)); err != nil {
			return err
		}
	}
	defer trackLayers(env)()
	if err := run("program", source); err != nil {
		return err
	}
	fb.off, fb.target = nil, "ui"
	return fb.SavePNG(out, nil)
}
