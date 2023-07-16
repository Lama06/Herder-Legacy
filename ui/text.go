package ui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type TextColorPalatte struct {
	Color      color.Color
	HoverColor color.Color
}

var defaultTextColorPalette = TextColorPalatte{
	Color:      color.RGBA{R: 87, G: 70, B: 123, A: 255},
	HoverColor: color.RGBA{R: 82, G: 73, B: 72, A: 255},
}

func (t TextColorPalatte) hoverColorOrDefault() color.Color {
	if t.HoverColor == nil {
		return t.Color
	}
	return t.HoverColor
}

type TextConfig struct {
	Position           Position
	Text               string
	CustomColorPalette bool
	ColorPalette       TextColorPalatte
}

type Text struct {
	position     Position
	text         string
	colorPalette TextColorPalatte
	hovered      bool
}

var _ Component = (*Text)(nil)

func NewText(config TextConfig) *Text {
	var colorPalette TextColorPalatte
	if config.CustomColorPalette {
		colorPalette = config.ColorPalette
	} else {
		colorPalette = defaultTextColorPalette
	}

	return &Text{
		position:     config.Position,
		text:         config.Text,
		colorPalette: colorPalette,
	}
}

func (t *Text) Position() Position {
	return t.position
}

func (t *Text) SetPosition(position Position) {
	t.position = position
}

func (t *Text) Text() string {
	return t.text
}

func (t *Text) SetText(text string) {
	t.text = text
}

func (t *Text) ColorPalette() TextColorPalatte {
	return t.colorPalette
}

func (t *Text) SetColorPalette(colorPalette TextColorPalatte) {
	t.colorPalette = colorPalette
}

func (t *Text) Update() {
	textBounds := text.BoundString(normalFontFace, t.text).Size()
	textWidth, textHeight := textBounds.X, textBounds.Y
	t.hovered = t.position.isHovered(float64(textWidth), float64(textHeight))
}

func (t *Text) Draw(screen *ebiten.Image) {
	textBounds := text.BoundString(normalFontFace, t.text).Size()
	textWidth, textHeight := textBounds.X, textBounds.Y
	eckeObenLinksX, eckeObenLinksY := t.position.eckeObenLinks(float64(textWidth), float64(textHeight))

	clr := t.colorPalette.Color
	if t.hovered {
		clr = t.colorPalette.hoverColorOrDefault()
	}

	lines := strings.Count(t.text, "\n") + 1
	text.Draw(
		screen,
		t.text,
		normalFontFace,
		int(eckeObenLinksX),
		int(eckeObenLinksY+float64(textHeight)/float64(lines)),
		clr,
	)
}
