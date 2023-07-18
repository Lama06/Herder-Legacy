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
	renderLayerFallendesUpgrade
	renderLayerPlattform
	renderLayerKanonenKugel
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
			var farbe color.Color
			if entity.hatRainbowModeColorChangeComponent {
				farbe = entity.rainbowModeColorChangeComponent.newColor
			} else {
				farbe = entity.rectComponent.farbe
			}

			vector.DrawFilledRect(
				screen,
				float32(entity.position.x),
				float32(entity.position.y),
				float32(entity.rectComponent.width),
				float32(entity.rectComponent.height),
				farbe,
				true,
			)
		}

		if entity.hatCircleComponent {
			var farbe color.Color
			if entity.hatRainbowModeColorChangeComponent {
				farbe = entity.rainbowModeColorChangeComponent.newColor
			} else {
				farbe = entity.circleComponent.farbe
			}

			vector.DrawFilledCircle(
				screen,
				float32(entity.position.x+entity.circleComponent.radius),
				float32(entity.position.y+entity.circleComponent.radius),
				float32(entity.circleComponent.radius),
				farbe,
				true,
			)
		}
	}
}
