package host

import (
	"fmt"
	"lizardscript/object"
	"math"
	"sort"
)

// Helpers to move values between Go and lizard-script.

func errorf(format string, args ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

// number reads an INTEGER or FLOAT.
func number(o object.Object) (float64, bool) {
	switch v := o.(type) {
	case *object.Integer:
		return float64(v.Value), true
	case *object.Float:
		return v.Value, true
	}
	return 0, false
}

// pointXY reads a POINT with at least two coordinates.
func pointXY(o object.Object) (x, y float64, ok bool) {
	p, isPoint := o.(*object.Point)
	if !isPoint || len(p.Coords) < 2 {
		return 0, 0, false
	}
	x, ok1 := number(p.Coords[0])
	y, ok2 := number(p.Coords[1])
	return x, y, ok1 && ok2
}

func colorOf(o object.Object) (rgba, bool) {
	c, ok := o.(*object.Color)
	if !ok {
		return rgba{}, false
	}
	return rgba{c.R, c.G, c.B, c.A}, true
}

func newInt(n int) object.Object { return &object.Integer{Value: int64(n)} }

func newPoint(x, y int) object.Object {
	return &object.Point{Coords: []object.Object{newInt(x), newInt(y)}}
}

func newFloat(f float64) object.Object { return &object.Float{Value: f} }

func newStr(s string) object.Object { return &object.String{Value: s} }

func newBool(b bool) object.Object {
	return &object.Boolean{Value: b}
}

var null = &object.Null{}

// fromGo builds a lizard value from maps, slices and scalars.
func fromGo(v any) object.Object {
	switch v := v.(type) {
	case nil:
		return null
	case object.Object:
		return v
	case bool:
		return newBool(v)
	case int:
		return newInt(v)
	case int64:
		return &object.Integer{Value: v}
	case float64:
		if v == math.Trunc(v) && math.Abs(v) < 1e15 {
			return &object.Integer{Value: int64(v)}
		}
		return newFloat(v)
	case string:
		return newStr(v)
	case []any:
		arr := &object.Array{}
		for _, e := range v {
			arr.Elements = append(arr.Elements, fromGo(e))
		}
		return arr
	case map[string]any:
		h := object.NewHash()
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			h.Set(&object.String{Value: k}, fromGo(v[k]))
		}
		return h
	}
	return newStr(fmt.Sprint(v))
}

// toGo turns a lizard value into something encoding/json can write.
func toGo(o object.Object, depth int) any {
	if depth > 40 {
		return "..."
	}
	switch v := o.(type) {
	case *object.Integer:
		return v.Value
	case *object.Float:
		return v.Value
	case *object.Boolean:
		return v.Value
	case *object.String:
		return v.Value
	case *object.Null, nil:
		return nil
	case *object.Color:
		return v.Inspect()
	case *object.Point:
		out := make([]any, len(v.Coords))
		for i, c := range v.Coords {
			out[i] = toGo(c, depth+1)
		}
		return out
	case *object.Array:
		out := make([]any, len(v.Elements))
		for i, e := range v.Elements {
			out[i] = toGo(e, depth+1)
		}
		return out
	case *object.Hash:
		out := map[string]any{}
		v.Each(func(p object.HashPair) {
			out[p.Key.Inspect()] = toGo(p.Value, depth+1)
		})
		return out
	}
	return o.Inspect()
}

// hashGet returns h[key], or nil.
func hashGet(h *object.Hash, key string) object.Object {
	if pair, ok := h.Pairs[(&object.String{Value: key}).HashKey()]; ok {
		return pair.Value
	}
	return nil
}

func hashSet(h *object.Hash, key string, v object.Object) {
	h.Set(&object.String{Value: key}, v)
}

// truthy follows lizard-script's rules for conditions.
func truthy(o object.Object) bool {
	switch v := o.(type) {
	case nil, *object.Null:
		return false
	case *object.Boolean:
		return v.Value
	case *object.Integer:
		return v.Value != 0
	case *object.Float:
		return v.Value != 0
	case *object.String:
		return v.Value != ""
	case *object.Array:
		return len(v.Elements) > 0
	case *object.Hash:
		return len(v.Order) > 0
	}
	return true
}

// builtinHash makes a hash of Go functions, so lizard calls them as fb.span().
func builtinHash(name string, fns map[string]object.BuiltinFunction) *object.Hash {
	h := object.NewHash()
	for k, fn := range fns {
		hashSet(h, k, &object.Builtin{Name: name + "." + k, Fn: fn})
	}
	sort.Slice(h.Order, func(i, j int) bool { return h.Pairs[h.Order[i]].Key.Inspect() < h.Pairs[h.Order[j]].Key.Inspect() })
	return h
}
