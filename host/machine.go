package host

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"linkz/lib"
	"linkz/lib/draw"
	"lizardscript/evaluator"
	"lizardscript/lexer"
	"lizardscript/object"
	"lizardscript/parser"
)

// What this machine can handle, measured once and kept in ~/.linkz/machine.liz
// (plain lizard data: `let machine = {...}`). The measuring is calibrate() in
// lib/gfx/calibrate.liz; this file only decides when to run it and where the
// result lives.

// AppVersion is the app's version, set by main. A new version measures again,
// since the drawing code may have got faster or slower.
var AppVersion = "dev"

// MachinePath is ~/.linkz/machine.liz.
func MachinePath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "machine.liz"), nil
}

// LoadMachine returns the kept measurement, or measures when there is none,
// when the file can't be read, or when it is from another app version (or
// force is set). measure gives the new file's text (a `let machine = {...}`
// program); it is written to path. The bool says whether it measured.
func LoadMachine(path string, force bool, measure func() (string, error)) (*object.Hash, bool, error) {
	if !force {
		if h, err := LoadSettingsFile(path); err == nil {
			if m, ok := hashGet(h, "machine").(*object.Hash); ok {
				if ver, ok := hashGet(m, "app_version").(*object.String); ok && ver.Value == AppVersion {
					return m, false, nil
				}
			}
		}
	}
	text, err := measure()
	if err != nil {
		return nil, true, err
	}
	if err := WriteFileAtomic(path, []byte(text)); err != nil {
		return nil, true, err
	}
	h, err := evalSettings(text)
	if err != nil {
		return nil, true, err
	}
	m, ok := hashGet(h, "machine").(*object.Hash)
	if !ok {
		return nil, true, fmt.Errorf("machine file has no machine")
	}
	return m, true, nil
}

// Machine is the startup entry: the kept measurement for this machine, which
// runs the benchmark when it has to. force re-measures (linkz calibrate).
func Machine(force bool, stdout io.Writer) (*object.Hash, error) {
	path, err := MachinePath()
	if err != nil {
		return nil, err
	}
	m, _, err := LoadMachine(path, force, func() (string, error) {
		return runCalibration("machine_source(calibrate())", stdout)
	})
	return m, err
}

// Recalibrate measures now, keeps the result in the machine file and returns
// it as JSON text (linkz calibrate prints it).
func Recalibrate(stdout io.Writer) (string, error) {
	path, err := MachinePath()
	if err != nil {
		return "", err
	}
	var json string
	_, _, err = LoadMachine(path, true, func() (string, error) {
		both, err := runCalibration("let m = calibrate()\nmachine_source(m) + machine_json(m)", stdout)
		if err != nil {
			return "", err
		}
		source, rest, _ := strings.Cut(both, "\n")
		json = rest
		return source + "\n", nil
	})
	return json, err
}

// runCalibration loads the libraries like Paint does, on a headless buffer,
// and evaluates expr, which must give a STRING.
func runCalibration(expr string, stdout io.Writer) (string, error) {
	prevOut := evaluator.Output
	defer func() { evaluator.Output = prevOut }()
	fb := NewFBUI(512, 512)
	env := object.NewEnvironment()
	evaluator.Output = stdout
	env.Define("fb", fb.Builtins())
	env.Define("time", timeHash())
	env.Define("os", osHash())
	env.Define("src", NewSrc(nil).Builtins())

	run := func(name, code string) (object.Object, error) {
		p := parser.New(lexer.New(code))
		prog := p.ParseProgram()
		if errs := p.Errors(); len(errs) > 0 {
			return nil, fmt.Errorf("%s: %s", name, errs[0])
		}
		res := evaluator.Eval(prog, env)
		if e, ok := res.(*object.Error); ok {
			return nil, fmt.Errorf("%s: %s", name, e.Message)
		}
		return res, nil
	}
	for _, name := range draw.Order {
		code, err := draw.Files.ReadFile(name)
		if err != nil {
			return "", err
		}
		if _, err := run(name, string(code)); err != nil {
			return "", err
		}
	}
	for _, name := range lib.Dir("gfx") {
		code, err := lib.Files.ReadFile(name)
		if err != nil {
			return "", err
		}
		if _, err := run(name, string(code)); err != nil {
			return "", err
		}
	}
	defer trackLayers(env)()
	res, err := run("calibrate", expr)
	if err != nil {
		return "", err
	}
	s, ok := res.(*object.String)
	if !ok {
		return "", fmt.Errorf("calibrate gave %s, not a STRING", res.Type())
	}
	return s.Value, nil
}
