package poker

import (
	"fmt"
	"strings"
	"time"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
)

func NewPokerRechnerScreen(
	herderLegacy herderlegacy.HerderLegacy,
	pokerRechnerBeendetCallback func() herderlegacy.Screen,
) herderlegacy.Screen {
	var (
		eigeneKarten    []karte
		mittelKarten    []karte
		abgelegteKarten []karte
	)

	return ui.NewListScreen(herderLegacy, ui.ListScreenConfig{
		Title: "Poker Rechner",
		Text:  "Berechne Wahrscheinlichkeiten für mögliche Poker Hände",
		Widgets: []ui.ListScreenWidget{
			ui.ListScreenButtonWidget{
				Text: "Eigene Karten bearbeiten",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()
					herderLegacy.OpenScreen(newKartenAuswahlScreen(
						herderLegacy,
						"Eigene Karten",
						eigeneKarten,
						func(neueEigeneKarten []karte) herderlegacy.Screen {
							eigeneKarten = neueEigeneKarten
							return thisScreen
						},
					))
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Mittelkarten bearbeiten",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()
					herderLegacy.OpenScreen(newKartenAuswahlScreen(
						herderLegacy,
						"Mittelkarten",
						mittelKarten,
						func(neueMittelkarten []karte) herderlegacy.Screen {
							mittelKarten = neueMittelkarten
							return thisScreen
						},
					))
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Abgelegte Karten bearbeiten",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()
					herderLegacy.OpenScreen(newKartenAuswahlScreen(
						herderLegacy,
						"Abgelegte Karten",
						abgelegteKarten,
						func(neueAbgelegteKarten []karte) herderlegacy.Screen {
							abgelegteKarten = neueAbgelegteKarten
							return thisScreen
						},
					))
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Chancen berechnen",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()

					if len(eigeneKarten) != 2 || !(len(mittelKarten) >= 3 && len(mittelKarten) <= 5) {
						herderLegacy.OpenScreen(ui.NewMessageScreen(herderLegacy, ui.MessageScreenConfig{
							Title:        "Fehler: Ungültige Kartenanzahlen",
							Text:         "Du musst genau 2 eigene Karten und zwischen 3 und 5 Mittelkarten angeben.",
							ContinueText: "Karten bearbeiten",
							ContinueAction: func() herderlegacy.Screen {
								return thisScreen
							},
						}))
						return
					}

					herderLegacy.OpenScreen(ui.NewMessageScreen(herderLegacy, ui.MessageScreenConfig{
						Title:        "Chancen berechnen",
						Text:         "Hinweis: Das Berechnen der Chancen kann etwas dauern.",
						ContinueText: "Chancen berechnen",
						ContinueAction: func() herderlegacy.Screen {
							start := time.Now()

							eigeneHandMöglichkeiten :=
								eigeneHandArtenMöglichkeitenBerechnen([2]karte(eigeneKarten), mittelKarten, abgelegteKarten)
							eigeneHandWahrscheinlichkeiten := handArtenWahrscheinlichkeitenBerechnen(eigeneHandMöglichkeiten)

							gegnerHandMöglichkeiten := gegnerHandArtenMöglichkeitenBerechnen([2]karte(eigeneKarten), mittelKarten, abgelegteKarten)
							gegnerHandWahrscheinlichkeiten := handArtenWahrscheinlichkeitenBerechnen(gegnerHandMöglichkeiten)

							benötigteZeit := time.Now().Sub(start)

							var text strings.Builder
							text.WriteString("Du hast folgende Chancen auf die Poker Kombinationen:")
							for handArt := handArtHöchsteKarte; handArt <= handArtRoyalFlush; handArt++ {
								text.WriteString(fmt.Sprintf("\n%v: %v%% (%v Möglichkeiten)",
									handArt, eigeneHandWahrscheinlichkeiten[handArt]*100, eigeneHandMöglichkeiten[handArt]))
							}
							text.WriteString("\n\nDeine Gegner haben (aus deiner Sicht) folgende Chancen auf die Poker Kombinationen:")
							for handArt := handArtHöchsteKarte; handArt <= handArtRoyalFlush; handArt++ {
								text.WriteString(fmt.Sprintf("\n%v: %v%% (%v Möglichkeiten)",
									handArt, gegnerHandWahrscheinlichkeiten[handArt]*100, gegnerHandMöglichkeiten[handArt]))
							}
							text.WriteString(fmt.Sprintf("\n\nBenötigte Sekunden zum Berechnen: %v", benötigteZeit.Seconds()))

							return ui.NewMessageScreen(herderLegacy, ui.MessageScreenConfig{
								Title:        "Chancen",
								Text:         text.String(),
								ContinueText: "Weiter",
								ContinueAction: func() herderlegacy.Screen {
									return thisScreen
								},
							})
						},
					}))
				},
			},
		},
		CancelText:   "Abbrechen",
		CancelAction: pokerRechnerBeendetCallback,
	})
}

