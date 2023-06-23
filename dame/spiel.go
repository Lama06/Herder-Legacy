package dame

import (
	"github.com/Lama06/Herder-Legacy/spiel"
	"github.com/hajimehoshi/ebiten/v2"
)

type dameSpiel struct {
	herderLegacy  spiel.HerderLegacy
	currentScreen screen
}

var _ spiel.Spiel = (*dameSpiel)(nil)

func (d *dameSpiel) Update() (beendet bool) {
	return d.currentScreen.update()
}

func (d *dameSpiel) Draw(screen *ebiten.Image) {
	d.currentScreen.draw(screen)
}

func NewDameSpiel(herderLegacy spiel.HerderLegacy) spiel.Spiel {
	dame := dameSpiel{
		herderLegacy: herderLegacy,
	}

	dame.currentScreen = newMenuScreen(&dame)

	return &dame
}

var _ spiel.Constructor = NewDameSpiel

type screen interface {
	update() (beendet bool)

	draw(screen *ebiten.Image)
}
