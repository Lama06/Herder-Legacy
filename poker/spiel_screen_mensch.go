package poker

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
)

func newSpielScreenEinsatzAuswahlScreen(
	herderLegacy herderlegacy.HerderLegacy,
	minEinsatz int,
	maxEinsatz int,
	callback func(einsatz int) herderlegacy.Screen,
) herderlegacy.Screen {
	var widgets []ui.ListScreenWidget
	for _, einsatzErhöhung := range [...]int{0, 1, 2, 3, 4, 5, 6, 8, 10, 12, 15, 20, 25, 30} {
		if minEinsatz+einsatzErhöhung > maxEinsatz {
			continue
		}
		einsatzErhöhung := einsatzErhöhung
		widgets = append(widgets, ui.ListScreenButtonWidget{
			Text: strconv.Itoa(minEinsatz + einsatzErhöhung),
			Callback: func() {
				herderLegacy.OpenScreen(callback(minEinsatz + einsatzErhöhung))
			},
		})
	}
	widgets = append(widgets, ui.ListScreenButtonWidget{
		Text: "ALL IN!",
		Callback: func() {
			herderLegacy.OpenScreen(callback(maxEinsatz))
		},
	})

	return ui.NewListScreen(herderLegacy, ui.ListScreenConfig{
		Title:   "Einsatz",
		Text:    "Wähle, wie viele Jetons du setzen willst:",
		Widgets: widgets,
	})
}

func spielScreenMenschKartenPosition(index int) (karteX float64, karteY float64, karteRotation float64) {
	const (
		y               = ui.Height - spielScreenKarteHöhe/2
		abstandX        = spielScreenKarteBreite / 2.0
		breiteGesamt    = abstandX + spielScreenKarteBreite
		abstandVomRandX = (ui.Width - breiteGesamt) / 2.0
		rotation        = 0
	)

	karteX = abstandVomRandX + float64(index)*abstandX
	return karteX, y, rotation
}

type spielScreenMensch struct {
	aufgegeben      bool
	jettons         int
	karten          [2]karte
	bewegendeKarten [2]*spielScreenBewegendeKarte
}

func newSpielScreenMensch(jettons int) *spielScreenMensch {
	return &spielScreenMensch{
		aufgegeben: false,
		jettons:    jettons,
	}
}

func (s *spielScreenMensch) kartenZiehen(karten [2]karte) {
	s.karten = karten
	for i := 0; i < 2; i++ {
		x, y, rotation := spielScreenMenschKartenPosition(i)
		s.bewegendeKarten[i] = &spielScreenBewegendeKarte{
			karte:          karten[i],
			targetRotation: rotation,
			currentX:       spielScreenStapelX,
			currentY:       spielScreenStapelY,
			targetX:        x,
			targetY:        y,
			autoAufdecken:  true,
		}
	}
}

func (s *spielScreenMensch) getKarten() [2]karte {
	return s.karten
}

func (s *spielScreenMensch) getBewegendeKarten() [2]*spielScreenBewegendeKarte {
	return s.bewegendeKarten
}

func (s *spielScreenMensch) jettonSpawnPunkt() (float64, float64) {
	offset := rand.Float64() * 200
	return ui.Width / 2, ui.Height + offset
}

func (s *spielScreenMensch) setAufgegeben(aufgegeben bool) {
	s.aufgegeben = aufgegeben
}

func (s *spielScreenMensch) hatAufgegeben() bool {
	return s.aufgegeben
}

func (s *spielScreenMensch) einsatzErmitteln(
	herderLegacy herderlegacy.HerderLegacy,
	status spielScreenStatus,
	wirdGewinnen bool,
	callback func(einsatz int),
) {
	previousScreen := herderLegacy.CurrentScreen()
	herderLegacy.OpenScreen(newSpielScreenEinsatzAuswahlScreen(
		herderLegacy,
		0,
		s.jettons,
		func(einsatz int) herderlegacy.Screen {
			s.jettons -= einsatz
			callback(einsatz)
			return previousScreen
		},
	))
}

func (s *spielScreenMensch) gehtMit(
	herderLegacy herderlegacy.HerderLegacy,
	einsatz int,
	callback func(gehtMit bool),
) {
	previousScreen := herderLegacy.CurrentScreen()
	herderLegacy.OpenScreen(ui.NewListScreen(herderLegacy, ui.ListScreenConfig{
		Title: "Mitgehen?",
		Text:  fmt.Sprintf("Es wurden %v Jetons gesetzt. Willst du mitgehen?", einsatz),
		Widgets: []ui.ListScreenWidget{
			ui.ListScreenButtonWidget{
				Text: "Rausgehen",
				Callback: func() {
					callback(false)
					herderLegacy.OpenScreen(previousScreen)
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Mitgehen",
				Callback: func() {
					s.jettons = max(0, s.jettons-einsatz)
					callback(true)
					herderLegacy.OpenScreen(previousScreen)
				},
			},
		},
	}))
}
