package dame

import (
	"fmt"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type lehrerInfoScreen struct {
	dame         *dameSpiel
	lehrer       lehrer
	title        *ui.Title
	info         *ui.Text
	zurueckKnopf *ui.Button
	weiterKnopf  *ui.Button
}

var _ screen = (*lehrerInfoScreen)(nil)

func newLehrerInfoScreen(dame *dameSpiel, lehrer lehrer) *lehrerInfoScreen {
	return &lehrerInfoScreen{
		dame:   dame,
		lehrer: lehrer,
		title: ui.NewTitle(ui.TitleConfig{
			Position: ui.NewCenteredPosition(ui.Width/2, 200),
			Text:     lehrer.name,
		}),
		info: ui.NewText(ui.TextConfig{
			Position: ui.Position{
				X:                ui.Width / 2,
				Y:                400,
				AnchorHorizontal: ui.HorizontalerAnchorMitte,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text: lehrer.info,
		}),
		zurueckKnopf: ui.NewButton(ui.ButtonConfig{
			Position: ui.Position{
				X:                ui.Width/2 - 20,
				Y:                ui.Height - 80,
				AnchorHorizontal: ui.HorizontalerAnchorRechts,
				AnchorVertikal:   ui.VertikalerAnchorUnten,
			},
			Text: "Gegen anderen Lehrer spielen",
			Callback: func() {
				dame.currentScreen = newLehrerAuswahlScreen(dame)
			},
		}),
		weiterKnopf: ui.NewButton(ui.ButtonConfig{
			Position: ui.Position{
				X:                ui.Width/2 + 20,
				Y:                ui.Height - 80,
				AnchorHorizontal: ui.HorizontalerAnchorLinks,
				AnchorVertikal:   ui.VertikalerAnchorUnten,
			},
			Text: fmt.Sprintf("Gegen %v spielen", lehrer.nameAkkusativOrDefault()),
			Callback: func() {
				dame.currentScreen = newSpielScreen(dame, lehrer)
			},
		}),
	}
}

func (l *lehrerInfoScreen) components() []ui.Component {
	return []ui.Component{l.title, l.info, l.zurueckKnopf, l.weiterKnopf}
}

func (l *lehrerInfoScreen) update() (beendet bool) {
	for _, component := range l.components() {
		component.Update()
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		l.dame.currentScreen = newLehrerAuswahlScreen(l.dame)
	}
	return false
}

func (l *lehrerInfoScreen) draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)
	for _, component := range l.components() {
		component.Draw(screen)
	}
}
