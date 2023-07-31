package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

func (w *World) debugPasswortLesen() {
	const debugPasswort = "debug"

	if w.debug {
		return
	}

	w.debugPasswort = ebiten.AppendInputChars(w.debugPasswort)
	if len(w.debugPasswort) > len(debugPasswort) {
		w.debugPasswort = w.debugPasswort[len(w.debugPasswort)-len(debugPasswort):]
	}

	if string(w.debugPasswort) == debugPasswort {
		w.debug = true
	}
}

func (w *World) pfadeVisualisierenDebug(screen *ebiten.Image) {
	if !w.debug {
		return
	}

	camera := w.findCamera()
	if camera == nil {
		return
	}

	for entity := range w.Entities {
		if entity.Level != camera.Level {
			continue
		}

		if !entity.HatMoveToPositionsComponent {
			continue
		}

		for i := entity.MoveToPositionsComponent.CurrentPosition; i < len(entity.MoveToPositionsComponent.Positions); i++ {
			screenPositionX, screenPositionY := calculateScreenPosition(entity.MoveToPositionsComponent.Positions[i], camera)
			vector.DrawFilledCircle(
				screen,
				float32(screenPositionX),
				float32(screenPositionY),
				10,
				colornames.Purple,
				true,
			)
		}
	}
}
