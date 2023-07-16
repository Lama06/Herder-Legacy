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
		{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalVorne,
		}: {},
		{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalVorne,
		}: {},
	}
	RichtungenDiagonal = map[Richtung]struct{}{
		{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalVorne,
		}: {},
		{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalVorne,
		}: {},
		{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalHinten,
		}: {},
		{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalHinten,
		}: {},
	}
	RichtungenSeiteDiagonalUndVorne = map[Richtung]struct{}{
		{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalKeine,
		}: {},
		{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalKeine,
		}: {},
		{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalVorne,
		}: {},
		{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalVorne,
		}: {},
		{
			Horizontal: richtungHorizontalKeine,
			Vertikal:   richtungVertikalVorne,
		}: {},
	}
	RichtungenAlle = map[Richtung]struct{}{
		{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalVorne,
		}: {},
		{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalKeine,
		}: {},
		{
			Horizontal: richtungHorizontalLinks,
			Vertikal:   richtungVertikalHinten,
		}: {},
		{
			Horizontal: richtungHorizontalKeine,
			Vertikal:   richtungVertikalVorne,
		}: {},
		{
			Horizontal: richtungHorizontalKeine,
			Vertikal:   richtungVertikalKeine,
		}: {},
		{
			Horizontal: richtungHorizontalKeine,
			Vertikal:   richtungVertikalHinten,
		}: {},
		{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalVorne,
		}: {},
		{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalKeine,
		}: {},
		{
			Horizontal: richtungHorizontalRechts,
			Vertikal:   richtungVertikalHinten,
		}: {},
	}
)
