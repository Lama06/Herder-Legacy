package world

import (
	"math"

	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type KeyboardMovementComponent struct {
	Unten, Oben, Links, Rechts ebiten.Key
	Speed                      float64
}

func (w *World) entitiesMitKeyboardSteuern() {
	for entity := range w.Entities {
		if !entity.HatKeyboardMovementComponent {
			continue
		}

		if ebiten.IsKeyPressed(entity.KeyboardMovementComponent.Oben) {
			entity.Position.Y -= entity.KeyboardMovementComponent.Speed
		}
		if ebiten.IsKeyPressed(entity.KeyboardMovementComponent.Unten) {
			entity.Position.Y += entity.KeyboardMovementComponent.Speed
		}
		if ebiten.IsKeyPressed(entity.KeyboardMovementComponent.Links) {
			entity.Position.X -= entity.KeyboardMovementComponent.Speed
		}
		if ebiten.IsKeyPressed(entity.KeyboardMovementComponent.Rechts) {
			entity.Position.X += entity.KeyboardMovementComponent.Speed
		}
	}
}

type TouchInputComponent struct {
	Speed        float64
	Sensitivität float64
}

func (w *World) entitiesMitTouchSteuern() {
	touchIds := ebiten.AppendTouchIDs(nil)
	if len(touchIds) != 1 {
		return
	}
	touchX, touchY := ebiten.TouchPosition(touchIds[0])

	for entity := range w.Entities {
		if !entity.HatTouchInputComponent {
			continue
		}

		xOffset := float64(touchX) - ui.Width/2
		yOffset := float64(touchY) - ui.Height/2

		xSpeed := xOffset * entity.TouchInputComponent.Sensitivität
		xSpeed = math.Max(-entity.TouchInputComponent.Speed, math.Min(entity.TouchInputComponent.Speed, xSpeed))
		ySpeed := yOffset * entity.TouchInputComponent.Sensitivität
		ySpeed = math.Max(-entity.TouchInputComponent.Speed, math.Min(entity.TouchInputComponent.Speed, ySpeed))

		entity.Position.X += xSpeed
		entity.Position.Y += ySpeed
	}
}
