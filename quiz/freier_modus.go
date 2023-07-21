package quiz

import (
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
)

func NewFreierModusScreen(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func() herderlegacy.Screen,
) herderlegacy.Screen {
	sekundenProFrage := 10.0

	return ui.NewListScreen(herderLegacy, ui.ListScreenConfig{
		Title:        "Quiz",
		Text:         "Wähle ein Quiz aus",
		CancelText:   "Schließen",
		CancelAction: nächsterScreen,
		Widgets: []ui.ListScreenWidget{
			ui.ListScreenButtonWidget{
				Text: "Hauptstädte der Nachbarländer Deutschlands",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()

					herderLegacy.OpenScreen(NewMultipleChoiceQuizScreen(
						herderLegacy,
						NewHauptstädtDerNachbarländerDeutschlandsQuizConfig(int(sekundenProFrage*60)),
						func(auswertung MultipleChoiceQuizAuswertung) herderlegacy.Screen {
							return thisScreen
						},
					))
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Hauptstädte Europas",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()

					herderLegacy.OpenScreen(NewMultipleChoiceQuizScreen(
						herderLegacy,
						NewHauptstädteEuropasQuizConfig(int(sekundenProFrage*60)),
						func(auswertung MultipleChoiceQuizAuswertung) herderlegacy.Screen {
							return thisScreen
						},
					))
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Hauptstädte Weltweit",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()

					herderLegacy.OpenScreen(NewMultipleChoiceQuizScreen(
						herderLegacy,
						NewHauptstädteInternationalQuizConfig(int(sekundenProFrage*60)),
						func(auswertung MultipleChoiceQuizAuswertung) herderlegacy.Screen {
							return thisScreen
						},
					))
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Einwohnerzahlen",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()

					herderLegacy.OpenScreen(NewRelationsQuizScreen(
						herderLegacy,
						NewEinwohnerQuizConfig(int(sekundenProFrage*60)),
						func(auswertung RelationQuizAuswertung) herderlegacy.Screen {
							return thisScreen
						},
					))
				},
			},
			ui.ListScreenSelectionWidget[float64]{
				Text:   "Antwortzeit pro Frage",
				Value:  sekundenProFrage,
				Values: []float64{1, 1.5, 2, 2.5, 3, 4, 6, 8, 10, 15, 20},
				Callback: func(neueAntwortzeit float64) {
					sekundenProFrage = neueAntwortzeit
				},
			},
		},
	})
}
