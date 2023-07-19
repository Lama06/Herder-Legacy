package breakout

import (
	"image/color"
	"math/rand"

	"golang.org/x/image/colornames"
)

type upgradeShuffleUpgrade struct {
	neueUpgrades int
}

var _ upgrade = upgradeShuffleUpgrade{}

func newRandomUpgradeShuffleUpgrade() upgrade {
	return upgradeShuffleUpgrade{
		neueUpgrades: 1 + rand.Intn(2),
	}
}

func (u upgradeShuffleUpgrade) farbe() color.Color {
	return colornames.Lightblue
}

func (u upgradeShuffleUpgrade) radius() float64 {
	return 10 * float64(u.neueUpgrades)
}

func (u upgradeShuffleUpgrade) collect(world *world, plattform *entity) {
	// Upgrades shuffeln
	for stein := range world.entities {
		if !stein.hatSteinComponent || !stein.hatRectComponent {
			continue
		}

		if !stein.steinComponent.hatUpgrade {
			continue
		}

		stein.steinComponent.upgrade = newRandomUpgrade()
		stein.rectComponent.farbe = stein.steinComponent.upgrade.farbe()
	}

	// Neue Upgrades hinzufügen
	for i := 0; i < u.neueUpgrades; i++ {
		var möglicheSteine []*entity
		for stein := range world.entities {
			if !stein.hatSteinComponent || !stein.hatRectComponent {
				continue
			}

			if stein.steinComponent.hatUpgrade {
				continue
			}

			möglicheSteine = append(möglicheSteine, stein)
		}

		if len(möglicheSteine) == 0 {
			break
		}

		randomStein := möglicheSteine[rand.Intn(len(möglicheSteine))]
		randomStein.steinComponent.hatUpgrade = true
		randomStein.steinComponent.upgrade = newRandomUpgrade()
		randomStein.rectComponent.farbe = randomStein.steinComponent.upgrade.farbe()
	}
}
