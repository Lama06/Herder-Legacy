package poker

import (
	"image/color"
	"sort"

	"golang.org/x/image/colornames"
)

func kartenNachWertSortieren(karten []karte) {
	sort.Slice(karten, func(i, j int) bool {
		if karten[i].wert == karten[j].wert {
			return karten[i].symbol > karten[j].symbol
		}
		return karten[i].wert > karten[j].wert
	})
}

func ersteWertwiederholungFinden(karten []karte, länge int) (
	kartenMitWiederholtemWert []karte,
	übrigeKarten []karte,
	gefunden bool,
) {
karten:
	for i := range karten {
		if i < länge-1 {
			continue
		}

		wiederholterWert := karten[i].wert
		for offset := 1; offset < länge; offset++ {
			if karten[i-offset].wert != wiederholterWert {
				continue karten
			}
		}

		kartenMitWiederholtemWert = karten[i-länge+1 : i+1]
		übrigeKarten = append(karten[:i-länge+1:i-länge+1], karten[i+1:]...)

		return kartenMitWiederholtemWert, übrigeKarten, true
	}

	return nil, nil, false
}

type handKartenAuswahl [7]karte

type handArt int8

const (
	handArtHöchsteKarte handArt = iota
	handArtEinPaar
	handArtZweiPaare
	handArtDrilling
	handArtStraße
	handArtFlush
	handArtFullHouse
	handArtVierling
	// TODO:
	// handArtStraightFlush
	// handArtRoyalFlush
)

func (h handArt) parser() handParser {
	switch h {
	case handArtHöchsteKarte:
		return parseHöchsteKarteHand
	case handArtEinPaar:
		return parseEinPaarHand
	case handArtZweiPaare:
		return parseZweiPaareHand
	case handArtDrilling:
		return parseDrillingsHand
	case handArtStraße:
		return parseStraßeHand
	case handArtFlush:
		return parseFlushHand
	case handArtFullHouse:
		return parseFullHouseHand
	case handArtVierling:
		return parseVierlingHand
	default:
		panic("unreachable")
	}
}

type hand interface {
	art() handArt

	karten() [5]karte

	displayName() string

	visualisierung(karte) color.Color
}

type handParser func(handKartenAuswahl) hand

func parseHand(karten handKartenAuswahl) hand {
	for art := handArtVierling; art >= handArtHöchsteKarte; art-- {
		if hand := art.parser()(karten); hand != nil {
			return hand
		}
	}
	return nil
}

func compareHände(hand1, hand2 hand) int {
	if hand1.art() < hand2.art() {
		return -1
	}
	if hand2.art() < hand1.art() {
		return 1
	}

	karten1 := hand1.karten()
	karten2 := hand2.karten()

	for i := range karten1 {
		if karten1[i].wert < karten2[i].wert {
			return -1
		}
		if karten2[i].wert < karten1[i].wert {
			return 1
		}
	}

	return 0
}

type höchsteKarteHand [5]karte

func parseHöchsteKarteHand(karten handKartenAuswahl) hand {
	kartenNachWertSortieren(karten[:])
	return höchsteKarteHand(karten[:5])
}

func (h höchsteKarteHand) art() handArt {
	return handArtHöchsteKarte
}

func (h höchsteKarteHand) karten() [5]karte {
	return h
}

func (h höchsteKarteHand) displayName() string {
	return "Höchste Karte: " + h[0].wert.String()
}

func (h höchsteKarteHand) visualisierung(karte karte) color.Color {
	switch karte {
	case h[0]:
		return colornames.Gold
	case h[1]:
		return colornames.Silver
	case h[2]:
		return colornames.Brown
	default:
		return nil
	}
}

type einPaarHand struct {
	paar      [2]karte
	beikarten [3]karte
}

func parseEinPaarHand(karten handKartenAuswahl) hand {
	kartenNachWertSortieren(karten[:])
	paar, übrig, gefunden := ersteWertwiederholungFinden(karten[:], 2)
	if !gefunden {
		return nil
	}
	return einPaarHand{
		paar:      [2]karte(paar),
		beikarten: [3]karte(übrig),
	}
}

