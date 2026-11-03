package pressure

import (
	"fmt"
	"math"
	"sync"
)

// Curve is the response curve of one kind of device. A Force Touch trackpad
// needs a firm press to register a click, so its pressure starts high and
// strokes come out thick; the curve cuts that off and bends the rest.
//
// The pressure the editor sees is clamp((raw - Low) / (High - Low), 0, 1) ** Gamma.
type Curve struct {
	Low, High, Gamma float64
}

// Identity changes nothing.
var Identity = Curve{Low: 0, High: 1, Gamma: 1}

// Limits of a valid curve: 0 <= Low < High <= 1 and minGamma <= Gamma <= maxGamma.
const (
	minGamma = 0.1
	maxGamma = 10
)

// Apply maps a raw pressure from 0 to 1 to the pressure the editor sees.
func (c Curve) Apply(raw float64) float64 {
	p := (raw - c.Low) / (c.High - c.Low)
	p = math.Min(1, math.Max(0, p))
	if c.Gamma == 1 {
		return p
	}
	return math.Pow(p, c.Gamma)
}

// Validate says why the curve can't be used, or returns nil.
func (c Curve) Validate() error {
	switch {
	case !(c.Low >= 0 && c.Low < 1):
		return fmt.Errorf("low must be from 0 up to (not including) 1, got %v", c.Low)
	case !(c.High > c.Low && c.High <= 1):
		return fmt.Errorf("high must be above low (%v) and at most 1, got %v", c.Low, c.High)
	case !(c.Gamma >= minGamma && c.Gamma <= maxGamma):
		return fmt.Errorf("gamma must be from %v to %v, got %v", minGamma, float64(maxGamma), c.Gamma)
	}
	return nil
}

// DefaultTrackpad is the trackpad curve before any setting. A Mac trackpad
// is near 90% of its range by the time a click registers; the numbers are a
// starting point for the user to tune.
var DefaultTrackpad = Curve{Low: 0.8, High: 1, Gamma: 1.5}

var (
	curveMu sync.Mutex
	curves  = map[string]Curve{"trackpad": DefaultTrackpad, "tablet": Identity}
)

// SetCurve sets the curve of "trackpad" or "tablet". A mouse has none.
func SetCurve(source string, c Curve) error {
	if source != "trackpad" && source != "tablet" {
		return fmt.Errorf("no pressure curve for %q", source)
	}
	if err := c.Validate(); err != nil {
		return err
	}
	curveMu.Lock()
	curves[source] = c
	curveMu.Unlock()
	return nil
}

// Shape is the pressure the editor sees for a raw pressure from a source:
// the source's curve applied, or the raw value for a source with none.
func Shape(source string, raw float64) float64 {
	curveMu.Lock()
	c, ok := curves[source]
	curveMu.Unlock()
	if !ok {
		return raw
	}
	return c.Apply(raw)
}
