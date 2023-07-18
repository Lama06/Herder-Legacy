package passwortdreher

import (
	"math/rand"
	"strconv"
	"unicode"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

const (
	spalteSpeed            = 6
	buchstabenProSpalte    = 13
	buchstabenAbstandY     = ui.Height / (buchstabenProSpalte - 1)
	möglicheFüllBuchstaben = "ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜ!?1234567890"
	mittelLinieHoehe       = 70
	mittelLinieY           = ui.Height/2 - mittelLinieHoehe/2
)

var (
	backgroundFarbe            = colornames.Black
	verbleibendeVersucheFarben = ui.TextColorPalatte{
		Color: colornames.Lightpink,
	}
	mitelLinieFarbe         = colornames.Pink
	falscherBuchstabeFarben = ui.TitleColorPalette{
		Color: colornames.Ghostwhite,
	}
	richtigerBuchstabeFarben = ui.TitleColorPalette{
		Color: colornames.Purple,
	}

	standardPasswörter = []string{
		"Passwort",
		"Passwort123",
		"Rüdiger",
		"Admin123",
		"HalloWelt",
		"AndreasWarHier",
		"ThomasWarHier",
		"Unknackbar",
	}
)

func randomFüllBuchstabe() rune {
	möglicheFüllBuchstabenRunes := []rune(möglicheFüllBuchstaben)
	return möglicheFüllBuchstabenRunes[rand.Intn(len(möglicheFüllBuchstabenRunes))]
}

type passworDreherScreen struct {
	herderLegacy       herderlegacy.HerderLegacy
	nächsterScreen     func(erfolg bool) herderlegacy.Screen
	möglichePasswörter []string

	spalten        []*spalte
	aktuelleSpalte int

	verbleibendeVersuche     int
	verbleibendeVersucheText *ui.Text

	aufgebenKnopf *ui.Button
}

var _ herderlegacy.Screen = (*passworDreherScreen)(nil)

func NewPasswortDreherScreen(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func(erfolg bool) herderlegacy.Screen,
	versuche int,
) herderlegacy.Screen {
	return NewPasswortDreherScreenMitCustomPasswort(
		herderLegacy,
		nächsterScreen,
		versuche,
		standardPasswörter...,
	)
}

func NewPasswortDreherScreenMitCustomPasswort(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func(erfolg bool) herderlegacy.Screen,
	versuche int,
	möglichePasswörter ...string,
) herderlegacy.Screen {
	passwort := []rune(möglichePasswörter[rand.Intn(len(möglichePasswörter))])

	spalten := make([]*spalte, len(passwort))
	spaltenAbstandX := ui.Width / float64(len(spalten))
	for i, buchstabe := range passwort {
		spalten[i] = newSpalte(float64(i)*spaltenAbstandX+spaltenAbstandX/2, buchstabe)
	}

	return &passworDreherScreen{
		herderLegacy:       herderLegacy,
		nächsterScreen:     nächsterScreen,
		möglichePasswörter: möglichePasswörter,

		spalten:        spalten,
		aktuelleSpalte: 0,

		verbleibendeVersuche: versuche,
		verbleibendeVersucheText: ui.NewText(ui.TextConfig{
			Position: ui.Position{
				X:                10,
				Y:                10,
				AnchorHorizontal: ui.HorizontalerAnchorLinks,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text:               "Verbleibende Versuche: " + strconv.Itoa(versuche),
			CustomColorPalette: true,
			ColorPalette:       verbleibendeVersucheFarben,
		}),

		aufgebenKnopf: ui.NewButton(ui.ButtonConfig{
			Position: ui.Position{
				X:                ui.Width - 10,
				Y:                10,
				AnchorHorizontal: ui.HorizontalerAnchorRechts,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text:               "Aufgeben",
			CustomColorPalette: true,
			ColorPalette:       ui.CancelButtonColorPalette,
			Callback: func() {
				herderLegacy.OpenScreen(nächsterScreen(false))
			},
		}),
	}
}

func (p *passworDreherScreen) Update() {
	var aktuelleSpalteErhöht bool
	for i, spalte := range p.spalten {
		spalte.update()
		if p.aktuelleSpalte == i && (inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
			inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			len(inpututil.AppendJustPressedTouchIDs(nil)) != 0) {
			if !spalte.istRichtigerBuchstabeAufMittellinie() {
				if p.verbleibendeVersuche >= 1 {
					p.herderLegacy.OpenScreen(NewPasswortDreherScreenMitCustomPasswort(
						p.herderLegacy,
						p.nächsterScreen,
						p.verbleibendeVersuche-1,
						p.möglichePasswörter...,
					))
					return
				}
				p.herderLegacy.OpenScreen(p.nächsterScreen(false))
				return
			}
			if p.aktuelleSpalte == len(p.spalten)-1 {
				p.herderLegacy.OpenScreen(p.nächsterScreen(true))
				return
			}
			spalte.fertig = true
			aktuelleSpalteErhöht = true
		}
	}
	if aktuelleSpalteErhöht {
		p.aktuelleSpalte++
	}

	p.verbleibendeVersucheText.Update()
	p.aufgebenKnopf.Update()
}

func (p *passworDreherScreen) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundFarbe)

	for _, spalte := range p.spalten {
		spalte.draw(screen)
	}

	vector.StrokeRect(
		screen,
		0,
		mittelLinieY,
		ui.Width,
		mittelLinieHoehe,
		7,
		mitelLinieFarbe,
		true,
	)

	p.verbleibendeVersucheText.Draw(screen)
	p.aufgebenKnopf.Draw(screen)
}

type spalte struct {
	buchstaben              []*ui.Title
	richtigerBuchstabeIndex int
	fertig                  bool
}

func newSpalte(x float64, richtigerBuchstabe rune) *spalte {
	richtigerBuchstabeIndex := rand.Intn(buchstabenProSpalte)
	buchstaben := make([]*ui.Title, buchstabenProSpalte)
	for i := range buchstaben {
		var buchstabe rune
		if i == richtigerBuchstabeIndex {
			buchstabe = unicode.ToUpper(richtigerBuchstabe)
		} else {
			buchstabe = randomFüllBuchstabe()
		}

		var farben ui.TitleColorPalette
		if i == richtigerBuchstabeIndex {
			farben = richtigerBuchstabeFarben
		} else {
			farben = falscherBuchstabeFarben
		}

		buchstaben[i] = ui.NewTitle(ui.TitleConfig{
			Position: ui.Position{
				X:                x,
				Y:                float64(i) * buchstabenAbstandY,
				AnchorHorizontal: ui.HorizontalerAnchorMitte,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text:               string(buchstabe),
			CustomColorPalette: true,
			ColorPalette:       farben,
		})
	}

	return &spalte{
		buchstaben:              buchstaben,
		richtigerBuchstabeIndex: richtigerBuchstabeIndex,
	}
}

func (s *spalte) draw(screen *ebiten.Image) {
	for _, buchstabe := range s.buchstaben {
		buchstabe.Draw(screen)
	}
}

func (s *spalte) update() {
	if s.fertig {
		return
	}

	for _, buchstabe := range s.buchstaben {
		buchstabe.Update()
		position := buchstabe.Position()
		position.Y += spalteSpeed
		if position.Y >= ui.Height {
			position.Y = -buchstabenAbstandY
		}
		buchstabe.SetPosition(position)
	}
}

func (s *spalte) istRichtigerBuchstabeAufMittellinie() bool {
	richtigerBuchstabePosition := s.buchstaben[s.richtigerBuchstabeIndex].Position()

	if richtigerBuchstabePosition.Y > mittelLinieY+mittelLinieHoehe {
		return false
	}

	if richtigerBuchstabePosition.Y+buchstabenAbstandY < mittelLinieY {
		return false
	}

	return true
}
