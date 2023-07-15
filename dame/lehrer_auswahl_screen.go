package dame

import (
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type lehrerAuswahlScreen struct {
	herderLegacy   herderlegacy.HerderLegacy
	nächsterScreen func() herderlegacy.Screen
	title          *ui.Title
	zurückKnopf    *ui.Button
	lehrerKnöpfe   []*ui.Button
}

var _ herderlegacy.Screen = (*lehrerAuswahlScreen)(nil)

func newLehrerAuswahlScreen(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func() herderlegacy.Screen,
) *lehrerAuswahlScreen {
	lehrerKnöpfe := make([]*ui.Button, len(alleLehrer))
	for i, lehrer := range alleLehrer {
		lehrer := lehrer
		lehrerKnöpfe[i] = ui.NewButton(ui.ButtonConfig{
			Position: ui.NewCenteredPosition(ui.Width/2, 200+80*float64(i)),
			Text:     lehrer.name,
			Callback: func() {
				herderLegacy.OpenScreen(newLehrerInfoScreen(herderLegacy, nächsterScreen, lehrer))
			},
		})
	}

	return &lehrerAuswahlScreen{
		herderLegacy:   herderLegacy,
		nächsterScreen: nächsterScreen,
		title: ui.NewTitle(ui.TitleConfig{
			Position: ui.NewCenteredPosition(ui.Width/2, 100),
			Text:     "Lehrer auswählen",
		}),
		zurückKnopf: ui.NewButton(ui.ButtonConfig{
			Position: ui.Position{
				X:                10,
				Y:                10,
				AnchorHorizontal: ui.HorizontalerAnchorLinks,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text: "Zurück",
			Callback: func() {
				herderLegacy.OpenScreen(nächsterScreen())
			},
		}),
		lehrerKnöpfe: lehrerKnöpfe,
	}
}

func (l *lehrerAuswahlScreen) components() []ui.Component {
	components := []ui.Component{l.title, l.zurückKnopf}
	for _, lehrerKnopf := range l.lehrerKnöpfe {
		components = append(components, lehrerKnopf)
	}
	return components
}

func (l *lehrerAuswahlScreen) Update() {
	for _, component := range l.components() {
		component.Update()
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		l.herderLegacy.OpenScreen(newMenüScreen(l.herderLegacy, l.nächsterScreen))
	}
}

func (l *lehrerAuswahlScreen) Draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)
	for _, component := range l.components() {
		component.Draw(screen)
	}
}
