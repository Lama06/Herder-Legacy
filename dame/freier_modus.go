package dame

import (
	"fmt"
	"strings"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
)

func NewFreierModusScreen(herderLegacy herderlegacy.HerderLegacy, dameBeendenCallback func() herderlegacy.Screen) herderlegacy.Screen {
	return ui.NewListScreen(herderLegacy, ui.ListScreenConfig{
		Title:        "Dame",
		CancelText:   "Dame beenden",
		CancelAction: dameBeendenCallback,
		Widgets: []ui.ListScreenWidget{
			ui.ListScreenButtonWidget{
				Text: "Info",
				Callback: func() {
					herderLegacy.OpenScreen(ui.NewMessageScreen(herderLegacy, ui.MessageScreenConfig{
						Title: "Dame",
						Text: `Wenn du gegen einen Lehrer in Dame gewinnst, wird dieser ein wenig Motivation, zu unterrichten, verlieren.
Du kannst auswählen gegen welchen Lehrer du antreten willst. 
Beachte aber, dass jeder mit seinen eigenen Regeln und unterschiedlicher Stragie spielt.
Hinweis: Teilweise kann es einige Sekunden dauern, um dem Zug des Lehrers zu berechnen.`,
						ContinueText: "Los gehts!",
						ContinueAction: func() herderlegacy.Screen {
							return NewFreierModusScreen(herderLegacy, dameBeendenCallback)
						},
					}))
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Gegen Lehrer Spielen",
				Callback: func() {
					herderLegacy.OpenScreen(newLehrerAuswahlScreen(herderLegacy, dameBeendenCallback))
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Mit eigenen Regeln spielen",
				Callback: func() {
					herderLegacy.OpenScreen(newDameEditorScreen(herderLegacy, dameBeendenCallback))
				},
			},
		},
	})
}

func newLehrerAuswahlScreen(
	herderLegacy herderlegacy.HerderLegacy,
	dameBeendenCallback func() herderlegacy.Screen,
) herderlegacy.Screen {
	widgets := make([]ui.ListScreenWidget, len(alleLehrer))
	for i, lehrer := range alleLehrer {
		lehrer := lehrer
		widgets[i] = ui.ListScreenButtonWidget{
			Text: lehrer.name,
			Callback: func() {
				herderLegacy.OpenScreen(newLehrerInfoScreen(herderLegacy, dameBeendenCallback, lehrer))
			},
		}
	}

	return ui.NewListScreen(herderLegacy, ui.ListScreenConfig{
		Title: "Lehrer auswählen",
		Text:  "Wähle einen Lehrer, gegen den du spielen willst, aus",

		CancelText: "Zurück",
		CancelAction: func() herderlegacy.Screen {
			return NewFreierModusScreen(herderLegacy, dameBeendenCallback)
		},

		Widgets: widgets,
	})
}

func newLehrerInfoScreen(
	herderLegacy herderlegacy.HerderLegacy,
	dameBeendenCallback func() herderlegacy.Screen,
	lehrer lehrer,
) herderlegacy.Screen {
	return ui.NewDecideScreen(herderLegacy, ui.DecideScreenConfig{
		Title: lehrer.name,
		Text:  lehrer.info,

		CancelText: "Gegen anderen Lehrer spielen",
		CancelAction: func() herderlegacy.Screen {
			return newLehrerAuswahlScreen(herderLegacy, dameBeendenCallback)
		},

		ConfirmText: fmt.Sprintf("Gegen %v spielen", lehrer.nameAkkusativOrDefault()),
		ConfirmAction: func() herderlegacy.Screen {
			return NewLehrerDameSpielScreen(herderLegacy, func(gewonnen bool) herderlegacy.Screen {
				return newGameOverScreen(herderLegacy, dameBeendenCallback, lehrer, gewonnen)
			}, lehrer.spielOptionen)
		},
	})
}

func newGameOverScreen(
	herderLegacy herderlegacy.HerderLegacy,
	dameBeendenCallback func() herderlegacy.Screen,
	lehrer lehrer,
	gewonnenn bool,
) herderlegacy.Screen {
	if gewonnenn {
		herderLegacy.AddVerhinderteStunden(3)
	}

	var title string
	if gewonnenn {
		title = "Gewonnen"
	} else {
		title = "Verloren"
	}

	var text string
	if gewonnenn {
		text = `Du hast gegen %nameAkk% gewonnen.
%pronomenSatzanfang% hat damit nicht gerechnet, weil %pronomenSatzmitte% ja die Regeln festgelegt hat.
Jetzt ist %name% schlecht gelaunt und hat auch weniger Motivation, zu unterrichten.
Durch deinen Sieg sind die Sommerferien 3 Stunden nach vorne gerutscht!`
	} else {
		text = `Du hast gegen %nameAkk% verloren.
Versuche es nocheinmal und gewinne, damit %pronomenSatzmitte% weniger motiviert ist, zu unterrichten.`
	}
	infoTextReplacer := strings.NewReplacer(
		"%name%", lehrer.name,
		"%nameAkk%", lehrer.nameAkkusativOrDefault(),
		"%pronomenSatzanfang%", lehrer.personalPronomenSatzanfang(),
		"%pronomenSatzmitte%", lehrer.personalPronomenSatzmitte(),
	)
	text = infoTextReplacer.Replace(text)

	return ui.NewDecideScreen(herderLegacy, ui.DecideScreenConfig{
		Title: title,
		Text:  text,

		CancelText:   "Dame beenden",
		CancelAction: dameBeendenCallback,

		ConfirmText: "Eine weitere Runde Dame spielen",
		ConfirmAction: func() herderlegacy.Screen {
			return newLehrerInfoScreen(herderLegacy, dameBeendenCallback, lehrer)
		},
	})
}
