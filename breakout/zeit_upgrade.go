package breakout

import (
	"image/color"
	"math"
	"math/rand"

	"golang.org/x/image/colornames"
)

type zeitUpgrade struct {
	duration int
	faktor   float64
}

var _ upgrade = zeitUpgrade{}

func newRandomZeitUpgrade() upgrade {
	schneller := rand.Float64() > 0.5
	var faktor float64
	if schneller {
		faktor = 1.5 + 0.5*rand.Float64()
	} else {
		faktor = 0.05 + 0.45*rand.Float64()
	}

	return zeitUpgrade{
		duration: 60 * (4 + rand.Intn(3)),
		faktor:   faktor,
	}
}

func (z zeitUpgrade) farbe() color.Color {
	if z.faktor > 1 {
		return colornames.Purple
	} else {
		return colornames.Hotpink
	}
}

func (z zeitUpgrade) radius() float64 {
	return 15 + 10*math.Abs(1-z.faktor)
}

func (z zeitUpgrade) collect(world *world, plattform *entity) {
	world.zeitUpgradeRemainingTime = z.duration
	world.zeitUpgradeFaktor = z.faktor
}

func (w *world) tickZeitUpgrade() {
	if w.zeitUpgradeRemainingTime <= 0 {
		return
	}
	w.zeitUpgradeRemainingTime--
}
