package window

// Presence remembers whether the window has focus and whether the cursor is
// over its content, so only changes become events. The zero value assumes a
// focused window with the cursor inside, so a normal start sends nothing.
type Presence struct {
	away    bool // not focused
	outside bool // cursor not over the content (or not focused)
}

// Update reads the current state and returns the events for what changed:
// "focus" when focus changes, then "leave" or "enter" when the cursor leaves
// or enters the content area. The cursor counts as outside while the window
// is not focused, and outside [0, w) x [0, h).
func (p *Presence) Update(focused bool, x, y, w, h int) []Event {
	var evs []Event
	if focused == p.away { // away is the opposite of focused: that is a change
		p.away = !focused
		evs = append(evs, Event{Type: "focus", Focused: focused})
	}
	outside := !focused || x < 0 || y < 0 || x >= w || y >= h
	if outside != p.outside {
		p.outside = outside
		if outside {
			evs = append(evs, Event{Type: "leave"})
		} else {
			evs = append(evs, Event{Type: "enter"})
		}
	}
	return evs
}
