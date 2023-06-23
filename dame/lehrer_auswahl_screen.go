package dame

import (
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type lehrerAuswahlScreen struct {
	dame          *dameSpiel
	title         *ui.Title
	zurueckKnopf  *ui.Button
	lehrerKnoepfe []*ui.Button
}

var _ screen = (*lehrerAuswahlScreen)(nil)

func newLehrerAuswahlScreen(dame *dameSpiel) *lehrerAuswahlScreen {
	lehrerKnoepfe := make([]*ui.Button, len(lehrerListe))
	for i, lehrer := range lehrerListe {
		lehrer := lehrer
		lehrerKnoepfe[i] = ui.NewButton(ui.ButtonConfig{
			Position: ui.NewCenteredPosition(ui.Width/2, 200+80*float64(i)),
			Text:     lehrer.name,
			Callback: func() {
				dame.currentScreen = newLehrerInfoScreen(dame, lehrer)
			},
		})
	}

	return &lehrerAuswahlScreen{
		dame: dame,
		title: ui.NewTitle(ui.TitleConfig{
			Position: ui.NewCenteredPosition(ui.Width/2, 100),
			Text:     "Lehrer auswählen",
		}),
		zurueckKnopf: ui.NewButton(ui.ButtonConfig{
			Position: ui.Position{
				X:                10,
				Y:                10,
				AnchorHorizontal: ui.HorizontalerAnchorLinks,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text: "Zurück",
			Callback: func() {
				dame.currentScreen = newMenuScreen(dame)
			},
		}),
		lehrerKnoepfe: lehrerKnoepfe,
	}
}

func (l *lehrerAuswahlScreen) components() []ui.Component {
	components := []ui.Component{l.title, l.zurueckKnopf}
	for _, lehrerKnopf := range l.lehrerKnoepfe {
		components = append(components, lehrerKnopf)
	}
	return components
}

func (l *lehrerAuswahlScreen) update() (beendet bool) {
	for _, component := range l.components() {
		component.Update()
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		l.dame.currentScreen = newMenuScreen(l.dame)
	}
	return false
}

func (l *lehrerAuswahlScreen) draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)
	for _, component := range l.components() {
		component.Draw(screen)
	}
}
