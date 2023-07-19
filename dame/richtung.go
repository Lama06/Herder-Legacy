package dame

type RichtungVertikal byte

const (
	RichtungVertikalVorne RichtungVertikal = iota
	RichtungVertikalKeine
	RichtungVertikalHinten
)

func (r RichtungVertikal) verschiebung(perspektive spieler) int {
	switch r {
	case RichtungVertikalKeine:
		return 0
	case RichtungVertikalVorne:
		switch perspektive {
		case spielerLehrer:
			return 1
		case spielerSchüler:
			return -1
		default:
			panic("unreachable")
		}
	case RichtungVertikalHinten:
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
	RichtungHorizontalLinks RichtungHorizontal = iota
	RichtungHorizontalKeine
	RichtungHorizontalRechts
)

func (r RichtungHorizontal) verschiebung(perspektive spieler) int {
	switch r {
	case RichtungHorizontalKeine:
		return 0
	case RichtungHorizontalLinks:
		switch perspektive {
		case spielerLehrer:
			return 1
		case spielerSchüler:
			return -1
		default:
			panic("unreachable")
		}
	case RichtungHorizontalRechts:
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
	RichtungVorneLinks = Richtung{
		Horizontal: RichtungHorizontalLinks,
		Vertikal:   RichtungVertikalVorne,
	}
	RichtungVorne = Richtung{
		Horizontal: RichtungHorizontalKeine,
		Vertikal:   RichtungVertikalVorne,
	}
	RichtungVorneRechts = Richtung{
		Horizontal: RichtungHorizontalRechts,
		Vertikal:   RichtungVertikalVorne,
	}
	RichtungLinks = Richtung{
		Horizontal: RichtungHorizontalLinks,
		Vertikal:   RichtungVertikalKeine,
	}
	RichtungRechts = Richtung{
		Horizontal: RichtungHorizontalRechts,
		Vertikal:   RichtungVertikalKeine,
	}
	RichtungHintenLinks = Richtung{
		Horizontal: RichtungHorizontalLinks,
		Vertikal:   RichtungVertikalHinten,
	}
	RichtungHinten = Richtung{
		Horizontal: RichtungHorizontalKeine,
		Vertikal:   RichtungVertikalHinten,
	}
	RichtungHintenRechts = Richtung{
		Horizontal: RichtungHorizontalRechts,
		Vertikal:   RichtungVertikalHinten,
	}
)

type Richtungen map[Richtung]struct{}

func (r Richtungen) contains(richtung Richtung) bool {
	_, ok := r[richtung]
	return ok
}

func (r Richtungen) set(richtung Richtung, present bool) {
	if present {
		r[richtung] = struct{}{}
	} else {
		delete(r, richtung)
	}
}

func (r Richtungen) clone() Richtungen {
	clone := make(Richtungen, len(r))
	for richtung := range r {
		clone[richtung] = struct{}{}
	}
	return clone
}

var (
	RichtungenDiagonalVorne = Richtungen{
		RichtungVorneLinks:  {},
		RichtungVorneRechts: {},
	}
	RichtungenDiagonal = Richtungen{
		RichtungVorneLinks:   {},
		RichtungVorneRechts:  {},
		RichtungHintenLinks:  {},
		RichtungHintenRechts: {},
	}
	RichtungenSeiteDiagonalUndVorne = Richtungen{
		RichtungVorneLinks:  {},
		RichtungVorne:       {},
		RichtungVorneRechts: {},
		RichtungLinks:       {},
		RichtungRechts:      {},
	}
	RichtungenAlle = Richtungen{
		RichtungVorneLinks:   {},
		RichtungVorne:        {},
		RichtungVorneRechts:  {},
		RichtungLinks:        {},
		RichtungRechts:       {},
		RichtungHintenLinks:  {},
		RichtungHinten:       {},
		RichtungHintenRechts: {},
	}
)
