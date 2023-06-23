package dame

import (
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"strings"
)

type gameOverScreen struct {
	dame         *dameSpiel
	lehrer       lehrer
	title        *ui.Title
	info         *ui.Text
	restartKnopf *ui.Button
	beendenKnopf *ui.Button
	beendet      bool
}

var _ screen = (*gameOverScreen)(nil)

func newGameOverScreen(dame *dameSpiel, lehrer lehrer, gewonnenn bool) *gameOverScreen {
	if gewonnenn {
		dame.herderLegacy.AddVerhinderteStunden(3)
	}

	var titleText string
	if gewonnenn {
		titleText = "Gewonnen"
	} else {
		titleText = "Verloren"
	}

	var infoText string
	if gewonnenn {
		infoText = `Du hast gegen %nameAkk% gewonnen.
%pronomenSatzanfang% hat damit nicht gerechnet, weil %pronomenSatzmitte% ja die Regeln festgelegt hat.
Jetzt ist %name% schlecht gelaunt und hat auch weniger Motivation, zu unterrichten.
Durch deinen Sieg sind die Sommerferien 3 Stunden nach vorne gerutscht!`
	} else {
		infoText = `Du hast gegen %nameAkk% verloren.
Versuche es nocheinmal und gewinne, damit %pronomenSatzmitte% weniger motiviert ist, zu unterrichten.`
	}
	infoTextReplacer := strings.NewReplacer(
		"%name%", lehrer.name,
		"%nameAkk%", lehrer.nameAkkusativOrDefault(),
		"%pronomenSatzanfang%", lehrer.personalPronomenSatzanfang(),
		"%pronomenSatzmitte%", lehrer.personalPronomenSatzmitte(),
	)
	infoText = infoTextReplacer.Replace(infoText)

	screen := gameOverScreen{
		dame:   dame,
		lehrer: lehrer,
		title: ui.NewTitle(ui.TitleConfig{
			Position: ui.NewCenteredPosition(ui.Width/2, 300),
			Text:     titleText,
		}),
		info: ui.NewText(ui.TextConfig{
			Position: ui.Position{
				X:                ui.Width / 2,
				Y:                400,
				AnchorHorizontal: ui.HorizontalerAnchorMitte,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text: infoText,
		}),
		restartKnopf: ui.NewButton(ui.ButtonConfig{
			Position: ui.Position{
				X:                ui.Width/2 - 20,
				Y:                ui.Height - 80,
				AnchorHorizontal: ui.HorizontalerAnchorRechts,
				AnchorVertikal:   ui.VertikalerAnchorUnten,
			},
			Text: "Ein weitere Runde Dame spielen",
			Callback: func() {
				dame.currentScreen = newLehrerInfoScreen(dame, lehrer)
			},
		}),
	}

	screen.beendenKnopf = ui.NewButton(ui.ButtonConfig{
		Position: ui.Position{
			X:                ui.Width/2 + 20,
			Y:                ui.Height - 80,
			AnchorHorizontal: ui.HorizontalerAnchorLinks,
			AnchorVertikal:   ui.VertikalerAnchorUnten,
		},
		Text: "Dame beenden",
		Callback: func() {
			screen.beendet = true
		},
	})

	return &screen
}

func (g *gameOverScreen) components() []ui.Component {
	return []ui.Component{g.title, g.info, g.restartKnopf, g.beendenKnopf}
}

func (g *gameOverScreen) update() (beendet bool) {
	for _, component := range g.components() {
		component.Update()
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		g.dame.currentScreen = newLehrerInfoScreen(g.dame, g.lehrer)
	}
	return g.beendet
}

func (g *gameOverScreen) draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)
	for _, component := range g.components() {
		component.Draw(screen)
	}
}
