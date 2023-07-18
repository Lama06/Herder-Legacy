package breakout

import (
	"math"

	"github.com/Lama06/Herder-Legacy/aabb"
	"github.com/Lama06/Herder-Legacy/ui"
)

type hitboxComponent struct {
	width, height float64
}

type amRandAbprallenComponent struct {
	oben, unten, links, rechts bool
}

type anHitboxenAbprallenComponent struct {
	minXSpeed, maxXSpeed float64
	minYSpeed, maxYSpeed float64
}

func (e *entity) hitbox() aabb.Aabb {
	if e.hatHitboxComponent {
		return aabb.Aabb{
			X:      e.position.x,
			Y:      e.position.y,
			Width:  e.hitboxComponent.width,
			Height: e.hitboxComponent.height,
		}
	}

	return aabb.Aabb{
		X:      e.position.x,
		Y:      e.position.y,
		Width:  0,
		Height: 0,
	}
}

func (w *world) amRandAbprallen() {
	for entity := range w.entities {
		if !entity.hatAmRandAbprallenComponent {
			continue
		}

		hitbox := entity.hitbox()

		if entity.amRandAbprallenComponent.oben && hitbox.Y < 0 {
			entity.position.y = 0
			if entity.hatVelocityComponent {
				entity.velocityComponent.velocityY *= -1
			}
		}

		if entity.amRandAbprallenComponent.links && hitbox.X < 0 {
			entity.position.x = 0
			if entity.hatVelocityComponent {
				entity.velocityComponent.velocityX *= -1
			}
		}

		if entity.amRandAbprallenComponent.unten && hitbox.Y+hitbox.Height > ui.Height {
			entity.position.y = ui.Height - hitbox.Height
			if entity.hatVelocityComponent {
				entity.velocityComponent.velocityY *= -1
			}
		}

		if entity.amRandAbprallenComponent.rechts && hitbox.X+hitbox.Width > ui.Width {
			entity.position.x = ui.Width - hitbox.Width
			if entity.hatVelocityComponent {
				entity.velocityComponent.velocityX *= -1
			}
		}
	}
}

func sekundärGeschwindigkeitNachAbprallen(
	hitboxCenter float64,
	otherHitboxCenter float64,
	otherHitboxSize float64,
	minSpeed float64,
	maxSpeed float64,
) float64 {
	abweichungVonMitte := hitboxCenter - otherHitboxCenter
	maximaleAbweichung := otherHitboxSize / 2
	abweichungProzent := abweichungVonMitte / maximaleAbweichung
	speed := abweichungProzent * maxSpeed
	if math.Abs(speed) < minSpeed {
		if speed == 0 {
			return minSpeed
		}
		return signum(speed) * minSpeed
	}
	return speed
}

