package breakout

import (
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type renderLayer int

const (
	renderLayerStein renderLayer = iota
	renderLayerKanonenKugel
	renderLayerFallendesUpgrade
	renderLayerPlattform
	renderLayerBall
)

type renderComponent struct {
	layer renderLayer
}

type rectComponent struct {
	width, height float64
	farbe         color.Color
}

type circleComponent struct {
	radius float64
	farbe  color.Color
}

func (w *world) renderObjects(screen *ebiten.Image) {
	var zuRenderndeEntities []*entity
	for entity := range w.entities {
		if !entity.hatRenderComponent {
			continue
		}

		zuRenderndeEntities = append(zuRenderndeEntities, entity)
	}

	sort.Slice(zuRenderndeEntities, func(i, j int) bool {
		return zuRenderndeEntities[i].renderComponent.layer < zuRenderndeEntities[j].renderComponent.layer
	})

	for _, entity := range zuRenderndeEntities {
		if entity.hatRectComponent {
			vector.DrawFilledRect(
				screen,
				float32(entity.position.x),
				float32(entity.position.y),
				float32(entity.rectComponent.width),
				float32(entity.rectComponent.height),
				entity.rectComponent.farbe,
				true,
			)
		}

		if entity.hatCircleComponent {
			vector.DrawFilledCircle(
				screen,
				float32(entity.position.x+entity.circleComponent.radius),
				float32(entity.position.y+entity.circleComponent.radius),
				float32(entity.circleComponent.radius),
				entity.circleComponent.farbe,
				true,
			)
		}
	}
}