func (e einPaarHand) art() handArt {
	return handArtEinPaar
}

func (e einPaarHand) karten() [5]karte {
	return [5]karte{e.paar[0], e.paar[1], e.beikarten[0], e.beikarten[1], e.beikarten[2]}
}

func (e einPaarHand) displayName() string {
	return "Ein Paar: " + e.paar[0].wert.String()
}

func (e einPaarHand) visualisierung(karte karte) color.Color {
	switch karte {
	case e.paar[0], e.paar[1]:
		return colornames.Gold
	case e.beikarten[0]:
		return colornames.Silver
	case e.beikarten[1]:
		return colornames.Brown
	default:
		return nil
	}
}

type zweiPaareHand struct {
	paar1, paar2 [2]karte
	beikarte     karte
}

func parseZweiPaareHand(karten handKartenAuswahl) hand {
	kartenNachWertSortieren(karten[:])
	paar1, übrig, gefunden := ersteWertwiederholungFinden(karten[:], 2)
	if !gefunden {
		return nil
	}
	paar2, übrig, gefunden := ersteWertwiederholungFinden(übrig, 2)
	if !gefunden {
		return nil
	}
	return zweiPaareHand{
		paar1:    [2]karte(paar1),
		paar2:    [2]karte(paar2),
		beikarte: übrig[0],
	}
}

func (z zweiPaareHand) art() handArt {
	return handArtZweiPaare
}

func (z zweiPaareHand) karten() [5]karte {
	return [5]karte{z.paar1[0], z.paar1[1], z.paar2[0], z.paar2[1], z.beikarte}
}

func (z zweiPaareHand) displayName() string {
	return "Zwei Paare: " + z.paar1[0].wert.String() + " und " + z.paar2[0].String()
}

func (z zweiPaareHand) visualisierung(karte karte) color.Color {
	switch karte {
	case z.paar1[0], z.paar1[1]:
		return colornames.Gold
	case z.paar2[0], z.paar2[1]:
		return colornames.Silver
	case z.beikarte:
		return colornames.Brown
	default:
		return nil
	}
}

type drillingHand struct {
	drilling  [3]karte
	beikarten [2]karte
}

func parseDrillingsHand(karten handKartenAuswahl) hand {
	kartenNachWertSortieren(karten[:])
	drilling, übrig, gefunden := ersteWertwiederholungFinden(karten[:], 3)
	if !gefunden {
		return nil
	}
	return drillingHand{
		drilling:  [3]karte(drilling),
		beikarten: [2]karte(übrig),
	}
}

func (d drillingHand) art() handArt {
	return handArtDrilling
}

func (d drillingHand) karten() [5]karte {
	return [5]karte{d.drilling[0], d.drilling[1], d.drilling[2], d.beikarten[0], d.beikarten[1]}
}

func (d drillingHand) displayName() string {
	return "Drei gleiche Karten: " + d.drilling[0].wert.String()
}

func (d drillingHand) visualisierung(karte karte) color.Color {
	switch karte {
	case d.drilling[0], d.drilling[1], d.drilling[2]:
		return colornames.Gold
	case d.beikarten[0]:
		return colornames.Silver
	case d.beikarten[1]:
		return colornames.Brown
	default:
		return nil
	}
}

type straßeHand [5]karte

func parseStraßeHand(karten handKartenAuswahl) hand {
	kartenNachWertSortieren(karten[:])
höchsterWertIndex:
	for höchstekarteIndex := 0; höchstekarteIndex <= 1; höchstekarteIndex++ {
		höchsteKarte := karten[höchstekarteIndex]
		for offset := 1; offset <= 4; offset++ {
			if karten[höchstekarteIndex+offset].wert != höchsteKarte.wert-wert(offset) {
				continue höchsterWertIndex
			}
		}
		return straßeHand(karten[höchstekarteIndex : höchstekarteIndex+5])
	}

	return nil
}

