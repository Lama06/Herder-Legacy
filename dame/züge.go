package dame

import (
	"fmt"
	"github.com/Lama06/Herder-Legacy/minimax"
	"strings"
)

type zugSchritt struct {
	von, zu                position
	hatGeschlagenePosition bool
	geschlagenePosition    position
	ergebnis               brett
}

func (schritt1 zugSchritt) equals(schritt2 zugSchritt) bool {
	if schritt1.zu != schritt2.zu {
		return false
	}

	if schritt1.hatGeschlagenePosition != schritt2.hatGeschlagenePosition {
		return false
	}

	if schritt1.geschlagenePosition != schritt2.geschlagenePosition {
		return false
	}

	if !schritt1.ergebnis.equals(schritt2.ergebnis) {
		return false
	}

	return true
}

func (z zugSchritt) String() string {
	return fmt.Sprintf(`{
von: %v
zu: %v
geschlagenePosition: %v
ergebnis: 
%v
}`, z.von, z.zu, z.geschlagenePosition, z.ergebnis)
}

type zug []zugSchritt

var _ minimax.Zug = zug{}

func (zug1 zug) equals(zug2 zug) bool {
	if len(zug1) != len(zug2) {
		return false
	}

	for schritt := 0; schritt < len(zug1); schritt++ {
		if !zug1[schritt].equals(zug2[schritt]) {
			return false
		}
	}

	return true
}

func (z zug) ergebnis() brett {
	return z[len(z)-1].ergebnis
}

func (z zug) Ergebnis() minimax.Brett {
	return z.ergebnis()
}

func (z zug) startPosition() position {
	return z[0].von
}

func (z zug) endPosition() position {
	return z[len(z)-1].zu
}

func (z zug) String() string {
	var builder strings.Builder
	builder.WriteString("[\n")
	for _, schritt := range z {
		builder.WriteString(schritt.String())
		builder.WriteRune('\n')
	}
	builder.WriteString("]")
	return builder.String()
}

type züge []zug

func (züge1 züge) equals(züge2 züge) bool {
	if len(züge1) != len(züge2) {
		return false
	}

züge1:
	for _, zug1 := range züge1 {
		for _, zug2 := range züge2 {
			if zug1.equals(zug2) {
				continue züge1
			}
		}

		return false
	}

	return true
}

func (z züge) String() string {
	var builder strings.Builder
	builder.WriteString("[\n")
	for _, zug := range z {
		builder.WriteString(zug.String())
		builder.WriteRune('\n')
	}
	builder.WriteString("]")
	return builder.String()
}

func (b brett) möglicheSteinBewegenZüge(perspektive spieler, startPosition position, regeln regeln) züge {
	if !startPosition.valid(b.zeilen, b.spalten) {
		return nil
	}
	startFeld := b.feld(startPosition)
	if startFeld != stein(perspektive) {
		return nil
	}

	var möglicheZüge züge

	for bewegenRichtung := range regeln.steinBewegenRichtungen {
		neuePosition := startPosition.add(
			bewegenRichtung.vertikal.verschiebung(perspektive),
			bewegenRichtung.horizontal.verschiebung(perspektive),
		)
		if !neuePosition.valid(b.zeilen, b.spalten) {
			continue
		}
		neuesFeld := b.feld(neuePosition)
		if neuesFeld != feldLeer {
			continue
		}

		neuerStein := stein(perspektive)
		if neuePosition.zeile == b.umwandlungsZeile(perspektive) {
			neuerStein = dame(perspektive)
		}

		neuesBrett := b.clone()
		neuesBrett.setFeld(startPosition, feldLeer)
		neuesBrett.setFeld(neuePosition, neuerStein)
		möglicherZug := zug{
			zugSchritt{
				von:                    startPosition,
				zu:                     neuePosition,
				hatGeschlagenePosition: false,
				ergebnis:               neuesBrett,
			},
		}
		möglicheZüge = append(möglicheZüge, möglicherZug)
	}

	return möglicheZüge
}

