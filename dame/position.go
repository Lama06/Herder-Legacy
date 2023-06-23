package dame

type position struct {
	zeile  int
	spalte int
}

func (p position) add(zeilen, spalten int) position {
	return position{
		zeile:  p.zeile + zeilen,
		spalte: p.spalte + spalten,
	}
}

func (p position) schwarz() bool {
	return p.zeile%2 == p.spalte%2
}

func (p position) valid(zeilen, spalten int) bool {
	return p.zeile >= 0 && p.spalte >= 0 && p.zeile < zeilen && p.spalte < spalten
}
