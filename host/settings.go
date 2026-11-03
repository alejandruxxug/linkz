package host

import (
	"fmt"
	"os"
	"path/filepath"

	"lizardscript/ast"
	"lizardscript/evaluator"
	"lizardscript/lexer"
	"lizardscript/object"
	"lizardscript/parser"

	"linkz/host/pressure"
)

// SettingsSteps is the step limit for evaluating settings.liz.
const SettingsSteps = 100_000

// SettingsPath is ~/.linkz/settings.liz (the folder is created if needed).
func SettingsPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "settings.liz"), nil
}

// LoadSettings evaluates the settings file once, at startup, and returns its
// variables as a hash: `let language = "es"` gives {"language": "es"}.
// A missing file is not a problem (empty hash, nil error). A file that does
// not parse, fails, or runs too long gives an empty hash and an error to warn
// about, so the editor starts with the defaults.
func LoadSettings() (*object.Hash, error) {
	path, err := SettingsPath()
	if err != nil {
		return object.NewHash(), err
	}
	return LoadSettingsFile(path)
}

// LoadSettingsFile is LoadSettings for a given path.
func LoadSettingsFile(path string) (*object.Hash, error) {
	body, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return object.NewHash(), nil
	}
	if err != nil {
		return object.NewHash(), err
	}
	h, err := evalSettings(string(body))
	if err != nil {
		return object.NewHash(), fmt.Errorf("%s: %w", path, err)
	}
	return h, nil
}

// evalSettings runs settings code in its own small environment (nothing but
// the language's builtins) and collects its top-level `let` variables.
func evalSettings(code string) (*object.Hash, error) {
	p := parser.New(lexer.New(code))
	prog := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		return nil, fmt.Errorf("%s", errs[0])
	}

	// Leave the caller's output, limit and hooks as they were.
	prevOut, prevMax := evaluator.Output, evaluator.MaxSteps
	prevStmt, prevLayer := evaluator.OnStatement, evaluator.OnLayer
	prevSteps := evaluator.Steps()
	defer func() {
		evaluator.Output, evaluator.MaxSteps = prevOut, prevMax
		evaluator.OnStatement, evaluator.OnLayer = prevStmt, prevLayer
		evaluator.SetSteps(prevSteps)
	}()
	evaluator.Output = os.Stderr
	evaluator.MaxSteps = SettingsSteps
	evaluator.OnStatement, evaluator.OnLayer = nil, nil
	evaluator.ResetSteps()

	env := object.NewEnvironment()
	if e, ok := Eval(prog, env).(*object.Error); ok {
		return nil, fmt.Errorf("%s", e.Message)
	}
	h := object.NewHash()
	for _, st := range prog.Statements {
		if let, ok := st.(*ast.LetStatement); ok {
			if v, ok := env.Get(let.Name.Value); ok {
				hashSet(h, let.Name.Value, v)
			}
		}
	}
	return h, nil
}

// settingBool reads a boolean setting; anything but true counts as false.
func settingBool(h *object.Hash, key string) bool {
	if h == nil {
		return false
	}
	b, ok := hashGet(h, key).(*object.Boolean)
	return ok && b.Value
}

// SkipLoadingAnim is the skip_loading_anim setting: `linkz edit` then behaves
// as if it had --no-splash.
func SkipLoadingAnim(h *object.Hash) bool { return settingBool(h, "skip_loading_anim") }

// pressureSettings are the settings that hold a device's pressure curve,
// by the source they shape.
var pressureSettings = map[string]string{
	"trackpad": "trackpad_pressure",
	"tablet":   "tablet_pressure",
}

// ApplyPressureCurves sets the pressure curves from the settings:
// `let trackpad_pressure = {"low": 0.8, "high": 1.0, "gamma": 1.5}`. Keys
// left out keep their default. If any curve is wrong none is applied and the
// error says which setting; the defaults stay.
func ApplyPressureCurves(h *object.Hash) error {
	if h == nil {
		return nil
	}
	next := map[string]pressure.Curve{}
	for source, name := range pressureSettings {
		v := hashGet(h, name)
		if v == nil {
			continue
		}
		c, err := curveFromHash(source, v)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		next[source] = c
	}
	for source, c := range next {
		if err := pressure.SetCurve(source, c); err != nil {
			return err
		}
	}
	return nil
}

// curveFromHash reads {"low", "high", "gamma"} on top of the default curve
// of source.
func curveFromHash(source string, v object.Object) (pressure.Curve, error) {
	c := pressure.Identity
	if source == "trackpad" {
		c = pressure.DefaultTrackpad
	}
	h, ok := v.(*object.Hash)
	if !ok {
		return c, fmt.Errorf("want a hash like {\"low\": 0.8}, got %s", v.Type())
	}
	fields := map[string]*float64{"low": &c.Low, "high": &c.High, "gamma": &c.Gamma}
	var err error
	h.Each(func(pair object.HashPair) {
		key := pair.Key.Inspect()
		field, known := fields[key]
		n, isNumber := number(pair.Value)
		switch {
		case err != nil:
		case !known:
			err = fmt.Errorf("unknown key %q: a curve has low, high and gamma", key)
		case !isNumber:
			err = fmt.Errorf("%s must be a number, got %s", key, pair.Value.Type())
		default:
			*field = n
		}
	})
	if err != nil {
		return c, err
	}
	return c, c.Validate()
}
