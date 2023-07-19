package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type ToggleConfig struct {
	Position Position
	Text     string
	Enabled  bool
	Callback func(bool)
}

type Toggle struct {
	config  ToggleConfig
	button  *Button
	enabled bool
}

func NewToggle(config ToggleConfig) *Toggle {
	toggle := Toggle{
		config: config,
		button: NewButton(ButtonConfig{
			Position: config.Position,
		}),
	}

	toggle.SetEnabled(config.Enabled)

	toggle.button.SetCallback(func() {
		toggle.SetEnabled(!toggle.enabled)
	})

	return &toggle
}

func (t *Toggle) Enabled() bool {
	return t.enabled
}

func (t *Toggle) SetEnabled(enabled bool) {
	t.enabled = enabled
	if enabled {
		t.button.SetText(t.config.Text + ": An")
		t.button.SetColorPalette(ButtonColorPalette{
			BackgroundColor:        colornames.Green,
			BackgroundColorHovered: colornames.Darkgreen,
			TextColor:              colornames.White,
		})
	} else {
		t.button.SetText(t.config.Text + ": Aus")
		t.button.SetColorPalette(ButtonColorPalette{
			BackgroundColor:        colornames.Red,
			BackgroundColorHovered: colornames.Darkred,
			TextColor:              colornames.White,
		})
	}
	if t.config.Callback != nil {
		t.config.Callback(enabled)
	}
}

func (t *Toggle) Position() Position {
	return t.button.Position()
}

func (t *Toggle) SetPosition(position Position) {
	t.button.SetPosition(position)
}

func (t *Toggle) Update() {
	t.button.Update()
}

func (t *Toggle) Draw(screen *ebiten.Image) {
	t.button.Draw(screen)
}
