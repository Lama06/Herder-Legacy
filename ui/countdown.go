package ui

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type CountdownColorPalette struct {
	Color      color.Color
	HoverColor color.Color

	LittleTimeColor      color.Color
	LittleTimeHoverColor color.Color
}

var defaultCountdownColorPalette = CountdownColorPalette{
	Color:                colornames.Green,
	HoverColor:           colornames.Darkgreen,
	LittleTimeColor:      colornames.Red,
	LittleTimeHoverColor: colornames.Darkred,
}

func (c CountdownColorPalette) hoverColorOrDefault() color.Color {
	if c.HoverColor == nil {
		return c.Color
	}
	return c.HoverColor
}

func (c CountdownColorPalette) littleTimeColorOrDefault() color.Color {
	if c.LittleTimeColor == nil {
		return c.hoverColorOrDefault()
	}
	return c.LittleTimeColor
}

func (c CountdownColorPalette) littleTimeHoverColorOrDefault() color.Color {
	if c.LittleTimeHoverColor == nil {
		return c.littleTimeColorOrDefault()
	}
	return c.LittleTimeHoverColor
}

type CountdownConfig struct {
	Position           Position
	StartZeit          int
	AbgelaufenCallback func()
	CustomColorPalette bool
	ColorPalette       CountdownColorPalette
}

type Countdown struct {
	colorPalette CountdownColorPalette

	title            *Title
	verbleibendeZeit int

	abgelaufenCallbackCalled bool
	abgelaufenCallback       func()
}

var _ Component = (*Countdown)(nil)

func NewCountdown(config CountdownConfig) *Countdown {
	if config.StartZeit < 0 {
		panic("Zeit negativ")
	}

	colorPalette := defaultCountdownColorPalette
	if config.CustomColorPalette {
		colorPalette = config.ColorPalette
	}

	return &Countdown{
		colorPalette: colorPalette,
		title: NewTitle(TitleConfig{
			Position: config.Position,
		}),
		verbleibendeZeit:         config.StartZeit,
		abgelaufenCallback:       config.AbgelaufenCallback,
		abgelaufenCallbackCalled: false,
	}
}

func (c *Countdown) Position() Position {
	return c.title.Position()
}

func (c *Countdown) SetPosition(position Position) {
	c.title.SetPosition(position)
}

func (c *Countdown) Abgelaufen() bool {
	return c.verbleibendeZeit == 0
}

func (c *Countdown) VerbleibendeZeit() int {
	return c.verbleibendeZeit
}

func (c *Countdown) SetVerbleibendeZeit(zeit int) {
	if zeit < 0 {
		panic("Zeit negativ")
	}
	c.verbleibendeZeit = zeit
	c.abgelaufenCallbackCalled = false
}

func (c *Countdown) AbgelaufenCallback() func() {
	return c.abgelaufenCallback
}

func (c *Countdown) SetAbgelaufenCallback(callback func()) {
	c.abgelaufenCallback = callback
}

func (c *Countdown) Update() {
	c.title.Update()

	if c.verbleibendeZeit > 0 {
		c.verbleibendeZeit--
		if c.verbleibendeZeit == 0 && !c.abgelaufenCallbackCalled && c.abgelaufenCallback != nil {
			c.abgelaufenCallback()
			c.abgelaufenCallbackCalled = true
		}
	}

	c.title.SetText(strconv.Itoa(c.verbleibendeZeit/60 + 1))
	if c.verbleibendeZeit > 3*60 {
		c.title.SetColorPalette(TitleColorPalette{
			Color:      c.colorPalette.Color,
			HoverColor: c.colorPalette.hoverColorOrDefault(),
		})
	} else {
		c.title.SetColorPalette(TitleColorPalette{
			Color:      c.colorPalette.littleTimeColorOrDefault(),
			HoverColor: c.colorPalette.littleTimeHoverColorOrDefault(),
		})
	}
}

func (c *Countdown) Draw(screen *ebiten.Image) {
	c.title.Draw(screen)
}
