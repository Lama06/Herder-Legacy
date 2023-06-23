package dame

import (
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type menuScreen struct {
	dame         *dameSpiel
	title        *ui.Title
	beschreibung *ui.Text
	spielenKnopf *ui.Button
	beendenKnopf *ui.Button
	geschlossen  bool
}

var _ screen = (*menuScreen)(nil)

func newMenuScreen(dame *dameSpiel) *menuScreen {
	screen := menuScreen{
		dame: dame,
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
				dame.currentScreen = newLehrerAuswahlScreen(dame)
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
			screen.geschlossen = true
		},
	})

	return &screen
}

func (m *menuScreen) components() []ui.Component {
	return []ui.Component{m.title, m.beschreibung, m.spielenKnopf, m.beendenKnopf}
}

func (m *menuScreen) update() (beendet bool) {
	for _, component := range m.components() {
		component.Update()
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		return true
	}
	return m.geschlossen
}

func (m *menuScreen) draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)
	for _, component := range m.components() {
		component.Draw(screen)
	}
}
