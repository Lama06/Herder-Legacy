package dame

import (
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func newDameEditorScreen(
	herderLegacy herderlegacy.HerderLegacy,
	dameBeendenCallback func() herderlegacy.Screen,
) herderlegacy.Screen {
	spielOptionen := SpielOptionen{
		StartBrett: MustParseBrett(
			"_l_l_l_l",
			"l_l_l_l_",
			"_l_l_l_l",
			"________",
			"________",
			"s_s_s_s_",
			"_s_s_s_s",
			"s_s_s_s_",
		),
		AiTiefe:   3,
		ZugRegeln: InternationaleZugRegeln.clone(),
	}

	return ui.NewListScreen(herderLegacy, ui.ListScreenConfig{
		Title:      "Dame Editor",
		Text:       "Spiele Dame mit deinen eigenen Regeln!",
		CancelText: "Zurück zum Hauptmenü",
		CancelAction: func() herderlegacy.Screen {
			return NewFreierModusScreen(herderLegacy, dameBeendenCallback)
		},
		Widgets: []ui.ListScreenWidget{
			ui.ListScreenButtonWidget{
				Text: "Brett bearbeiten",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()

					herderLegacy.OpenScreen(newBrettEditorScreen(
						herderLegacy,
						spielOptionen.StartBrett.clone(),
						func(brett Brett) herderlegacy.Screen {
							spielOptionen.StartBrett = brett
							return thisScreen
						},
					))
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Zug Regeln bearbeiten",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()

					herderLegacy.OpenScreen(newZugRegelnEditorScreen(
						herderLegacy,
						spielOptionen.ZugRegeln.clone(),
						func(zugRegeln ZugRegeln) herderlegacy.Screen {
							spielOptionen.ZugRegeln = zugRegeln
							return thisScreen
						},
					))
				},
			},
			ui.ListScreenSelectionWidget[int]{
				Text:   "Ai Stärke",
				Value:  spielOptionen.AiTiefe,
				Values: []int{1, 2, 3, 4, 5, 6, 7},
				Callback: func(aiTiefe int) {
					spielOptionen.AiTiefe = aiTiefe
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Spielen",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()

					herderLegacy.OpenScreen(NewLehrerDameSpielScreen(
						herderLegacy,
						func(gewonnen bool) herderlegacy.Screen {
							var text string
							if gewonnen {
								text = "Gewonnen"
							} else {
								text = "Verloren"
							}
							return ui.NewDecideScreen(herderLegacy, ui.DecideScreenConfig{
								Title:      text,
								CancelText: "Beenden",
								CancelAction: func() herderlegacy.Screen {
									return NewFreierModusScreen(herderLegacy, dameBeendenCallback)
								},
								ConfirmText: "Noch eine Runde",
								ConfirmAction: func() herderlegacy.Screen {
									return thisScreen
								},
							})
						},
						spielOptionen,
					))
				},
			},
		},
	})
}

const (
	brettEditorScreenBrettX         = 0
	brettEditorScreenBrettY         = 100
	brettEditorScreenBrettMaxBreite = ui.Width
	brettEditorScreenBrettMaxHöhe   = ui.Height - brettEditorScreenBrettY
)

type brettEditorScreen struct {
	zurückKnopf *ui.Button
	sizeÄndern  *ui.Button
	feldAuswahl *ui.Selection[feld]

	brett Brett
}

