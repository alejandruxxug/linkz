package window

import (
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
