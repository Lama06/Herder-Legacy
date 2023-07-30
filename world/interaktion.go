package world

import (
	"github.com/Lama06/Herder-Legacy/aabb"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InteraktionComponent struct {
	OffsetX, OffsetY float64
	Width, Height    float64
	Callback         func()
}

func (w *World) interaktionenHandeln() {
	camera := w.findCamera()
	if camera == nil {
		return
	}

	for entity := range w.Entities {
		if !entity.HatInteraktionComponent {
			continue
		}

		screenX, screenY := calculateScreenPosition(entity, camera)
		hitbox := aabb.Aabb{
			X:      screenX + entity.InteraktionComponent.OffsetX,
			Y:      screenY + entity.InteraktionComponent.OffsetY,
			Width:  entity.InteraktionComponent.Width,
			Height: entity.InteraktionComponent.Height,
		}

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			mouseX, mouseY := ebiten.CursorPosition()
			if hitbox.IsInside(float64(mouseX), float64(mouseY)) {
				entity.InteraktionComponent.Callback()
			}
			return
		}

		if touchIds := inpututil.AppendJustReleasedTouchIDs(nil); len(touchIds) == 1 {
			touchX, touchY := inpututil.TouchPositionInPreviousTick(touchIds[0])
			if hitbox.IsInside(float64(touchX), float64(touchY)) {
				entity.InteraktionComponent.Callback()
			}
			return
		}
	}
}
