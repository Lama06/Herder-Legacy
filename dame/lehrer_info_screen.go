package dame

import (
	"fmt"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type lehrerInfoScreen struct {
	herderLegacy   herderlegacy.HerderLegacy
	nächsterScreen func() herderlegacy.Screen
	title          *ui.Title
	info           *ui.Text
	zurückKnopf    *ui.Button
	weiterKnopf    *ui.Button
}

var _ herderlegacy.Screen = (*lehrerInfoScreen)(nil)

func newLehrerInfoScreen(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func() herderlegacy.Screen,
	lehrer lehrer,
) *lehrerInfoScreen {
	return &lehrerInfoScreen{
		herderLegacy:   herderLegacy,
		nächsterScreen: nächsterScreen,
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
		zurückKnopf: ui.NewButton(ui.ButtonConfig{
			Position: ui.Position{
				X:                ui.Width/2 - 20,
				Y:                ui.Height - 80,
				AnchorHorizontal: ui.HorizontalerAnchorRechts,
				AnchorVertikal:   ui.VertikalerAnchorUnten,
			},
			Text: "Gegen anderen Lehrer spielen",
			Callback: func() {
				herderLegacy.OpenScreen(newLehrerAuswahlScreen(herderLegacy, nächsterScreen))
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
				herderLegacy.OpenScreen(newSpielScreen(herderLegacy, func(gewonnen bool) herderlegacy.Screen {
					return newGameOverScreen(herderLegacy, nächsterScreen, lehrer, gewonnen)
				}, lehrer.spielOptionen))
			},
		}),
	}
}

func (l *lehrerInfoScreen) components() []ui.Component {
	return []ui.Component{l.title, l.info, l.zurückKnopf, l.weiterKnopf}
}

func (l *lehrerInfoScreen) Update() {
	for _, component := range l.components() {
		component.Update()
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		l.herderLegacy.OpenScreen(newLehrerAuswahlScreen(l.herderLegacy, l.nächsterScreen))
	}
}

func (l *lehrerInfoScreen) Draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)
	for _, component := range l.components() {
		component.Draw(screen)
	}
}
