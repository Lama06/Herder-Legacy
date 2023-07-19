package breakout

import (
	"image/color"
	"math/rand"
	"sort"

	"golang.org/x/image/colornames"
)

type automaticInputUpgrade struct {
	duration int
}

var _ upgrade = automaticInputUpgrade{}

func newRandomAutomaticInputUpgrade() upgrade {
	return automaticInputUpgrade{
		duration: 60 * (9 + rand.Intn(7)),
	}
}

func (a automaticInputUpgrade) farbe() color.Color {
	return colornames.Lightgreen
}

func (a automaticInputUpgrade) radius() float64 {
	return 15
}

func (a automaticInputUpgrade) collect(world *world, plattform *entity) {
	world.automaticInputUpgradeRemainingTime = a.duration
}

func (w *world) performAutomaticInput() {
	if w.automaticInputUpgradeRemainingTime <= 0 {
		return
	}
	w.automaticInputUpgradeRemainingTime--

	for inputEntity := range w.entities {
		if !inputEntity.hatMovesWithInputComponent {
			continue
		}

		var possibleEntities []*entity
		for possibleEntity := range w.entities {
			if !possibleEntity.affectsAutomaticInput || !possibleEntity.hatVelocityComponent {
				continue
			}
			possibleEntities = append(possibleEntities, possibleEntity)
		}

		if len(possibleEntities) == 0 {
			continue
		}

		if inputEntity.moveWithInputComponent.x {
			sort.Slice(possibleEntities, func(i, j int) bool {
				timeI := (inputEntity.position.y - possibleEntities[i].position.y) / possibleEntities[i].velocityComponent.velocityY
				timeJ := (inputEntity.position.y - possibleEntities[j].position.y) / possibleEntities[j].velocityComponent.velocityY
				if timeI < 0 && timeJ < 0 {
					return timeI < timeJ
				}
				if timeI < 0 {
					return false
				}
				if timeJ < 0 {
					return true
				}
				return timeI < timeJ
			})
			inputEntity.position.x = possibleEntities[0].position.x + inputEntity.moveWithInputComponent.offsetX
		}

		if inputEntity.moveWithInputComponent.y {
			sort.Slice(possibleEntities, func(i, j int) bool {
				timeI := (inputEntity.position.x - possibleEntities[i].position.x) / possibleEntities[i].velocityComponent.velocityX
				timeJ := (inputEntity.position.x - possibleEntities[j].position.x) / possibleEntities[j].velocityComponent.velocityX
				if timeI < 0 && timeJ < 0 {
					return timeI < timeJ
				}
				if timeI < 0 {
					return false
				}
				if timeJ < 0 {
					return true
				}
				return timeI < timeJ
			})
			inputEntity.position.y = possibleEntities[0].position.y + inputEntity.moveWithInputComponent.offsetY
		}
	}
}
