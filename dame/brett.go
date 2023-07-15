package dame

import (
	"errors"
	"fmt"
	"github.com/Lama06/Herder-Legacy/minimax"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"strings"
)

type brett struct {
	zeilen, spalten int
	felder          [][]feld
}

var _ minimax.Brett = brett{}

func newLeeresBrett(zeilen, spalten int) brett {
	felder := make([][]feld, zeilen)
	for zeile := 0; zeile < zeilen; zeile++ {
		felder[zeile] = make([]feld, spalten)
	}
	return brett{
		zeilen:  zeilen,
		spalten: spalten,
		felder:  felder,
	}
}

func parseBrett(zeilen ...string) (brett, error) {
	zeilenAnzahl := len(zeilen)
	if zeilenAnzahl == 0 {
		return brett{}, errors.New("Brett muss mindestens eine Zeile haben")
	}

	spaltenAnzahl := len(zeilen[0])
	for _, zeile := range zeilen {
		if len(zeile) != spaltenAnzahl {
			return brett{}, errors.New("alle Zeilen müssen gleich lang sein")
		}
	}

	parsedBrett := newLeeresBrett(zeilenAnzahl, spaltenAnzahl)
	for zeile, zeileText := range zeilen {
		for spalte, feldZeichen := range zeileText {
			feld, err := parseFeld(feldZeichen)
			if err != nil {
				return brett{}, fmt.Errorf("feld konnte nicht erkannt werden: %w", err)
			}
			parsedBrett.setFeld(position{zeile: zeile, spalte: spalte}, feld)
		}
	}
	return parsedBrett, nil
}

func mustParseBrett(zeilen ...string) brett {
	brett, err := parseBrett(zeilen...)
	if err != nil {
		panic(err)
	}
	return brett
}

func (b brett) clone() brett {
	felder := make([][]feld, b.zeilen)
	for zeile := 0; zeile < b.zeilen; zeile++ {
		felder[zeile] = make([]feld, b.spalten)
		for spalte := 0; spalte < b.spalten; spalte++ {
			felder[zeile][spalte] = b.felder[zeile][spalte]
		}
	}
	return brett{
		zeilen:  b.zeilen,
		spalten: b.spalten,
		felder:  felder,
	}
}

func (brett1 brett) equals(brett2 brett) bool {
	if brett1.zeilen != brett2.zeilen || brett1.spalten != brett2.spalten {
		return false
	}

	for zeile := 0; zeile < brett1.zeilen; zeile++ {
		for spalte := 0; spalte < brett1.spalten; spalte++ {
			if brett1.felder[zeile][spalte] != brett2.felder[zeile][spalte] {
				return false
			}
		}
	}

	return true
}

func (b brett) String() string {
	var result strings.Builder
	for zeile := 0; zeile < b.zeilen; zeile++ {
		if zeile != 0 {
			result.WriteRune('\n')
		}

		for spalte := 0; spalte < b.spalten; spalte++ {
			feld := b.feld(position{zeile: zeile, spalte: spalte})
			result.WriteRune(feld.zeichen())
		}
	}
	return result.String()
}

func (b brett) feld(pos position) feld {
	if !pos.valid(b.zeilen, b.spalten) {
		panic("invalid position")
	}

	return b.felder[pos.zeile][pos.spalte]
}

func (b brett) setFeld(pos position, neuesFeld feld) {
	b.felder[pos.zeile][pos.spalte] = neuesFeld
}

func (b brett) umwandlungsZeile(perspektive spieler) int {
	switch perspektive {
	case spielerLehrer:
		return b.zeilen - 1
	case spielerSchüler:
		return 0
	default:
		panic("unreachable")
	}
}

func (b brett) felderZählen(gesucht feld) int {
	var anzahl int
	for zeile := 0; zeile < b.zeilen; zeile++ {
		for spalte := 0; spalte < b.spalten; spalte++ {
			if b.feld(position{zeile: zeile, spalte: spalte}) == gesucht {
				anzahl++
			}
		}
	}
	return anzahl
}

func (b brett) gewonnen(perspektive spieler, regeln regeln) bool {
	return len(b.möglicheZüge(perspektive.gegner(), regeln, false)) == 0
}

