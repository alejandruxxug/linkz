package host

import (
	"fmt"
	"image"

	"linkz/lib"
	"lizardscript/evaluator"
	"lizardscript/lexer"
	"lizardscript/object"
	"lizardscript/parser"
	"os"
)

// SplashScale is how many window pixels one banner pixel takes.
const SplashScale = 3

// Splash holds the startup banner's pieces, drawn by lib/ui/splash.liz into
// offscreen buffers. Go only stacks and scales them: it never decides what a
// piece looks like.
type Splash struct {
	fb  *FB
	env *object.Environment

	W, H       int // banner size in banner pixels
	FPS        int
	LizardAt   image.Point // where the lizard walks in place
	Lizards    int
	Vines      int // vine stages; 0 if the art has none
	Bloom      int // bloom frames; 0 if the art has none
	logo       *Buf
	lizard     []*Buf
	vine       []*Buf
	bloom      []*Buf
	errorLayer *Buf
}

// splashLibs are loaded, in this order, into the splash's own environment.
var splashLibs = []string{"gfx/font_data.liz", "gfx/font.liz", "ui/splash_art.liz", "ui/splash.liz"}

// NewSplash loads the splash libraries and has lizard render every piece.
func NewSplash() (*Splash, error) {
	s := &Splash{fb: NewFB(2, 2), env: object.NewEnvironment()}
	evaluator.Output = os.Stderr
	s.env.Define("fb", s.fb.Builtins())
	for _, name := range splashLibs {
		code, err := lib.Files.ReadFile(name)
		if err != nil {
			return nil, err
		}
		if _, err := s.run(name, string(code)); err != nil {
			return nil, err
		}
	}
	info, err := s.run("splash_render()", "splash_render()")
	if err != nil {
		return nil, err
	}
	h, ok := info.(*object.Hash)
	if !ok {
		return nil, fmt.Errorf("splash_render must give a hash")
	}
	size, _ := hashGet(h, "size").(*object.Point)
	if size == nil {
		return nil, fmt.Errorf("splash_render: no size")
	}
	s.W, s.H, _ = pointInts(size)
	s.FPS = hashInt(h, "fps")
	if at, ok := hashGet(h, "lizard_at").(*object.Point); ok {
		x, y, _ := pointInts(at)
		s.LizardAt = image.Pt(x, y)
	}
	s.Lizards, s.Vines, s.Bloom = hashInt(h, "lizards"), hashInt(h, "vines"), hashInt(h, "bloom")
	if s.W <= 0 || s.H <= 0 || s.FPS <= 0 || s.Lizards <= 0 {
		return nil, fmt.Errorf("splash_render gave a bad info hash")
	}
	s.logo = s.fb.Kept("logo")
	s.lizard = s.kept("lizard", s.Lizards)
	s.vine = s.kept("vine", s.Vines)
	s.bloom = s.kept("bloom", s.Bloom)
	return s, nil
}

func (s *Splash) kept(prefix string, n int) []*Buf {
	var out []*Buf
	for i := 0; i < n; i++ {
		out = append(out, s.fb.Kept(fmt.Sprintf("%s%d", prefix, i)))
	}
	return out
}

func (s *Splash) run(name, code string) (object.Object, error) {
	p := parser.New(lexer.New(code))
	prog := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		return nil, fmt.Errorf("%s: %s", name, errs[0])
	}
	res := evaluator.Eval(prog, s.env)
	if e, ok := res.(*object.Error); ok {
		return nil, fmt.Errorf("%s: %s", name, e.Message)
	}
	return res, nil
}

// ShowError has lizard draw msg under the tagline; Frame then includes it.
func (s *Splash) ShowError(msg string) {
	s.env.Define("__msg", newStr(msg))
	if _, err := s.run("splash_error", "splash_error(__msg)"); err != nil {
		fmt.Fprintln(os.Stderr, "splash:", err)
		return
	}
	s.errorLayer = s.fb.Kept("error")
}

// WindowSize is the splash window's size in window pixels.
func (s *Splash) WindowSize() (int, int) { return s.W * SplashScale, s.H * SplashScale }

// Frame stacks the pieces for animation frame n and scales them to the window
// size: lizard (walking in place), logo and tagline, vine stage, bloom frame, error line. A
// stage or bloom index out of range counts as none (-1).
func (s *Splash) Frame(n, vine, bloom int, withError bool) *image.RGBA {
	banner := image.NewRGBA(image.Rect(0, 0, s.W, s.H))
	overlay(banner, s.lizard[n%s.Lizards], s.LizardAt)
	overlay(banner, s.logo, image.Point{})
	if vine >= 0 && vine < len(s.vine) {
		overlay(banner, s.vine[vine], image.Point{})
	}
	if bloom >= 0 && bloom < len(s.bloom) {
		overlay(banner, s.bloom[bloom], image.Point{})
	}
	if withError && s.errorLayer != nil {
		overlay(banner, s.errorLayer, image.Point{})
	}
	return scaleUp(banner, SplashScale)
}

// overlay draws src over dst with its top left at at (source-over, both
// premultiplied). Whatever falls outside dst is cut.
func overlay(dst *image.RGBA, src *Buf, at image.Point) {
	if src == nil {
		return
	}
	r := src.Img.Rect.Add(at).Intersect(dst.Rect)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			s := src.Img.Pix[(y-at.Y)*src.Img.Stride+(x-at.X)*4:]
			if s[3] == 0 {
				continue
			}
			d := dst.Pix[y*dst.Stride+x*4:]
			inv := 255 - int(s[3])
			for i := 0; i < 4; i++ {
				d[i] = uint8(min(255, int(s[i])+(int(d[i])*inv+127)/255))
			}
		}
	}
}

// scaleUp repeats every pixel k by k times.
func scaleUp(src *image.RGBA, k int) *image.RGBA {
	w, h := src.Rect.Dx(), src.Rect.Dy()
	out := image.NewRGBA(image.Rect(0, 0, w*k, h*k))
	for y := 0; y < h*k; y++ {
		for x := 0; x < w*k; x++ {
			copy(out.Pix[y*out.Stride+x*4:y*out.Stride+x*4+4], src.Pix[(y/k)*src.Stride+(x/k)*4:])
		}
	}
	return out
}

func hashInt(h *object.Hash, key string) int {
	if n, ok := hashGet(h, key).(*object.Integer); ok {
		return int(n.Value)
	}
	return 0
}

func pointInts(p *object.Point) (int, int, bool) {
	x, y, ok := pointXY(p)
	return int(x), int(y), ok
}
