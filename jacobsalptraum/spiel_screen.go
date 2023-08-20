package jacobsalptraum

import (
	"math"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/minimax"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func clickedPosition() (x, y int, ok bool) {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y = ebiten.CursorPosition()
		return x, y, true
	}

	if touchIds := inpututil.AppendJustReleasedTouchIDs(nil); len(touchIds) == 1 {
		x, y = inpututil.TouchPositionInPreviousTick(touchIds[0])
		return x, y, true
	}

	return 0, 0, false
}

var spielScreenBrettDrawOptions = brettDrawOptions{
	x:      0,
	y:      0,
	width:  ui.Width,
	height: ui.Height,
}

type spielScreen struct {
	herderLegacy herderlegacy.HerderLegacy
	callback     func(gewonnen bool) herderlegacy.Screen

	menschSpieler Spieler
	aiStärke      int
	regeln        Regeln

	brett Brett

	istPositionAusgewält bool
	ausgewähltePosition  position

	bewegenderStein  *bewegenderStein
	zugSpieler       Spieler
	zugErgebnisBrett Brett

	aufgebenKnopf *ui.Button
}

type SpielConfig struct {
	AiStärke      int
	StartBrett    Brett
	Regeln        Regeln
	MenschSpieler Spieler
}

func NewSpielScreen(
	herderLegacy herderlegacy.HerderLegacy,
	config SpielConfig,
	callback func(gewonnen bool) herderlegacy.Screen,
) herderlegacy.Screen {
	return &spielScreen{
		herderLegacy: herderLegacy,
		callback:     callback,

		menschSpieler: config.MenschSpieler,
		aiStärke:      config.AiStärke,
		regeln:        config.Regeln,

		brett: config.StartBrett.clone(),

		aufgebenKnopf: ui.NewButton(ui.ButtonConfig{
			Position: ui.Position{
				X:                10,
				Y:                10,
				AnchorHorizontal: ui.HorizontalerAnchorLinks,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text:               "Aufgeben",
			CustomColorPalette: true,
			ColorPalette:       ui.CancelButtonColorPalette,
			Callback: func() {
				herderLegacy.OpenScreen(callback(false))
			},
		}),
	}
}

func (s *spielScreen) Draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)
	s.brett.draw(spielScreenBrettDrawOptions, screen, s.istPositionAusgewält, s.ausgewähltePosition)
	if s.bewegenderStein != nil {
		s.bewegenderStein.draw(screen)
	}
	s.aufgebenKnopf.Draw(screen)
}

func (s *spielScreen) vierGewinntZugAusführen(zug vierGewinntZug) {
	s.zugErgebnisBrett = zug.ergebnis

	startX, startY := s.brett.brettPositionToScreenPosition(spielScreenBrettDrawOptions, position{
		zeile:  0,
		spalte: zug.position.spalte,
	})
	zielX, zielY := s.brett.brettPositionToScreenPosition(spielScreenBrettDrawOptions, zug.position)
	s.bewegenderStein = &bewegenderStein{
		aktuellesX: startX,
		aktuellesY: startY,
		zielX:      zielX,
		zielY:      zielY,
		size:       s.brett.calculateFeldSize(spielScreenBrettDrawOptions),
		img:        feldVierGewinntStein.image(),
	}

	s.zugSpieler = SpielerVierGewinnt
}

func (s *spielScreen) schachZugAusführen(zug schachZug) {
	zugFeld := s.brett.zeilen[zug.start.zeile][zug.start.spalte]
	s.brett.zeilen[zug.start.zeile][zug.start.spalte] = feldLeer

	s.zugErgebnisBrett = zug.ergebnis

	startX, startY := s.brett.brettPositionToScreenPosition(spielScreenBrettDrawOptions, zug.start)
	zielX, zielY := s.brett.brettPositionToScreenPosition(spielScreenBrettDrawOptions, zug.ziel)
	s.bewegenderStein = &bewegenderStein{
		aktuellesX: startX,
		aktuellesY: startY,
		zielX:      zielX,
		zielY:      zielY,
		size:       s.brett.calculateFeldSize(spielScreenBrettDrawOptions),
		img:        zugFeld.image(),
	}

	s.zugSpieler = SpielerSchach
}

func (s *spielScreen) zugAuswühren(zug zug) {
	switch konkreterZug := zug.(type) {
	case vierGewinntZug:
		s.vierGewinntZugAusführen(konkreterZug)
	case schachZug:
		s.schachZugAusführen(konkreterZug)
	}
}

func (s *spielScreen) Update() {
	s.aufgebenKnopf.Update()

	if s.bewegenderStein != nil {
		s.bewegenderStein.update()
		if s.bewegenderStein.angekommen() {
			if s.zugSpieler == s.menschSpieler {
				s.brett = s.zugErgebnisBrett
				computerZug, ok := minimax.BesterNächsterZug(s.brett, s.regeln, s.menschSpieler.gegner(), s.aiStärke)
				if !ok {
					s.bewegenderStein = nil
					return
				}
				s.zugAuswühren(computerZug)
				return
			}

			s.bewegenderStein = nil
			s.brett = s.zugErgebnisBrett
		}
		return
	}

	gewinner, ok := s.brett.gewinner(s.regeln)
	if ok {
		s.herderLegacy.OpenScreen(s.callback(gewinner == s.menschSpieler))
		return
	}

	mausX, mausY, istKlick := clickedPosition()
	if !istKlick {
		return
	}
	mausPosition, ok := s.brett.screenPositionToBrettPosition(spielScreenBrettDrawOptions, float64(mausX), float64(mausY))
	if !ok {
		return
	}

	switch s.menschSpieler {
	case SpielerVierGewinnt:
		for _, zug := range s.brett.möglicheVierGewinntZüge() {
			if zug.position.spalte == mausPosition.spalte {
				s.vierGewinntZugAusführen(zug)
				return
			}
		}
	case SpielerSchach:
		if !s.istPositionAusgewält {
			s.istPositionAusgewält = true
			s.ausgewähltePosition = mausPosition
			return
		}

		s.istPositionAusgewält = false

		for _, zug := range s.brett.möglicheSchachZüge() {
			if zug.start == s.ausgewähltePosition && zug.ziel == mausPosition {
				s.schachZugAusführen(zug)
				return
			}
		}
	}
}

type bewegenderStein struct {
	aktuellesX, aktuellesY float64
	zielX, zielY           float64
	size                   float64
	img                    *ebiten.Image
}

func (b *bewegenderStein) angekommen() bool {
	const toleranz = 0.00001
	return math.Abs(b.aktuellesX-b.zielX) <= toleranz && math.Abs(b.aktuellesY-b.zielY) <= toleranz
}

func (b *bewegenderStein) update() {
	const maxSpeed = 6
	b.aktuellesX += math.Max(-maxSpeed, math.Min(maxSpeed, b.zielX-b.aktuellesX))
	b.aktuellesY += math.Max(-maxSpeed, math.Min(maxSpeed, b.zielY-b.aktuellesY))
}

func (b *bewegenderStein) draw(screen *ebiten.Image) {
	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Scale(b.size/float64(b.img.Bounds().Dx()), b.size/float64(b.img.Bounds().Dy()))
	drawOptions.GeoM.Translate(b.aktuellesX, b.aktuellesY)
	screen.DrawImage(b.img, &drawOptions)
}
