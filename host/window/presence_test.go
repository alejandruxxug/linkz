package window

import (
	"reflect"
	"testing"
)

func TestPresenceUpdate(t *testing.T) {
	type step struct {
		focused bool
		x, y    int
		want    []Event
	}
	focus := func(b bool) Event { return Event{Type: "focus", Focused: b} }
	leave, enter := Event{Type: "leave"}, Event{Type: "enter"}
	steps := []step{
		{true, 10, 10, nil},                         // normal start: nothing
		{true, 20, 30, nil},                         // moving inside: nothing
		{true, 639, 419, nil},                       // last pixel is inside
		{true, 640, 100, []Event{leave}},            // cursor leaves on the right
		{true, 700, 100, nil},                       // still outside
		{true, -1, 100, nil},                        // outside on the other side
		{true, 5, 5, []Event{enter}},                // back in
		{false, 5, 5, []Event{focus(false), leave}}, // focus lost: cursor counts as out
		{false, 6, 6, nil},                          // nothing while unfocused
		{true, 6, 6, []Event{focus(true), enter}},   // focus back, cursor in
		{false, 900, 6, []Event{focus(false), leave}},
		{true, 900, 6, []Event{focus(true)}}, // focus back, cursor still out
		{true, 6, 6, []Event{enter}},
		{true, 6, -1, []Event{leave}}, // above the window
	}
	var p Presence
	for i, s := range steps {
		got := p.Update(s.focused, s.x, s.y, 640, 420)
		if !reflect.DeepEqual(got, s.want) {
			t.Errorf("step %d (%v at %d,%d): got %+v, want %+v", i, s.focused, s.x, s.y, got, s.want)
		}
	}
	// A window that starts unfocused reports it on the first tick.
	var q Presence
	got := q.Update(false, 0, 0, 640, 420)
	if !reflect.DeepEqual(got, []Event{focus(false), leave}) {
		t.Errorf("unfocused start: %+v", got)
	}
}
