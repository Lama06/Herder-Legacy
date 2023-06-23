package ui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	buttonPadding = 20

	buttonScaleNotHovered       = 1
	buttonScaleHovered          = 1.2
	buttonScaleMaxChangePerTick = 0.03
)

var (
	buttonBackgroundColor        = color.RGBA{R: 18, G: 53, B: 91, A: 255}
	hoveredButtonBackgroundColor = color.RGBA{R: 134, G: 22, B: 87, A: 255}
	buttonTextColor              = color.RGBA{R: 212, G: 245, B: 245, A: 255}
	hoveredButtonTextColor       = color.RGBA{R: 212, G: 245, B: 245, A: 255}

	disabledButtonBackgroundColor        = color.RGBA{R: 42, G: 59, B: 82, A: 255}
	disabledHoveredButtonBackgroundColor = color.RGBA{R: 29, G: 37, B: 48, A: 255}
	disabledButtonTextColor              = color.RGBA{R: 212, G: 245, B: 245, A: 255}
	disabledHoveredButtonTextColor       = color.RGBA{R: 212, G: 245, B: 245, A: 255}
)

type ButtonConfig struct {
	Position Position
	Text     string
	Callback func()
	Disabled bool
}

type Button struct {
	position Position
	text     string
	callback func()
	disabled bool

	currentScale  float64
	hovered       bool
	imgNotHovered *ebiten.Image
	imgHovered    *ebiten.Image
}

var _ Component = (*Button)(nil)

func NewButton(config ButtonConfig) *Button {
	button := Button{
		position: config.Position,
		text:     config.Text,
		callback: config.Callback,
		disabled: config.Disabled,

		currentScale: buttonScaleNotHovered,
	}

	button.updateImages()

	return &button
}

func (b *Button) textBounds() (textWidth, textHeight int) {
	bounds := text.BoundString(normalFontFace, b.text)
	return bounds.Dx(), bounds.Dy()
}

func (b *Button) buttonSize() (buttonWidth, buttonHeight int) {
	textWidth, textHeight := b.textBounds()
	return textWidth + buttonPadding*2, textHeight + buttonPadding*2
}

func (b *Button) createImage(backgroundColor, textColor color.Color) *ebiten.Image {
	_, textHeight := b.textBounds()
	buttonWidth, buttonHeight := b.buttonSize()
	img := ebiten.NewImage(buttonWidth, buttonHeight)
	img.Fill(backgroundColor)
	textX := buttonPadding
	lines := strings.Count(b.text, "\n") + 1
	textY := buttonPadding + textHeight/lines
	text.Draw(img, b.text, normalFontFace, textX, textY, textColor)
	return img
}

func (b *Button) backgroundColor() color.Color {
	if b.disabled {
		return disabledButtonBackgroundColor
	}
	return buttonBackgroundColor
}

func (b *Button) backgroundColorHovered() color.Color {
	if b.disabled {
		return disabledHoveredButtonBackgroundColor
	}
	return hoveredButtonBackgroundColor
}

func (b *Button) textColor() color.Color {
	if b.disabled {
		return disabledButtonTextColor
	}
	return buttonTextColor
}

func (b *Button) textColorHovered() color.Color {
	if b.disabled {
		return disabledHoveredButtonTextColor
	}
	return hoveredButtonTextColor
}

func (b *Button) updateImages() {
	b.imgNotHovered = b.createImage(b.backgroundColor(), b.textColor())
	b.imgHovered = b.createImage(b.backgroundColorHovered(), b.textColorHovered())
}

func (b *Button) Position() Position {
	return b.position
}

func (b *Button) SetPosition(position Position) {
	b.position = position
}

func (b *Button) Text() string {
	return b.text
}

func (b *Button) SetText(text string) {
	b.text = text
	b.updateImages()
}

func (b *Button) Disabled() bool {
	return b.disabled
}

func (b *Button) SetDisabled(disabled bool) {
	b.disabled = disabled
	b.updateImages()
}

func (b *Button) Update() {
	buttonWidth, buttonHeight := b.buttonSize()

	if b.position.isClicked(float64(buttonWidth), float64(buttonHeight)) && b.callback != nil {
		b.callback()
	}

	b.hovered = b.position.isHovered(float64(buttonWidth), float64(buttonHeight))
	if b.hovered {
		if b.currentScale < buttonScaleHovered {
			diff := buttonScaleHovered - b.currentScale
			if diff > buttonScaleMaxChangePerTick {
				diff = buttonScaleMaxChangePerTick
			}
			b.currentScale += diff
		}
	} else {
		if b.currentScale > buttonScaleNotHovered {
			diff := b.currentScale - buttonScaleNotHovered
			if diff > buttonScaleMaxChangePerTick {
				diff = buttonScaleMaxChangePerTick
			}
			b.currentScale -= diff
		}
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	img := b.imgNotHovered
	if b.hovered {
		img = b.imgHovered
	}

	buttonWidth, buttonHeight := b.buttonSize()
	eckeObenLinksX, eckeObenLinksY := b.position.eckeObenLinks(float64(buttonWidth), float64(buttonHeight))
	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Scale(b.currentScale, b.currentScale)
	drawOptions.GeoM.Translate(eckeObenLinksX, eckeObenLinksY)
	drawOptions.GeoM.Translate(
		-float64(b.position.AnchorHorizontal)*(b.currentScale-buttonScaleNotHovered)*float64(buttonWidth),
		-float64(b.position.AnchorVertikal)*(b.currentScale-buttonScaleNotHovered)*float64(buttonHeight),
	)
	screen.DrawImage(img, &drawOptions)
}
