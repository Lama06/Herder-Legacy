package leitstellenspiel

import (
	"fmt"
	"image/color"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
)

func NewLeitstellenspielBotScreen(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func() herderlegacy.Screen,
) herderlegacy.Screen {
	return newBootScreen(herderLegacy, nächsterScreen)
}

type bootScreen struct {
	herderLegacy   herderlegacy.HerderLegacy
	nächsterScreen func() herderlegacy.Screen

	title     *ui.Title
	info      *ui.Text
	countdown *ui.Countdown
}

func newBootScreen(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func() herderlegacy.Screen,
) *bootScreen {
	return &bootScreen{
		herderLegacy:   herderLegacy,
		nächsterScreen: nächsterScreen,

		title: ui.NewTitle(ui.TitleConfig{
			Position:           ui.NewCenteredPosition(ui.Width/2, 100),
			CustomColorPalette: true,
			ColorPalette: ui.TitleColorPalette{
				Color: colornames.White,
			},
			Text: "MySHN: Intialisieren...",
		}),
		info: ui.NewText(ui.TextConfig{
			Position: ui.Position{
				X:                ui.Width - 5,
				Y:                ui.Height - 5,
				AnchorHorizontal: ui.HorizontalerAnchorRechts,
				AnchorVertikal:   ui.VertikalerAnchorUnten,
			},
			Text:               "Boot Menü: F12 drücken",
			CustomColorPalette: true,
			ColorPalette: ui.TextColorPalatte{
				Color: colornames.Lightgray,
			},
		}),
		countdown: ui.NewCountdown(ui.CountdownConfig{
			Position:  ui.NewCenteredPosition(ui.Width/2, ui.Height/2),
			StartZeit: 6 * 60,
			AbgelaufenCallback: func() {
				herderLegacy.OpenScreen(nächsterScreen())
			},
			CustomColorPalette: true,
			ColorPalette: ui.CountdownColorPalette{
				Color: colornames.White,
			},
		}),
	}
}

func (b *bootScreen) components() []ui.Component {
	return []ui.Component{b.title, b.info, b.countdown}
}

func (b *bootScreen) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		b.herderLegacy.OpenScreen(newPasswortEingabeScreen(b.herderLegacy, b.nächsterScreen))
	}

	for _, component := range b.components() {
		component.Update()
	}
}

func (b *bootScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0, G: 162, B: 237, A: 255})
	for _, component := range b.components() {
		component.Draw(screen)
	}
}

type passwortEingabeScreen struct {
	herderLegacy   herderlegacy.HerderLegacy
	nächsterScreen func() herderlegacy.Screen

	info           *ui.Text
	passwortWidget *ui.Text
	passwort       []rune
}

func newPasswortEingabeScreen(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func() herderlegacy.Screen,
) *passwortEingabeScreen {
	return &passwortEingabeScreen{
		herderLegacy:   herderLegacy,
		nächsterScreen: nächsterScreen,

		info: ui.NewText(ui.TextConfig{
			Position: ui.Position{
				X:                ui.Width / 2,
				Y:                ui.Height/2 - 10,
				AnchorHorizontal: ui.HorizontalerAnchorMitte,
				AnchorVertikal:   ui.VertikalerAnchorUnten,
			},
			Text:               "Admin Passwort eingeben und mit Enter bestätigen",
			CustomColorPalette: true,
			ColorPalette: ui.TextColorPalatte{
				Color: colornames.White,
			},
		}),
		passwortWidget: ui.NewText(ui.TextConfig{
			Position: ui.Position{
				X:                ui.Width / 2,
				Y:                ui.Height/2 + 10,
				AnchorHorizontal: ui.HorizontalerAnchorMitte,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text:               "Admin",
			CustomColorPalette: true,
			ColorPalette: ui.TextColorPalatte{
				Color: colornames.White,
			},
		}),
		passwort: []rune("Admin"),
	}
}

func (p *passwortEingabeScreen) components() []ui.Component {
	return []ui.Component{p.info, p.passwortWidget}
}

func (p *passwortEingabeScreen) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(p.passwort) != 0 {
		p.passwort = p.passwort[:len(p.passwort)-1]
	}
	p.passwort = ebiten.AppendInputChars(p.passwort)
	p.passwortWidget.SetText("Passwort: " + string(p.passwort))

	if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		if string(p.passwort) == "" {
			p.herderLegacy.OpenScreen(newDesktopScreen(p.herderLegacy, p.nächsterScreen))
		} else {
			p.herderLegacy.OpenScreen(p.nächsterScreen())
		}
	}

	for _, component := range p.components() {
		component.Update()
	}
}

func (p *passwortEingabeScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0, G: 162, B: 237, A: 255})
	for _, component := range p.components() {
		component.Draw(screen)
	}
}

func newDesktopScreen(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func() herderlegacy.Screen,
) herderlegacy.Screen {
	return ui.NewDecideScreen(herderLegacy, ui.DecideScreenConfig{
		Title:        "Windows 7 Home Premium",
		CancelText:   "Herunterfahren",
		CancelAction: nächsterScreen,
		ConfirmText:  "Leitstellenspiel Pro Tool 3001",
		ConfirmAction: func() herderlegacy.Screen {
			herderLegacy.AddVerhinderteStunden(3)
			return newConsoleScreen(herderLegacy, nächsterScreen)
		},
	})
}

const consoleAccountDelay = 10

type consoleScreen struct {
	konsoleText           *ui.Text
	herunterfahrenKnopf   *ui.Button
	aktuelleAccountNumber int
	nächsterAccountZeit   int
}

func newConsoleScreen(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func() herderlegacy.Screen,
) herderlegacy.Screen {
	return &consoleScreen{
		konsoleText: ui.NewText(ui.TextConfig{
			Position: ui.Position{
				X:                ui.Width / 2,
				Y:                ui.Height - 10,
				AnchorHorizontal: ui.HorizontalerAnchorMitte,
				AnchorVertikal:   ui.VertikalerAnchorUnten,
			},
			Text:               "Erstelle Allianz Mission",
			CustomColorPalette: true,
			ColorPalette: ui.TextColorPalatte{
				Color: colornames.Lightgreen,
			},
		}),
		herunterfahrenKnopf: ui.NewButton(ui.ButtonConfig{
			Position: ui.Position{
				X:                0,
				Y:                ui.Height,
				AnchorHorizontal: ui.HorizontalerAnchorLinks,
				AnchorVertikal:   ui.VertikalerAnchorUnten,
			},
			Text:               "Herunterfahren",
			CustomColorPalette: true,
			ColorPalette:       ui.CancelButtonColorPalette,
			Callback: func() {
				herderLegacy.OpenScreen(nächsterScreen())
			},
		}),
		aktuelleAccountNumber: 0,
		nächsterAccountZeit:   consoleAccountDelay,
	}
}

func (c *consoleScreen) components() []ui.Component {
	return []ui.Component{c.konsoleText, c.herunterfahrenKnopf}
}

func (c *consoleScreen) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Black)
	for _, component := range c.components() {
		component.Draw(screen)
	}
}

func (c *consoleScreen) Update() {
	if c.nächsterAccountZeit <= 0 {
		c.aktuelleAccountNumber++
		c.konsoleText.SetText(fmt.Sprintf("%v\nAlamiere Account: Wehner.%v", c.konsoleText.Text(), c.aktuelleAccountNumber))
		c.nächsterAccountZeit = consoleAccountDelay
	} else {
		c.nächsterAccountZeit--
	}

	for _, component := range c.components() {
		component.Update()
	}
}
