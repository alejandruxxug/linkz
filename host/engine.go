package host

import (
	"os"

	"lizardscript/ast"
	"lizardscript/evaluator"
	"lizardscript/object"
	"lizardscript/vm"
)

// Engine names the way programs run: the bytecode VM, or the tree-walking
// evaluator. Both give the same results; the VM is several times faster.
type Engine string

const (
	EngineVM   Engine = "vm"
	EngineTree Engine = "tree"
)

// CurrentEngine is the engine every program runs on. It is the VM unless
// the LINKZ_ENGINE environment variable says "tree", the fallback if the VM
// ever misbehaves.
var CurrentEngine = engineFromEnv()

func engineFromEnv() Engine {
	if Engine(os.Getenv("LINKZ_ENGINE")) == EngineTree {
		return EngineTree
	}
	return EngineVM
}

// Eval runs program in env on the current engine, with the contract of
// evaluator.Eval: the program's value, or an *object.Error. It is the only
// place in Linkz that starts a program.
func Eval(program *ast.Program, env *object.Environment) object.Object {
	if CurrentEngine == EngineTree {
		return evaluator.Eval(program, env)
	}
	return vm.Run(program, env)
}
