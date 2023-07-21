package quiz

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
)

var (
	relationsQuizLetzterGewinnerBtnPos = ui.NewCenteredPosition(300, ui.Height/2)
	relationsQuizVergleichBtnPos       = ui.NewCenteredPosition(ui.Width-300, ui.Height/2)
)

func moveBtnTowards(btn *ui.Button, targetPosition ui.Position) {
	const maxSpeed = 15

	speedX := targetPosition.X - btn.Position().X
	speedX = math.Max(-maxSpeed, math.Min(maxSpeed, speedX))
	newX := btn.Position().X + speedX

	speedY := targetPosition.Y - btn.Position().Y
	speedY = math.Max(-maxSpeed, math.Min(maxSpeed, speedY))
	newY := btn.Position().Y + speedY

	btn.SetPosition(ui.Position{
		X:                newX,
		Y:                newY,
		AnchorHorizontal: btn.Position().AnchorHorizontal,
		AnchorVertikal:   btn.Position().AnchorVertikal,
	})
}

type RelationQuizConfig struct {
	Name         string
	Frage        string
	ZeitProFrage int
	Werte        map[string]int
}

type RelationQuizAuswertung struct {
	RichtigeAntworten int
	FalscheAntworten  int
}

func NewRelationsQuizScreen(
	herderLegacy herderlegacy.HerderLegacy,
	config RelationQuizConfig,
	quizBeendetCallback func(RelationQuizAuswertung) herderlegacy.Screen,
) herderlegacy.Screen {

	relationen := make([]string, 0, len(config.Werte))
	for relation := range config.Werte {
		relationen = append(relationen, relation)
	}

	rand.Shuffle(len(relationen), func(i, j int) {
		relationen[i], relationen[j] = relationen[j], relationen[i]
	})

	return ui.NewMessageScreen(herderLegacy, ui.MessageScreenConfig{
		Title:        config.Name,
		Text:         "Auf die Plätze. Fertig. Los!",
		ContinueText: "Bereit!",
		ContinueAction: func() herderlegacy.Screen {
			return newRelationsQuizFrageScreen(
				herderLegacy,
				config,
				false,
				"",
				ui.Position{},
				relationen[0],
				relationsQuizLetzterGewinnerBtnPos,
				relationen[1],
				relationen[2:],
				RelationQuizAuswertung{},
				quizBeendetCallback,
			)
		},
	})
}

type relationsQuizFrageScreen struct {
	herderLegacy           herderlegacy.HerderLegacy
	config                 RelationQuizConfig
	letzterGewinner        string
	vergleich              string
	verbleibendeVergleiche []string
	auswertung             RelationQuizAuswertung
	quizBeendetCallback    func(RelationQuizAuswertung) herderlegacy.Screen

	verbleibendeZeit int

	aufgebenKnopf *ui.Button
	countdown     *ui.Title
	statistik     *ui.Text
	frage         *ui.Title

	letzterVerlierenButton *ui.Button
	letzterGewinnerButton  *ui.Button
	vergleichsButton       *ui.Button
}

