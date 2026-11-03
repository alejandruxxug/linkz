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

// Sample is one pressure reading with the time it arrived.
type Sample struct {
	At       time.Time
	Pressure float64
}

const maxSamples = 64

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
	mu      sync.Mutex
	state   State
	samples []Sample
)

// update changes the state under the lock and logs a sample.
func update(f func(s *State)) {
	mu.Lock()
	defer mu.Unlock()
	f(&state)
	state.OK = true
	samples = append(samples, Sample{time.Now(), state.Pressure})
	if len(samples) > maxSamples {
		samples = append(samples[:0], samples[len(samples)-maxSamples:]...)
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

// Samples returns a copy of the most recent readings, oldest first.
func Samples() []Sample {
	mu.Lock()
	defer mu.Unlock()
	return append([]Sample(nil), samples...)
}
