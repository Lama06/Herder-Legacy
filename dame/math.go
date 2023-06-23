package dame

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func signum(f float64) float64 {
	switch {
	case f == 0:
		return 0
	case f > 0:
		return 1
	case f < 0:
		return -1
	default:
		panic("unreachable")
	}
}
