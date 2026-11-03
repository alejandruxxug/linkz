package window

import (
	"math"
	"testing"

	"linkz/host/pressure"
)

func TestDragSamples(t *testing.T) {
	native := []pressure.Sample{
		{Phase: pressure.PhaseDown, X: 1, Y: 1},
		{Phase: pressure.PhaseDrag, X: 10, Y: 20, Pressure: 0.3, T: 1},
		{Phase: pressure.PhaseDrag, X: 11, Y: 22, Pressure: 0.5, T: 2},
		{Phase: pressure.PhaseDrag, X: 12, Y: 25, Pressure: 0.7, T: 3},
	}
	got, last, ok := dragSamples(native, 13, 26)
	if !ok || len(got) != 3 || last.Pressure != 0.7 {
		t.Fatalf("got %+v ok=%v", got, ok)
	}
	if got[0].X != 10 || got[1].Y != 22 {
		t.Errorf("middle points changed: %+v", got)
	}
	// The last sample sits on the event's own position, keeping its pressure.
	if got[2].X != 13 || got[2].Y != 26 || got[2].Pressure != 0.7 {
		t.Errorf("last = %+v", got[2])
	}
	if _, _, ok := dragSamples(nil, 0, 0); ok {
		t.Error("no native samples should report none")
	}
	if s, ok := lastPhase(native, pressure.PhaseDown); !ok || s.X != 1 {
		t.Errorf("lastPhase = %+v %v", s, ok)
	}
}

// A native reading goes through its device's curve on its way to an event:
// the default trackpad curve is {0.8, 1, 1.5}, a mouse has none.
func TestEventsUseTheResponseCurve(t *testing.T) {
	trackpad := pressure.Sample{Phase: pressure.PhaseDrag, Pressure: 0.9, Source: "trackpad"}
	want := math.Pow(0.5, 1.5)
	if p := penOf(trackpad, false); math.Abs(p.Pressure-want) > 1e-12 {
		t.Errorf("penOf pressure = %v, want %v", p.Pressure, want)
	}
	got, _, _ := dragSamples([]pressure.Sample{trackpad}, 0, 0)
	if len(got) != 1 || math.Abs(got[0].Pressure-want) > 1e-12 {
		t.Errorf("drag sample = %+v, want pressure %v", got, want)
	}
	mouse := pressure.Sample{Pressure: 0.9, Source: "mouse"}
	if p := penOf(mouse, false); p.Pressure != 0.9 {
		t.Errorf("a mouse changed: %v", p.Pressure)
	}
}
