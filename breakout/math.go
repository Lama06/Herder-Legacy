package breakout

import "math"

func signum(f float64) float64 {
	switch {
	case f > 0:
		return 1
	case f < 0:
		return -1
	default:
		return 0
	}
}

func clamp(min, x, max float64) float64 {
	switch {
	case x < min:
		return min
	case x > max:
		return max
	default:
		return x
	}
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x1-x2, 2) + math.Pow(y1-y2, 2))
}
