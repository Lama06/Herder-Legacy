package jacobsalptraum

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"slices"

	"github.com/Lama06/Herder-Legacy/minimax"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

type Brett struct {
	breite, höhe int
	zeilen       [][]feld
}

var _ minimax.Brett = Brett{}

func parseBrett(zeilenText ...string) (Brett, error) {
	zeilen := make([][]feld, len(zeilenText))
	var breite int
	for zeileIndex, zeileText := range zeilenText {
		zeileBuchstaben := []rune(zeileText)

		if zeileIndex == 0 {
			breite = len(zeileBuchstaben)
		} else if len(zeileBuchstaben) != breite {
			return Brett{}, errors.New("nicht alle Zeilen haben gleich viele Spalten")
		}

		zeile := make([]feld, breite)
		for spalteIndex := range zeile {
			feld, ok := parseFeld(zeileBuchstaben[spalteIndex])
			if !ok {
				return Brett{}, errors.New(fmt.Sprintf("ungültiges feld: %c", zeileBuchstaben[spalteIndex]))
			}
			zeile[spalteIndex] = feld
		}
		zeilen[zeileIndex] = zeile
	}
	return Brett{
		breite: breite,
		höhe:   len(zeilenText),
		zeilen: zeilen,
	}, nil
}

func MustParseBrett(zeilenText ...string) Brett {
	brett, err := parseBrett(zeilenText...)
	if err != nil {
		panic(err)
	}
	return brett
}

var StandardBrett = MustParseBrett(
	"________",
	"________",
	"________",
	"________",
	"________",
	"________",
	"BBBBBBBB",
	"TSLDKLST",
)

func (b Brett) clone() Brett {
	zeilen := make([][]feld, b.höhe)
	for zeile := range zeilen {
		zeilen[zeile] = slices.Clone(b.zeilen[zeile])
	}

	return Brett{
		breite: b.breite,
		höhe:   b.höhe,
		zeilen: zeilen,
	}
}

func (b Brett) vierGewinntSteineZählen() int {
	var anzahl int
	for spalte := 0; spalte < b.breite; spalte++ {
		for zeile := 0; zeile < b.höhe; zeile++ {
			if b.zeilen[zeile][spalte] == feldVierGewinntStein {
				anzahl++
			}
		}
	}
	return anzahl
}

func (b Brett) vierGewinntWiederholungenInZeilenZählen(länge int) int {
	var anzahl int
	for zeile := 0; zeile < b.höhe; zeile++ {
		var längeBisher int
		for spalte := 0; spalte < b.breite; spalte++ {
			feld := b.zeilen[zeile][spalte]
			if feld == feldVierGewinntStein {
				längeBisher++
				if längeBisher == länge {
					anzahl++
					längeBisher = 0
				}
				continue
			}
			längeBisher = 0
		}
	}
	return anzahl
}

func (b Brett) vierGewinntWiederholungenInSpaltenZählen(länge int) int {
	var anzahl int
	for spalte := 0; spalte < b.breite; spalte++ {
		var längeBisher int
		for zeile := 0; zeile < b.höhe; zeile++ {
			feld := b.zeilen[zeile][spalte]
			if feld == feldVierGewinntStein {
				längeBisher++
				if längeBisher == länge {
					anzahl++
					längeBisher = 0
				}
				continue
			}
			längeBisher = 0
		}
	}
	return anzahl
}

func (b Brett) vierGewinntWiederholungenDiagonalZählen(länge int) int {
	var anzahl int

	diagonalRichtungen := []struct{ x, y int }{
		{x: 1, y: -1}, {x: 1, y: 1},
	}

	for startZeile := 0; startZeile < b.höhe; startZeile++ {
		for startSpalte := 0; startSpalte < b.breite; startSpalte++ {
		richtungen:
			for _, richtung := range diagonalRichtungen {
				for offset := 0; offset < länge; offset++ {
					zeile := startZeile + richtung.y*offset
					spalte := startSpalte + richtung.x*offset
					if !(position{zeile: zeile, spalte: spalte}).valid(b.breite, b.höhe) {
						continue richtungen
					}
					feld := b.zeilen[zeile][spalte]
					if feld != feldVierGewinntStein {
						continue richtungen
					}
				}

				anzahl++
			}
		}
	}

	return anzahl
}

func (b Brett) vierGewinntWiederholungenZählen(länge int) int {
	return b.vierGewinntWiederholungenInZeilenZählen(länge) +
		b.vierGewinntWiederholungenInSpaltenZählen(länge) +
		b.vierGewinntWiederholungenDiagonalZählen(länge)
}

func (b Brett) gewinner(regeln Regeln) (Spieler, bool) {
	if b.vierGewinntSteineZählen() >= regeln.MaximaleVierGewinntSteine {
		return SpielerSchach, true
	}
	if b.vierGewinntWiederholungenZählen(4) >= 1 {
		return SpielerVierGewinnt, true
	}
	return false, false
}

