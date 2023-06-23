package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
	"strings"
)

var (
	textColor      = color.RGBA{R: 87, G: 70, B: 123, A: 255}
	textHoverColor = color.RGBA{R: 82, G: 73, B: 72, A: 255}
)

type TextConfig struct {
	Position Position
	Text     string
}

type Text struct {
	position Position
	text     string
	hovered  bool
}

var _ Component = (*Text)(nil)

func NewText(config TextConfig) *Text {
	return &Text{
		position: config.Position,
		text:     config.Text,
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

func (t *Text) Update() {
	textBounds := text.BoundString(normalFontFace, t.text).Size()
	textWidth, textHeight := textBounds.X, textBounds.Y
	t.hovered = t.position.isHovered(float64(textWidth), float64(textHeight))
}

func (t *Text) Draw(screen *ebiten.Image) {
	textBounds := text.BoundString(normalFontFace, t.text).Size()
	textWidth, textHeight := textBounds.X, textBounds.Y
	eckeObenLinksX, eckeObenLinksY := t.position.eckeObenLinks(float64(textWidth), float64(textHeight))

	clr := textColor
	if t.hovered {
		clr = textHoverColor
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
