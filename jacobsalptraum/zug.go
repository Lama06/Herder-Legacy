package jacobsalptraum

import (
	"github.com/Lama06/Herder-Legacy/minimax"
)

type zug interface {
	minimax.Zug
}

type schachZug struct {
	start    position
	ziel     position
	ergebnis Brett
}

func (s schachZug) MinimaxErgebnis() minimax.Brett {
	return s.ergebnis
}

type vierGewinntZug struct {
	position position
	ergebnis Brett
}

func (v vierGewinntZug) MinimaxErgebnis() minimax.Brett {
	return v.ergebnis
}

var schachFigurenZugInfos = map[feld]struct {
	bewegungen            []struct{ x, y int }
	bewegungsWiederholung bool
}{
	feldBauer: {
		bewegungen: []struct{ x, y int }{
			{x: 0, y: -1},
		},
		bewegungsWiederholung: false,
	},
	feldLäufer: {
		bewegungen: []struct{ x, y int }{
			{x: -1, y: -1}, {x: -1, y: 1},
			{x: 1, y: -1}, {x: 1, y: 1},
		},
		bewegungsWiederholung: true,
	},
	feldSpringer: {
		bewegungen: []struct{ x, y int }{
			{x: -2, y: 1}, {x: -2, y: -1},
			{x: -1, y: 2}, {x: -1, y: -2},
			{x: 1, y: 2}, {x: 1, y: -2},
			{x: 2, y: 1}, {x: 2, y: -1},
		},
		bewegungsWiederholung: false,
	},
	feldTurm: {
		bewegungen: []struct{ x, y int }{
			{x: 0, y: 1}, {x: 0, y: -1},
			{x: 1, y: 0}, {x: -1, y: 0},
		},
		bewegungsWiederholung: true,
	},
	feldDame: {
		bewegungen: []struct{ x, y int }{
			{x: -1, y: -1}, {x: -1, y: 0}, {x: -1, y: 1},
			{x: 0, y: -1}, {x: 0, y: 0}, {x: 0, y: 1},
			{x: 1, y: -1}, {x: 1, y: 0}, {x: 1, y: 1},
		},
		bewegungsWiederholung: true,
	},
	feldKönig: {
		bewegungen: []struct{ x, y int }{
			{x: -1, y: -1}, {x: -1, y: 0}, {x: -1, y: 1},
			{x: 0, y: -1}, {x: 0, y: 0}, {x: 0, y: 1},
			{x: 1, y: -1}, {x: 1, y: 0}, {x: 1, y: 1},
		},
		bewegungsWiederholung: false,
	},
}

func (b Brett) zeileFürNeuenVierGewinntSteinFinden(spalte int) (int, bool) {
	for zeile := 0; zeile < b.höhe; zeile++ {
		if b.zeilen[zeile][spalte] != feldLeer {
			if zeile == 0 {
				return 0, false
			}
			return zeile - 1, true
		}
	}
	return b.höhe - 1, true
}

func (b Brett) möglicheVierGewinntZüge() []vierGewinntZug {
	var züge []vierGewinntZug
	for spalte := 0; spalte < b.breite; spalte++ {
		zeile, ok := b.zeileFürNeuenVierGewinntSteinFinden(spalte)
		if !ok {
			continue
		}

		neuesBrett := b.clone()
		neuesBrett.zeilen[zeile][spalte] = feldVierGewinntStein

		züge = append(züge, vierGewinntZug{
			position: position{zeile: zeile, spalte: spalte},
			ergebnis: neuesBrett,
		})
	}
	return züge
}

func (b Brett) vierGewinntSteineRunterfallenLassen() {
	for zeile := b.höhe - 2; zeile >= 0; zeile-- {
	spalten:
		for spalte := 0; spalte < b.breite; spalte++ {
			if b.zeilen[zeile][spalte] != feldVierGewinntStein {
				continue
			}

			b.zeilen[zeile][spalte] = feldLeer

			for neueZeile := zeile + 1; neueZeile < b.höhe; neueZeile++ {
				if b.zeilen[neueZeile][spalte] != feldLeer {
					b.zeilen[neueZeile-1][spalte] = feldVierGewinntStein
					continue spalten
				}
			}

			b.zeilen[b.höhe-1][spalte] = feldVierGewinntStein
		}
	}
}

func (b Brett) möglicheSchachZüge() []schachZug {
	var züge []schachZug
	for zeile := 0; zeile < b.höhe; zeile++ {
		for spalte := 0; spalte < b.breite; spalte++ {
			feld := b.zeilen[zeile][spalte]
			schachFigurInfo, istSchachFigur := schachFigurenZugInfos[feld]
			if !istSchachFigur {
				continue
			}

			bewegungen := schachFigurInfo.bewegungen
			if feld == feldBauer && zeile == b.höhe-2 && b.zeilen[zeile-1][spalte] == feldLeer {
				bewegungen = []struct{ x, y int }{{x: 0, y: -1}, {x: 0, y: -2}}
			}

			for _, bewegung := range bewegungen {
				for bewegungsWiederholung := 1; bewegungsWiederholung == 1 || schachFigurInfo.bewegungsWiederholung; bewegungsWiederholung++ {
					neueZeile := zeile + bewegung.y*bewegungsWiederholung
					neueSpalte := spalte + bewegung.x*bewegungsWiederholung
					neuePosition := position{zeile: neueZeile, spalte: neueSpalte}
					if !neuePosition.valid(b.breite, b.höhe) {
						break
					}
					neuesFeld := b.zeilen[neueZeile][neueSpalte]
					if neuesFeld != feldLeer {
						break
					}
					neuesBrett := b.clone()
					neuesBrett.zeilen[zeile][spalte] = feldLeer
					neuesBrett.zeilen[neueZeile][neueSpalte] = feld
					neuesBrett.vierGewinntSteineRunterfallenLassen()
					züge = append(züge, schachZug{
						start:    position{zeile: zeile, spalte: spalte},
						ziel:     neuePosition,
						ergebnis: neuesBrett,
					})
				}
			}
		}
	}
	return züge
}
