package breakout

import "github.com/hajimehoshi/ebiten/v2"

type moveWithInputComponent struct {
	x, y                 bool
	offsetX, offsetY     float64
	maxSpeedX, maxSpeedY float64
}

func cursorPosition() (x, y int, ok bool) {
	touchIds := ebiten.AppendTouchIDs(nil)
	if len(touchIds) == 1 {
		touchX, touchY := ebiten.TouchPosition(touchIds[0])
		return touchX, touchY, true
	}

	mausX, mausY := ebiten.CursorPosition()
	if mausX != 0 || mausY != 0 {
		return mausX, mausY, true
	}

	return 0, 0, false
}

func (w *world) moveWithInput() {
	if w.automaticInputUpgradeRemainingTime > 0 {
		return
	}

	cursorX, cursorY, ok := cursorPosition()
	if !ok {
		return
	}

	var maxSpeedMultiplier float64
	if w.fasterInputUpgradeRemainingTime > 0 {
		maxSpeedMultiplier = w.fasterInputUpgradeMultiplier
	} else {
		maxSpeedMultiplier = 1
	}

	for entity := range w.entities {
		if !entity.hatMovesWithInputComponent {
			continue
		}

		if entity.moveWithInputComponent.x {
			targetX := float64(cursorX) + entity.moveWithInputComponent.offsetX
			speedX := clamp(
				-entity.moveWithInputComponent.maxSpeedX*maxSpeedMultiplier,
				targetX-entity.position.x,
				entity.moveWithInputComponent.maxSpeedX*maxSpeedMultiplier,
			)
			entity.position.x += speedX
		}

		if entity.moveWithInputComponent.y {
			targetY := float64(cursorY) + entity.moveWithInputComponent.offsetY
			speedY := clamp(
				-entity.moveWithInputComponent.maxSpeedY*maxSpeedMultiplier,
				targetY-entity.position.y,
				entity.moveWithInputComponent.maxSpeedY*maxSpeedMultiplier,
			)
			entity.position.y += speedY
		}
	}
}