func (s straßeHand) art() handArt {
	return handArtStraße
}

func (s straßeHand) karten() [5]karte {
	return s
}

func (s straßeHand) displayName() string {
	return "Straße"
}

func (s straßeHand) visualisierung(karte karte) color.Color {
	switch karte {
	case s[0]:
		return colornames.Gold
	case s[1], s[2], s[3], s[4]:
		return colornames.Silver
	default:
		return nil
	}
}

type flushHand [5]karte

func parseFlushHand(karten handKartenAuswahl) hand {
	kartenNachSymbol := make(map[symbol][]karte)
	for _, karte := range karten {
		kartenNachSymbol[karte.symbol] = append(kartenNachSymbol[karte.symbol], karte)
	}
	for _, kartenMitSymbol := range kartenNachSymbol {
		if len(kartenMitSymbol) != 5 {
			continue
		}

		kartenNachWertSortieren(kartenMitSymbol)

		return flushHand(kartenMitSymbol)
	}
	return nil
}

func (f flushHand) art() handArt {
	return handArtFlush
}

func (f flushHand) karten() [5]karte {
	return f
}

func (f flushHand) displayName() string {
	return "Flush: 5 mal " + f[0].symbol.String()
}

func (f flushHand) visualisierung(karte karte) color.Color {
	switch karte {
	case f[0], f[1], f[2], f[3], f[4]:
		return colornames.Gold
	default:
		return nil
	}
}

type fullHouseHand struct {
	drilling [3]karte
	paar     [2]karte
}

func parseFullHouseHand(karten handKartenAuswahl) hand {
	kartenNachWertSortieren(karten[:])
	drilling, übrig, gefunden := ersteWertwiederholungFinden(karten[:], 3)
	if !gefunden {
		return nil
	}
	paar, _, gefunden := ersteWertwiederholungFinden(übrig, 2)
	if !gefunden {
		return nil
	}
	return fullHouseHand{
		drilling: [3]karte(drilling),
		paar:     [2]karte(paar),
	}
}

func (f fullHouseHand) art() handArt {
	return handArtFullHouse
}

func (f fullHouseHand) karten() [5]karte {
	return [5]karte{f.drilling[0], f.drilling[1], f.drilling[2], f.paar[0], f.paar[1]}
}

func (f fullHouseHand) displayName() string {
	return "Full House: dreimal " + f.drilling[0].wert.String() + " und zweimal " + f.paar[0].wert.String()
}

func (f fullHouseHand) visualisierung(karte karte) color.Color {
	switch karte {
	case f.drilling[0], f.drilling[1], f.drilling[2]:
		return colornames.Gold
	case f.paar[0], f.paar[1]:
		return colornames.Silver
	default:
		return nil
	}
}

type vierlingHand struct {
	vierling [4]karte
	beikarte karte
}

func parseVierlingHand(karten handKartenAuswahl) hand {
	kartenNachWertSortieren(karten[:])
	vierling, übrig, gefunden := ersteWertwiederholungFinden(karten[:], 4)
	if !gefunden {
		return nil
	}
	return vierlingHand{
		vierling: [4]karte(vierling),
		beikarte: übrig[0],
	}
}

func (v vierlingHand) art() handArt {
	return handArtVierling
}

func (v vierlingHand) karten() [5]karte {
	return [5]karte{v.vierling[0], v.vierling[1], v.vierling[2], v.vierling[3], v.beikarte}
}

func (v vierlingHand) displayName() string {
	return "Vierling: 4 x " + v.vierling[0].wert.String()
}

func (v vierlingHand) visualisierung(karte karte) color.Color {
	switch karte {
	case v.vierling[0], v.vierling[1], v.vierling[2], v.vierling[3]:
		return colornames.Gold
	case v.beikarte:
		return colornames.Silver
	default:
		return nil
	}
}
