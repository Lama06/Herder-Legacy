package dame

import (
	"math"
	"math/rand"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/minimax"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	spielScreenBrettX         = 300
	spielScreenBrettY         = 0
	spielScreenBrettMaxBreite = ui.Width - spielScreenBrettX
	spielScreenBrettMaxHöhe   = ui.Height
)

type SpielOptionen struct {
	StartBrett Brett
	AiTiefe    int
	ZugRegeln  ZugRegeln
}

type spielScreen struct {
	herderLegacy   herderlegacy.HerderLegacy
	nächsterScreen func(gewonnen bool) herderlegacy.Screen
	optionen       SpielOptionen

	brett Brett

	aufgebenKnopf *ui.Button

	hatAusgewähltePosition bool
	ausgewähltePosition    position

	verbeibendeZugSchritte []zugSchritt
	aktuellerZugSchritt    *zugSchritt
	zugStein               *bewegenderStein

	geschlageneSteine []*bewegenderStein
}

var _ herderlegacy.Screen = (*spielScreen)(nil)

func newSpielScreen(
	herderLegacy herderlegacy.HerderLegacy,
	nächsterScreen func(gewonnen bool) herderlegacy.Screen,
	optionen SpielOptionen,
) *spielScreen {
	return &spielScreen{
		herderLegacy:   herderLegacy,
		nächsterScreen: nächsterScreen,
		brett:          optionen.StartBrett.clone(),
		optionen:       optionen,
		aufgebenKnopf: ui.NewButton(ui.ButtonConfig{
			Position: ui.Position{
				X:                10,
				Y:                10,
				AnchorHorizontal: ui.HorizontalerAnchorLinks,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text: "Aufgeben",
			Callback: func() {
				herderLegacy.OpenScreen(nächsterScreen(false))
			},
		}),
	}
}

func (s *spielScreen) zugAusführen() (zugLäuft bool) {
	if len(s.verbeibendeZugSchritte) == 0 && s.aktuellerZugSchritt == nil {
		// Gerade wird kein Zug ausgeführt
		return false
	}

	if s.zugStein != nil && !s.zugStein.angekommen() {
		// Ein Zug Schritt ist noch am Laufen
		return true
	}

	// Nächsten Zug Schritt beginnen

	if s.aktuellerZugSchritt != nil {
		s.brett = s.aktuellerZugSchritt.ergebnis
	}

	if len(s.verbeibendeZugSchritte) == 0 {
		// Kein verbleibender Schritt mehr, Zug beendet
		s.aktuellerZugSchritt = nil
		s.zugStein = nil
		return false
	}

	nächsterZugSchritt := s.verbeibendeZugSchritte[0]
	s.aktuellerZugSchritt = &nächsterZugSchritt
	s.verbeibendeZugSchritte = s.verbeibendeZugSchritte[1:]

	feldSize := s.brett.feldSize(spielScreenBrettMaxBreite, spielScreenBrettMaxHöhe)

	if nächsterZugSchritt.hatGeschlagenePosition {
		geschlagenesFeld := s.brett.feld(nächsterZugSchritt.geschlagenePosition)
		geschlagenScreenX, geschlagenScreenY := s.brett.feldPosition(
			nächsterZugSchritt.geschlagenePosition,
			spielScreenBrettX,
			spielScreenBrettY,
			spielScreenBrettMaxBreite,
			spielScreenBrettMaxHöhe,
		)
		s.geschlageneSteine = append(s.geschlageneSteine, &bewegenderStein{
			feld:     geschlagenesFeld,
			currentX: geschlagenScreenX + feldSize/2,
			currentY: geschlagenScreenY + feldSize/2,
			targetX:  float64(rand.Intn(spielScreenBrettX)),
			targetY:  float64(rand.Intn(ui.Height)),
			radius:   feldSize * 0.35,
			speed:    6,
		})

		s.brett.setFeld(nächsterZugSchritt.geschlagenePosition, feldLeer)
	}

	vonFeld := s.brett.feld(nächsterZugSchritt.von)
	vonScreenX, vonScreenY := s.brett.feldPosition(
		nächsterZugSchritt.von,
		spielScreenBrettX,
		spielScreenBrettY,
		spielScreenBrettMaxBreite,
		spielScreenBrettMaxHöhe,
	)
	zuScreenX, zuScreenY := s.brett.feldPosition(
		nächsterZugSchritt.zu,
		spielScreenBrettX,
		spielScreenBrettY,
		spielScreenBrettMaxBreite,
		spielScreenBrettMaxHöhe,
	)
	s.zugStein = &bewegenderStein{
		feld:     vonFeld,
		currentX: vonScreenX + feldSize/2,
		currentY: vonScreenY + feldSize/2,
		targetX:  zuScreenX + feldSize/2,
		targetY:  zuScreenY + feldSize/2,
		radius:   float64(feldSize) * 0.35,
		speed:    2,
	}

	s.brett.setFeld(nächsterZugSchritt.von, feldLeer)

	return true
}

func (s *spielScreen) handleClick(mausX, mausY int) {
	mausPosition, ok := s.brett.screenPositionToBrettPosition(
		float64(mausX),
		float64(mausY),
		spielScreenBrettX,
		spielScreenBrettY,
		spielScreenBrettMaxBreite,
		spielScreenBrettMaxHöhe,
	)
	if !ok {
		s.hatAusgewähltePosition = false
		return
	}

	if !s.hatAusgewähltePosition {
		mausFeld := s.brett.feld(mausPosition)
		feldEigentümer, feldHatEigentümer := mausFeld.eigentümer()
		if !feldHatEigentümer || feldEigentümer == spielerLehrer {
			s.hatAusgewähltePosition = false
			return
		}

		s.hatAusgewähltePosition = true
		s.ausgewähltePosition = mausPosition
		return
	}

	möglicheZüge := s.brett.möglicheZügeMitStartPosition(s.ausgewähltePosition, s.optionen.ZugRegeln)
	s.hatAusgewähltePosition = false

	for _, möglicherZug := range möglicheZüge {
		if möglicherZug.endPosition() == mausPosition {
			lehrerAiZug, lehrerZugGefunden := minimax.BesterNächsterZug(
				möglicherZug.ergebnis(),
				s.optionen.ZugRegeln,
				spielerLehrer,
				s.optionen.AiTiefe,
			)
			if !lehrerZugGefunden {
				s.brett = möglicherZug.ergebnis()
				return
			}
			lehrerZug := lehrerAiZug.(zug)

			s.aktuellerZugSchritt = nil
			s.zugStein = nil
			s.verbeibendeZugSchritte = append(möglicherZug, lehrerZug...)
			return
		}
	}
}

func (s *spielScreen) Update() {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		s.herderLegacy.OpenScreen(s.nächsterScreen(false))
		return
	}

	if s.brett.gewonnen(spielerSchüler, s.optionen.ZugRegeln) {
		s.herderLegacy.OpenScreen(s.nächsterScreen(true))
		return
	}

	if s.brett.gewonnen(spielerLehrer, s.optionen.ZugRegeln) {
		s.herderLegacy.OpenScreen(s.nächsterScreen(false))
		return
	}

	s.aufgebenKnopf.Update()

	for _, verlorenerStein := range s.geschlageneSteine {
		verlorenerStein.update()
	}

	if s.zugStein != nil {
		s.zugStein.update()
	}

	if zugLäuft := s.zugAusführen(); zugLäuft {
		return
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		s.handleClick(ebiten.CursorPosition())
		return
	}

	for _, touchId := range inpututil.AppendJustReleasedTouchIDs(nil) {
		s.handleClick(inpututil.TouchPositionInPreviousTick(touchId))
		return
	}
}