func (b brett) möglicheSteinSchlagenZüge(
	perspektive spieler,
	startPosition position,
	regeln regeln,
	weiterschlagen bool,
) züge {
	if !startPosition.valid(b.zeilen, b.spalten) {
		return nil
	}
	startFeld := b.feld(startPosition)
	if startFeld != stein(perspektive) {
		return nil
	}

	var möglicheZüge züge

	for schlagenRichtung := range regeln.steinSchlagenRichtungen(weiterschlagen) {
		schlagenPosition := startPosition.add(
			schlagenRichtung.vertikal.verschiebung(perspektive),
			schlagenRichtung.horizontal.verschiebung(perspektive),
		)
		if !schlagenPosition.valid(b.zeilen, b.spalten) {
			continue
		}
		schlagenFeld := b.feld(schlagenPosition)
		if schlagenFeld == feldLeer || schlagenFeld == stein(perspektive) || schlagenFeld == dame(perspektive) {
			continue
		}

		neuePosition := schlagenPosition.add(
			schlagenRichtung.vertikal.verschiebung(perspektive),
			schlagenRichtung.horizontal.verschiebung(perspektive),
		)
		if !neuePosition.valid(b.zeilen, b.spalten) {
			continue
		}
		neuesFeld := b.feld(neuePosition)
		if neuesFeld != feldLeer {
			continue
		}

		neuerStein := stein(perspektive)
		if neuePosition.zeile == b.umwandlungsZeile(perspektive) {
			neuerStein = dame(perspektive)
		}

		neuesBrett := b.clone()
		neuesBrett.setFeld(startPosition, feldLeer)
		neuesBrett.setFeld(schlagenPosition, feldLeer)
		neuesBrett.setFeld(neuePosition, neuerStein)

		ersterZugSchritt := zugSchritt{
			von:                    startPosition,
			zu:                     neuePosition,
			hatGeschlagenePosition: true,
			geschlagenePosition:    schlagenPosition,
			ergebnis:               neuesBrett,
		}

		weiterschlagenZüge := neuesBrett.möglicheSteinSchlagenZüge(perspektive, neuePosition, regeln, true)
		if len(weiterschlagenZüge) == 0 {
			möglicherZug := zug{ersterZugSchritt}
			möglicheZüge = append(möglicheZüge, möglicherZug)
		} else {
			for _, weiterschlagenZug := range weiterschlagenZüge {
				möglicherZug := append(zug{ersterZugSchritt}, weiterschlagenZug...)
				möglicheZüge = append(möglicheZüge, möglicherZug)
			}
		}
	}

	return möglicheZüge
}

func (b brett) möglicheDameBewegenZüge(perspektive spieler, startPosition position, regeln regeln) züge {
	if !startPosition.valid(b.zeilen, b.spalten) {
		return nil
	}
	startFeld := b.feld(startPosition)
	if startFeld != dame(perspektive) {
		return nil
	}

	var möglicheZüge züge

bewegenRichtungen:
	for bewegenRichtung := range regeln.dameBewegenRichtungen {
		for anzahlSchritte := 1; anzahlSchritte < maxInt(b.zeilen, b.spalten); anzahlSchritte++ {
			neuePosition := startPosition.add(
				bewegenRichtung.vertikal.verschiebung(perspektive)*anzahlSchritte,
				bewegenRichtung.horizontal.verschiebung(perspektive)*anzahlSchritte,
			)
			if !neuePosition.valid(b.zeilen, b.spalten) {
				continue bewegenRichtungen
			}
			neuesFeld := b.feld(neuePosition)
			if neuesFeld != feldLeer {
				continue bewegenRichtungen
			}

			neuesBrett := b.clone()
			neuesBrett.setFeld(startPosition, feldLeer)
			neuesBrett.setFeld(neuePosition, dame(perspektive))
			möglicherZug := zug{
				zugSchritt{
					von:                    startPosition,
					zu:                     neuePosition,
					hatGeschlagenePosition: false,
					ergebnis:               neuesBrett,
				},
			}
			möglicheZüge = append(möglicheZüge, möglicherZug)
		}
	}

	return möglicheZüge
}

