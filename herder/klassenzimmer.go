package herder

import (
	"math/rand"

	"github.com/Lama06/Herder-Legacy/world"
)

func createBasisKlassenzimmer(
	w *world.World,
	level world.Level,
	rng *rand.Rand,
) {
	createBoden(w, level, bodenArtParkett, 0, 0, 20, 40)
}
