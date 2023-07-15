package dame

import "github.com/Lama06/Herder-Legacy/herderlegacy"

func NewFreierModusScreen(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func() herderlegacy.Screen,
) herderlegacy.Screen {
	return newMenüScreen(herderLegacy, nächsterScreen)
}

func NewLehrerDameSpielScreen(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func(gewonnen bool) herderlegacy.Screen,
	lehrer lehrer,
) herderlegacy.Screen {
	return newSpielScreen(herderLegacy, nächsterScreen, lehrer)
}
