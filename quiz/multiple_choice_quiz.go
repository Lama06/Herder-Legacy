package quiz

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
)

const multipleChoiceQuizAnzahlAntworten = 4

type MultipleChoiceQuizFrage struct {
	Frage            string
	Antwort          string
	FalscheAntworten []string
}

type MultipleChoiceQuizConfig struct {
	Name         string
	ZeitProFrage int
	Fragen       []MultipleChoiceQuizFrage
}

type MultipleChoiceQuizAuswertung struct {
	RichtigeAntworten int
	FalscheAntworten  int
}

func NewMultipleChoiceQuizScreen(
	herderLegacy herderlegacy.HerderLegacy,
	config MultipleChoiceQuizConfig,
	quizBeendetCallback func(MultipleChoiceQuizAuswertung) herderlegacy.Screen,
) herderlegacy.Screen {
	if len(config.Fragen) < multipleChoiceQuizAnzahlAntworten {
		panic("zu wenig Antworten")
	}

	rand.Shuffle(len(config.Fragen), func(i, j int) {
		config.Fragen[i], config.Fragen[j] = config.Fragen[j], config.Fragen[i]
	})

	return ui.NewMessageScreen(herderLegacy, ui.MessageScreenConfig{
		Title:        "Quiz: " + config.Name,
		Text:         "Auf die Plätze. Fertig. Los!",
		ContinueText: "Bereit!",
		ContinueAction: func() herderlegacy.Screen {
			return newMultipleChoiceFrageScreen(
				herderLegacy,
				config,
				config.Fragen[0],
				config.Fragen[1:],
				MultipleChoiceQuizAuswertung{},
				quizBeendetCallback,
			)
		},
	})
}

type multipleChoiceFrageScreen struct {
	herderLegacy        herderlegacy.HerderLegacy
	config              MultipleChoiceQuizConfig
	verbleibendeFragen  []MultipleChoiceQuizFrage
	auswertung          MultipleChoiceQuizAuswertung
	quizBeendetCallback func(MultipleChoiceQuizAuswertung) herderlegacy.Screen

	verbleibendeZeit int

	aufgebenKnopf *ui.Button
	countdown     *ui.Title
	statistik     *ui.Text
	frage         *ui.Title
	antwortKnöpfe []*ui.Button
}

var _ herderlegacy.Screen = (*multipleChoiceFrageScreen)(nil)