func newKartenAuswahlScreen(
	herderLegacy herderlegacy.HerderLegacy,
	title string,
	karten []karte,
	nächsterScreen func([]karte) herderlegacy.Screen,
) herderlegacy.Screen {
	widgets := make([]ui.ListScreenWidget, len(karten))
	for i, karte := range karten {
		widgets[i] = ui.ListScreenButtonWidget{
			Text: karte.String(),
			Callback: func() {
				karten = append(karten[:i], karten[i+1:]...)
				herderLegacy.OpenScreen(newKartenAuswahlScreen(
					herderLegacy,
					title,
					karten,
					nächsterScreen,
				))
			},
		}
	}

	widgets = append(widgets, ui.ListScreenButtonWidget{
		Text: "Hinzufügen",
		Callback: func() {
			ausgewähltesSymbol := symbolKreuz
			ausgewählterWert := wert2

			herderLegacy.OpenScreen(ui.NewListScreen(herderLegacy, ui.ListScreenConfig{
				Title: "Karte hinzufügen",
				Text:  "",
				Widgets: []ui.ListScreenWidget{
					ui.ListScreenSelectionWidget[symbol]{
						Text:   "Symbol",
						Value:  ausgewähltesSymbol,
						Values: []symbol{symbolKreuz, symbolKaro, symbolPik, symbolHerz},
						Callback: func(neuesSymbol symbol) {
							ausgewähltesSymbol = neuesSymbol
						},
					},
					ui.ListScreenSelectionWidget[wert]{
						Text:  "Wert",
						Value: ausgewählterWert,
						Values: []wert{wert2, wert3, wert4, wert5, wert6, wert7, wert8, wert9, wert10,
							wertBube, wertDame, wertKönig, wertAss},
						Callback: func(neuerWert wert) {
							ausgewählterWert = neuerWert
						},
					},
				},
				CancelText: "Hinzufügen",
				CancelAction: func() herderlegacy.Screen {
					karten = append(karten, karte{symbol: ausgewähltesSymbol, wert: ausgewählterWert})
					return newKartenAuswahlScreen(
						herderLegacy,
						title,
						karten,
						nächsterScreen,
					)
				},
			}))
		},
	})

	return ui.NewListScreen(herderLegacy, ui.ListScreenConfig{
		Title:      title,
		Widgets:    widgets,
		CancelText: "Speichern",
		CancelAction: func() herderlegacy.Screen {
			return nächsterScreen(karten)
		},
	})
}

func eigeneHandArtenMöglichkeitenBerechnen(
	eigeneKarten [2]karte,
	mittelkarten []karte,
	abgelegteKarten []karte,
) map[handArt]int {
	stapel := vollständigerKartenStapel.clone()
	delete(stapel, eigeneKarten[0])
	delete(stapel, eigeneKarten[1])
	for _, mittelkarte := range mittelkarten {
		delete(stapel, mittelkarte)
	}
	for _, abgelegteKarte := range abgelegteKarten {
		delete(stapel, abgelegteKarte)
	}

	anzahlenMöglicherHandArten := make(map[handArt]int)
	switch len(mittelkarten) {
	case 5:
		verfügbareKarten := [7]karte{eigeneKarten[0], eigeneKarten[1],
			mittelkarten[0], mittelkarten[1], mittelkarten[2], mittelkarten[3], mittelkarten[4]}
		anzahlenMöglicherHandArten[parseHand(verfügbareKarten).art()] = 1
	case 4:
		for _, mittelkarte5 := range stapel.karten() {
			verfügbareKarten := [7]karte{eigeneKarten[0], eigeneKarten[1],
				mittelkarten[0], mittelkarten[1], mittelkarten[2], mittelkarten[3], mittelkarte5}
			anzahlenMöglicherHandArten[parseHand(verfügbareKarten).art()]++
		}
	case 3:
		for _, mittelkarte4 := range stapel.karten() {
			delete(stapel, mittelkarte4)
			for _, mittelkarte5 := range stapel.karten() {
				verfügbareKarten := [7]karte{eigeneKarten[0], eigeneKarten[1],
					mittelkarten[0], mittelkarten[1], mittelkarten[2], mittelkarte4, mittelkarte5}
				anzahlenMöglicherHandArten[parseHand(verfügbareKarten).art()]++
			}
			stapel[mittelkarte4] = struct{}{}
		}
	default:
		panic("Es gibt immer 3, 4 oder 5 Mittelkarten")
	}
	return anzahlenMöglicherHandArten
}

func gegnerHandArtenMöglichkeitenBerechnen(
	eigeneKarten [2]karte,
	mittelkarten []karte,
	abgelegteKarten []karte,
) map[handArt]int {
	stapel := vollständigerKartenStapel.clone()
	delete(stapel, eigeneKarten[0])
	delete(stapel, eigeneKarten[1])
	for _, mittelkarte := range mittelkarten {
		delete(stapel, mittelkarte)
	}
	for _, abgelegteKarte := range abgelegteKarten {
		delete(stapel, abgelegteKarte)
	}

	anzahlenMöglicherHandArten := make(map[handArt]int)
	for _, gegnerKarte1 := range stapel.karten() {
		delete(stapel, gegnerKarte1)
		for _, gegnerKarte2 := range stapel.karten() {
			for handArt, anzahlMöglichkeiten := range eigeneHandArtenMöglichkeitenBerechnen(
				[2]karte{gegnerKarte1, gegnerKarte2},
				mittelkarten,
				append(abgelegteKarten, eigeneKarten[0], eigeneKarten[1]),
			) {
				anzahlenMöglicherHandArten[handArt] += anzahlMöglichkeiten
			}
		}
		stapel[gegnerKarte1] = struct{}{}
	}
	return anzahlenMöglicherHandArten
}

func handArtenWahrscheinlichkeitenBerechnen(
	anzahlenMöglicherHandArten map[handArt]int,
) map[handArt]float64 {
	var anzahlMöglichkeitenGesamt int
	for _, anzahlMöglichkeitenFürHandArt := range anzahlenMöglicherHandArten {
		anzahlMöglichkeitenGesamt += anzahlMöglichkeitenFürHandArt
	}

	wahrscheinlichkeiten := make(map[handArt]float64, len(anzahlenMöglicherHandArten))
	for handArt, anzahl := range anzahlenMöglicherHandArten {
		wahrscheinlichkeiten[handArt] = float64(anzahl) / float64(anzahlMöglichkeitenGesamt)
	}

	return wahrscheinlichkeiten
}
