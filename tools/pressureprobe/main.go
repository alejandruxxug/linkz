// pressureprobe shows live trackpad and pen pressure. It is a throwaway dev tool:
// Ebitengine draws here ONLY for that reason. Linkz itself never draws
// with Ebitengine.
package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"linkz/host/pressure"
)

type game struct {
	started bool
	ok      bool
	last    pressure.State
}

func (g *game) Update() error {
	if !g.started {
		g.started = true
		g.ok = pressure.Start()
		fmt.Println("monitor installed:", g.ok)
		if !g.ok {
			fmt.Println("pressure monitor unavailable")
		}
	}
	s := pressure.Latest()
	if s != g.last {
		g.last = s
		if s.Pressure > 0 {
			fmt.Printf("pressure=%.2f stage=%d source=%s tilt=%.2f,%.2f rotation=%.0f eraser=%v\n",
				s.Pressure, s.Stage, s.Source, s.Tilt[0], s.Tilt[1], s.Rotation, s.Eraser)
		}
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	s := pressure.Latest()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("monitor=%v source=%s pressure=%.2f stage=%d\ntilt=%.2f,%.2f rot=%.0f eraser=%v",
		g.ok, s.Source, s.Pressure, s.Stage, s.Tilt[0], s.Tilt[1], s.Rotation, s.Eraser))
	vector.DrawFilledRect(screen, 20, 60, 400, 30, color.Gray{40}, false)
	vector.DrawFilledRect(screen, 20, 60, float32(400*s.Pressure), 30, color.RGBA{230, 120, 40, 255}, false)
}

func (g *game) Layout(w, h int) (int, int) { return 440, 120 }

func main() {
	ebiten.SetWindowSize(880, 240)
	ebiten.SetWindowTitle("pressure probe")
	if err := ebiten.RunGame(&game{}); err != nil {
		log.Fatal(err)
	}
}