func (b Brett) bewertung(perspektive Spieler, regeln Regeln) int {
	const (
		verlorenBewertung           = -1000
		gewonnenBewertung           = 1000
		dreierWiederholungBewertung = 1
	)

	gewinner, gibtGewinner := b.gewinner(regeln)
	if gibtGewinner {
		if gewinner == perspektive {
			return gewonnenBewertung
		}
		return verlorenBewertung
	}

	vierGewinntBewertung := dreierWiederholungBewertung * b.vierGewinntWiederholungenZählen(3)

	switch perspektive {
	case SpielerVierGewinnt:
		return vierGewinntBewertung
	case SpielerSchach:
		return -vierGewinntBewertung
	default:
		panic("unreachable")
	}
}

func (b Brett) MinimaxBewertung(perspektive minimax.Spieler, minimaxRegeln minimax.Regeln) int {
	return b.bewertung(perspektive.(Spieler), minimaxRegeln.(Regeln))
}

func (b Brett) MinimaxMöglicheZüge(perspektive minimax.Spieler, minimaxRegeln minimax.Regeln) []minimax.Zug {
	_, gibtGewinner := b.gewinner(minimaxRegeln.(Regeln))
	if gibtGewinner {
		return nil
	}

	var minimaxZüge []minimax.Zug
	switch perspektive.(Spieler) {
	case SpielerVierGewinnt:
		züge := b.möglicheVierGewinntZüge()
		minimaxZüge = make([]minimax.Zug, len(züge))
		for i, zug := range züge {
			minimaxZüge[i] = zug
		}
	case SpielerSchach:
		züge := b.möglicheSchachZüge()
		minimaxZüge = make([]minimax.Zug, len(züge))
		for i, zug := range züge {
			minimaxZüge[i] = zug
		}
	}
	return minimaxZüge
}

type brettDrawOptions struct {
	x, y, width, height float64
}

func (b Brett) calculateFeldSize(options brettDrawOptions) float64 {
	return math.Min(
		options.width/float64(b.breite),
		options.height/float64(b.höhe),
	)
}

func (b Brett) calculateAbstandX(options brettDrawOptions) float64 {
	return (options.width - b.calculateFeldSize(options)*float64(b.breite)) / 2
}

func (b Brett) calculateAbstandY(options brettDrawOptions) float64 {
	return (options.height - b.calculateFeldSize(options)*float64(b.höhe)) / 2
}

func (b Brett) draw(
	options brettDrawOptions,
	screen *ebiten.Image,
	istFeldAusgewählt bool,
	ausgewähltesFeld position,
) {
	zielFelder := make(map[position]struct{})
	if istFeldAusgewählt {
		for _, zug := range b.möglicheSchachZüge() {
			if zug.start == ausgewähltesFeld {
				zielFelder[zug.ziel] = struct{}{}
			}
		}
	}

	feldSize := b.calculateFeldSize(options)

	for zeile := 0; zeile < b.höhe; zeile++ {
		for spalte := 0; spalte < b.breite; spalte++ {
			pos := position{zeile: zeile, spalte: spalte}
			_, istZiel := zielFelder[pos]
			screenX, screenY := b.brettPositionToScreenPosition(options, pos)

			var hintergrundFarbe color.Color
			switch {
			case istFeldAusgewählt && pos == ausgewähltesFeld:
				hintergrundFarbe = colornames.Purple
			case istZiel:
				hintergrundFarbe = colornames.Pink
			default:
				hintergrundFarbe = pos.farbe()
			}
			vector.DrawFilledRect(
				screen,
				float32(screenX), float32(screenY),
				float32(feldSize), float32(feldSize),
				hintergrundFarbe,
				true,
			)

			img := b.zeilen[zeile][spalte].image()
			if img == nil {
				continue
			}

			var drawOptions ebiten.DrawImageOptions
			drawOptions.GeoM.Scale(feldSize/float64(img.Bounds().Dx()), feldSize/float64(img.Bounds().Dy()))
			drawOptions.GeoM.Translate(screenX, screenY)
			screen.DrawImage(img, &drawOptions)
		}
	}
}

func (b Brett) screenPositionToBrettPosition(options brettDrawOptions, screenX, screenY float64) (position, bool) {
	abstandX, abstandY := b.calculateAbstandX(options), b.calculateAbstandY(options)
	startX, startY := options.x+abstandX, options.y+abstandY
	feldSize := b.calculateFeldSize(options)

	zeile := int((screenY - startY) / feldSize)
	spalte := int((screenX - startX) / feldSize)
	pos := position{zeile: zeile, spalte: spalte}
	return pos, pos.valid(b.breite, b.höhe)
}

func (b Brett) brettPositionToScreenPosition(options brettDrawOptions, pos position) (x, y float64) {
	abstandX, abstandY := b.calculateAbstandX(options), b.calculateAbstandY(options)
	startX, startY := options.x+abstandX, options.y+abstandY
	feldSize := b.calculateFeldSize(options)
	return startX + feldSize*float64(pos.spalte), startY + feldSize*float64(pos.zeile)
}
