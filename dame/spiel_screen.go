package dame

import (
	"github.com/Lama06/Herder-Legacy/ai"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math"
	"math/rand"
)

const (
	spielScreenBrettX         = 300
	spielScreenBrettY         = 0
	spielScreenBrettMaxBreite = ui.Width - spielScreenBrettX
	spielScreenBrettMaxHoehe  = ui.Height
)

type spielScreen struct {
	dame   *dameSpiel
	brett  brett
	lehrer lehrer

	aufgebenKnopf *ui.Button

	ausgewaehltePosition *position

	verbeibendeZugSchritte []zugSchritt
	aktuellerZugSchritt    *zugSchritt
	zugStein               *bewegenderStein

	geschlageneSteine []*bewegenderStein
}

var _ screen = (*spielScreen)(nil)

func newSpielScreen(dame *dameSpiel, lehrer lehrer) *spielScreen {
	return &spielScreen{
		dame:   dame,
		brett:  lehrer.anfangsBrett.clone(),
		lehrer: lehrer,
		aufgebenKnopf: ui.NewButton(ui.ButtonConfig{
			Position: ui.Position{
				X:                10,
				Y:                10,
				AnchorHorizontal: ui.HorizontalerAnchorLinks,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text: "Aufgeben",
			Callback: func() {
				dame.currentScreen = newGameOverScreen(dame, lehrer, false)
			},
		}),
	}
}

func (s *spielScreen) zugAusfuehren() (zugLaeuft bool) {
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

	naechsterZugSchritt := s.verbeibendeZugSchritte[0]
	s.aktuellerZugSchritt = &naechsterZugSchritt
	s.verbeibendeZugSchritte = s.verbeibendeZugSchritte[1:]

	feldSize := s.brett.feldSize(spielScreenBrettMaxBreite, spielScreenBrettMaxHoehe)

	if naechsterZugSchritt.geschlagenePosition != nil {
		geschlagenesFeld := s.brett.feld(*naechsterZugSchritt.geschlagenePosition)
		geschlagenScreenX, geschlagenScreenY := s.brett.feldPosition(
			*naechsterZugSchritt.geschlagenePosition,
			spielScreenBrettX,
			spielScreenBrettY,
			spielScreenBrettMaxBreite,
			spielScreenBrettMaxHoehe,
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

		s.brett.setFeld(*naechsterZugSchritt.geschlagenePosition, feldLeer)
	}

	vonFeld := s.brett.feld(naechsterZugSchritt.von)
	vonScreenX, vonScreenY := s.brett.feldPosition(
		naechsterZugSchritt.von,
		spielScreenBrettX,
		spielScreenBrettY,
		spielScreenBrettMaxBreite,
		spielScreenBrettMaxHoehe,
	)
	zuScreenX, zuScreenY := s.brett.feldPosition(
		naechsterZugSchritt.zu,
		spielScreenBrettX,
		spielScreenBrettY,
		spielScreenBrettMaxBreite,
		spielScreenBrettMaxHoehe,
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

	s.brett.setFeld(naechsterZugSchritt.von, feldLeer)

	return true
}

func (s *spielScreen) handleClick(mausX, mausY int) {
	mausPosition, ok := s.brett.screenPositionToBrettPosition(
		float64(mausX),
		float64(mausY),
		spielScreenBrettX,
		spielScreenBrettY,
		spielScreenBrettMaxBreite,
		spielScreenBrettMaxHoehe,
	)
	if !ok {
		s.ausgewaehltePosition = nil
		return
	}

	if s.ausgewaehltePosition == nil {
		mausFeld := s.brett.feld(mausPosition)
		feldEigentuemer, feldHatEigentuemer := mausFeld.eigentuemer()
		if !feldHatEigentuemer || feldEigentuemer == spielerLehrer {
			s.ausgewaehltePosition = nil
			return
		}

		s.ausgewaehltePosition = &mausPosition
		return
	}

	moeglicheZuege := s.brett.moeglicheZuegeMitStartPosition(*s.ausgewaehltePosition, s.lehrer.regeln)
	s.ausgewaehltePosition = nil

	for _, moeglicherZug := range moeglicheZuege {
		if moeglicherZug.endPosition() == mausPosition {
			lehrerAiZug, lehrerZugGefunden := ai.BesterNaechsterZug(
				moeglicherZug.ergebnis(),
				s.lehrer.regeln,
				spielerLehrer,
				s.lehrer.aiTiefe,
			)
			if !lehrerZugGefunden {
				s.brett = moeglicherZug.ergebnis()
				return
			}
			lehrerZug := lehrerAiZug.(zug)

			for _, schritt := range moeglicherZug {
				if schritt.geschlagenePosition != nil {
					feldSize := s.brett.feldSize(spielScreenBrettMaxBreite, spielScreenBrettMaxHoehe)

					geschlagenesFeld := s.brett.feld(*schritt.geschlagenePosition)

					geschlagenScreenX, geschlagenScreenY := s.brett.feldPosition(
						*schritt.geschlagenePosition,
						spielScreenBrettX,
						spielScreenBrettY,
						spielScreenBrettMaxBreite,
						spielScreenBrettMaxHoehe,
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
				}
			}

			s.brett = moeglicherZug.ergebnis()
			s.aktuellerZugSchritt = nil
			s.zugStein = nil
			s.verbeibendeZugSchritte = lehrerZug

			return
		}
	}
}

func (s *spielScreen) update() (beendet bool) {
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		s.dame.currentScreen = newGameOverScreen(s.dame, s.lehrer, false)
		return false
	}

	if s.brett.gewonnen(spielerSchueler, s.lehrer.regeln) {
		s.dame.currentScreen = newGameOverScreen(s.dame, s.lehrer, true)
		return false
	}

	if s.brett.gewonnen(spielerLehrer, s.lehrer.regeln) {
		s.dame.currentScreen = newGameOverScreen(s.dame, s.lehrer, false)
		return false
	}

	s.aufgebenKnopf.Update()

	for _, verlorenerStein := range s.geschlageneSteine {
		verlorenerStein.update()
	}

	if s.zugStein != nil {
		s.zugStein.update()
	}

	if zugLaeuft := s.zugAusfuehren(); zugLaeuft {
		return false
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		s.handleClick(ebiten.CursorPosition())
		return false
	}

	for _, touchId := range inpututil.AppendJustReleasedTouchIDs(nil) {
		s.handleClick(inpututil.TouchPositionInPreviousTick(touchId))
		return false
	}

	return false
}

func (s *spielScreen) draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)
	s.brett.draw(
		screen,
		spielScreenBrettX,
		spielScreenBrettY,
		spielScreenBrettMaxBreite,
		spielScreenBrettMaxHoehe,
		s.ausgewaehltePosition,
		s.lehrer.regeln,
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
