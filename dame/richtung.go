package dame

type richtungVertikal byte

const (
	richtungVertikalVorne richtungVertikal = iota
	richtungVertikalKeine
	richtungVertikalHinten
)

func (r richtungVertikal) verschiebung(perspektive spieler) int {
	switch r {
	case richtungVertikalKeine:
		return 0
	case richtungVertikalVorne:
		switch perspektive {
		case spielerLehrer:
			return 1
		case spielerSchueler:
			return -1
		default:
			panic("unreachable")
		}
	case richtungVertikalHinten:
		switch perspektive {
		case spielerLehrer:
			return -1
		case spielerSchueler:
			return 1
		default:
			panic("unreachable")
		}
	default:
		panic("unreachable")
	}
}

type richtungHorizontal byte

const (
	richtungHorizontalLinks richtungHorizontal = iota
	richtungHorizontalKeine
	richtungHorizontalRechts
)

func (r richtungHorizontal) verschiebung(perspektive spieler) int {
	switch r {
	case richtungHorizontalKeine:
		return 0
	case richtungHorizontalLinks:
		switch perspektive {
		case spielerLehrer:
			return 1
		case spielerSchueler:
			return -1
		default:
			panic("unreachable")
		}
	case richtungHorizontalRechts:
		switch perspektive {
		case spielerLehrer:
			return -1
		case spielerSchueler:
			return 1
		default:
			panic("unreachable")
		}
	default:
		panic("unreachable")
	}
}

type richtung struct {
	horizontal richtungHorizontal
	vertikal   richtungVertikal
}

var (
	richtungenDiagonalVorne = map[richtung]struct{}{
		richtung{
			horizontal: richtungHorizontalLinks,
			vertikal:   richtungVertikalVorne,
		}: {},
		richtung{
			horizontal: richtungHorizontalRechts,
			vertikal:   richtungVertikalVorne,
		}: {},
	}
	richtungenDiagonal = map[richtung]struct{}{
		richtung{
			horizontal: richtungHorizontalLinks,
			vertikal:   richtungVertikalVorne,
		}: {},
		richtung{
			horizontal: richtungHorizontalRechts,
			vertikal:   richtungVertikalVorne,
		}: {},
		richtung{
			horizontal: richtungHorizontalLinks,
			vertikal:   richtungVertikalHinten,
		}: {},
		richtung{
			horizontal: richtungHorizontalRechts,
			vertikal:   richtungVertikalHinten,
		}: {},
	}
	richtungenSeiteDiagonalUndVorne = map[richtung]struct{}{
		richtung{
			horizontal: richtungHorizontalLinks,
			vertikal:   richtungVertikalKeine,
		}: {},
		richtung{
			horizontal: richtungHorizontalRechts,
			vertikal:   richtungVertikalKeine,
		}: {},
		richtung{
			horizontal: richtungHorizontalLinks,
			vertikal:   richtungVertikalVorne,
		}: {},
		richtung{
			horizontal: richtungHorizontalRechts,
			vertikal:   richtungVertikalVorne,
		}: {},
		richtung{
			horizontal: richtungHorizontalKeine,
			vertikal:   richtungVertikalVorne,
		}: {},
	}
	richtungenAlle = map[richtung]struct{}{
		richtung{
			horizontal: richtungHorizontalLinks,
			vertikal:   richtungVertikalVorne,
		}: {},
		richtung{
			horizontal: richtungHorizontalLinks,
			vertikal:   richtungVertikalKeine,
		}: {},
		richtung{
			horizontal: richtungHorizontalLinks,
			vertikal:   richtungVertikalHinten,
		}: {},
		richtung{
			horizontal: richtungHorizontalKeine,
			vertikal:   richtungVertikalVorne,
		}: {},
		richtung{
			horizontal: richtungHorizontalKeine,
			vertikal:   richtungVertikalKeine,
		}: {},
		richtung{
			horizontal: richtungHorizontalKeine,
			vertikal:   richtungVertikalHinten,
		}: {},
		richtung{
			horizontal: richtungHorizontalRechts,
			vertikal:   richtungVertikalVorne,
		}: {},
		richtung{
			horizontal: richtungHorizontalRechts,
			vertikal:   richtungVertikalKeine,
		}: {},
		richtung{
			horizontal: richtungHorizontalRechts,
			vertikal:   richtungVertikalHinten,
		}: {},
	}
)
