package herderlegacy

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type HerderLegacy interface {
	SpielerName() string

	VerhinderteStunden() float64

	AddVerhinderteStunden(stunden float64)

	AudioContext() *audio.Context

	CurrentScreen() Screen

	OpenScreen(screen Screen)
}

type Screen interface {
	Update()

	Draw(screen *ebiten.Image)
}