func newBrettEditorScreen(
	herderLegacy herderlegacy.HerderLegacy,
	brett Brett,
	nächsterScreen func(brett Brett) herderlegacy.Screen,
) *brettEditorScreen {
	screen := brettEditorScreen{
		feldAuswahl: ui.NewSelection(herderLegacy, ui.SelectionConfig[feld]{
			Position: ui.Position{
				X:                ui.Width - 20,
				Y:                20,
				AnchorHorizontal: ui.HorizontalerAnchorRechts,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text:     "Setzen",
			Value:    feldLeer,
			Values:   []feld{feldLeer, feldSteinSchüler, feldDameSchüler, feldSteinLehrer, feldDameLehrer},
			ToString: feld.anzeigeName,
		}),

		brett: brett,
	}

	screen.zurückKnopf = ui.NewButton(ui.ButtonConfig{
		Position: ui.Position{
			X:                20,
			Y:                20,
			AnchorHorizontal: ui.HorizontalerAnchorLinks,
			AnchorVertikal:   ui.VertikalerAnchorOben,
		},
		Text:               "Speichern",
		CustomColorPalette: true,
		ColorPalette:       ui.CancelButtonColorPalette,
		Callback: func() {
			herderLegacy.OpenScreen(nächsterScreen(screen.brett))
		},
	})

	screen.sizeÄndern = ui.NewButton(ui.ButtonConfig{
		Position: ui.Position{
			X:                ui.Width / 2,
			Y:                20,
			AnchorHorizontal: ui.HorizontalerAnchorMitte,
			AnchorVertikal:   ui.VertikalerAnchorOben,
		},
		Text: "Größe ändern",
		Callback: func() {
			thisScreen := herderLegacy.CurrentScreen()

			herderLegacy.OpenScreen(newBrettGrößeAuswahlScreen(
				herderLegacy,
				screen.brett.zeilen,
				screen.brett.spalten,
				func(zeilen, spalten int) herderlegacy.Screen {
					screen.brett = newLeeresBrett(zeilen, spalten)
					return thisScreen
				},
			))
		},
	})

	return &screen
}

func (b *brettEditorScreen) handleClick(clickX, clickY int) {
	clickedPosition, ok := b.brett.screenPositionToBrettPosition(
		float64(clickX),
		float64(clickY),
		brettEditorScreenBrettX,
		brettEditorScreenBrettY,
		brettEditorScreenBrettMaxBreite,
		brettEditorScreenBrettMaxHöhe,
	)
	if !ok {
		return
	}

	b.brett.setFeld(clickedPosition, b.feldAuswahl.Value())
}

func (b *brettEditorScreen) components() []ui.Component {
	return []ui.Component{b.zurückKnopf, b.sizeÄndern, b.feldAuswahl}
}

func (b *brettEditorScreen) Update() {
	for _, component := range b.components() {
		component.Update()
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		b.handleClick(ebiten.CursorPosition())
	} else if touchIds := inpututil.AppendJustReleasedTouchIDs(nil); len(touchIds) == 1 {
		b.handleClick(inpututil.TouchPositionInPreviousTick(touchIds[0]))
	}
}

func (b *brettEditorScreen) Draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)

	for _, component := range b.components() {
		component.Draw(screen)
	}

	b.brett.draw(
		screen,
		brettEditorScreenBrettX,
		brettEditorScreenBrettY,
		brettEditorScreenBrettMaxBreite,
		brettEditorScreenBrettMaxHöhe,
		false,
		position{},
		ZugRegeln{},
	)
}

func newBrettGrößeAuswahlScreen(
	herderLegacy herderlegacy.HerderLegacy,
	zeilenVorher, spaltenVorher int,
	nächsterScreen func(zeilen, spalten int) herderlegacy.Screen,
) herderlegacy.Screen {
	möglicheGrößen := []int{5, 6, 7, 8, 9, 10, 12, 16, 32, 64, 128}
	zeilen, spalten := zeilenVorher, spaltenVorher

	return ui.NewListScreen(herderLegacy, ui.ListScreenConfig{
		Title:      "Größe ändern",
		CancelText: "Speichern",
		CancelAction: func() herderlegacy.Screen {
			return nächsterScreen(zeilen, spalten)
		},
		Widgets: []ui.ListScreenWidget{
			ui.ListScreenSelectionWidget[int]{
				Text:   "Zeilen",
				Value:  zeilen,
				Values: möglicheGrößen,
				Callback: func(neueZeilen int) {
					zeilen = neueZeilen
				},
			},
			ui.ListScreenSelectionWidget[int]{
				Text:   "Spalten",
				Value:  spalten,
				Values: möglicheGrößen,
				Callback: func(neueSpalten int) {
					spalten = neueSpalten
				},
			},
		},
	})
}

