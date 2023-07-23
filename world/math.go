package world

import "math"

func floatsRoughlyEqual(a, b float64) bool {
	const tolerance = 0.00000001
	return math.Abs(a-b) <= tolerance
}
