package breakout

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
