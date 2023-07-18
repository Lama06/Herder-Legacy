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
		duration: 60 * (5 + rand.Intn(6)),
	}
}

func (a automaticInputUpgrade) farbe() color.Color {
	return colornames.Aquamarine
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
			if !possibleEntity.affectsAutomaticInput {
				continue
			}
			possibleEntities = append(possibleEntities, possibleEntity)
		}

		sort.Slice(possibleEntities, func(i, j int) bool {
			distanceToI := distance(
				inputEntity.position.x,
				inputEntity.position.y,
				possibleEntities[i].position.x,
				possibleEntities[i].position.y,
			)
			distanceToJ := distance(
				inputEntity.position.x,
				inputEntity.position.y,
				possibleEntities[j].position.x,
				possibleEntities[j].position.y,
			)
			return distanceToI < distanceToJ
		})

		if inputEntity.moveWithInputComponent.x {
			inputEntity.position.x = possibleEntities[0].position.x + inputEntity.moveWithInputComponent.offsetX
		}

		if inputEntity.moveWithInputComponent.y {
			inputEntity.position.y = possibleEntities[0].position.y + inputEntity.moveWithInputComponent.offsetY
		}
	}
}
