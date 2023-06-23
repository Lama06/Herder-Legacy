package spiel

import "github.com/hajimehoshi/ebiten/v2"

type HerderLegacy interface {
	SpielerName() string

	VerhinderteStunden() float64

	AddVerhinderteStunden(stunden float64)
}

type Spiel interface {
	Update() (beendet bool)

	Draw(screen *ebiten.Image)
}

type Constructor func(HerderLegacy) Spiel
