package ui

import (
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ListScreenWidget interface {
	abstand() float64

	createComponent(herderLegacy herderlegacy.HerderLegacy, position Position) Component
}

type ListScreenButtonWidget struct {
	Text     string
	Callback func()

	CustomColorPalette bool
	ColorPalette       ButtonColorPalette
}

func (l ListScreenButtonWidget) abstand() float64 {
	return 80
}

func (l ListScreenButtonWidget) createComponent(herderLegacy herderlegacy.HerderLegacy, position Position) Component {
	return NewButton(ButtonConfig{
		Position:           position,
		Text:               l.Text,
		CustomColorPalette: l.CustomColorPalette,
		ColorPalette:       l.ColorPalette,
		Callback:           l.Callback,
	})
}

type ListScreenToggleWidget struct {
	Text     string
	Callback func(bool)
	Enabled  bool
}

func (l ListScreenToggleWidget) abstand() float64 {
	return 80
}

func (l ListScreenToggleWidget) createComponent(herderLegacy herderlegacy.HerderLegacy, position Position) Component {
	return NewToggle(ToggleConfig{
		Position: position,
		Text:     l.Text,
		Enabled:  l.Enabled,
		Callback: l.Callback,
	})
}

type ListScreenSelectionWidget[T any] struct {
	Text     string
	Value    T
	Values   []T
	Callback func(T)
	ToString func(T) string
}

func (l ListScreenSelectionWidget[T]) abstand() float64 {
	return 80
}

func (l ListScreenSelectionWidget[T]) createComponent(herderLegacy herderlegacy.HerderLegacy, position Position) Component {
	return NewSelection(herderLegacy, SelectionConfig[T]{
		Position: position,
		Text:     l.Text,
		Value:    l.Value,
		Values:   l.Values,
		Callback: l.Callback,
		ToString: l.ToString,
	})
}

type ListScreenConfig struct {
	Title   string
	Text    string
	Widgets []ListScreenWidget

	CancelText   string
	CancelAction func() herderlegacy.Screen
}

func (l ListScreenConfig) cancelTextOrDefault() string {
	if l.CancelText == "" {
		return "Zurück"
	}
	return l.CancelText
}

func (l ListScreenConfig) abbrechbar() bool {
	return l.CancelAction != nil
}

type listScreen struct {
	herderLegacy herderlegacy.HerderLegacy
	config       ListScreenConfig

	title        *Title
	text         *Text
	cancelButton *Button
	widgets      []Component
}

func NewListScreen(herderLegacy herderlegacy.HerderLegacy, config ListScreenConfig) herderlegacy.Screen {
	const widgetStartY = 175.0

	spalten := 1
	widgetY := widgetStartY
	for _, widget := range config.Widgets {
		if widgetY+widget.abstand() >= Height {
			spalten++
			widgetY = widgetStartY + widget.abstand()
			continue
		}
		widgetY += widget.abstand()
	}

	widgets := make([]Component, len(config.Widgets))
	widgetIndex := 0
	widgetY = widgetStartY
spalten:
	for spalte := 0; spalte < spalten; spalte++ {
		widgetX := (Width / float64(spalten+1)) * float64(spalte+1)
		for _, widget := range config.Widgets[widgetIndex:] {
			if widgetY+widget.abstand() >= Height {
				widgetY = widgetStartY
				continue spalten
			}
			widgetY += widget.abstand() / 2
			widgets[widgetIndex] = widget.createComponent(herderLegacy, NewCenteredPosition(widgetX, widgetY))
			widgetY += widget.abstand() / 2
			widgetIndex++
		}
	}

	var cancelButton *Button
	if config.abbrechbar() {
		cancelButton = NewButton(ButtonConfig{
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

	return &listScreen{
		herderLegacy: herderLegacy,
		config:       config,

		title: NewTitle(TitleConfig{
			Position: NewCenteredPosition(Width/2, 100),
			Text:     config.Title,
		}),
		text: NewText(TextConfig{
			Position: Position{
				X:                Width / 2,
				Y:                140,
				AnchorHorizontal: HorizontalerAnchorMitte,
				AnchorVertikal:   VertikalerAnchorOben,
			},
			Text: config.Text,
		}),
		cancelButton: cancelButton,
		widgets:      widgets,
	}
}

func (l *listScreen) components() []Component {
	components := append(l.widgets, l.title, l.text)
	if l.cancelButton != nil {
		components = append(components, l.cancelButton)
	}
	return components
}

func (l *listScreen) Update() {
	if l.config.abbrechbar() && inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		l.herderLegacy.OpenScreen(l.config.CancelAction())
		return
	}

	for _, component := range l.components() {
		component.Update()
	}
}

func (l *listScreen) Draw(screen *ebiten.Image) {
	screen.Fill(BackgroundColor)

	for _, component := range l.components() {
		component.Draw(screen)
	}
}
