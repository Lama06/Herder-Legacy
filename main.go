package main

import (
	_ "image/jpeg"
	_ "image/png"
	"log"
	"runtime"

	"github.com/Lama06/Herder-Legacy/breakout"
	"github.com/Lama06/Herder-Legacy/dame"
	"github.com/Lama06/Herder-Legacy/dialog"
	"github.com/Lama06/Herder-Legacy/flappyoinky"
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/jacobsalptraum"
	"github.com/Lama06/Herder-Legacy/leitstellenspiel"
	"github.com/Lama06/Herder-Legacy/openworld"
	"github.com/Lama06/Herder-Legacy/passwortdreher"
	"github.com/Lama06/Herder-Legacy/poker"
	"github.com/Lama06/Herder-Legacy/quiz"
	"github.com/Lama06/Herder-Legacy/sodoku"
	"github.com/Lama06/Herder-Legacy/stabwelle"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type herderLegacy struct {
	audioContext       *audio.Context
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

func (h *herderLegacy) AudioContext() *audio.Context {
	return h.audioContext
}

func (h *herderLegacy) CurrentScreen() herderlegacy.Screen {
	return h.currentScreen
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
	herderLegacy := herderLegacy{
		audioContext: audio.NewContext(44100),
	}

	var newMenuScreen func() herderlegacy.Screen
	newMenuScreen = func() herderlegacy.Screen {
		return dialog.NewDialogScreen(
			&herderLegacy,
			"Herr Weber",
			"Was willst du spielen",
			dialog.NewAntwort("Openworld", func() herderlegacy.Screen {
				return openworld.NewOpenWorldScreen(&herderLegacy)
			}),
			dialog.NewAntwort("Breakout", func() herderlegacy.Screen {
				return breakout.NewFreierModusScreen(
					&herderLegacy,
					newMenuScreen,
				)
			}),
			dialog.NewAntwort("Passwort knacken", func() herderlegacy.Screen {
				return passwortdreher.NewPasswortDreherScreen(
					&herderLegacy,
					func(erfolg bool) herderlegacy.Screen {
						return newMenuScreen()
					},
					3,
				)
			}),
			dialog.NewAntwort("Stabwellen", func() herderlegacy.Screen {
				return stabwelle.NewStabwelleScreen(
					&herderLegacy,
					func(erfolg bool) herderlegacy.Screen {
						return newMenuScreen()
					},
					3,
				)
			}),
			dialog.NewAntwort("Hauptstadt Quiz", func() herderlegacy.Screen {
				return quiz.NewFreierModusScreen(
					&herderLegacy,
					newMenuScreen,
				)
			}),
			dialog.NewAntwort("Dame", func() herderlegacy.Screen {
				return dame.NewFreierModusScreen(
					&herderLegacy,
					newMenuScreen,
				)
			}),
			dialog.NewAntwort("Flappyoinky", func() herderlegacy.Screen {
				return flappyoinky.NewFlappyOinkyScreen(
					&herderLegacy,
					func(i int) herderlegacy.Screen {
						return newMenuScreen()
					},
				)
			}),
			dialog.NewAntwort("Lss", func() herderlegacy.Screen {
				return leitstellenspiel.NewLeitstellenspielBotScreen(
					&herderLegacy,
					func() herderlegacy.Screen {
						return newMenuScreen()
					},
				)
			}),
			dialog.NewAntwort("Poker", func() herderlegacy.Screen {
				return poker.NewSpielScreen(&herderLegacy, 20, func(jettons int) herderlegacy.Screen {
					return newMenuScreen()
				})
			}),
			dialog.NewAntwort("Poker Rechner", func() herderlegacy.Screen {
				return poker.NewPokerRechnerScreen(&herderLegacy, newMenuScreen)
			}),
			dialog.NewAntwort("Schach", func() herderlegacy.Screen {
				return jacobsalptraum.NewSpielScreen(
					&herderLegacy,
					jacobsalptraum.SpielConfig{
						AiStärke:      5,
						StartBrett:    jacobsalptraum.StandardBrett,
						Regeln:        jacobsalptraum.StandardRegeln,
						MenschSpieler: jacobsalptraum.SpielerVierGewinnt,
					},
					func(gewonnen bool) herderlegacy.Screen {
						log.Println(gewonnen)
						return newMenuScreen()
					},
				)
			}),
			dialog.NewAntwort("Sodoku", func() herderlegacy.Screen {
				return sodoku.NewSodokuScreen(&herderLegacy)
			}),
		)
	}

	herderLegacy.OpenScreen(dialog.NewDialogScreen(
		&herderLegacy,
		"Herr Weber",
		`Willkommen zu Herrder Games, dem neuem Projekt von WuW`,
		dialog.NewAntwort("Weiter", newMenuScreen),
	))

	err := ebiten.RunGame(&herderLegacy)
	if err != nil {
		panic(err)
	}
}
