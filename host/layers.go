package host

import (
	"lizardscript/evaluator"
	"lizardscript/object"
	"lizardscript/token"
)

// trackLayers makes here() report the layer path of the running call, for
// hosts that have no editor (render, paint): blend() and the SVG groups need
// it. It returns a function that puts the old hook back. The editor has its
// own, in Src.Run, which also tracks lines.
func trackLayers(env *object.Environment) (restore func()) {
	prev := evaluator.OnLayer
	var layers []string
	evaluator.OnLayer = func(name string, enter bool, tok token.Token) {
		if enter {
			layers = append(layers, name)
		} else if len(layers) > 0 {
			layers = layers[:len(layers)-1]
		}
	}
	env.Define("here", &object.Builtin{Name: "here", Fn: func(args ...object.Object) object.Object {
		path := &object.Array{}
		for _, l := range layers {
			path.Elements = append(path.Elements, newStr(l))
		}
		h := object.NewHash()
		hashSet(h, "line", null)
		hashSet(h, "layer", path)
		return h
	}})
	return func() { evaluator.OnLayer = prev }
}
