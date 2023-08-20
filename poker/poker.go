package poker

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/Lama06/Herder-Legacy/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

type symbol int8

const (
	symbolPik symbol = iota
	symbolHerz
	symbolKaro
	symbolKreuz
	anzahlSymbole
)

func parseSymbol(text string) (symbol symbol, ok bool) {
	for s := symbolPik; s <= symbolKreuz; s++ {
		if text == s.String() {
			return s, true
		}
	}
	return 0, false
}

func (s symbol) fileName() string {
	return strings.ToLower(s.String())
}

func (s symbol) String() string {
	switch s {
	case symbolHerz:
		return "Herz"
	case symbolKreuz:
		return "Kreuz"
	case symbolPik:
		return "Pik"
	case symbolKaro:
		return "Karo"
	default:
		panic("unreachable")
	}
}

type wert int8

const (
	wert2 wert = iota + 2
	wert3
	wert4
	wert5
	wert6
	wert7
	wert8
	wert9
	wert10
	wertBube
	wertDame
	wertKönig
	wertAss
)

func parseWert(text string) (wert wert, ok bool) {
	for w := wert2; w <= wertAss; w++ {
		if text == w.String() {
			return w, true
		}
	}
	return 0, false
}

func (w wert) fileName() string {
	return strings.ToLower(w.String())
}

func (w wert) String() string {
	if w >= wert2 && w <= wert10 {
		return strconv.Itoa(int(w))
	}

	switch w {
	case wertBube:
		return "Bube"
	case wertDame:
		return "Dame"
	case wertKönig:
		return "König"
	case wertAss:
		return "Ass"
	default:
		panic("unreachable")
	}
}

type karte struct {
	wert   wert
	symbol symbol
}

func parseKarte(text string) (ergebnis karte, ok bool) {
	words := strings.Split(text, " ")
	if len(words) != 2 {
		return karte{}, false
	}
	symbolText := words[0]
	wertText := words[1]

	symbol, ok := parseSymbol(symbolText)
	if !ok {
		return karte{}, false
	}

	wert, ok := parseWert(wertText)
	if !ok {
		return karte{}, false
	}

	return karte{wert: wert, symbol: symbol}, true
}

func mustParseKarte(text string) karte {
	ergebnis, ok := parseKarte(text)
	if !ok {
		panic("failed to parse karte")
	}
	return ergebnis
}

func (k karte) image() *ebiten.Image {
	return assets.RequireImage(fmt.Sprintf("spielkarten/%v/%v.png", k.symbol.fileName(), k.wert.fileName()))
}

func (k karte) String() string {
	return k.symbol.String() + " " + k.wert.String()
}

type kartenStapel map[karte]struct{}

var vollständigerKartenStapel kartenStapel

func init() {
	vollständigerKartenStapel = make(kartenStapel)
	for s := symbolPik; s <= symbolKreuz; s++ {
		for w := wert2; w <= wertAss; w++ {
			vollständigerKartenStapel[karte{symbol: s, wert: w}] = struct{}{}
		}
	}
}

func (k kartenStapel) clone() kartenStapel {
	neuesDeck := make(kartenStapel, len(k))
	for karte := range k {
		neuesDeck[karte] = struct{}{}
	}
	return neuesDeck
}

func (k kartenStapel) karten() []karte {
	liste := make([]karte, 0, len(k))
	for karte := range k {
		liste = append(liste, karte)
	}
	return liste
}

func (k kartenStapel) karteZiehen() karte {
	liste := k.karten()
	karte := liste[rand.Intn(len(liste))]
	delete(k, karte)
	return karte
}
