package breakout

type velocityComponent struct {
	velocityX, velocityY float64
}

func (w *world) moveWithVelocity() {
	for entity := range w.entities {
		if !entity.hatVelocityComponent {
			continue
		}

		entity.position.x += entity.velocityComponent.velocityX
		entity.position.y += entity.velocityComponent.velocityY
	}
}