func (b brett) bewertung(perspektive spieler, regeln regeln) int {
	const (
		gewonnenBewertung = 1000
		steinBewertung    = 1
		dameBewertung     = 3
	)

	if b.gewonnen(perspektive, regeln) {
		return gewonnenBewertung
	}

	if b.gewonnen(perspektive.gegner(), regeln) {
		return -gewonnenBewertung
	}

	perspektiveBewertung := steinBewertung*b.felderZählen(stein(perspektive)) + dameBewertung*b.felderZählen(dame(perspektive))
	gegnerBewertung := steinBewertung*b.felderZählen(stein(perspektive.gegner())) + dameBewertung*b.felderZählen(dame(perspektive.gegner()))

	return perspektiveBewertung - gegnerBewertung
}

func (b brett) Bewertung(perspektive minimax.Spieler, aiRegeln minimax.Regeln) int {
	return b.bewertung(perspektive.(spieler), aiRegeln.(regeln))
}

func (b brett) feldSize(maxBrettBreite, maxBrettHoehe float64) float64 {
	return minFloat64(maxBrettBreite/float64(b.spalten), maxBrettHoehe/float64(b.zeilen))
}

func (b brett) brettSize(maxBrettBreite, maxBrettHoehe float64) (brettBreite, brettHoehe float64) {
	feldSize := b.feldSize(maxBrettBreite, maxBrettHoehe)
	return feldSize * float64(b.spalten), feldSize * float64(b.zeilen)
}

func (b brett) brettAbstand(maxBrettBreite, maxBrettHoehe float64) (brettAbstandX, brettAbstandY float64) {
	brettBreite, brettHoehe := b.brettSize(maxBrettBreite, maxBrettHoehe)
	return (maxBrettBreite - brettBreite) / 2, (maxBrettHoehe - brettHoehe) / 2
}

func (b brett) feldPosition(
	pos position,
	brettX, brettY, maxBrettBreite, maxBrettHoehe float64,
) (feldX, feldY float64) {
	feldSize := b.feldSize(maxBrettBreite, maxBrettHoehe)
	brettAbstandX, brettAbstandY := b.brettAbstand(maxBrettBreite, maxBrettHoehe)
	return brettX + brettAbstandX + feldSize*float64(pos.spalte), brettY + brettAbstandY + feldSize*float64(pos.zeile)
}

func (b brett) draw(
	screen *ebiten.Image,
	brettX, brettY, maxBrettBreite, maxBrettHoehe float64,
	hatAusgewähltePosition bool, ausgewähltePosition position,
	regeln regeln,
) {
	zugEndPositionen := make(map[position]struct{})
	if hatAusgewähltePosition {
		for _, moeglicherZug := range b.möglicheZügeMitStartPosition(ausgewähltePosition, regeln) {
			zugEndPositionen[moeglicherZug.endPosition()] = struct{}{}
		}
	}

	feldSize := b.feldSize(maxBrettBreite, maxBrettHoehe)

	for zeile := 0; zeile < b.zeilen; zeile++ {
		for spalte := 0; spalte < b.spalten; spalte++ {
			pos := position{zeile: zeile, spalte: spalte}
			feld := b.feld(pos)
			feldX, feldY := b.feldPosition(pos, brettX, brettY, maxBrettBreite, maxBrettHoehe)

			var clr color.Color
			if hatAusgewähltePosition && ausgewähltePosition == pos {
				clr = colornames.Purple
			} else if _, istZugEndPosition := zugEndPositionen[pos]; istZugEndPosition {
				clr = colornames.Pink
			} else if pos.schwarz() {
				clr = colornames.Gray
			} else {
				clr = colornames.Wheat
			}
			vector.DrawFilledRect(
				screen,
				float32(feldX),
				float32(feldY),
				float32(feldSize),
				float32(feldSize),
				clr,
				true,
			)

			feld.draw(screen, feldX, feldY, feldSize)
		}
	}
}

func (b brett) screenPositionToBrettPosition(
	screenX, screenY float64,
	brettX, brettY, maxBrettBreite, maxBrettHoehe float64,
) (pos position, ok bool) {
	feldSize := b.feldSize(maxBrettBreite, maxBrettHoehe)
	brettAbstandX, brettAbstandY := b.brettAbstand(maxBrettBreite, maxBrettHoehe)
	pos = position{
		zeile:  int(math.Floor((screenY - (brettY + brettAbstandY)) / feldSize)),
		spalte: int(math.Floor((screenX - (brettX + brettAbstandX)) / feldSize)),
	}
	if !pos.valid(b.zeilen, b.spalten) {
		return position{}, false
	}
	return pos, true
}