func newRelationsQuizFrageScreen(
	herderLegacy herderlegacy.HerderLegacy,
	config RelationQuizConfig,

	hatLetztenVerlierer bool,
	letzterVerlierer string,
	letzterVerliererButtonPos ui.Position,

	letzterGewinner string,
	letzterGewinnerButtonPos ui.Position,

	vergleich string,

	verbleibendeVergleiche []string,
	auswertung RelationQuizAuswertung,
	quizBeendetCallback func(RelationQuizAuswertung) herderlegacy.Screen,
) *relationsQuizFrageScreen {
	screen := relationsQuizFrageScreen{
		herderLegacy:           herderLegacy,
		config:                 config,
		letzterGewinner:        letzterGewinner,
		vergleich:              vergleich,
		verbleibendeVergleiche: verbleibendeVergleiche,
		auswertung:             auswertung,
		quizBeendetCallback:    quizBeendetCallback,

		verbleibendeZeit: config.ZeitProFrage,

		aufgebenKnopf: ui.NewButton(ui.ButtonConfig{
			Position: ui.Position{
				X:                20,
				Y:                20,
				AnchorHorizontal: ui.HorizontalerAnchorLinks,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text:               "Aufgeben",
			CustomColorPalette: true,
			ColorPalette:       ui.CancelButtonColorPalette,
			Callback: func() {
				auswertung.FalscheAntworten += len(verbleibendeVergleiche)
				herderLegacy.OpenScreen(quizBeendetCallback(auswertung))
			},
		}),
		countdown: ui.NewTitle(ui.TitleConfig{
			Position: ui.Position{
				X:                ui.Width / 2,
				Y:                20,
				AnchorHorizontal: ui.HorizontalerAnchorMitte,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
		}),
		statistik: ui.NewText(ui.TextConfig{
			Position: ui.Position{
				X:                ui.Width - 20,
				Y:                20,
				AnchorHorizontal: ui.HorizontalerAnchorRechts,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text: fmt.Sprintf(
				"Richtig: %v\nFalsch: %v\nVerbleibend: %v",
				auswertung.RichtigeAntworten,
				auswertung.FalscheAntworten,
				len(verbleibendeVergleiche),
			),
		}),
		frage: ui.NewTitle(ui.TitleConfig{
			Position: ui.NewCenteredPosition(ui.Width/2, 250),
			Text:     config.Frage,
		}),
	}

	if hatLetztenVerlierer {
		screen.letzterVerlierenButton = ui.NewButton(ui.ButtonConfig{
			Position: letzterVerliererButtonPos,
			Text:     letzterVerlierer,
			Disabled: true,
		})
	}

	screen.letzterGewinnerButton = ui.NewButton(ui.ButtonConfig{
		Position: letzterGewinnerButtonPos,
		Text:     letzterGewinner,
		Silent:   true,
		Callback: func() {
			screen.handleSelect(letzterGewinner, screen.letzterGewinnerButton, vergleich, screen.vergleichsButton)
		},
	})

	screen.vergleichsButton = ui.NewButton(ui.ButtonConfig{
		Position: ui.NewCenteredPosition(ui.Width+100, ui.Height/2),
		Text:     vergleich,
		Silent:   true,
		Callback: func() {
			screen.handleSelect(vergleich, screen.vergleichsButton, letzterGewinner, screen.letzterGewinnerButton)
		},
	})

	return &screen
}

func (r *relationsQuizFrageScreen) handleSelect(
	ausgewählt string,
	ausgewähltButton *ui.Button,
	alternative string,
	alternativeButton *ui.Button,
) {
	gewonnen := r.config.Werte[ausgewählt] > r.config.Werte[alternative]

	if gewonnen {
		r.auswertung.RichtigeAntworten++
		richtigSound.Rewind()
		richtigSound.Play()
	} else {
		r.auswertung.FalscheAntworten++
		falschSound.Rewind()
		falschSound.Play()
	}

	if len(r.verbleibendeVergleiche) == 0 {
		r.herderLegacy.OpenScreen(r.quizBeendetCallback(r.auswertung))
		return
	}

	if gewonnen {
		r.herderLegacy.OpenScreen(newRelationsQuizFrageScreen(
			r.herderLegacy,
			r.config,

			true,
			alternative,
			alternativeButton.Position(),

			ausgewählt,
			ausgewähltButton.Position(),

			r.verbleibendeVergleiche[0],
			r.verbleibendeVergleiche[1:],
			r.auswertung,
			r.quizBeendetCallback,
		))
		return
	}
	r.herderLegacy.OpenScreen(newRelationsQuizFrageScreen(
		r.herderLegacy,
		r.config,

		true,
		ausgewählt,
		ausgewähltButton.Position(),

		alternative,
		alternativeButton.Position(),

		r.verbleibendeVergleiche[0],
		r.verbleibendeVergleiche[1:],
		r.auswertung,
		r.quizBeendetCallback,
	))
}

func (r *relationsQuizFrageScreen) components() []ui.Component {
	components := []ui.Component{
		r.aufgebenKnopf,
		r.countdown,
		r.statistik,
		r.frage,
		r.letzterGewinnerButton,
		r.vergleichsButton,
	}
	if r.letzterVerlierenButton != nil {
		components = append(components, r.letzterVerlierenButton)
	}
	return components
}

func (r *relationsQuizFrageScreen) Draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)
	for _, component := range r.components() {
		component.Draw(screen)
	}
}

func (r *relationsQuizFrageScreen) Update() {
	for _, component := range r.components() {
		component.Update()
	}

	if r.verbleibendeZeit <= 0 {
		if r.config.Werte[r.letzterGewinner] > r.config.Werte[r.vergleich] {
			r.vergleichsButton.Callback()()
		} else {
			r.letzterGewinnerButton.Callback()()
		}
		return
	}
	r.verbleibendeZeit--

	r.countdown.SetText(strconv.Itoa(r.verbleibendeZeit/60 + 1))
	if r.verbleibendeZeit > 3*60 {
		r.countdown.SetColorPalette(ui.TitleColorPalette{
			Color:      colornames.Green,
			HoverColor: colornames.Darkgreen,
		})
	} else {
		r.countdown.SetColorPalette(ui.TitleColorPalette{
			Color:      colornames.Red,
			HoverColor: colornames.Darkred,
		})
	}

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyLeft):
		r.letzterGewinnerButton.Callback()()
	case inpututil.IsKeyJustPressed(ebiten.KeyRight):
		r.vergleichsButton.Callback()()
	}

	moveBtnTowards(r.letzterGewinnerButton, relationsQuizLetzterGewinnerBtnPos)
	moveBtnTowards(r.vergleichsButton, relationsQuizVergleichBtnPos)
	if r.letzterVerlierenButton != nil {
		moveBtnTowards(r.letzterVerlierenButton, ui.Position{
			X:                r.letzterVerlierenButton.Position().X,
			Y:                math.Inf(1),
			AnchorHorizontal: r.letzterVerlierenButton.Position().AnchorHorizontal,
			AnchorVertikal:   r.letzterVerlierenButton.Position().AnchorVertikal,
		})
	}
}
