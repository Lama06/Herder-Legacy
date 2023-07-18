package ui

import (
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DecideScreenConfig struct {
	Title string
	Text  string

	CancelText   string
	CancelAction func() herderlegacy.Screen

	ConfirmText   string
	ConfirmAction func() herderlegacy.Screen
}

func (d DecideScreenConfig) cancelTextOrDefault() string {
	if d.CancelText == "" {
		return "Abbrechen"
	}
	return d.CancelText
}

func (d DecideScreenConfig) confirmTextOrDefault() string {
	if d.ConfirmText == "" {
		return "Weiter"
	}
	return d.ConfirmText
}

type decideScreen struct {
	herderLegacy herderlegacy.HerderLegacy
	config       DecideScreenConfig

	title      *Title
	text       *Text
	cancelBtn  *Button
	confirmBtn *Button
}

var _ herderlegacy.Screen = (*decideScreen)(nil)

func NewDecideScreen(herderLegacy herderlegacy.HerderLegacy, config DecideScreenConfig) herderlegacy.Screen {
	return &decideScreen{
		herderLegacy: herderLegacy,
		config:       config,

		title: NewTitle(TitleConfig{
			Position: NewCenteredPosition(Width/2, 100),
			Text:     config.Title,
		}),
		text: NewText(TextConfig{
			Position: Position{
				X:                Width / 2,
				Y:                175,
				AnchorHorizontal: HorizontalerAnchorMitte,
				AnchorVertikal:   VertikalerAnchorOben,
			},
			Text: config.Text,
		}),
		cancelBtn: NewButton(ButtonConfig{
			Position: Position{
				X:                Width/2 - 25,
				Y:                Height - 100,
				AnchorHorizontal: HorizontalerAnchorRechts,
				AnchorVertikal:   VertikalerAnchorUnten,
			},
			Text:               config.cancelTextOrDefault(),
			CustomColorPalette: true,
			ColorPalette:       CancelButtonColorPalette,
			Callback: func() {
				herderLegacy.OpenScreen(config.CancelAction())
			},
		}),
		confirmBtn: NewButton(ButtonConfig{
			Position: Position{
				X:                Width/2 + 25,
				Y:                Height - 100,
				AnchorHorizontal: HorizontalerAnchorLinks,
				AnchorVertikal:   VertikalerAnchorUnten,
			},
			Text: config.confirmTextOrDefault(),
			Callback: func() {
				herderLegacy.OpenScreen(config.ConfirmAction())
			},
		}),
	}
}

func (d *decideScreen) components() []Component {
	return []Component{d.title, d.text, d.cancelBtn, d.confirmBtn}
}

func (d *decideScreen) Update() {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		d.herderLegacy.OpenScreen(d.config.CancelAction())
		return
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		d.herderLegacy.OpenScreen(d.config.ConfirmAction())
		return
	}

	for _, component := range d.components() {
		component.Update()
	}
}

func (d *decideScreen) Draw(screen *ebiten.Image) {
	screen.Fill(BackgroundColor)

	for _, component := range d.components() {
		component.Draw(screen)
	}
}
