package main

import (
	_ "embed"
	"github.com/Lama06/Herder-Legacy/dame"
	"github.com/Lama06/Herder-Legacy/spiel"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"runtime"
)

type herderLegacy struct {
	currentSpiel       spiel.Spiel
	verhinderteStunden float64
}

func (h *herderLegacy) SpielerName() string {
	return "Spieler"
}

func (h *herderLegacy) VerhinderteStunden() float64 {
	return h.verhinderteStunden
}

func (h *herderLegacy) AddVerhinderteStunden(stunden float64) {
	h.verhinderteStunden += stunden
}

func (h *herderLegacy) Update() error {
	if runtime.GOOS == "js" && !ebiten.IsFullscreen() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ebiten.SetFullscreen(true)
	}
	if h.currentSpiel.Update() {
		return ebiten.Termination
	}
	return nil
}

func (h *herderLegacy) Draw(screen *ebiten.Image) {
	h.currentSpiel.Draw(screen)
}

func (h *herderLegacy) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ui.Width, ui.Height
}

func main() {
	ebiten.SetWindowTitle("Herder Legacy")
	ebiten.SetFullscreen(true)
	herderLegacy := herderLegacy{}
	herderLegacy.currentSpiel = dame.NewDameSpiel(&herderLegacy)
	err := ebiten.RunGame(&herderLegacy)
	if err != nil {
		panic(err)
	}
}
