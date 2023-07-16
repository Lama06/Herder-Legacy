package ui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var (
	titleColor      = color.RGBA{R: 87, G: 70, B: 123, A: 255}
	titleHoverColor = color.RGBA{R: 112, G: 248, B: 186, A: 255}
)

type TitleConfig struct {
	Position Position
	Text     string
}

type Title struct {
	position Position
	text     string
	hovered  bool
}

var _ Component = (*Title)(nil)

func NewTitle(config TitleConfig) *Title {
	return &Title{
		position: config.Position,
		text:     config.Text,
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

func (t *Title) Update() {
	textBounds := text.BoundString(titleFontFace, t.text).Size()
	textWidth, textHeight := textBounds.X, textBounds.Y
	t.hovered = t.position.isHovered(float64(textWidth), float64(textHeight))
}

func (t *Title) Draw(screen *ebiten.Image) {
	textBounds := text.BoundString(titleFontFace, t.text).Size()
	textWidth, textHeight := textBounds.X, textBounds.Y
	eckeObenLinksX, eckeObenLinksY := t.position.eckeObenLinks(float64(textWidth), float64(textHeight))

	clr := titleColor
	if t.hovered {
		clr = titleHoverColor
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