func (s *spielScreen) Draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)
	s.brett.draw(
		screen,
		spielScreenBrettX,
		spielScreenBrettY,
		spielScreenBrettMaxBreite,
		spielScreenBrettMaxHöhe,
		s.hatAusgewähltePosition,
		s.ausgewähltePosition,
		s.optionen.ZugRegeln,
	)
	for _, verlorenerStein := range s.geschlageneSteine {
		verlorenerStein.draw(screen)
	}
	if s.zugStein != nil {
		s.zugStein.draw(screen)
	}
	s.aufgebenKnopf.Draw(screen)
}

type bewegenderStein struct {
	feld               feld
	currentX, currentY float64
	targetX, targetY   float64
	radius             float64
	speed              float64
}

func (b *bewegenderStein) angekommen() bool {
	return b.currentX == b.targetX && b.currentY == b.targetY
}

func (b *bewegenderStein) update() {
	changeX := b.targetX - b.currentX
	if math.Abs(changeX) > b.speed {
		changeX = signum(changeX) * b.speed
		b.currentX += changeX
	} else {
		b.currentX = b.targetX
	}

	changeY := b.targetY - b.currentY
	if math.Abs(changeY) > b.speed {
		changeY = signum(changeY) * b.speed
		b.currentY += changeY
	} else {
		b.currentY = b.targetY
	}
}

func (b *bewegenderStein) draw(screen *ebiten.Image) {
	farbe, ok := b.feld.farbe()
	if !ok {
		return
	}
	vector.DrawFilledCircle(screen, float32(b.currentX), float32(b.currentY), float32(b.radius), farbe, true)
}
