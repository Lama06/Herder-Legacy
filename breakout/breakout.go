package breakout

import (
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/hajimehoshi/ebiten/v2"
)

func NewFreierModusScreen(
	herderLegacy herderlegacy.HerderLegacy,
	breakoutBeendenCallback func() herderlegacy.Screen,
) herderlegacy.Screen {
	return nil
}

type breakoutScreen struct {
	herderLegacy   herderlegacy.HerderLegacy
	nächsterScreen func(gewonnen bool) herderlegacy.Screen
	world          *world
}

var _ herderlegacy.Screen = (*breakoutScreen)(nil)

func NewBreakoutScreen(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func(gewonnen bool) herderlegacy.Screen,
) herderlegacy.Screen {
	return &breakoutScreen{
		herderLegacy:   herderLegacy,
		nächsterScreen: nächsterScreen,
		world:          NewStandardLevel(StandardLevelHardConfig),
	}
}

func (b *breakoutScreen) Update() {
	b.world.update()
}

func (b *breakoutScreen) Draw(screen *ebiten.Image) {
	b.world.draw(screen)
}
