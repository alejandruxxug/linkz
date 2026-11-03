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
