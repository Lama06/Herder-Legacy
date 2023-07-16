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
	optionen SpielOptionen,
	nächsterScreen func(gewonnen bool) herderlegacy.Screen,
) herderlegacy.Screen {
	return newSpielScreen(herderLegacy, nächsterScreen, optionen)
}
