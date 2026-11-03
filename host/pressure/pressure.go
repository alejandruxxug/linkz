// Package pressure reads Force Touch trackpad pressure.
//
// On macOS it listens to the NSEvents the app already receives, so it works
// next to Ebitengine without touching the window. On other systems there is
// no pressure and Latest always reports ok=false.
package pressure

import (
	"sync"
	"time"
)

// Reading is one pressure reading with the time it arrived.
type Reading struct {
	At       time.Time
	Pressure float64
}

const maxReadings = 64

// State is the latest pen or finger reading.
type State struct {
	Pressure float64    // 0 to 1
	Stage    int        // trackpad force stage: 0 up, 1 click, 2 deep press
	Tilt     [2]float64 // tablet tilt, -1 to 1 on each axis
	Rotation float64    // tablet pen rotation in degrees
	Source   string     // "trackpad", "tablet" or "mouse"
	Eraser   bool       // the tablet's eraser end is in use
	OK       bool       // false until something was read
}

var (
	mu       sync.Mutex
	state    State
	readings []Reading
)

// update changes the state under the lock and logs a sample.
func update(f func(s *State)) {
	mu.Lock()
	defer mu.Unlock()
	f(&state)
	state.OK = true
	readings = append(readings, Reading{time.Now(), state.Pressure})
	if len(readings) > maxReadings {
		readings = append(readings[:0], readings[len(readings)-maxReadings:]...)
	}
}

// setEraser records the pen end from a proximity event (no sample).
func setEraser(e bool) {
	mu.Lock()
	state.Eraser = e
	mu.Unlock()
}

// Latest returns the last reading. OK is false when nothing has been read yet
// or the system has no pressure support.
func Latest() State {
	mu.Lock()
	defer mu.Unlock()
	return state
}

// Readings returns a copy of the most recent readings, oldest first.
func Readings() []Reading {
	mu.Lock()
	defer mu.Unlock()
	return append([]Reading(nil), readings...)
}

// Phase says which part of a gesture a sample belongs to.
const (
	PhaseDown = "down"
	PhaseDrag = "drag"
	PhaseUp   = "up"
	PhasePen  = "pen" // tablet point with no button down (hover)
)

// Sample is one native pointer event. Unlike Latest, nothing is skipped: the
// monitor sees every drag, even when the app only polls once per tick.
// X and Y are logical points in the window's content view, measured from the
// top left. T is seconds since the app started.
type Sample struct {
	T        float64
	X, Y     float64
	Pressure float64
	Tilt     [2]float64
	Rotation float64
	Source   string // "mouse", "trackpad" or "tablet"
	Phase    string
}

// maxQueue bounds the samples kept between two Drain calls.
const maxQueue = 4096

var (
	queue []Sample
	begin = time.Now()
)

// Now is the app clock in seconds, the same clock Sample.T uses.
func Now() float64 { return time.Since(begin).Seconds() }

// record adds a sample, dropping the oldest ones past maxQueue.
func record(s Sample) {
	mu.Lock()
	defer mu.Unlock()
	queue = append(queue, s)
	if len(queue) > maxQueue {
		queue = append(queue[:0], queue[len(queue)-maxQueue:]...)
	}
}

// Drain returns every sample since the last call, oldest first, and clears
// them. It is empty when the system has no native events.
func Drain() []Sample {
	mu.Lock()
	defer mu.Unlock()
	out := queue
	queue = nil
	return out
}

// FlipY converts AppKit's bottom-left y to a y measured from the top of a
// view that is height points tall.
func FlipY(y, height float64) float64 { return height - y }

// snapshot builds a sample from the state the monitor just updated. x and y
// are AppKit's bottom-left origin; height is the content view's height.
func snapshot(t, x, y, height float64, phase string) Sample {
	s := Latest()
	if phase == PhaseUp {
		s.Pressure = 0
	}
	src := s.Source
	if src == "" {
		src = "mouse"
	}
	return Sample{T: t, X: x, Y: FlipY(y, height), Pressure: s.Pressure, Tilt: s.Tilt,
		Rotation: s.Rotation, Source: src, Phase: phase}
}
