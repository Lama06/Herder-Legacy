package breakout

import (
	"image/color"

	"golang.org/x/image/colornames"
)

type hiddenUpgrade struct {
	upgrade upgrade
}

func newRandomHiddenUpgrade() upgrade {
	return hiddenUpgrade{
		upgrade: newRandomUpgrade(),
	}
}

var _ upgrade = hiddenUpgrade{}

func (h hiddenUpgrade) farbe() color.Color {
	return colornames.White
}

func (h hiddenUpgrade) radius() float64 {
	return h.upgrade.radius()
}

func (h hiddenUpgrade) collect(world *world, plattform *entity) {
	h.upgrade.collect(world, plattform)
}
