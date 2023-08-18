package herder

import (
	"math/rand"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/world"
)

func CreateHerder(herderLegacy herderlegacy.HerderLegacy) *world.World {
	rng := rand.New(rand.NewSource(42))

	w := world.NewEmptyWorld()
	createBasisKlassenzimmer(w, levelB001, rng)
	return w
}