func (b brett) möglicheDameSchlagenZüge(
	perspektive spieler,
	startPosition position,
	regeln regeln,
	weiterschlagen bool,
) züge {
	if !startPosition.valid(b.zeilen, b.spalten) {
		return nil
	}
	startFeld := b.feld(startPosition)
	if startFeld != dame(perspektive) {
		return nil
	}

	var möglicheZüge züge

schlagenRichtungen:
	for schlagenRichtung := range regeln.dameSchlagenRichtungen(weiterschlagen) {
		for anzahlSchritte := 1; anzahlSchritte < maxInt(b.zeilen, b.spalten); anzahlSchritte++ {
			schlagenPosition := startPosition.add(
				schlagenRichtung.vertikal.verschiebung(perspektive)*anzahlSchritte,
				schlagenRichtung.horizontal.verschiebung(perspektive)*anzahlSchritte,
			)
			if !schlagenPosition.valid(b.zeilen, b.spalten) {
				continue schlagenRichtungen
			}
			schlagenFeld := b.feld(schlagenPosition)
			if schlagenFeld == stein(perspektive) || schlagenFeld == dame(perspektive) {
				continue schlagenRichtungen
			}
			if schlagenFeld == feldLeer {
				continue
			}

			neuePosition := schlagenPosition.add(
				schlagenRichtung.vertikal.verschiebung(perspektive),
				schlagenRichtung.horizontal.verschiebung(perspektive),
			)
			if !neuePosition.valid(b.zeilen, b.spalten) {
				continue schlagenRichtungen
			}
			neuesFeld := b.feld(neuePosition)
			if neuesFeld != feldLeer {
				continue schlagenRichtungen
			}

			neuesBrett := b.clone()
			neuesBrett.setFeld(startPosition, feldLeer)
			neuesBrett.setFeld(schlagenPosition, feldLeer)
			neuesBrett.setFeld(neuePosition, dame(perspektive))

			ersterZugSchritt := zugSchritt{
				von:                    startPosition,
				zu:                     neuePosition,
				hatGeschlagenePosition: true,
				geschlagenePosition:    schlagenPosition,
				ergebnis:               neuesBrett,
			}

			weiterschlagenZüge := neuesBrett.möglicheDameSchlagenZüge(perspektive, neuePosition, regeln, true)
			if len(weiterschlagenZüge) == 0 {
				möglicherZug := zug{ersterZugSchritt}
				möglicheZüge = append(möglicheZüge, möglicherZug)
			} else {
				for _, weiterschlagenZug := range weiterschlagenZüge {
					möglicherZug := append(zug{ersterZugSchritt}, weiterschlagenZug...)
					möglicheZüge = append(möglicheZüge, möglicherZug)
				}
			}

			continue schlagenRichtungen
		}
	}

	return möglicheZüge
}

func (b brett) möglicheZüge(perspektive spieler, regeln regeln, gewonnenÜberprüfen bool) züge {
	if gewonnenÜberprüfen && b.gewonnen(perspektive, regeln) {
		return nil
	}

	var möglicheZüge züge

	for zeile := 0; zeile < b.zeilen; zeile++ {
		for spalte := 0; spalte < b.spalten; spalte++ {
			startPosition := position{
				zeile:  zeile,
				spalte: spalte,
			}
			möglicheZüge = append(möglicheZüge, b.möglicheSteinSchlagenZüge(perspektive, startPosition, regeln, false)...)
			möglicheZüge = append(möglicheZüge, b.möglicheDameSchlagenZüge(perspektive, startPosition, regeln, false)...)
		}
	}

	if regeln.schlagZwang && len(möglicheZüge) != 0 {
		return möglicheZüge
	}

	for zeile := 0; zeile < b.zeilen; zeile++ {
		for spalte := 0; spalte < b.spalten; spalte++ {
			startPosition := position{
				zeile:  zeile,
				spalte: spalte,
			}
			möglicheZüge = append(möglicheZüge, b.möglicheSteinBewegenZüge(perspektive, startPosition, regeln)...)
			möglicheZüge = append(möglicheZüge, b.möglicheDameBewegenZüge(perspektive, startPosition, regeln)...)
		}
	}

	return möglicheZüge
}

func (b brett) möglicheZügeMitStartPosition(startPosition position, regeln regeln) züge {
	if !startPosition.valid(b.zeilen, b.spalten) {
		return nil
	}
	feld := b.feld(startPosition)
	eigentümer, ok := feld.eigentümer()
	if !ok {
		return nil
	}
	möglicheZüge := b.möglicheZüge(eigentümer, regeln, true)
	var möglicheZügeMitStartPosition züge
	for _, möglicherZug := range möglicheZüge {
		if möglicherZug.startPosition() == startPosition {
			möglicheZügeMitStartPosition = append(möglicheZügeMitStartPosition, möglicherZug)
		}
	}
	return möglicheZügeMitStartPosition
}

func (b brett) MöglicheZüge(perspektive minimax.Spieler, aiRegeln minimax.Regeln) []minimax.Zug {
	züge := b.möglicheZüge(perspektive.(spieler), aiRegeln.(regeln), true)
	aiZüge := make([]minimax.Zug, len(züge))
	for i, zug := range züge {
		aiZüge[i] = zug
	}
	return aiZüge
}