func newMultipleChoiceFrageScreen(
	herderLegacy herderlegacy.HerderLegacy,
	config MultipleChoiceQuizConfig,
	frage MultipleChoiceQuizFrage,
	verbleibendeFragen []MultipleChoiceQuizFrage,
	auswertung MultipleChoiceQuizAuswertung,
	quizBeendetCallback func(MultipleChoiceQuizAuswertung) herderlegacy.Screen,
) *multipleChoiceFrageScreen {
	screen := multipleChoiceFrageScreen{
		herderLegacy:        herderLegacy,
		config:              config,
		verbleibendeFragen:  verbleibendeFragen,
		auswertung:          auswertung,
		quizBeendetCallback: quizBeendetCallback,

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
				auswertung.FalscheAntworten += len(verbleibendeFragen)
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
				len(verbleibendeFragen),
			),
		}),
		frage: ui.NewTitle(ui.TitleConfig{
			Position: ui.NewCenteredPosition(ui.Width/2, 250),
			Text:     frage.Frage,
		}),
	}

	screen.antwortKnöpfe = []*ui.Button{
		ui.NewButton(ui.ButtonConfig{
			Text:   frage.Antwort,
			Silent: true,
			Callback: func() {
				richtigSound.Rewind()
				richtigSound.Play()

				auswertung.RichtigeAntworten++

				if len(verbleibendeFragen) == 0 {
					herderLegacy.OpenScreen(quizBeendetCallback(auswertung))
					return
				}

				herderLegacy.OpenScreen(newMultipleChoiceFrageScreen(
					herderLegacy,
					config,
					verbleibendeFragen[0],
					verbleibendeFragen[1:],
					auswertung,
					quizBeendetCallback,
				))
			},
		}),
	}
	var möglicheFalscheAntworten []string
	for _, andereFrage := range config.Fragen {
		if andereFrage.Antwort == frage.Antwort {
			continue
		}
		möglicheFalscheAntworten = append(möglicheFalscheAntworten, andereFrage.Antwort)
	}
	rand.Shuffle(len(möglicheFalscheAntworten), func(i, j int) {
		möglicheFalscheAntworten[i], möglicheFalscheAntworten[j] = möglicheFalscheAntworten[j], möglicheFalscheAntworten[i]
	})
	möglicheFalscheAntworten = append(frage.FalscheAntworten, möglicheFalscheAntworten...)
	for i := 0; i < multipleChoiceQuizAnzahlAntworten-1; i++ {
		screen.antwortKnöpfe = append(screen.antwortKnöpfe, ui.NewButton(ui.ButtonConfig{
			Text:   möglicheFalscheAntworten[i],
			Silent: true,
			Callback: func() {
				falschSound.Rewind()
				falschSound.Play()

				auswertung.FalscheAntworten++

				if len(verbleibendeFragen) == 0 {
					herderLegacy.OpenScreen(quizBeendetCallback(auswertung))
					return
				}

				herderLegacy.OpenScreen(newMultipleChoiceFrageScreen(
					herderLegacy,
					config,
					verbleibendeFragen[0],
					verbleibendeFragen[1:],
					auswertung,
					quizBeendetCallback,
				))
			},
		}))
	}
	rand.Shuffle(multipleChoiceQuizAnzahlAntworten, func(i, j int) {
		screen.antwortKnöpfe[i], screen.antwortKnöpfe[j] = screen.antwortKnöpfe[j], screen.antwortKnöpfe[i]
	})
	for i, antwortKnopf := range screen.antwortKnöpfe {
		antwortKnopf.SetPosition(ui.NewCenteredPosition(
			ui.Width/2,
			ui.Height-multipleChoiceQuizAnzahlAntworten*80+float64(i)*80,
		))
	}

	return &screen
}

func (f *multipleChoiceFrageScreen) components() []ui.Component {
	components := []ui.Component{f.aufgebenKnopf, f.countdown, f.statistik, f.frage}
	for _, antwortKnopf := range f.antwortKnöpfe {
		components = append(components, antwortKnopf)
	}
	return components
}

func (f *multipleChoiceFrageScreen) Update() {
	if f.verbleibendeZeit <= 0 {
		f.auswertung.FalscheAntworten++

		if len(f.verbleibendeFragen) == 0 {
			f.herderLegacy.OpenScreen(f.quizBeendetCallback(f.auswertung))
			return
		}

		f.herderLegacy.OpenScreen(newMultipleChoiceFrageScreen(
			f.herderLegacy,
			f.config,
			f.verbleibendeFragen[0],
			f.verbleibendeFragen[1:],
			f.auswertung,
			f.quizBeendetCallback,
		))
		return
	}
	f.verbleibendeZeit--

	f.countdown.SetText(strconv.Itoa(f.verbleibendeZeit/60 + 1))
	if f.verbleibendeZeit > 3*60 {
		f.countdown.SetColorPalette(ui.TitleColorPalette{
			Color:      colornames.Green,
			HoverColor: colornames.Darkgreen,
		})
	} else {
		f.countdown.SetColorPalette(ui.TitleColorPalette{
			Color:      colornames.Red,
			HoverColor: colornames.Darkred,
		})
	}

	for i, antwortKnopf := range f.antwortKnöpfe {
		if inpututil.IsKeyJustPressed(ebiten.Key1 + ebiten.Key(i)) {
			antwortKnopf.Callback()()
		}
	}

	for _, component := range f.components() {
		component.Update()
	}
}

func (f *multipleChoiceFrageScreen) Draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)
	for _, component := range f.components() {
		component.Draw(screen)
	}
}
