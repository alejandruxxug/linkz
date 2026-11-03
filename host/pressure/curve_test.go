package pressure

import (
	"math"
	"strings"
	"testing"
)

func TestCurveApply(t *testing.T) {
	for _, tc := range []struct {
		name  string
		curve Curve
		raw   float64
		want  float64
	}{
		{"identity", Identity, 0.37, 0.37},
		{"below low", DefaultTrackpad, 0.5, 0},
		{"at low", DefaultTrackpad, 0.8, 0},
		{"halfway, bent", DefaultTrackpad, 0.9, math.Pow(0.5, 1.5)},
		{"at high", DefaultTrackpad, 1, 1},
		{"above high", Curve{Low: 0, High: 0.5, Gamma: 1}, 0.75, 1},
		{"gamma below 1 raises", Curve{Low: 0, High: 1, Gamma: 0.5}, 0.25, 0.5},
		{"negative raw", Identity, -0.2, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.curve.Apply(tc.raw); math.Abs(got-tc.want) > 1e-12 {
				t.Errorf("Apply(%v) = %v, want %v", tc.raw, got, tc.want)
			}
		})
	}
}

func TestCurveValidate(t *testing.T) {
	for _, tc := range []struct {
		name  string
		curve Curve
		err   string // empty: valid
	}{
		{"default trackpad", DefaultTrackpad, ""},
		{"edges", Curve{Low: 0, High: 1, Gamma: 10}, ""},
		{"negative low", Curve{Low: -0.1, High: 1, Gamma: 1}, "low must be"},
		{"low is 1", Curve{Low: 1, High: 1, Gamma: 1}, "low must be"},
		{"high not above low", Curve{Low: 0.5, High: 0.5, Gamma: 1}, "high must be above low"},
		{"high past 1", Curve{Low: 0, High: 1.2, Gamma: 1}, "high must be above low"},
		{"gamma too small", Curve{Low: 0, High: 1, Gamma: 0.05}, "gamma must be"},
		{"gamma too big", Curve{Low: 0, High: 1, Gamma: 11}, "gamma must be"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.curve.Validate()
			if (err == nil) != (tc.err == "") || err != nil && !strings.Contains(err.Error(), tc.err) {
				t.Errorf("Validate() = %v, want %q", err, tc.err)
			}
		})
	}
}

func TestShapeAndSetCurve(t *testing.T) {
	t.Cleanup(func() {
		SetCurve("trackpad", DefaultTrackpad)
		SetCurve("tablet", Identity)
	})
	if got := Shape("trackpad", 0.5); got != 0 {
		t.Errorf("default trackpad Shape(0.5) = %v, want 0", got)
	}
	if got := Shape("mouse", 0.5); got != 0.5 {
		t.Errorf("a mouse has no curve: Shape = %v, want 0.5", got)
	}
	if err := SetCurve("tablet", Curve{Low: 0.2, High: 1, Gamma: 1}); err != nil {
		t.Fatal(err)
	}
	if got := Shape("tablet", 0.6); math.Abs(got-0.5) > 1e-12 {
		t.Errorf("tablet Shape(0.6) = %v, want 0.5", got)
	}
	if err := SetCurve("tablet", Curve{Low: 0.9, High: 0.1, Gamma: 1}); err == nil {
		t.Error("an invalid curve was accepted")
	}
	if got := Shape("tablet", 0.6); math.Abs(got-0.5) > 1e-12 {
		t.Errorf("a rejected curve changed the tablet: Shape(0.6) = %v", got)
	}
	if err := SetCurve("mouse", Identity); err == nil {
		t.Error("a mouse has no curve to set")
	}
}
