package host

import (
	"fmt"
	"io"

	"linkz/lib"
	"linkz/lib/draw"
	"lizardscript/evaluator"
	"lizardscript/object"
)

// Paint runs a lizard program headless with the fb primitives on a UI buffer
// of w by h pixels, then saves that buffer to out. The drawing libraries and
// the graphics libraries are loaded first. It is how the lizard graphics
// libraries get tested without a window.
func Paint(source, out string, w, h int, stdout io.Writer) error {
	if w < 1 || h < 1 || w > MaxCanvasSide || h > MaxCanvasSide {
		return fmt.Errorf("paint: size must be between (1, 1) and (%d, %d), got (%d, %d)", MaxCanvasSide, MaxCanvasSide, w, h)
	}
	fb := NewFBUI(w, h)
	env := object.NewEnvironment()
	evaluator.Output = stdout
	env.Define("fb", fb.Builtins())
	env.Define("fs", fsHash())
	env.Define("img", imgHash())
	env.Define("time", timeHash())
	env.Define("clip", clipHash())
	env.Define("os", osHash())

	// label starts the error message; the source is named so that errors
	// from the libraries can be told from the user's program (named "").
	runSrc := func(label, srcName, code string, attr *attribution) error {
		if _, err := run(code, srcName, env, attr); err != nil {
			return fmt.Errorf("%s: %w", label, err)
		}
		return nil
	}
	for _, name := range draw.Order {
		code, err := draw.Files.ReadFile(name)
		if err != nil {
			return err
		}
		if err := runSrc(name, "draw/"+name, string(code), nil); err != nil {
			return err
		}
	}
	for _, name := range lib.Dir("gfx") {
		code, err := lib.Files.ReadFile(name)
		if err != nil {
			return err
		}
		if err := runSrc(name, "gfx/"+name, string(code), nil); err != nil {
			return err
		}
	}
	defer limitSteps()()
	defer trackLayers(env)()
	if err := runSrc("program", "", source, &attribution{}); err != nil {
		return err
	}
	fb.off, fb.target = nil, "ui"
	return fb.SavePNG(out, nil)
}
