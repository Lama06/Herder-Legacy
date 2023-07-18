package breakout

import (
	"image/color"
	"image/color/palette"
	"math/rand"

	"golang.org/x/image/colornames"
)

var rainbowColors = palette.Plan9

type rainbowUpgrade struct {
	duration        int
	changeFrequency int
}

func newRandomRainbowUpgrade() upgrade {
	return rainbowUpgrade{
		duration:        60 * (5 + rand.Intn(10)),
		changeFrequency: 30 + rand.Intn(200),
	}
}

var _ upgrade = rainbowUpgrade{}

func (r rainbowUpgrade) farbe() color.Color {
	return colornames.Darkorange
}

func (r rainbowUpgrade) radius() float64 {
	return 15
}

func (r rainbowUpgrade) collect(world *world, plattform *entity) {
	world.rainbowUpgradeRemainingTime = r.duration
	world.rainbowUpgradeChangeFrequency = r.changeFrequency
}

type rainbowModeColorChangeComponent struct {
	newColor color.Color
}

func (w *world) changeColorsInRainbowMode() {
	if w.rainbowUpgradeRemainingTime <= 0 {
		for entity := range w.entities {
			entity.hatRainbowModeColorChangeComponent = false
		}
		return
	}
	w.rainbowUpgradeRemainingTime--

	if w.rainbowUpgradeRemainingTime%w.rainbowUpgradeChangeFrequency != 0 {
		return
	}

	for entity := range w.entities {
		if entity.hatSteinComponent && entity.steinComponent.hatUpgrade {
			continue
		}

		neueFarbe := rainbowColors[rand.Intn(len(rainbowColors))]

		entity.hatRainbowModeColorChangeComponent = true
		entity.rainbowModeColorChangeComponent = rainbowModeColorChangeComponent{
			newColor: neueFarbe,
		}
	}
}