func (w *world) anHitboxenAbprallen() {
	const hitboxTeilProzent = 0.15

	for firstEntity := range w.entities {
		if !firstEntity.hatAnHitboxenAbprallenComponent || !firstEntity.hatVelocityComponent {
			continue
		}

		firstHitbox := firstEntity.hitbox()

		for secondEntity := range w.entities {
			if firstEntity == secondEntity {
				continue
			}

			if !secondEntity.hatHitboxComponent {
				continue
			}

			secondHitbox := secondEntity.hitbox()
			secondHitboxOben := secondHitbox.VonObenZerschneiden(hitboxTeilProzent)
			secondHitboxUnten := secondHitbox.VonUntenZerschneiden(hitboxTeilProzent)
			secondHitboxLinks := secondHitbox.VonLinksZerschneiden(hitboxTeilProzent)
			secondHitboxRechts := secondHitbox.VonRechtsZerschneiden(hitboxTeilProzent)

			if !firstHitbox.KollidiertMit(secondHitbox) {
				continue
			}

			if secondEntity.hatSteinComponent && firstEntity.istBall {
				delete(w.entities, secondEntity)
				fallendesUpgradeSpawnen(w, secondEntity)
			}

			switch {
			case firstHitbox.KollidiertMit(secondHitboxOben) && firstHitbox.KollidiertMit(secondHitboxLinks):
				firstEntity.velocityComponent.velocityX = -firstEntity.anHitboxenAbprallenComponent.maxXSpeed
				firstEntity.velocityComponent.velocityY = -firstEntity.anHitboxenAbprallenComponent.maxYSpeed

				newX := secondHitbox.X - firstHitbox.Width
				xAdjustment := newX - firstHitbox.X

				newY := secondHitbox.Y - firstHitbox.Height
				yAdjustment := newY - firstHitbox.Y

				if math.Abs(xAdjustment) < math.Abs(yAdjustment) {
					firstEntity.position.x += xAdjustment
				} else {
					firstEntity.position.y += yAdjustment
				}
			case firstHitbox.KollidiertMit(secondHitboxOben) && firstHitbox.KollidiertMit(secondHitboxRechts):
				firstEntity.velocityComponent.velocityX = firstEntity.anHitboxenAbprallenComponent.maxXSpeed
				firstEntity.velocityComponent.velocityY = -firstEntity.anHitboxenAbprallenComponent.maxYSpeed

				newX := secondHitbox.MaxX()
				xAdjustment := newX - firstHitbox.X

				newY := secondHitbox.Y - firstHitbox.Height
				yAdjustment := newY - firstHitbox.Y

				if math.Abs(xAdjustment) < math.Abs(yAdjustment) {
					firstEntity.position.x += xAdjustment
				} else {
					firstEntity.position.y += yAdjustment
				}
			case firstHitbox.KollidiertMit(secondHitboxUnten) && firstHitbox.KollidiertMit(secondHitboxLinks):
				firstEntity.velocityComponent.velocityX = -firstEntity.anHitboxenAbprallenComponent.maxXSpeed
				firstEntity.velocityComponent.velocityY = firstEntity.anHitboxenAbprallenComponent.maxYSpeed

				newX := secondHitbox.X - firstHitbox.Width
				xAdjustment := newX - firstHitbox.X

				newY := secondHitbox.MaxY()
				yAdjustment := newY - firstHitbox.Y

				if math.Abs(xAdjustment) < math.Abs(yAdjustment) {
					firstEntity.position.x += xAdjustment
				} else {
					firstEntity.position.y += yAdjustment
				}
			case firstHitbox.KollidiertMit(secondHitboxUnten) && firstHitbox.KollidiertMit(secondHitboxRechts):
				firstEntity.velocityComponent.velocityX = firstEntity.anHitboxenAbprallenComponent.maxXSpeed
				firstEntity.velocityComponent.velocityY = firstEntity.anHitboxenAbprallenComponent.maxYSpeed

				newX := secondHitbox.MaxX()
				xAdjustment := newX - firstHitbox.X

				newY := secondHitbox.MaxY()
				yAdjustment := newY - firstHitbox.Y

				if math.Abs(xAdjustment) < math.Abs(yAdjustment) {
					firstEntity.position.x += xAdjustment
				} else {
					firstEntity.position.y += yAdjustment
				}
			case firstHitbox.KollidiertMit(secondHitboxOben) || firstHitbox.KollidiertMit(secondHitboxUnten):
				firstEntity.velocityComponent.velocityY = -signum(firstEntity.velocityComponent.velocityY) * firstEntity.anHitboxenAbprallenComponent.maxYSpeed
				firstEntity.velocityComponent.velocityX = sekundärGeschwindigkeitNachAbprallen(
					firstHitbox.CenterX(),
					secondHitbox.CenterX(),
					secondHitbox.Width,
					firstEntity.anHitboxenAbprallenComponent.minXSpeed,
					firstEntity.anHitboxenAbprallenComponent.maxXSpeed,
				)

				if firstHitbox.KollidiertMit(secondHitboxOben) {
					firstEntity.position.y = secondEntity.position.y - firstHitbox.Height
				}

				if firstHitbox.KollidiertMit(secondHitboxUnten) {
					firstEntity.position.y = secondEntity.position.y + secondEntity.hitboxComponent.height
				}
			case firstHitbox.KollidiertMit(secondHitboxLinks) || firstHitbox.KollidiertMit(secondHitboxRechts):
				firstEntity.velocityComponent.velocityX = -signum(firstEntity.velocityComponent.velocityX) * firstEntity.anHitboxenAbprallenComponent.maxXSpeed
				firstEntity.velocityComponent.velocityY = sekundärGeschwindigkeitNachAbprallen(
					firstHitbox.CenterY(),
					secondHitbox.CenterY(),
					secondHitbox.Height,
					firstEntity.anHitboxenAbprallenComponent.minYSpeed,
					firstEntity.anHitboxenAbprallenComponent.maxYSpeed,
				)

				if firstHitbox.KollidiertMit(secondHitboxLinks) {
					firstEntity.position.x = secondEntity.position.x - firstHitbox.Width
				}

				if firstHitbox.KollidiertMit(secondHitboxRechts) {
					firstEntity.position.x = secondEntity.position.x + secondEntity.hitboxComponent.width
				}
			}
		}
	}
}