func newZugRegelnEditorScreen(
	herderLegacy herderlegacy.HerderLegacy,
	zugRegeln ZugRegeln,
	nächsterScreen func(ZugRegeln) herderlegacy.Screen,
) herderlegacy.Screen {
	return ui.NewListScreen(herderLegacy, ui.ListScreenConfig{
		Title:      "Zug Regeln",
		Text:       "Stelle ein in welche Richtungen man mit Steinen bzw. Damen Züge ausführen darf",
		CancelText: "Speichern",
		CancelAction: func() herderlegacy.Screen {
			return nächsterScreen(zugRegeln)
		},
		Widgets: []ui.ListScreenWidget{
			ui.ListScreenToggleWidget{
				Text:    "Schlagzwang",
				Enabled: true,
				Callback: func(schlagZwang bool) {
					zugRegeln.SchlagZwang = schlagZwang
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Normaler Stein",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()

					herderLegacy.OpenScreen(newFigurZugRegelnScreen(
						herderLegacy,
						"Normaler Stein",
						zugRegeln.SteinBewegenRichtungen.clone(),
						zugRegeln.SteinSchlagenRichtungenAnfang.clone(),
						zugRegeln.SteinSchlagenRichtungenWeiterschlagen.clone(),
						func(
							neueBewegenRichtungen, neueSchlagenRichtungen, neueWeiterschlagenRichtungen Richtungen,
						) herderlegacy.Screen {
							zugRegeln.SteinBewegenRichtungen = neueBewegenRichtungen
							zugRegeln.SteinSchlagenRichtungenAnfang = neueSchlagenRichtungen
							zugRegeln.SteinSchlagenRichtungenWeiterschlagen = neueWeiterschlagenRichtungen
							return thisScreen
						},
					))
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Dame",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()

					herderLegacy.OpenScreen(newFigurZugRegelnScreen(
						herderLegacy,
						"Dame",
						zugRegeln.DameBewegenRichtungen.clone(),
						zugRegeln.DameSchlagenRichtungenAnfang.clone(),
						zugRegeln.DameSchlagenRichtungenWeiterschlagen.clone(),
						func(
							neueBewegenRichtungen, neueSchlagenRichtungen, neueWeiterschlagenRichtungen Richtungen,
						) herderlegacy.Screen {
							zugRegeln.DameBewegenRichtungen = neueBewegenRichtungen
							zugRegeln.DameSchlagenRichtungenAnfang = neueSchlagenRichtungen
							zugRegeln.DameSchlagenRichtungenWeiterschlagen = neueWeiterschlagenRichtungen
							return thisScreen
						},
					))
				},
			},
		},
	})
}

