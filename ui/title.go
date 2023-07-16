package ui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type TitleColorPalette struct {
	Color      color.Color
	HoverColor color.Color
}

var defaultTitleColorPalette = TitleColorPalette{
	Color:      color.RGBA{R: 87, G: 70, B: 123, A: 255},
	HoverColor: color.RGBA{R: 112, G: 248, B: 186, A: 255},
}

func (t TitleColorPalette) hoverColorOrDefault() color.Color {
	if t.HoverColor == nil {
		return t.Color
	}
	return t.HoverColor
}

type TitleConfig struct {
	Position           Position
	Text               string
	CustomColorPalette bool
	ColorPalette       TitleColorPalette
}

type Title struct {
	position     Position
	text         string
	colorPalette TitleColorPalette
	hovered      bool
}

var _ Component = (*Title)(nil)

func NewTitle(config TitleConfig) *Title {
	var colorPalette TitleColorPalette
	if config.CustomColorPalette {
		colorPalette = config.ColorPalette
	} else {
		colorPalette = defaultTitleColorPalette
	}

	return &Title{
		position:     config.Position,
		text:         config.Text,
		colorPalette: colorPalette,
	}
}

func (t *Title) Position() Position {
	return t.position
}

func (t *Title) SetPosition(position Position) {
	t.position = position
}

func (t *Title) Text() string {
	return t.text
}

func (t *Title) SetText(text string) {
	t.text = text
}

func (t *Title) ColorPalette() TitleColorPalette {
	return t.colorPalette
}

func (t *Title) SetColorPalette(colorPalette TitleColorPalette) {
	t.colorPalette = colorPalette
}

func (t *Title) Update() {
	textBounds := text.BoundString(titleFontFace, t.text).Size()
	textWidth, textHeight := textBounds.X, textBounds.Y
	t.hovered = t.position.isHovered(float64(textWidth), float64(textHeight))
}

func (t *Title) Draw(screen *ebiten.Image) {
	textBounds := text.BoundString(titleFontFace, t.text).Size()
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
		titleFontFace,
		int(eckeObenLinksX),
		int(eckeObenLinksY+float64(textHeight)/float64(lines)),
		clr,
	)
}
