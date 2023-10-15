package sodoku

type position struct {
	zeile, spalte int
}

func (p position) nächste() (position, bool) {
	if p.zeile == 8 && p.spalte == 8 {
		return position{}, false
	}
	if p.spalte == 8 {
		return position{spalte: 0, zeile: p.zeile + 1}, true
	}
	return position{spalte: p.spalte + 1, zeile: p.zeile}, true
}

type sodoku [9][9]byte

func (s sodoku) möglicheZahlen() [9][9]map[byte]struct{} {
	var möglicheZahlen [9][9]map[byte]struct{}
	for zeile := 0; zeile < 9; zeile++ {
		for spalte := 0; spalte < 9; spalte++ {
			if s[zeile][spalte] != 0 {
				continue
			}

			möglicheZahlen[zeile][spalte] = make(map[byte]struct{}, 9)
			for zahl := 1; zahl <= 9; zahl++ {
				möglicheZahlen[zeile][spalte][byte(zahl)] = struct{}{}
			}
		}
	}

	for startZeile := 0; startZeile < 9; startZeile++ {
		for startSpalte := 0; startSpalte < 9; startSpalte++ {
			zahl := s[startZeile][startSpalte]
			if zahl == 0 {
				continue
			}

			for spalte := 0; spalte < 9; spalte++ {
				if s[startZeile][spalte] != 0 {
					continue
				}
				delete(möglicheZahlen[startZeile][spalte], zahl)
			}

			for zeile := 0; zeile < 9; zeile++ {
				if s[zeile][startSpalte] != 0 {
					continue
				}
				delete(möglicheZahlen[zeile][startSpalte], zahl)
			}

			quadratX := (startSpalte / 3) * 3
			quadratY := (startZeile / 3) * 3
			for zeile := quadratY; zeile < quadratY+3; quadratY++ {
				for spalte := quadratX; spalte < quadratX+3; quadratX++ {
					if s[zeile][spalte] != 0 {
						continue
					}
					delete(möglicheZahlen[zeile][spalte], zahl)
				}
			}
		}
	}

	return möglicheZahlen
}

func (s sodoku) teilweiseVereinfachen() (sodoku, bool) {
	möglicheZahlen := s.möglicheZahlen()

	var vereinfacht bool
	for zeile := 0; zeile < 9; zeile++ {
		for spalte := 0; spalte < 9; spalte++ {
			if s[zeile][spalte] != 0 {
				continue
			}

			if len(möglicheZahlen[zeile][spalte]) == 1 {
				for zahl := range möglicheZahlen[zeile][spalte] {
					s[zeile][spalte] = zahl
					vereinfacht = true
				}
			}
		}
	}

	return s, vereinfacht
}

func (s sodoku) vereinfachen() sodoku {
	for {
		var vereinfacht bool
		if s, vereinfacht = s.teilweiseVereinfachen(); !vereinfacht {
			return s
		}
	}
}

func (s sodoku) hatFehlerInSpalten() bool {
	zahlen := make(map[byte]struct{}, 9)
	for spalte := 0; spalte < 9; spalte++ {
		for zeile := 0; zeile < 9; zeile++ {
			zahl := s[zeile][spalte]
			if zahl == 0 {
				continue
			}
			if _, doppelt := zahlen[zahl]; doppelt {
				return true
			}
			zahlen[zahl] = struct{}{}
		}
		clear(zahlen)
	}
	return false
}

func (s sodoku) hatFehlerInZeilen() bool {
	zahlen := make(map[byte]struct{}, 9)
	for zeile := 0; zeile < 9; zeile++ {
		for spalte := 0; spalte < 9; spalte++ {
			zahl := s[zeile][spalte]
			if zahl == 0 {
				continue
			}
			if _, doppelt := zahlen[zahl]; doppelt {
				return true
			}
			zahlen[zahl] = struct{}{}
		}
		clear(zahlen)
	}
	return false
}

func (s sodoku) hatFelherInQuadrat() bool {
	zahlen := make(map[byte]struct{}, 9)
	for quadratX := 0; quadratX < 3; quadratX++ {
		for quadratY := 0; quadratY < 3; quadratY++ {
			for spalte := quadratX * 3; spalte < quadratX*3+3; spalte++ {
				for zeile := quadratY * 3; zeile < quadratY*3+3; zeile++ {
					zahl := s[zeile][spalte]
					if zahl == 0 {
						continue
					}
					if _, doppelt := zahlen[zahl]; doppelt {
						return true
					}
					zahlen[zahl] = struct{}{}
				}
			}
			clear(zahlen)
		}
	}
	return false
}

func (s sodoku) hatFehler() bool {
	return s.hatFehlerInZeilen() || s.hatFehlerInSpalten() || s.hatFelherInQuadrat()
}

func (s sodoku) nächstesFreiesFeld(pos position) (position, bool) {
	for s[pos.zeile][pos.spalte] != 0 {
		var ok bool
		pos, ok = pos.nächste()
		if !ok {
			return position{}, false
		}
	}
	return pos, true
}

func (s sodoku) rekursivLösen(pos position) (sodoku, bool) {
	for zahl := 1; zahl <= 9; zahl++ {
		s[pos.zeile][pos.spalte] = byte(zahl)
		if s.hatFehler() {
			continue
		}
		nächstesFeld, ok := s.nächstesFreiesFeld(pos)
		if !ok {
			return s, true
		}
		lösung, ok := s.rekursivLösen(nächstesFeld)
		if ok {
			return lösung, true
		}
	}
	return sodoku{}, false
}

func (s sodoku) lösen() (sodoku, bool) {
	ersteFreiePosition, ok := s.nächstesFreiesFeld(position{0, 0})
	if !ok {
		return sodoku{}, false
	}
	return s.rekursivLösen(ersteFreiePosition)
}
