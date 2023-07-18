package ui

import (
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SelectScreenAuswahlMöglichkeit struct {
	Text   string
	Action func() herderlegacy.Screen
}

type SelectScreenConfig struct {
	Title string
	Text  string

	AuswahlMöglichkeiten []SelectScreenAuswahlMöglichkeit

	CancelText   string
	CancelAction func() herderlegacy.Screen
}

func (s SelectScreenConfig) cancelTextOrDefault() string {
	if s.CancelText == "" {
		return "Zurück"
	}
	return s.CancelText
}

func (s SelectScreenConfig) abbrechbar() bool {
	return s.CancelAction != nil
}

type selectScreen struct {
	herderlegacy herderlegacy.HerderLegacy
	config       SelectScreenConfig

	title         *Title
	text          *Text
	auswahlKnöpfe []*Button
	cancelKnopf   *Button
}

var _ herderlegacy.Screen = (*selectScreen)(nil)

func NewSelectScreen(herderLegacy herderlegacy.HerderLegacy, config SelectScreenConfig) herderlegacy.Screen {
	auswahlKnöpfe := make([]*Button, len(config.AuswahlMöglichkeiten))
	for i, auswahlMöglichkeit := range config.AuswahlMöglichkeiten {
		auswahlMöglichkeit := auswahlMöglichkeit
		auswahlKnöpfe[i] = NewButton(ButtonConfig{
			Position: NewCenteredPosition(Width/2, 300+80*float64(i)),
			Text:     auswahlMöglichkeit.Text,
			Callback: func() {
				herderLegacy.OpenScreen(auswahlMöglichkeit.Action())
			},
		})
	}

	var cancelKnopf *Button
	if config.abbrechbar() {
		cancelKnopf = NewButton(ButtonConfig{
			Position: Position{
				X:                20,
				Y:                20,
				AnchorHorizontal: HorizontalerAnchorLinks,
				AnchorVertikal:   VertikalerAnchorOben,
			},
			Text:               config.cancelTextOrDefault(),
			CustomColorPalette: true,
			ColorPalette:       CancelButtonColorPalette,
			Callback: func() {
				herderLegacy.OpenScreen(config.CancelAction())
			},
		})
	}

	return &selectScreen{
		herderlegacy: herderLegacy,
		config:       config,

		title: NewTitle(TitleConfig{
			Position: NewCenteredPosition(Width/2, 100),
			Text:     config.Title,
		}),
		text: NewText(TextConfig{
			Position: Position{
				X:                Width / 2,
				Y:                200,
				AnchorHorizontal: HorizontalerAnchorMitte,
				AnchorVertikal:   VertikalerAnchorOben,
			},
			Text: config.Text,
		}),
		auswahlKnöpfe: auswahlKnöpfe,
		cancelKnopf:   cancelKnopf,
	}
}

func (s *selectScreen) components() []Component {
	components := []Component{s.title, s.text}
	for _, auswahlKnopf := range s.auswahlKnöpfe {
		components = append(components, auswahlKnopf)
	}
	if s.cancelKnopf != nil {
		components = append(components, s.cancelKnopf)
	}
	return components
}

func (s *selectScreen) Update() {
	if s.config.abbrechbar() && inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		s.herderlegacy.OpenScreen(s.config.CancelAction())
		return
	}

	for _, component := range s.components() {
		component.Update()
	}
}

func (s *selectScreen) Draw(screen *ebiten.Image) {
	screen.Fill(BackgroundColor)

	for _, component := range s.components() {
		component.Draw(screen)
	}
}