func newFigurZugRegelnScreen(
	herderLegacy herderlegacy.HerderLegacy,
	name string,
	bewegenRichtungen Richtungen,
	schlagenRichtungen Richtungen,
	weiterschlagenRichtungen Richtungen,
	nächsterScreen func(
		bewegenRichtungen Richtungen,
		schlagenRichtungen Richtungen,
		weiterschlagenRichtungen Richtungen,
	) herderlegacy.Screen,
) herderlegacy.Screen {
	return ui.NewListScreen(herderLegacy, ui.ListScreenConfig{
		Title:      name,
		CancelText: "Speichern",
		CancelAction: func() herderlegacy.Screen {
			return nächsterScreen(bewegenRichtungen, schlagenRichtungen, weiterschlagenRichtungen)
		},
		Widgets: []ui.ListScreenWidget{
			ui.ListScreenButtonWidget{
				Text: "Bewegen Richtungen",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()

					herderLegacy.OpenScreen(newRichtungenEditorScreen(
						herderLegacy,
						"Bewegen Richtungen",
						bewegenRichtungen.clone(),
						func(neueBewegenRichtungen Richtungen) herderlegacy.Screen {
							bewegenRichtungen = neueBewegenRichtungen
							return thisScreen
						},
					))
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Schlagen Richtungen am Anfang",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()

					herderLegacy.OpenScreen(newRichtungenEditorScreen(
						herderLegacy,
						"Schlagen Richtungen am Anfang",
						schlagenRichtungen.clone(),
						func(neueSchlagenRichtungen Richtungen) herderlegacy.Screen {
							schlagenRichtungen = neueSchlagenRichtungen
							return thisScreen
						},
					))
				},
			},
			ui.ListScreenButtonWidget{
				Text: "Schlagen Richtungen bei Weiterschlagen",
				Callback: func() {
					thisScreen := herderLegacy.CurrentScreen()

					herderLegacy.OpenScreen(newRichtungenEditorScreen(
						herderLegacy,
						"Schlagen Richtungen bei Weiterschlagen",
						weiterschlagenRichtungen.clone(),
						func(neueWeiterSchlagenRichtungen Richtungen) herderlegacy.Screen {
							weiterschlagenRichtungen = neueWeiterSchlagenRichtungen
							return thisScreen
						},
					))
				},
			},
		},
	})
}

func newRichtungenEditorScreen(
	herderLegacy herderlegacy.HerderLegacy,
	name string,
	richtungen Richtungen,
	nächsterScreen func(Richtungen) herderlegacy.Screen,
) herderlegacy.Screen {
	return ui.NewListScreen(herderLegacy, ui.ListScreenConfig{
		Title:      name,
		CancelText: "Speichern",
		CancelAction: func() herderlegacy.Screen {
			return nächsterScreen(richtungen)
		},
		Widgets: []ui.ListScreenWidget{
			ui.ListScreenToggleWidget{
				Text:    "Vorne Links",
				Enabled: richtungen.contains(RichtungVorneLinks),
				Callback: func(enabled bool) {
					richtungen.set(RichtungVorneLinks, enabled)
				},
			},
			ui.ListScreenToggleWidget{
				Text:    "Vorne",
				Enabled: richtungen.contains(RichtungVorne),
				Callback: func(enabled bool) {
					richtungen.set(RichtungVorne, enabled)
				},
			},
			ui.ListScreenToggleWidget{
				Text:    "Vorne Rechts",
				Enabled: richtungen.contains(RichtungVorneRechts),
				Callback: func(enabled bool) {
					richtungen.set(RichtungVorneRechts, enabled)
				},
			},

			ui.ListScreenToggleWidget{
				Text:    "Links",
				Enabled: richtungen.contains(RichtungLinks),
				Callback: func(enabled bool) {
					richtungen.set(RichtungLinks, enabled)
				},
			},
			ui.ListScreenToggleWidget{
				Text:    "Rechts",
				Enabled: richtungen.contains(RichtungRechts),
				Callback: func(enabled bool) {
					richtungen.set(RichtungRechts, enabled)
				},
			},

			ui.ListScreenToggleWidget{
				Text:    "Hinten Links",
				Enabled: richtungen.contains(RichtungHintenLinks),
				Callback: func(enabled bool) {
					richtungen.set(RichtungHintenLinks, enabled)
				},
			},
			ui.ListScreenToggleWidget{
				Text:    "Hinten",
				Enabled: richtungen.contains(RichtungHinten),
				Callback: func(enabled bool) {
					richtungen.set(RichtungHinten, enabled)
				},
			},
			ui.ListScreenToggleWidget{
				Text:    "Hinten Rechts",
				Enabled: richtungen.contains(RichtungHintenRechts),
				Callback: func(enabled bool) {
					richtungen.set(RichtungHintenRechts, enabled)
				},
			},
		},
	})
}
