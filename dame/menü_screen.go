package dame

import (
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type menüScreen struct {
	herderLegacy   herderlegacy.HerderLegacy
	nächsterScreen func() herderlegacy.Screen
	title          *ui.Title
	beschreibung   *ui.Text
	spielenKnopf   *ui.Button
	beendenKnopf   *ui.Button
}

var _ herderlegacy.Screen = (*menüScreen)(nil)

func newMenüScreen(herderLegacy herderlegacy.HerderLegacy, nächsterScreen func() herderlegacy.Screen) *menüScreen {
	screen := menüScreen{
		herderLegacy:   herderLegacy,
		nächsterScreen: nächsterScreen,
		title: ui.NewTitle(ui.TitleConfig{
			Position: ui.NewCenteredPosition(ui.Width/2, 100),
			Text:     "Dame",
		}),
		beschreibung: ui.NewText(ui.TextConfig{
			Position: ui.NewCenteredPosition(ui.Width/2, 200),
			Text: `Wenn du gegen einen Lehrer in Dame gewinnst, wird dieser ein wenig Motivation, zu unterrichten, verlieren.
Du kannst auswählen gegen welchen Lehrer du antreten willst. 
Beachte aber, dass jeder mit seinen eigenen Regeln und unterschiedlicher Stragie spielt.
Hinweis: Teilweise kann es einige Sekunden dauern, um dem Zug des Lehrers zu berechnen.`,
		}),
		spielenKnopf: ui.NewButton(ui.ButtonConfig{
			Position: ui.NewCenteredPosition(ui.Width/2, ui.Height-100),
			Text:     "Lehrer auswählen",
			Callback: func() {
				herderLegacy.OpenScreen(newLehrerAuswahlScreen(herderLegacy, nächsterScreen))
			},
		}),
	}

	screen.beendenKnopf = ui.NewButton(ui.ButtonConfig{
		Position: ui.Position{
			X:                10,
			Y:                10,
			AnchorHorizontal: ui.HorizontalerAnchorLinks,
			AnchorVertikal:   ui.VertikalerAnchorOben,
		},
		Text: "Schließen",
		Callback: func() {
			herderLegacy.OpenScreen(nächsterScreen())
		},
	})

	return &screen
}

func (m *menüScreen) components() []ui.Component {
	return []ui.Component{m.title, m.beschreibung, m.spielenKnopf, m.beendenKnopf}
}

func (m *menüScreen) Update() {
	for _, component := range m.components() {
		component.Update()
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		m.herderLegacy.OpenScreen(m.nächsterScreen())
	}
}

func (m *menüScreen) Draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)
	for _, component := range m.components() {
		component.Draw(screen)
	}
}
