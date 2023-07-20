package ui

import (
	"bytes"
	_ "embed"
	"image/color"
	"strings"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
)

const (
	buttonPadding = 20

	buttonScaleNotHovered       = 1
	buttonScaleHovered          = 1.2
	buttonScaleMaxChangePerTick = 0.03
)

var (
	//go:embed button.mp3
	buttonClickSoundData []byte
	buttonClickSound     *audio.Player
)

func initButtonClickSound(herderLegacy herderlegacy.HerderLegacy) {
	context := herderLegacy.AudioContext()
	stream, err := mp3.DecodeWithSampleRate(context.SampleRate(), bytes.NewReader(buttonClickSoundData))
	if err != nil {
		panic(err)
	}
	buttonClickSound, err = context.NewPlayer(stream)
	if err != nil {
		panic(err)
	}
}

type ButtonColorPalette struct {
	BackgroundColor        color.Color
	BackgroundColorHovered color.Color
	TextColor              color.Color
	TextColorHovered       color.Color

	BackgroundColorDisabled        color.Color
	BackgroundColorHoveredDisabled color.Color
	TextColorDisabled              color.Color
	TextColorHoveredDisabled       color.Color
}

func (b ButtonColorPalette) backgroundColorHoveredOrDefault() color.Color {
	if b.BackgroundColorHovered == nil {
		return b.BackgroundColor
	}
	return b.BackgroundColorHovered
}

func (b ButtonColorPalette) textColorHoveredOrDefault() color.Color {
	if b.TextColorHovered == nil {
		return b.TextColor
	}
	return b.TextColorHovered
}

func (b ButtonColorPalette) backgroundColorDisabledOrDefault() color.Color {
	if b.BackgroundColorDisabled == nil {
		return b.BackgroundColor
	}
	return b.BackgroundColorDisabled
}

func (b ButtonColorPalette) backgroundColorHoveredDisabledOrDefault() color.Color {
	if b.BackgroundColorHoveredDisabled == nil {
		return b.backgroundColorDisabledOrDefault()
	}
	return b.BackgroundColorHoveredDisabled
}

func (b ButtonColorPalette) textColorDisabledOrDefault() color.Color {
	if b.TextColorDisabled == nil {
		return b.TextColor
	}
	return b.TextColorDisabled
}

func (b ButtonColorPalette) textColorHoveredDisabledOrDefault() color.Color {
	if b.TextColorHoveredDisabled == nil {
		return b.textColorDisabledOrDefault()
	}
	return b.TextColorHoveredDisabled
}

func (b ButtonColorPalette) backgroundColor(hovered, disabled bool) color.Color {
	if hovered {
		if disabled {
			return b.backgroundColorHoveredDisabledOrDefault()
		}
		return b.backgroundColorHoveredOrDefault()
	}
	if disabled {
		return b.backgroundColorDisabledOrDefault()
	}
	return b.BackgroundColor
}

func (b ButtonColorPalette) textColor(hovered, disabled bool) color.Color {
	if hovered {
		if disabled {
			return b.textColorHoveredDisabledOrDefault()
		}
		return b.textColorHoveredOrDefault()
	}
	if disabled {
		return b.textColorDisabledOrDefault()
	}
	return b.TextColor
}

var (
	CancelButtonColorPalette = ButtonColorPalette{
		BackgroundColor:        colornames.Red,
		BackgroundColorHovered: colornames.Darkred,
		TextColor:              colornames.Whitesmoke,
		TextColorHovered:       colornames.White,
	}
	defaultButtonColorPalette = ButtonColorPalette{
		BackgroundColor:        color.RGBA{R: 18, G: 53, B: 91, A: 255},
		BackgroundColorHovered: color.RGBA{R: 134, G: 22, B: 87, A: 255},
		TextColor:              color.RGBA{R: 212, G: 245, B: 245, A: 255},
		TextColorHovered:       color.RGBA{R: 212, G: 245, B: 245, A: 255},

		BackgroundColorDisabled:        color.RGBA{R: 42, G: 59, B: 82, A: 255},
		BackgroundColorHoveredDisabled: color.RGBA{R: 29, G: 37, B: 48, A: 255},
		TextColorDisabled:              color.RGBA{R: 212, G: 245, B: 245, A: 255},
		TextColorHoveredDisabled:       color.RGBA{R: 212, G: 245, B: 245, A: 255},
	}
)

type ButtonConfig struct {
	Position           Position
	Text               string
	CustomColorPalette bool
	ColorPalette       ButtonColorPalette
	Callback           func()
	Disabled           bool
}

type Button struct {
	position     Position
	text         string
	colorPalette ButtonColorPalette
	callback     func()
	disabled     bool

	currentScale  float64
	hovered       bool
	imgNotHovered *ebiten.Image
	imgHovered    *ebiten.Image
}

var _ Component = (*Button)(nil)

func NewButton(config ButtonConfig) *Button {
	var colorPalette ButtonColorPalette
	if config.CustomColorPalette {
		colorPalette = config.ColorPalette
	} else {
		colorPalette = defaultButtonColorPalette
	}

	button := Button{
		position:     config.Position,
		text:         config.Text,
		colorPalette: colorPalette,
		callback:     config.Callback,
		disabled:     config.Disabled,

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

func (b *Button) updateImages() {
	b.imgNotHovered = b.createImage(
		b.colorPalette.backgroundColor(false, b.disabled),
		b.colorPalette.textColor(false, b.disabled),
	)
	b.imgHovered = b.createImage(
		b.colorPalette.backgroundColor(true, b.disabled),
		b.colorPalette.textColor(true, b.disabled),
	)
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

func (b *Button) ColorPalette() ButtonColorPalette {
	return b.colorPalette
}

func (b *Button) SetColorPalette(colorPalette ButtonColorPalette) {
	b.colorPalette = colorPalette
	b.updateImages()
}

func (b *Button) Callback() func() {
	return b.callback
}

func (b *Button) SetCallback(callback func()) {
	b.callback = callback
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

	if b.position.isClicked(float64(buttonWidth), float64(buttonHeight)) && !b.disabled {
		buttonClickSound.Rewind()
		buttonClickSound.Play()

		if b.callback != nil {
			b.callback()
		}
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
