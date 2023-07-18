package dialog

import (
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Antwort struct {
	Text   string
	Screen func() herderlegacy.Screen
}

func NewAntwort(text string, screen func() herderlegacy.Screen) Antwort {
	return Antwort{
		Text:   text,
		Screen: screen,
	}
}

type dialogScreen struct {
	title *ui.Title

	text              string
	textWidget        *ui.Text
	angezeigteZeichen int

	antworten     []Antwort
	antwortKnöpfe []*ui.Button
}

func NewDialogScreen(
	herderLegacy herderlegacy.HerderLegacy,
	person string,
	text string,
	antworten ...Antwort,
) herderlegacy.Screen {
	antwortKnoepfe := make([]*ui.Button, len(antworten))
	for i, antwort := range antworten {
		antwort := antwort
		antwortKnoepfe[i] = ui.NewButton(ui.ButtonConfig{
			Position: ui.NewCenteredPosition(ui.Width/2, ui.Height-float64(len(antworten))*80+float64(i)*80),
			Text:     antwort.Text,
			Callback: func() {
				herderLegacy.OpenScreen(antwort.Screen())
			},
			Disabled: true,
		})
	}

	return &dialogScreen{
		title: ui.NewTitle(ui.TitleConfig{
			Position: ui.NewCenteredPosition(ui.Width/2, 100),
			Text:     person + ":",
		}),

		text: text,
		textWidget: ui.NewText(ui.TextConfig{
			Position: ui.Position{
				X:                ui.Width / 2,
				Y:                175,
				AnchorHorizontal: ui.HorizontalerAnchorMitte,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text: "",
		}),
		angezeigteZeichen: 0,

		antworten:     antworten,
		antwortKnöpfe: antwortKnoepfe,
	}
}

func (d *dialogScreen) Update() {
	d.title.Update()
	d.textWidget.Update()
	for _, antwortKnopf := range d.antwortKnöpfe {
		antwortKnopf.Update()
	}

	if d.angezeigteZeichen == len(d.text) {
		for _, antwortKnopf := range d.antwortKnöpfe {
			antwortKnopf.SetDisabled(false)
		}
		return
	}

	d.angezeigteZeichen++
	d.textWidget.SetText(d.text[0:d.angezeigteZeichen])
}

func (d *dialogScreen) Draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)
	d.title.Draw(screen)
	d.textWidget.Draw(screen)
	for _, antwortKnopf := range d.antwortKnöpfe {
		antwortKnopf.Draw(screen)
	}
}
