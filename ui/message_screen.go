package ui

import (
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MessageScreenConfig struct {
	Title          string
	Text           string
	ContinueText   string
	ContinueAction func() herderlegacy.Screen
}

func (m MessageScreenConfig) continueTextOrDefault() string {
	if m.ContinueText == "" {
		return "Weiter"
	}
	return m.ContinueText
}

type messageScreen struct {
	herderLegacy herderlegacy.HerderLegacy
	config       MessageScreenConfig

	title       *Title
	text        *Text
	continueBtn *Button
}

var _ herderlegacy.Screen = (*messageScreen)(nil)

func NewMessageScreen(herderLegacy herderlegacy.HerderLegacy, config MessageScreenConfig) herderlegacy.Screen {
	return &messageScreen{
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
		continueBtn: NewButton(ButtonConfig{
			Position: NewCenteredPosition(Width/2, Height-100),
			Text:     config.continueTextOrDefault(),
			Callback: func() {
				herderLegacy.OpenScreen(config.ContinueAction())
			},
		}),
	}
}

func (m *messageScreen) components() []Component {
	return []Component{m.title, m.text, m.continueBtn}
}

func (m *messageScreen) Update() {
	if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		m.herderLegacy.OpenScreen(m.config.ContinueAction())
		return
	}

	for _, component := range m.components() {
		component.Update()
	}
}

func (m *messageScreen) Draw(screen *ebiten.Image) {
	screen.Fill(BackgroundColor)

	for _, component := range m.components() {
		component.Draw(screen)
	}
}
