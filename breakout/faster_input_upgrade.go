package breakout

import (
	"image/color"
	"math/rand"

	"golang.org/x/image/colornames"
)

type fasterInputUpgrade struct {
	duration   int
	multiplier float64
}

var _ upgrade = fasterInputUpgrade{}

func newRandomFasterInputUpgrade() upgrade {
	return fasterInputUpgrade{
		duration:   60 * (6 + rand.Intn(10)),
		multiplier: 1.5 + 1.5*rand.Float64(),
	}
}

func (f fasterInputUpgrade) farbe() color.Color {
	return colornames.Wheat
}

func (f fasterInputUpgrade) radius() float64 {
	return f.multiplier * 10
}

func (f fasterInputUpgrade) collect(world *world, plattform *entity) {
	world.fasterInputUpgradeRemainingTime = f.duration
	world.fasterInputUpgradeMultiplier = f.multiplier
}

func (w *world) tickFasterInputUpgradeTimer() {
	if w.fasterInputUpgradeRemainingTime <= 0 {
		return
	}
	w.fasterInputUpgradeRemainingTime--
}
