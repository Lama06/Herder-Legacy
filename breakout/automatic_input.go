package breakout

import (
	"image/color"
	"math"
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
			if !possibleEntity.affectsAutomaticInput {
				continue
			}
			possibleEntities = append(possibleEntities, possibleEntity)
		}

		if inputEntity.moveWithInputComponent.x {
			sort.Slice(possibleEntities, func(i, j int) bool {
				distanceToI := math.Abs(inputEntity.position.y - possibleEntities[i].position.y)
				distanceToJ := math.Abs(inputEntity.position.y - possibleEntities[j].position.y)
				return distanceToI < distanceToJ
			})
			inputEntity.position.x = possibleEntities[0].position.x + inputEntity.moveWithInputComponent.offsetX
		}

		if inputEntity.moveWithInputComponent.y {
			sort.Slice(possibleEntities, func(i, j int) bool {
				distanceToI := math.Abs(inputEntity.position.x - possibleEntities[i].position.x)
				distanceToJ := math.Abs(inputEntity.position.x - possibleEntities[j].position.x)
				return distanceToI < distanceToJ
			})
			inputEntity.position.y = possibleEntities[0].position.y + inputEntity.moveWithInputComponent.offsetY
		}
	}
}
