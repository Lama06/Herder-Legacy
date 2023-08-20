package dame

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
