package stabwelle

import (
	"strconv"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

const (
	stäbeAnzahl             = 10
	stabAktivUmrandungWidth = 20

	stäbeAbstandX    = 300
	stabWidth        = (ui.Width - (stäbeAbstandX * 2)) / (stäbeAnzahl*2 - 1)
	stabHeight       = ui.Height * (2.0 / 3.0)
	stabLückeHeight  = stabHeight / 5
	stabTeilHeight   = (stabHeight - stabLückeHeight) / 2
	stabObererTeilY  = ui.Height/2 - stabLückeHeight/2 - stabTeilHeight
	stabUntererTeilY = ui.Height/2 + stabLückeHeight/2
	stabMaxYOffset   = (ui.Height - stabHeight) / 2

	stabYSpeed = stabMaxYOffset / 50

	mittelLinieHoehe = 20
	mittelLinieY     = ui.Height/2 - mittelLinieHoehe/2
)

var (
	backgroundFarbe            = colornames.Black
	verbleibendeVersucheFarben = ui.TextColorPalatte{
		Color: colornames.Lightpink,
	}
	aufgebenKnopfFarbe = ui.ButtonColorPalette{
		BackgroundColor:        colornames.Crimson,
		TextColor:              colornames.Whitesmoke,
		BackgroundColorHovered: colornames.Red,
	}
	mitelLinieFarbe    = colornames.Purple
	stabFarbe          = colornames.Hotpink
	stabUmrandungFarbe = colornames.Blanchedalmond
)

type stabwelleScreen struct {
	herderLegacy   herderlegacy.HerderLegacy
	nächsterScreen func(erfolg bool) herderlegacy.Screen

	verbleibendeVersuche     int
	verbleibendeVersucheText *ui.Text

	aufgebenKnopf *ui.Button

	aktuellerStab int
	stäbe         []*stab
}

var _ herderlegacy.Screen = (*stabwelleScreen)(nil)

func NewStabwelleScreen(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func(erfolg bool) herderlegacy.Screen,
	versuche int,
) herderlegacy.Screen {
	stäbe := make([]*stab, stäbeAnzahl)
	for i := range stäbe {
		stäbe[i] = &stab{
			x:         stäbeAbstandX + float64(i)*stabWidth*2,
			yOffset:   float64(i) * (stabMaxYOffset / stäbeAnzahl),
			yRichtung: 1,
		}
	}

	return &stabwelleScreen{
		herderLegacy:   herderLegacy,
		nächsterScreen: nächsterScreen,

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
			ColorPalette:       aufgebenKnopfFarbe,
			Callback: func() {
				herderLegacy.OpenScreen(nächsterScreen(false))
			},
		}),

		aktuellerStab: 0,
		stäbe:         stäbe,
	}
}

func (s *stabwelleScreen) Update() {
	s.verbleibendeVersucheText.Update()
	s.aufgebenKnopf.Update()

	var aktuellenStabErhöht bool
	for i, stab := range s.stäbe {
		verloren, fertig := stab.update(i == s.aktuellerStab)
		if verloren {
			if s.verbleibendeVersuche-1 >= 1 {
				s.herderLegacy.OpenScreen(NewStabwelleScreen(s.herderLegacy, s.nächsterScreen, s.verbleibendeVersuche-1))
				return
			}
			s.herderLegacy.OpenScreen(s.nächsterScreen(false))
			return
		}
		if fertig {
			if s.aktuellerStab == stäbeAnzahl-1 {
				s.herderLegacy.OpenScreen(s.nächsterScreen(true))
				return
			}
			aktuellenStabErhöht = true
		}
	}

	if aktuellenStabErhöht {
		s.aktuellerStab++
	}
}

func (s *stabwelleScreen) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundFarbe)

	for i, stab := range s.stäbe {
		stab.draw(screen, i == s.aktuellerStab)
	}

	vector.StrokeLine(
		screen,
		0,
		mittelLinieY+mittelLinieHoehe/2,
		ui.Width,
		mittelLinieY+mittelLinieHoehe/2,
		mittelLinieHoehe,
		mitelLinieFarbe,
		true,
	)

	s.verbleibendeVersucheText.Draw(screen)
	s.aufgebenKnopf.Draw(screen)
}

type stab struct {
	x         float64
	yOffset   float64
	yRichtung int
	fertig    bool
}

func (s *stab) obererTeilYPosition() float64 {
	return stabObererTeilY + s.yOffset
}

func (s *stab) untererTeilYPosition() float64 {
	return stabUntererTeilY + s.yOffset
}

func (s *stab) kollidiertTeilMitMittelLinie(teilY float64) bool {
	if mittelLinieY > teilY+stabTeilHeight {
		return false
	}

	if mittelLinieY+mittelLinieHoehe < teilY {
		return false
	}

	return true
}

func (s *stab) kollidiertMitMittellinie() bool {
	return s.kollidiertTeilMitMittelLinie(s.obererTeilYPosition()) || s.kollidiertTeilMitMittelLinie(s.untererTeilYPosition())
}

func (s *stab) bewegen() {
	if s.fertig {
		return
	}

	s.yOffset += float64(s.yRichtung) * stabYSpeed
	if s.yOffset <= -stabMaxYOffset {
		s.yOffset = -stabMaxYOffset
		s.yRichtung = 1
	}
	if s.yOffset >= stabMaxYOffset {
		s.yOffset = stabMaxYOffset
		s.yRichtung = -1
	}
}

func (s *stab) update(aktiv bool) (verloren bool, fertig bool) {
	s.bewegen()

	if aktiv && (inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustPressedTouchIDs(nil)) != 0 ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace)) {
		if s.kollidiertMitMittellinie() {
			return true, false
		}
		s.fertig = true
		return false, true
	}

	return false, false
}

func (s *stab) drawTeil(screen *ebiten.Image, teilY float64, aktiv bool) {
	vector.DrawFilledRect(
		screen,
		float32(s.x),
		float32(teilY),
		stabWidth,
		stabTeilHeight,
		stabFarbe,
		true,
	)

	if aktiv {
		vector.StrokeRect(
			screen,
			float32(s.x),
			float32(teilY),
			stabWidth,
			stabTeilHeight,
			stabAktivUmrandungWidth,
			stabUmrandungFarbe,
			true,
		)
	}
}

func (s *stab) draw(screen *ebiten.Image, aktiv bool) {
	s.drawTeil(screen, s.obererTeilYPosition(), aktiv)
	s.drawTeil(screen, s.untererTeilYPosition(), aktiv)
}
