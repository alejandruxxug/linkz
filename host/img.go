package host

import (
	"fmt"
	"image"
	_ "image/jpeg" // registers the JPEG decoder for image.Decode
	_ "image/png"  // registers the PNG decoder for image.Decode
	"io"
	"os"

	"linkz/lib/draw"
	"linkz/lib/trace"
	"lizardscript/evaluator"
	"lizardscript/object"
)

// The image primitive: only the format and size work of tracing (decode a
// file, shrink it). Finding colors, regions and curves is lizard-script, in
// lib/trace.

const defaultMaxSide = 256

// pixelRGB is an opaque pixel.
type pixelRGB struct{ R, G, B uint8 }

// imgHash is `img`. img.load(path, max_side = 256) decodes a PNG or JPEG,
// shrinks it with a box filter until its longer side is at most max_side
// (never enlarges it), and gives {"size": (w, h), "pixels": [COLOR, ...]}.
// The pixels are one flat array, row by row from the top-left: the pixel at
// column x, row y is pixels[y * w + x]. A flat array is what a lizard loop
// reads fastest (one index, no inner array to fetch). Transparent pixels are
// flattened onto white, so every color is opaque.
func imgHash() *object.Hash {
	return builtinHash("img", map[string]object.BuiltinFunction{
		"load": func(args ...object.Object) object.Object {
			if len(args) < 1 || len(args) > 2 {
				return errorf("wrong number of arguments: want 1 or 2 (path, max_side), got %d", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return errorf("argument 1 must be a STRING, got %s", args[0].Type())
			}
			side := defaultMaxSide
			if len(args) == 2 {
				n, ok := args[1].(*object.Integer)
				if !ok || n.Value < 1 || n.Value > 8192 {
					return errorf("max_side must be a whole number from 1 to 8192")
				}
				side = int(n.Value)
			}
			w, h, px, err := loadImage(path.Value, side)
			if err != nil {
				return errorf("img.load: %v", err)
			}
			pixels := &object.Array{Elements: make([]object.Object, len(px))}
			for i, c := range px {
				pixels.Elements[i] = &object.Color{R: c.R, G: c.G, B: c.B, A: 255}
			}
			out := object.NewHash()
			hashSet(out, "size", newPoint(w, h))
			hashSet(out, "pixels", pixels)
			return out
		},
	})
}

// loadImage decodes the file and returns it at most `side` pixels on its
// longer side, row-major, over white.
func loadImage(path string, side int) (w, h int, px []pixelRGB, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("%s: %w", path, err)
	}
	b := img.Bounds()
	sw, sh := b.Dx(), b.Dy()
	if sw == 0 || sh == 0 {
		return 0, 0, nil, fmt.Errorf("%s: empty image", path)
	}
	w, h = sw, sh
	if longer := max(sw, sh); longer > side {
		w = max(1, (sw*side+longer/2)/longer)
		h = max(1, (sh*side+longer/2)/longer)
	}
	px = make([]pixelRGB, w*h)
	for y := 0; y < h; y++ {
		y0, y1 := y*sh/h, max(y*sh/h+1, ((y+1)*sh+h-1)/h)
		for x := 0; x < w; x++ {
			x0, x1 := x*sw/w, max(x*sw/w+1, ((x+1)*sw+w-1)/w)
			var r, g, bl, n uint64
			for sy := y0; sy < y1 && sy < sh; sy++ {
				for sx := x0; sx < x1 && sx < sw; sx++ {
					cr, cg, cb, ca := img.At(b.Min.X+sx, b.Min.Y+sy).RGBA() // 16 bit, premultiplied
					// Over white: premultiplied color + (1 - alpha) * white.
					r += uint64(cr + (0xffff - ca))
					g += uint64(cg + (0xffff - ca))
					bl += uint64(cb + (0xffff - ca))
					n++
				}
			}
			px[y*w+x] = pixelRGB{uint8(r / n >> 8), uint8(g / n >> 8), uint8(bl / n >> 8)}
		}
	}
	return w, h, px, nil
}

// Trace loads the tracing library and runs trace(path, options) on an image
// file (options: the keys of trace_options, as Go values). It returns the
// generated .lzdraw source and the stats hash as Go values.
func Trace(path string, options map[string]any, out io.Writer) (string, map[string]any, error) {
	evaluator.Output = out
	evaluator.ResetMath()
	env := object.NewEnvironment()
	env.Define("img", imgHash())
	env.Define("time", timeHash())
	for _, name := range draw.Order {
		code, err := draw.Files.ReadFile(name)
		if err != nil {
			return "", nil, err
		}
		if _, err := run(string(code), "draw/"+name, env, nil); err != nil {
			return "", nil, fmt.Errorf("library %s: %w", name, err)
		}
	}
	for _, name := range trace.Order {
		code, err := trace.Files.ReadFile(name)
		if err != nil {
			return "", nil, err
		}
		if _, err := run(string(code), "trace/"+name, env, nil); err != nil {
			return "", nil, fmt.Errorf("library %s: %w", name, err)
		}
	}
	env.Define("trace_input_path", newStr(path))
	env.Define("trace_input_options", fromGo(options))
	res, err := run("trace(trace_input_path, trace_input_options)", "trace", env, nil)
	if err != nil {
		return "", nil, err
	}
	h, ok := res.(*object.Hash)
	if !ok {
		return "", nil, fmt.Errorf("trace gave %s, not a HASH", res.Type())
	}
	code, ok := hashGet(h, "code").(*object.String)
	if !ok {
		return "", nil, fmt.Errorf("trace gave no code")
	}
	stats, _ := toGo(hashGet(h, "stats"), 0).(map[string]any)
	return code.Value, stats, nil
}
