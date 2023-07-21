package world_test

import (
	"math"

	"github.com/Lama06/Herder-Legacy/world"
)

func floatsRoughlyEqual(a, b float64) bool {
	const tolerance = 0.00000001
	return math.Abs(a-b) <= tolerance
}

func positionsAreRoughlyEqual(a, b world.Position) bool {
	return floatsRoughlyEqual(a.X, b.X) && floatsRoughlyEqual(a.Y, b.Y)
}
