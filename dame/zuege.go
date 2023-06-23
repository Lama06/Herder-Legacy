package dame

import (
	"fmt"
	"github.com/Lama06/Herder-Legacy/ai"
	"strings"
)

type zugSchritt struct {
	von, zu             position
	geschlagenePosition *position
	ergebnis            brett
}

func (schritt1 zugSchritt) equals(schritt2 zugSchritt) bool {
	if schritt1.zu != schritt2.zu {
		return false
	}

	if schritt1.geschlagenePosition == nil && schritt2.geschlagenePosition != nil {
		return false
	}

	if schritt1.geschlagenePosition != nil &&
		(schritt2.geschlagenePosition == nil || *schritt2.geschlagenePosition != *schritt1.geschlagenePosition) {
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

var _ ai.Zug = zug{}

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

func (z zug) Ergebnis() ai.Brett {
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

type zuege []zug

func (zuege1 zuege) equals(zuege2 zuege) bool {
	if len(zuege1) != len(zuege2) {
		return false
	}

zuege1:
	for _, zug1 := range zuege1 {
		for _, zug2 := range zuege2 {
			if zug1.equals(zug2) {
				continue zuege1
			}
		}

		return false
	}

	return true
}

func (z zuege) String() string {
	var builder strings.Builder
	builder.WriteString("[\n")
	for _, zug := range z {
		builder.WriteString(zug.String())
		builder.WriteRune('\n')
	}
	builder.WriteString("]")
	return builder.String()
}

func (b brett) moeglicheSteinBewegenZuege(perspektive spieler, startPosition position, regeln regeln) zuege {
	if !startPosition.valid(b.zeilen, b.spalten) {
		return nil
	}
	startFeld := b.feld(startPosition)
	if startFeld != stein(perspektive) {
		return nil
	}

	var moeglicheZuege zuege

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
		moeglicherZug := zug{
			zugSchritt{
				von:                 startPosition,
				zu:                  neuePosition,
				geschlagenePosition: nil,
				ergebnis:            neuesBrett,
			},
		}
		moeglicheZuege = append(moeglicheZuege, moeglicherZug)
	}

	return moeglicheZuege
}

func (b brett) moeglicheSteinSchlagenZuege(
	perspektive spieler,
	startPosition position,
	regeln regeln,
	weiterschlagen bool,
) zuege {
	if !startPosition.valid(b.zeilen, b.spalten) {
		return nil
	}
	startFeld := b.feld(startPosition)
	if startFeld != stein(perspektive) {
		return nil
	}

	var moeglicheZuege zuege

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
			von:                 startPosition,
			zu:                  neuePosition,
			geschlagenePosition: &schlagenPosition,
			ergebnis:            neuesBrett,
		}

		weiterschlagenZuege := neuesBrett.moeglicheSteinSchlagenZuege(perspektive, neuePosition, regeln, true)
		if len(weiterschlagenZuege) == 0 {
			moeglicherZug := zug{ersterZugSchritt}
			moeglicheZuege = append(moeglicheZuege, moeglicherZug)
		} else {
			for _, weiterschlagenZug := range weiterschlagenZuege {
				moeglicherZug := append(zug{ersterZugSchritt}, weiterschlagenZug...)
				moeglicheZuege = append(moeglicheZuege, moeglicherZug)
			}
		}
	}

	return moeglicheZuege
}

func (b brett) moeglicheDameBewegenZuege(perspektive spieler, startPosition position, regeln regeln) zuege {
	if !startPosition.valid(b.zeilen, b.spalten) {
		return nil
	}
	startFeld := b.feld(startPosition)
	if startFeld != dame(perspektive) {
		return nil
	}

	var moeglicheZuege zuege

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
			moeglicherZug := zug{
				zugSchritt{
					von:                 startPosition,
					zu:                  neuePosition,
					geschlagenePosition: nil,
					ergebnis:            neuesBrett,
				},
			}
			moeglicheZuege = append(moeglicheZuege, moeglicherZug)
		}
	}

	return moeglicheZuege
}

