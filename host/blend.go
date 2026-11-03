package host

import (
	"lizardscript/object"
	"math"
)

// Layer blending. A layer with a blend mode or an opacity is painted into its
// own transparent buffer (fb.push_layer), then mixed onto what is below
// (fb.pop_layer, or fb.composite for a buffer kept under a name). Lizard
// decides which layers need it; this file does the per-pixel math, the second
// place where Go decides a color (see Architecture, principles).

// blendFn mixes the backdrop channel cb with the source channel cs (both 0 to 1).
type blendFn func(cb, cs float64) float64

var blendModes = map[string]blendFn{
	"normal":   func(cb, cs float64) float64 { return cs },
	"multiply": func(cb, cs float64) float64 { return cb * cs },
	"screen":   func(cb, cs float64) float64 { return cb + cs - cb*cs },
	"overlay": func(cb, cs float64) float64 {
		if cb <= 0.5 {
			return 2 * cb * cs
		}
		return 1 - 2*(1-cb)*(1-cs)
	},
	"darken":  math.Min,
	"lighten": math.Max,
	"add":     func(cb, cs float64) float64 { return math.Min(1, cb+cs) },
}

// compositeBuf mixes src over dst with a blend mode and an opacity (0 to 1).
// Both buffers are premultiplied and the same coordinates. This is the
// W3C compositing formula: the source color becomes
// (1 - ab) * Cs + ab * B(Cb, Cs), then goes source-over.
func compositeBuf(dst, src *Buf, mode blendFn, opacity float64) {
	area := src.Img.Rect.Intersect(dst.Img.Rect).Intersect(dst.clip)
	for y := area.Min.Y; y < area.Max.Y; y++ {
		for x := area.Min.X; x < area.Max.X; x++ {
			s := src.Img.Pix[y*src.Img.Stride+x*4:]
			if s[3] == 0 {
				continue
			}
			d := dst.Img.Pix[y*dst.Img.Stride+x*4:]
			as := float64(s[3]) / 255
			ab := float64(d[3]) / 255
			a := as * opacity
			var out [3]float64
			for i := 0; i < 3; i++ {
				cs := float64(s[i]) / 255 / as
				cb := 0.0
				if ab > 0 {
					cb = float64(d[i]) / 255 / ab
				}
				cr := (1-ab)*cs + ab*mode(cb, math.Min(1, cs))
				out[i] = a*cr + (1-a)*float64(d[i])/255
			}
			for i := 0; i < 3; i++ {
				d[i] = uint8(math.Round(math.Min(1, out[i]) * 255))
			}
			d[3] = uint8(math.Round((a + ab*(1-a)) * 255))
		}
	}
}

// blendArgs reads (mode, opacity).
func blendArgs(mode, opacity object.Object) (blendFn, float64, *object.Error) {
	m, ok := mode.(*object.String)
	op, okOp := number(opacity)
	if !ok || !okOp {
		return nil, 0, errorf("want (STRING mode, number opacity)")
	}
	fn := blendModes[m.Value]
	if fn == nil {
		return nil, 0, errorf("unknown blend mode %q", m.Value)
	}
	return fn, math.Max(0, math.Min(1, op)), nil
}

// pushLayer is fb.push_layer(): drawing goes to a new transparent buffer the
// size of the current target (with its clip) until fb.pop_layer. Layers nest.
func (f *FB) pushLayer(args ...object.Object) object.Object {
	if len(args) != 0 {
		return errorf("wrong number of arguments: want 0, got %d", len(args))
	}
	cur := f.current()
	f.layers = append(f.layers, f.off)
	f.off = NewBuf(cur.W(), cur.H())
	f.off.clip = cur.clip
	return null
}

// popLayer is fb.pop_layer(mode, opacity): goes back to the buffer below the
// last push_layer and mixes the layer onto it.
func (f *FB) popLayer(args ...object.Object) object.Object {
	if len(args) != 2 {
		return errorf("wrong number of arguments: want 2 (mode, opacity), got %d", len(args))
	}
	if len(f.layers) == 0 {
		return errorf("pop_layer without a push_layer")
	}
	mode, opacity, err := blendArgs(args[0], args[1])
	if err != nil {
		return err
	}
	layer := f.off
	f.off = f.layers[len(f.layers)-1]
	f.layers = f.layers[:len(f.layers)-1]
	f.mark()
	compositeBuf(f.current(), layer, mode, opacity)
	return null
}

// composite is fb.composite(name, mode, opacity): mixes a buffer stored with
// fb.keep onto the current target.
func (f *FB) composite(args ...object.Object) object.Object {
	if len(args) != 3 {
		return errorf("wrong number of arguments: want 3 (name, mode, opacity), got %d", len(args))
	}
	name, ok := args[0].(*object.String)
	if !ok {
		return errorf("name must be a STRING")
	}
	mode, opacity, err := blendArgs(args[1], args[2])
	if err != nil {
		return err
	}
	src := f.Kept(name.Value)
	if src == nil {
		return errorf("no kept buffer called %q", name.Value)
	}
	f.mark()
	compositeBuf(f.current(), src, mode, opacity)
	return null
}
