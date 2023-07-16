package main

import (
	"github.com/Lama06/Herder-Legacy/dame"
	"github.com/Lama06/Herder-Legacy/dialog"
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"runtime"
)

type herderLegacy struct {
	currentScreen      herderlegacy.Screen
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

func (h *herderLegacy) OpenScreen(screen herderlegacy.Screen) {
	h.currentScreen = screen
}

func (h *herderLegacy) Update() error {
	if runtime.GOOS == "js" && !ebiten.IsFullscreen() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ebiten.SetFullscreen(true)
	}
	h.currentScreen.Update()
	return nil
}

func (h *herderLegacy) Draw(screen *ebiten.Image) {
	h.currentScreen.Draw(screen)
}

func (h *herderLegacy) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ui.Width, ui.Height
}

func main() {
	ebiten.SetWindowTitle("Herder Legacy")
	ebiten.SetFullscreen(true)
	herderLegacy := herderLegacy{}
	herderLegacy.currentScreen = dialog.NewDialogScreen(
		&herderLegacy,
		"Herr Weber",
		"Hallo",
		dialog.NewAntwort("Tschau", func() herderlegacy.Screen {
			return dame.NewFreierModusScreen(&herderLegacy, func() herderlegacy.Screen {
				return dialog.NewDialogScreen(
					&herderLegacy,
					"Herr Weber",
					"Gut gemacht",
				)
			})
		}),
		dialog.NewAntwort("Hallo", func() herderlegacy.Screen {
			return dame.NewLehrerDameSpielScreen(
				&herderLegacy,
				dame.SpielOptionen{
					StartBrett: dame.MustParseBrett(
						"_l_l_l_l",
						"l_l_l_l_",
						"_l_l_l_l",
						"________",
						"________",
						"s_s_s_s_",
						"_s_s_s_s",
						"s_s_s_s_",
					),
					ZugRegeln: dame.InternationaleZugRegeln,
					AiTiefe:   5,
				},
				func(gewonnen bool) herderlegacy.Screen {
					if gewonnen {
						return dialog.NewDialogScreen(
							&herderLegacy,
							"Herr Weber",
							"Gut gemacht",
						)
					}
					return dialog.NewDialogScreen(
						&herderLegacy,
						"Herr Weber",
						"Verdammt",
						dialog.NewAntwort("Tschau", func() herderlegacy.Screen {
							return dame.NewFreierModusScreen(&herderLegacy, func() herderlegacy.Screen {
								return dialog.NewDialogScreen(
									&herderLegacy,
									"Herr Weber",
									"Gut gemacht",
								)
							})
						}),
					)
				},
			)
		}),
	)
	err := ebiten.RunGame(&herderLegacy)
	if err != nil {
		panic(err)
	}
}
