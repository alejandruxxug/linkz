package host

import (
	"lizardscript/object"
	"math"
)

// Gradients. Lizard works out what a gradient is (see lib/draw/gradient.liz:
// where it starts, which way it runs, its stops and its steps); this file only
// evaluates it for each pixel along a span, like a shader. It is the one place
// where Go decides a color, because doing it per pixel in lizard is too slow.

type gradStop struct {
	t float64
	c rgba
}

// grad is a gradient ready to evaluate. A linear one projects a pixel onto
// the axis; a radial one measures its distance from the origin.
type grad struct {
	radial bool
	ox, oy float64 // linear: where t is 0; radial: the center
	ax, ay float64 // linear: the axis divided by its squared length
	scale  float64 // radial: 1 / radius
	stops  []gradStop
	steps  int // 0 or 1: smooth; n >= 2: n flat bands
}

// parseGrad reads the hash lizard built: {"kind": "linear" or "radial",
// "origin": point, "axis": point, "scale": number, "stops": [[t, color], ...],
// "steps": int}.
func parseGrad(o object.Object) (*grad, bool) {
	h, ok := o.(*object.Hash)
	if !ok {
		return nil, false
	}
	g := &grad{}
	kind, _ := hashGet(h, "kind").(*object.String)
	if kind == nil || (kind.Value != "linear" && kind.Value != "radial") {
		return nil, false
	}
	g.radial = kind.Value == "radial"
	var ok1, ok2, ok3 bool
	g.ox, g.oy, ok1 = pointXY(hashGet(h, "origin"))
	if g.radial {
		g.scale, ok2 = number(hashGet(h, "scale"))
		ok3 = true
	} else {
		g.ax, g.ay, ok3 = pointXY(hashGet(h, "axis"))
		ok2 = true
	}
	if !ok1 || !ok2 || !ok3 {
		return nil, false
	}
	list, ok := hashGet(h, "stops").(*object.Array)
	if !ok || len(list.Elements) == 0 {
		return nil, false
	}
	for _, item := range list.Elements {
		pair, ok := item.(*object.Array)
		if !ok || len(pair.Elements) < 2 {
			return nil, false
		}
		t, ok1 := number(pair.Elements[0])
		c, ok2 := colorOf(pair.Elements[1])
		if !ok1 || !ok2 {
			return nil, false
		}
		g.stops = append(g.stops, gradStop{t, c})
	}
	if s, ok := number(hashGet(h, "steps")); ok && s >= 2 {
		g.steps = int(s)
	}
	return g, true
}

// colorAt is the color of the gradient at t (0 to 1, already clamped).
func (g *grad) colorAt(t float64) rgba {
	st := g.stops
	if t <= st[0].t {
		return st[0].c
	}
	j := 0
	for j < len(st) && st[j].t <= t {
		j++
	}
	if j == len(st) {
		return st[len(st)-1].c
	}
	a, b := st[j-1], st[j]
	f := (t - a.t) / (b.t - a.t)
	mix := func(p, q uint8) uint8 {
		return uint8(math.Round(float64(p) + (float64(q)-float64(p))*f))
	}
	return rgba{mix(a.c.r, b.c.r), mix(a.c.g, b.c.g), mix(a.c.b, b.c.b), mix(a.c.a, b.c.a)}
}

// at is the color at the point (x, y) of the gradient's own space.
func (g *grad) at(x, y float64) rgba {
	var t float64
	if g.radial {
		t = math.Hypot(x-g.ox, y-g.oy) * g.scale
	} else {
		t = (x-g.ox)*g.ax + (y-g.oy)*g.ay
	}
	t = math.Max(0, math.Min(1, t))
	if g.steps >= 2 {
		band := math.Min(math.Floor(t*float64(g.steps)), float64(g.steps-1))
		t = band / float64(g.steps-1)
	}
	return g.colorAt(t)
}

// GradSpan fills the pixels [x0, x1) of row y like Span, but each pixel takes
// its color from g, evaluated at the pixel's middle minus (dx, dy). The ends
// are soft in the same way.
func (b *Buf) GradSpan(y int, x0, x1 float64, g *grad, dx, dy float64) {
	if y < b.clip.Min.Y || y >= b.clip.Max.Y {
		return
	}
	x0 = math.Max(x0, float64(b.clip.Min.X))
	x1 = math.Min(x1, float64(b.clip.Max.X))
	if !(x1 > x0) {
		return
	}
	first := int(math.Floor(x0))
	last := int(math.Ceil(x1)) - 1
	row := b.Img.Pix[y*b.Img.Stride:]
	py := float64(y) + 0.5 - dy
	for x := first; x <= last; x++ {
		cov := math.Min(float64(x+1), x1) - math.Max(float64(x), x0)
		blend(row[x*4:x*4+4], g.at(float64(x)+0.5-dx, py), cov)
	}
}

// spanGradient is fb.span_gradient(y, x0, x1, grad).
func (f *FB) spanGradient(args ...object.Object) object.Object {
	if len(args) != 4 {
		return errorf("wrong number of arguments: want 4 (y, x0, x1, gradient), got %d", len(args))
	}
	y, ok1 := number(args[0])
	x0, ok2 := number(args[1])
	x1, ok3 := number(args[2])
	g, ok4 := parseGrad(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return errorf("want (number y, number x0, number x1, gradient)")
	}
	f.mark()
	f.current().GradSpan(int(math.Floor(y)), x0, x1, g, 0, 0)
	return null
}

// spansGradient is fb.spans_gradient(list, offset, grad): many [y, x0, x1]
// runs at an offset. The gradient lives in the space of the runs, before the
// offset, so a moved shape keeps its gradient.
func (f *FB) spansGradient(args ...object.Object) object.Object {
	if len(args) != 3 {
		return errorf("wrong number of arguments: want 3 (list, offset, gradient), got %d", len(args))
	}
	list, ok1 := args[0].(*object.Array)
	dx, dy, ok2 := pointXY(args[1])
	g, ok3 := parseGrad(args[2])
	if !ok1 || !ok2 || !ok3 {
		return errorf("want (ARRAY of [y, x0, x1], POINT offset, gradient)")
	}
	f.mark()
	b := f.current()
	for _, item := range list.Elements {
		run, ok := item.(*object.Array)
		if !ok || len(run.Elements) != 3 {
			return errorf("each span must be [y, x0, x1]")
		}
		y, ok1 := number(run.Elements[0])
		x0, ok2 := number(run.Elements[1])
		x1, ok3 := number(run.Elements[2])
		if !ok1 || !ok2 || !ok3 {
			return errorf("each span must be [y, x0, x1] of numbers")
		}
		b.GradSpan(int(math.Floor(y+dy)), x0+dx, x1+dx, g, dx, dy)
	}
	return null
}
