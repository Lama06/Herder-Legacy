package breakout

type velocityComponent struct {
	velocityX, velocityY float64
}

func (w *world) moveWithVelocity() {
	for entity := range w.entities {
		if !entity.hatVelocityComponent {
			continue
		}

		var faktor float64
		if w.zeitUpgradeRemainingTime > 0 {
			faktor = w.zeitUpgradeFaktor
		} else {
			faktor = 1
		}

		entity.position.x += entity.velocityComponent.velocityX * faktor
		entity.position.y += entity.velocityComponent.velocityY * faktor
	}
}
