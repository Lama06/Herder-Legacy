package dame

type RichtungVertikal byte

const (
	richtungVertikalVorne RichtungVertikal = iota
	richtungVertikalKeine
	richtungVertikalHinten
)

func (r RichtungVertikal) verschiebung(perspektive spieler) int {
	switch r {
	case richtungVertikalKeine:
		return 0
	case richtungVertikalVorne:
		switch perspektive {
		case spielerLehrer:
			return 1
		case spielerSchüler:
			return -1
		default:
			panic("unreachable")
		}
	case richtungVertikalHinten:
		switch perspektive {
		case spielerLehrer:
			return -1
		case spielerSchüler:
			return 1
		default:
			panic("unreachable")
		}
	default:
		panic("unreachable")
	}
}

type RichtungHorizontal byte

const (
	richtungHorizontalLinks RichtungHorizontal = iota
	richtungHorizontalKeine
	richtungHorizontalRechts
)

func (r RichtungHorizontal) verschiebung(perspektive spieler) int {
	switch r {
	case richtungHorizontalKeine:
		return 0
	case richtungHorizontalLinks:
		switch perspektive {
		case spielerLehrer:
			return 1
		case spielerSchüler:
			return -1
		default:
			panic("unreachable")
		}
	case richtungHorizontalRechts:
		switch perspektive {
		case spielerLehrer:
			return -1
		case spielerSchüler:
			return 1
		default:
			panic("unreachable")
		}
	default:
		panic("unreachable")
	}
}

type Richtung struct {
	Horizontal RichtungHorizontal
	Vertikal   RichtungVertikal
}

var (
	RichtungenDiagonalVorne = map[Richtung]struct{}{
		Richtung{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalVorne,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalVorne,
		}: {},
	}
	RichtungenDiagonal = map[Richtung]struct{}{
		Richtung{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalVorne,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalVorne,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalHinten,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalHinten,
		}: {},
	}
	RichtungenSeiteDiagonalUndVorne = map[Richtung]struct{}{
		Richtung{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalKeine,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalKeine,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalVorne,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalVorne,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalKeine,
			Vertikal:   richtungVertikalVorne,
		}: {},
	}
	RichtungenAlle = map[Richtung]struct{}{
		Richtung{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalVorne,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalKeine,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalHinten,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalKeine,
			Vertikal:   richtungVertikalVorne,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalKeine,
			Vertikal:   richtungVertikalKeine,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalKeine,
			Vertikal:   richtungVertikalHinten,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalVorne,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalKeine,
		}: {},
		Richtung{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalHinten,
		}: {},
	}
)
