package breakout

import (
	"image/color"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func NewFreierModusScreen(
	herderLegacy herderlegacy.HerderLegacy,
	breakoutBeendenCallback func() herderlegacy.Screen,
) herderlegacy.Screen {
	return ui.NewSelectScreen(herderLegacy, ui.SelectScreenConfig{
		Title: "Breakout",
		Text:  "Herr Hammdorfs Lieblingsspiel!",

		CancelText:   "Breakout beenden",
		CancelAction: breakoutBeendenCallback,

		AuswahlMöglichkeiten: []ui.SelectScreenAuswahlMöglichkeit{
			{
				Text: "Level auswählen",
				Action: func() herderlegacy.Screen {
					return newLevelSelectScreen(herderLegacy, breakoutBeendenCallback)
				},
			},
			{
				Text: "Anleitung",
				Action: func() herderlegacy.Screen {
					return newAnleitungScreen(herderLegacy, breakoutBeendenCallback)
				},
			},
			{
				Text: "Hilfe zu Upgrades",
				Action: func() herderlegacy.Screen {
					return newHilfeZuUpgradesScreen(herderLegacy, breakoutBeendenCallback)
				},
			},
			{
				Text: "Tipps",
				Action: func() herderlegacy.Screen {
					return newTippsScreen(herderLegacy, breakoutBeendenCallback)
				},
			},
		},
	})
}

func newAnleitungScreen(
	herderLegacy herderlegacy.HerderLegacy,
	breakoutBeendenCallback func() herderlegacy.Screen,
) herderlegacy.Screen {
	return ui.NewMessageScreen(herderLegacy, ui.MessageScreenConfig{
		Title: "Anleitung",
		Text: `Das Ziel des Spieles besteht darin, alle Steine zu zerstören, denn nach deren Zerstörung gewinnt man.
Steine zerstört man, indem man sie mit dem Ball abschießt.
Man muss aufpassen, dass der Ball nicht aus den Seiten herausfliegt, was man verhindern kann,
indem man seine Plattform mit der Maus bzw. dem Finger bewegt und mit ihr den Ball wegschießt.
Sollte der Ball aus dem Spielfeld fliegen, verliert man.`,
		ContinueText: "Alles verstanden! Zurück zum Hauptmenü",
		ContinueAction: func() herderlegacy.Screen {
			return NewFreierModusScreen(herderLegacy, breakoutBeendenCallback)
		},
	})
}

func newTippsScreen(
	herderLegacy herderlegacy.HerderLegacy,
	breakoutBeendenCallback func() herderlegacy.Screen,
) herderlegacy.Screen {
	return ui.NewMessageScreen(herderLegacy, ui.MessageScreenConfig{
		Title: "Tipps",
		Text: `Die Kollisionen sind nicht realistisch.
Wenn der Ball beispielsweise von oben auf ein Objekt trifft, prallt er nach links oder rechts ab,
jenachdem wie weit er das Objekt links bzw. rechts trifft. 
Das kannst du dir in zweierlei Hinsicht zu Nutze machen.
Einerseits kannst vorraussagen, wie der Ball fliegen wird und so deine Plattform vorrausplanend dorthin bewegen.
Andererseits kannst du gezielt z.B. Upgrades abschießen, indem du deine Plattform korrekt ausrichtest.`,
		ContinueText: "Alles verstanden! Zurück zum Hauptmenü",
		ContinueAction: func() herderlegacy.Screen {
			return NewFreierModusScreen(herderLegacy, breakoutBeendenCallback)
		},
	})
}

func newHilfeZuUpgradesScreen(
	herderLegacy herderlegacy.HerderLegacy,
	breakoutBeendenCallback func() herderlegacy.Screen,
) herderlegacy.Screen {
	return ui.NewMessageScreen(herderLegacy, ui.MessageScreenConfig{
		Title: "Erklärung der Upgrades",
		Text: `Farbig markierte Steine enthalten Upgrades.
Wenn du mit einem Ball eins triffst, fällt es in Form einer farbigen Kugel herunter.
Durch deren Aufsammeln, schaltest du das Upgrade frei.

Folgende Upgrades gibt es:

Königsblau: Deine Plattform wird zu einer Kanone, deren Kugeln Steine zerstören
Dunkelorange: Das Spielfeld wird zu einem Regenbogen
Gelb: Deine Plattform kann schneller bewegt werden
Hellgrün: Du kannst entspannen, während deine Plattform automatisch Kugeln wegschießt und Upgrades sammelt
Weiß: In diesem Upgrade kann alles drinstecken
Lila: Die Zeit wird kurz beschleunigt
Pink: Das Spiel läuft in Zeitlupe
Hellblau: Die Upgradesteine werden zufällig neu belegt und es kommen ein paar neue dazu

Alle Upgrades sind zeitlich begrenzt.
Die Dauer und Intensität der Upgrades sind zufällig.`,
		ContinueText: "Alles verstanden! Zurück zum Hauptmenü",
		ContinueAction: func() herderlegacy.Screen {
			return NewFreierModusScreen(herderLegacy, breakoutBeendenCallback)
		},
	})
}

func newLevelSelectScreen(
	herderLegacy herderlegacy.HerderLegacy,
	breakoutBeendenCallback func() herderlegacy.Screen,
) herderlegacy.Screen {
	auswahlMoeglichkeiten := make([]ui.SelectScreenAuswahlMöglichkeit, len(breakoutLevelListe))
	for i, breakoutLevel := range breakoutLevelListe {
		breakoutLevel := breakoutLevel
		auswahlMoeglichkeiten[i] = ui.SelectScreenAuswahlMöglichkeit{
			Text: breakoutLevel.name,
			Action: func() herderlegacy.Screen {
				return NewBreakoutScreen(herderLegacy, breakoutLevel.worldCreator(), func(gewonnen bool) herderlegacy.Screen {
					return newGameOverScreen(herderLegacy, breakoutBeendenCallback, gewonnen)
				})
			},
		}
	}

	return ui.NewSelectScreen(herderLegacy, ui.SelectScreenConfig{
		Title: "Level auswählen",
		Text:  "Wähle ein Breakoutlevel aus:",

		CancelText: "Zurück",
		CancelAction: func() herderlegacy.Screen {
			return NewFreierModusScreen(herderLegacy, breakoutBeendenCallback)
		},

		AuswahlMöglichkeiten: auswahlMoeglichkeiten,
	})
}

func newGameOverScreen(
	herderLegacy herderlegacy.HerderLegacy,
	breakoutBeendenCallback func() herderlegacy.Screen,
	gewonnen bool,
) herderlegacy.Screen {
	var title string
	if gewonnen {
		title = "Gewonnen!"
	} else {
		title = "Verloren"
	}

	var text string
	if gewonnen {
		text = "Du hast alle Steine zerstört und gewonnen"
	} else {
		text = "Du hast alle Bälle verloren und deswegen verloren"
	}

	return ui.NewDecideScreen(herderLegacy, ui.DecideScreenConfig{
		Title: title,
		Text:  text,

		CancelText:   "Breakout beenden",
		CancelAction: breakoutBeendenCallback,

		ConfirmText: "Zurück zur Levelauswahl",
		ConfirmAction: func() herderlegacy.Screen {
			return newLevelSelectScreen(herderLegacy, breakoutBeendenCallback)
		},
	})
}

type breakoutScreen struct {
	herderLegacy   herderlegacy.HerderLegacy
	nächsterScreen func(gewonnen bool) herderlegacy.Screen
	world          *world
	aufgebenKnopf  *ui.Button
}

var _ herderlegacy.Screen = (*breakoutScreen)(nil)

func NewBreakoutScreen(
	herderLegacy herderlegacy.HerderLegacy,
	world *world,
	nächsterScreen func(gewonnen bool) herderlegacy.Screen,
) herderlegacy.Screen {
	return &breakoutScreen{
		herderLegacy:   herderLegacy,
		nächsterScreen: nächsterScreen,
		world:          world,
		aufgebenKnopf: ui.NewButton(ui.ButtonConfig{
			Position: ui.Position{
				X:                15,
				Y:                ui.Height - 15,
				AnchorHorizontal: ui.HorizontalerAnchorLinks,
				AnchorVertikal:   ui.VertikalerAnchorUnten,
			},
			Text:               "Aufgeben",
			CustomColorPalette: true,
			ColorPalette: ui.ButtonColorPalette{
				BackgroundColor: color.RGBA{
					R: colornames.Red.R,
					G: colornames.Red.G,
					B: colornames.Red.B,
					A: 80,
				},
				BackgroundColorHovered: colornames.Darkred,
				TextColor:              color.White,
			},
			Callback: func() {
				herderLegacy.OpenScreen(nächsterScreen(false))
			},
		}),
	}
}

func (b *breakoutScreen) Update() {
	b.world.update()
	if b.world.gewonnen() {
		b.herderLegacy.OpenScreen(b.nächsterScreen(true))
		return
	}
	if b.world.verloren() {
		b.herderLegacy.OpenScreen(b.nächsterScreen(false))
	}
	b.aufgebenKnopf.Update()
}

func (b *breakoutScreen) Draw(screen *ebiten.Image) {
	b.world.draw(screen)
	b.aufgebenKnopf.Draw(screen)
}
