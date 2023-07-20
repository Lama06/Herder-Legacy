package ui

import (
	"image/color"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
)

const (
	Width  = 1920.0
	Height = 1080.0
)

var BackgroundColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}

func Init(herderLegacy herderlegacy.HerderLegacy) {
	initButtonClickSound(herderLegacy)
}
