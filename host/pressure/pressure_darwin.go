//go:build darwin

package pressure

import (
	"sync"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

// NSEventType values we care about.
const (
	typeLeftMouseDown    = 1
	typeLeftMouseUp      = 2
	typeLeftMouseDragged = 6
	typeTabletPoint      = 23
	typeTabletProximity  = 24
	typePressure         = 34

	subtypeTabletPoint = 1 // NSEventSubtypeTabletPoint
	deviceEraser       = 3 // NSPointingDeviceTypeEraser
)

// nsPoint is an NSPoint (two CGFloats).
type nsPoint struct{ X, Y float64 }

// sawForce turns true once a pressure event arrives: only a Force Touch
// trackpad sends them, so later plain mouse events are trackpad drags.
var sawForce bool

// monitor is the id returned by addLocalMonitor; kept so it is never freed.
var (
	startOnce sync.Once
	monitor   objc.ID
	handler   objc.Block
)

// Start installs the event monitor. Call it once, on the main thread, after
// the app is running (Ebitengine's first Update is fine). It reports whether
// the monitor was installed.
func Start() bool {
	startOnce.Do(install)
	return monitor != 0
}

func install() {
	if _, err := purego.Dlopen("/System/Library/Frameworks/AppKit.framework/AppKit", purego.RTLD_GLOBAL|purego.RTLD_NOW); err != nil {
		return
	}
	var (
		selType     = objc.RegisterName("type")
		selSubtype  = objc.RegisterName("subtype")
		selPressure = objc.RegisterName("pressure")
		selStage    = objc.RegisterName("stage")
		selTilt     = objc.RegisterName("tilt")
		selRotation = objc.RegisterName("rotation")
		selEntering = objc.RegisterName("isEnteringProximity")
		selDevice   = objc.RegisterName("pointingDeviceType")
	)
	// The block must hand the event back so Ebitengine still gets it.
	handler = objc.NewBlock(func(_ objc.Block, ev objc.ID) objc.ID {
		t := objc.Send[uint64](ev, selType)
		switch t {
		case typeTabletProximity:
			if objc.Send[bool](ev, selEntering) {
				setEraser(objc.Send[uint64](ev, selDevice) == deviceEraser)
			} else {
				update(func(s *State) { *s = State{Eraser: s.Eraser, OK: true} })
			}
		case typeTabletPoint:
			readTablet(ev, selPressure, selTilt, selRotation)
		case typeLeftMouseDown, typeLeftMouseDragged:
			if objc.Send[int16](ev, selSubtype) == subtypeTabletPoint {
				readTablet(ev, selPressure, selTilt, selRotation)
				break
			}
			p := float64(objc.Send[float32](ev, selPressure))
			update(func(s *State) {
				s.Pressure, s.Tilt, s.Rotation = p, [2]float64{}, 0
				s.Source = "mouse"
				if sawForce {
					s.Source = "trackpad"
				}
				if s.Stage == 0 {
					s.Stage = 1
				}
			})
		case typePressure:
			// -stage is only valid on pressure events.
			sawForce = true
			p := float64(objc.Send[float32](ev, selPressure))
			st := int(objc.Send[int64](ev, selStage))
			update(func(s *State) { s.Pressure, s.Stage, s.Source = p, st, "trackpad" })
		case typeLeftMouseUp:
			update(func(s *State) {
				s.Pressure, s.Stage, s.Tilt, s.Rotation = 0, 0, [2]float64{}, 0
			})
		}
		return ev
	})
	mask := uint64(1<<typeLeftMouseDown | 1<<typeLeftMouseUp | 1<<typeLeftMouseDragged |
		1<<typeTabletPoint | 1<<typeTabletProximity | 1<<typePressure)
	monitor = objc.ID(objc.GetClass("NSEvent")).Send(
		objc.RegisterName("addLocalMonitorForEventsMatchingMask:handler:"), mask, handler)
}

// readTablet stores pressure, tilt and rotation from a tablet event.
func readTablet(ev objc.ID, selPressure, selTilt, selRotation objc.SEL) {
	p := float64(objc.Send[float32](ev, selPressure))
	tl := objc.Send[nsPoint](ev, selTilt)
	rot := float64(objc.Send[float32](ev, selRotation))
	update(func(s *State) {
		s.Pressure, s.Tilt, s.Rotation, s.Source = p, [2]float64{tl.X, tl.Y}, rot, "tablet"
		s.Stage = 0
	})
}