func (b brett) moeglicheDameSchlagenZuege(
	perspektive spieler,
	startPosition position,
	regeln regeln,
	weiterschlagen bool,
) zuege {
	if !startPosition.valid(b.zeilen, b.spalten) {
		return nil
	}
	startFeld := b.feld(startPosition)
	if startFeld != dame(perspektive) {
		return nil
	}

	var moeglicheZuege zuege

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
				von:                 startPosition,
				zu:                  neuePosition,
				geschlagenePosition: &schlagenPosition,
				ergebnis:            neuesBrett,
			}

			weiterschlagenZuege := neuesBrett.moeglicheDameSchlagenZuege(perspektive, neuePosition, regeln, true)
			if len(weiterschlagenZuege) == 0 {
				moeglicherZug := zug{ersterZugSchritt}
				moeglicheZuege = append(moeglicheZuege, moeglicherZug)
			} else {
				for _, weiterschlagenZug := range weiterschlagenZuege {
					moeglicherZug := append(zug{ersterZugSchritt}, weiterschlagenZug...)
					moeglicheZuege = append(moeglicheZuege, moeglicherZug)
				}
			}

			continue schlagenRichtungen
		}
	}

	return moeglicheZuege
}

func (b brett) moeglicheZuege(perspektive spieler, regeln regeln, gewonnenUeberpruefen bool) zuege {
	if gewonnenUeberpruefen && b.gewonnen(perspektive, regeln) {
		return nil
	}

	var moeglicheZuege zuege

	for zeile := 0; zeile < b.zeilen; zeile++ {
		for spalte := 0; spalte < b.spalten; spalte++ {
			startPosition := position{
				zeile:  zeile,
				spalte: spalte,
			}
			moeglicheZuege = append(moeglicheZuege, b.moeglicheSteinSchlagenZuege(perspektive, startPosition, regeln, false)...)
			moeglicheZuege = append(moeglicheZuege, b.moeglicheDameSchlagenZuege(perspektive, startPosition, regeln, false)...)
		}
	}

	if regeln.schlagZwang && len(moeglicheZuege) != 0 {
		return moeglicheZuege
	}

	for zeile := 0; zeile < b.zeilen; zeile++ {
		for spalte := 0; spalte < b.spalten; spalte++ {
			startPosition := position{
				zeile:  zeile,
				spalte: spalte,
			}
			moeglicheZuege = append(moeglicheZuege, b.moeglicheSteinBewegenZuege(perspektive, startPosition, regeln)...)
			moeglicheZuege = append(moeglicheZuege, b.moeglicheDameBewegenZuege(perspektive, startPosition, regeln)...)
		}
	}

	return moeglicheZuege
}

func (b brett) moeglicheZuegeMitStartPosition(startPosition position, regeln regeln) zuege {
	if !startPosition.valid(b.zeilen, b.spalten) {
		return nil
	}
	feld := b.feld(startPosition)
	eigentuemer, ok := feld.eigentuemer()
	if !ok {
		return nil
	}
	moeglicheZuege := b.moeglicheZuege(eigentuemer, regeln, true)
	var moeglicheZuegeMitStartPosition zuege
	for _, moeglicherZug := range moeglicheZuege {
		if moeglicherZug.startPosition() == startPosition {
			moeglicheZuegeMitStartPosition = append(moeglicheZuegeMitStartPosition, moeglicherZug)
		}
	}
	return moeglicheZuegeMitStartPosition
}

func (b brett) MoeglicheZuege(perspektive ai.Spieler, aiRegeln ai.Regeln) []ai.Zug {
	zuege := b.moeglicheZuege(perspektive.(spieler), aiRegeln.(regeln), true)
	aiZuege := make([]ai.Zug, len(zuege))
	for i, zug := range zuege {
		aiZuege[i] = zug
	}
	return aiZuege
}
