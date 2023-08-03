package poker

import (
	"image/color"
	"math"

	"github.com/Lama06/Herder-Legacy/assets"
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func istKlick() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustReleasedTouchIDs(nil)) != 0
}

type spielScreenStatus int

const (
	spielStatusKartenAnfang spielScreenStatus = iota
	spielStatusKartenWerdenGezogen
	spielStatusVerdeckteMittelkarten
	spielStatusMittelkartenWerdenAufgedeckt
	spielStatus3AufgedeckteMittelKarten
	spielStatusVierteKarteWirdGezogen
	spielStatus4AufgedeckteMittelkarten
	spielStatusFünfteKarteWirdGezogen
	spielStatus5AufgedeckteMittelkarten
	spielStatusKartenWerdenAufgedeckt
	spielStatusSiegerermittlung
)

const (
	spielScreenKarteFormat            = 489.0 / 338.0
	spielScreenKarteBreite            = 200.0
	spielScreenKarteHöhe              = spielScreenKarteBreite * spielScreenKarteFormat
	spielScreenKarteMaxSpeed          = 4
	spielScreenKarteMaxRotationSpeed  = math.Pi / (3 * 60)
	spielScreenKarteMaxAufdeckenSpeed = 1.0 / 60.0

	spielScreenStapelY = 100.0
	spielScreenStapelX = (ui.Width - spielScreenKarteBreite) / 2.0

	spielScreenAnzahlMittelkarten          = 5
	spielScreenMittelKartenY               = spielScreenStapelY + spielScreenKarteHöhe*1.5
	spielScreenMittelKartenAbstandX        = spielScreenKarteBreite * 1.5
	spielScreenMittelKartenBreiteGesamt    = (spielScreenAnzahlMittelkarten-1)*spielScreenMittelKartenAbstandX + spielScreenKarteBreite
	spielScreenMittelKartenAbstandVomRandX = (ui.Width - spielScreenMittelKartenBreiteGesamt) / 2

	spielScreenAnzahlEigenerKarten         = 2
	spielScreenEigeneKartenY               = ui.Height - spielScreenKarteHöhe/2
	spielScreenEigeneKartenAbstandX        = spielScreenKarteBreite / 2.0
	spielScreenEigeneKartenBreiteGesamt    = (spielScreenAnzahlEigenerKarten-1)*spielScreenEigeneKartenAbstandX + spielScreenKarteBreite
	spielScreenEigeneKartenAbstandVomRandX = (ui.Width - spielScreenEigeneKartenBreiteGesamt) / 2.0

	spielScreenAnzahlGegnerKarten          = 2
	spielScreenGegnerKartenRotation        = math.Pi / 2
	spielScreenGegnerKartenAbstandY        = spielScreenKarteBreite / 2
	spielScreenGegnerKartenHöheGesamt      = (spielScreenAnzahlGegnerKarten-1)*spielScreenGegnerKartenAbstandY + spielScreenKarteBreite
	spielScreenGegnerKartenAbstandVomRandY = (ui.Height - spielScreenGegnerKartenHöheGesamt) / 2

	spielScreenLinkerGegnerKartenX  = spielScreenKarteHöhe / 2
	spielScreenRechterGegnerKartenX = ui.Width + spielScreenKarteHöhe/2
)

func spielScreenMittelKarteX(karteIndex int) float64 {
	return spielScreenMittelKartenAbstandVomRandX + float64(karteIndex)*spielScreenMittelKartenAbstandX
}

func spielScreenEigeneKarteX(karteIndex int) float64 {
	return spielScreenEigeneKartenAbstandVomRandX + float64(karteIndex)*spielScreenEigeneKartenAbstandX
}

func spielScreenGegnerKarteY(karteIndex int) float64 {
	return spielScreenGegnerKartenAbstandVomRandY + float64(karteIndex)*spielScreenGegnerKartenAbstandY
}

type spielScreen struct {
	status spielScreenStatus
	info   *ui.Title

	eigeneKarten          [2]karte
	eigeneBewegendeKarten [2]*spielScreenBewegendeKarte

	linkerGegnerKarten          [2]karte
	linkerGegnerBewegendeKarten [2]*spielScreenBewegendeKarte

	rechterGegnerKarten          [2]karte
	rechterGegnerBewegendeKarten [2]*spielScreenBewegendeKarte

	mittelkarten          [5]karte
	bewegendeMittelkarten [5]*spielScreenBewegendeKarte
}

var _ herderlegacy.Screen = (*spielScreen)(nil)

func NewSpielScreen() herderlegacy.Screen {
	stapel := vollständigerKartenStapel.clone()

	return &spielScreen{
		status: spielStatusKartenAnfang,
		info: ui.NewTitle(ui.TitleConfig{
			Position:           ui.NewCenteredPosition(ui.Width/2, ui.Height/3),
			Text:               "Klicken um Karten zu ziehen",
			CustomColorPalette: false,
			ColorPalette:       ui.TitleColorPalette{},
		}),

		eigeneKarten:        [2]karte{stapel.karteZiehen(), stapel.karteZiehen()},
		linkerGegnerKarten:  [2]karte{stapel.karteZiehen(), stapel.karteZiehen()},
		rechterGegnerKarten: [2]karte{stapel.karteZiehen(), stapel.karteZiehen()},
		mittelkarten:        [5]karte{stapel.karteZiehen(), stapel.karteZiehen(), stapel.karteZiehen(), stapel.karteZiehen(), stapel.karteZiehen()},
	}
}

func (s *spielScreen) bewegendeKarten() []*spielScreenBewegendeKarte {
	return append(
		s.bewegendeMittelkarten[:],
		s.eigeneBewegendeKarten[0], s.eigeneBewegendeKarten[1],
		s.linkerGegnerBewegendeKarten[0], s.linkerGegnerBewegendeKarten[1],
		s.rechterGegnerBewegendeKarten[0], s.rechterGegnerBewegendeKarten[1],
	)
}

func (s *spielScreen) components() []ui.Component {
	components := []ui.Component{s.info}
	for _, bewegendeKarte := range s.bewegendeKarten() {
		components = append(components, bewegendeKarte)
	}
	return components
}

func (s *spielScreen) menschHand() hand {
	return parseHand([7]karte(append(s.mittelkarten[:], s.eigeneKarten[:]...)))
}

func (s *spielScreen) linksHand() hand {
	return parseHand([7]karte(append(s.mittelkarten[:], s.linkerGegnerKarten[:]...)))
}

func (s *spielScreen) rechtsHand() hand {
	return parseHand([7]karte(append(s.mittelkarten[:], s.rechterGegnerKarten[:]...)))
}

func (s *spielScreen) gewinner() (mensch, links, rechts bool) {
	menschHand := s.menschHand()
	linksHand := s.linksHand()
	rechtsHand := s.rechtsHand()

	switch {
	case compareHände(linksHand, menschHand) == -1 && compareHände(rechtsHand, menschHand) == -1:
		return true, false, false
	default:
		fallthrough
	case compareHände(rechtsHand, linksHand) == -1 && compareHände(menschHand, linksHand) == -1:
		return false, true, false
	case compareHände(linksHand, rechtsHand) == -1 && compareHände(menschHand, rechtsHand) == -1:
		return false, false, true
	}
}

func (s *spielScreen) Update() {
	for _, component := range s.components() {
		component.Update()
	}

	switch s.status {
	case spielStatusKartenAnfang:
		if !istKlick() {
			return
		}

		s.status = spielStatusKartenWerdenGezogen
		s.info.SetText("")
		for i, karte := range s.eigeneKarten {
			s.eigeneBewegendeKarten[i] = &spielScreenBewegendeKarte{
				currentX:      spielScreenStapelX,
				currentY:      spielScreenStapelY,
				targetX:       spielScreenEigeneKarteX(i),
				targetY:       spielScreenEigeneKartenY,
				autoAufdecken: true,
				karte:         karte,
			}
		}
		for i, karte := range s.linkerGegnerKarten {
			s.linkerGegnerBewegendeKarten[i] = &spielScreenBewegendeKarte{
				targetRotation: spielScreenGegnerKartenRotation,
				currentX:       spielScreenStapelX,
				currentY:       spielScreenStapelY,
				targetX:        spielScreenLinkerGegnerKartenX,
				targetY:        spielScreenGegnerKarteY(i),
				karte:          karte,
			}
		}
		for i, karte := range s.rechterGegnerKarten {
			s.rechterGegnerBewegendeKarten[i] = &spielScreenBewegendeKarte{
				targetRotation: spielScreenGegnerKartenRotation,
				currentX:       spielScreenStapelX,
				currentY:       spielScreenStapelY,
				targetX:        spielScreenRechterGegnerKartenX,
				targetY:        spielScreenGegnerKarteY(i),
				karte:          karte,
			}
		}
		for i := 0; i < 3; i++ {
			s.bewegendeMittelkarten[i] = &spielScreenBewegendeKarte{
				currentX: spielScreenStapelX,
				currentY: spielScreenStapelY,
				targetX:  spielScreenMittelKarteX(i),
				targetY:  spielScreenMittelKartenY,
				karte:    s.mittelkarten[i],
			}
		}
	case spielStatusKartenWerdenGezogen:
		for _, eigeneKarte := range s.eigeneBewegendeKarten {
			if !eigeneKarte.animationBeendet() {
				return
			}
		}

		for _, mittelKarte := range s.bewegendeMittelkarten {
			if mittelKarte != nil && !mittelKarte.animationBeendet() {
				return
			}
		}

		s.status = spielStatusVerdeckteMittelkarten
	case spielStatusVerdeckteMittelkarten:
		if !istKlick() {
			return
		}

		s.status = spielStatusMittelkartenWerdenAufgedeckt
		for i := 0; i < 3; i++ {
			s.bewegendeMittelkarten[i].targetAufgedecktStatus = true
		}
	case spielStatusMittelkartenWerdenAufgedeckt:
		for i := 0; i < 3; i++ {
			if !s.bewegendeMittelkarten[i].animationBeendet() {
				return
			}
		}

		s.status = spielStatus3AufgedeckteMittelKarten
	case spielStatus3AufgedeckteMittelKarten:
		if !istKlick() {
			return
		}

		s.status = spielStatusVierteKarteWirdGezogen
		s.bewegendeMittelkarten[3] = &spielScreenBewegendeKarte{
			currentX:      spielScreenStapelX,
			currentY:      spielScreenStapelY,
			targetX:       spielScreenMittelKarteX(3),
			targetY:       spielScreenMittelKartenY,
			karte:         s.mittelkarten[3],
			autoAufdecken: true,
		}
	case spielStatusVierteKarteWirdGezogen:
		if !s.bewegendeMittelkarten[3].animationBeendet() {
			return
		}

		s.status = spielStatus4AufgedeckteMittelkarten
	case spielStatus4AufgedeckteMittelkarten:
		if !istKlick() {
			return
		}

		s.status = spielStatusFünfteKarteWirdGezogen
		s.bewegendeMittelkarten[4] = &spielScreenBewegendeKarte{
			currentX:      spielScreenStapelX,
			currentY:      spielScreenStapelY,
			targetX:       spielScreenMittelKarteX(4),
			targetY:       spielScreenMittelKartenY,
			karte:         s.mittelkarten[4],
			autoAufdecken: true,
		}
	case spielStatusFünfteKarteWirdGezogen:
		if !s.bewegendeMittelkarten[4].animationBeendet() {
			return
		}

		s.status = spielStatus5AufgedeckteMittelkarten
	case spielStatus5AufgedeckteMittelkarten:
		if !istKlick() {
			return
		}

		for _, linkerGegnerKarte := range s.linkerGegnerBewegendeKarten {
			linkerGegnerKarte.targetAufgedecktStatus = true
		}
		for _, rechterGegnerKarte := range s.rechterGegnerBewegendeKarten {
			rechterGegnerKarte.targetAufgedecktStatus = true
		}

		s.status = spielStatusKartenWerdenAufgedeckt
	case spielStatusKartenWerdenAufgedeckt:
		for _, linkerGegnerKarte := range s.linkerGegnerBewegendeKarten {
			if !linkerGegnerKarte.animationBeendet() {
				return
			}
		}
		for _, rechterGegnerKarte := range s.rechterGegnerBewegendeKarten {
			if !rechterGegnerKarte.animationBeendet() {
				return
			}
		}

		s.status = spielStatusSiegerermittlung
	case spielStatusSiegerermittlung:
		var gewinnerHand hand
		menschGewonnen, linksGewonnen, rechtsGewonnen := s.gewinner()
		switch {
		case menschGewonnen:
			gewinnerHand = s.menschHand()
		case linksGewonnen:
			gewinnerHand = s.linksHand()
		case rechtsGewonnen:
			gewinnerHand = s.rechtsHand()
		}

		s.info.SetText(gewinnerHand.displayName())

		for _, karte := range gewinnerHand.karten() {
			if gewinnerHand.visualisierung(karte) == nil {
				continue
			}

			for _, bewegendeKarte := range s.bewegendeKarten() {
				if bewegendeKarte.karte != karte {
					continue
				}

				bewegendeKarte.hatUmrandung = true
				bewegendeKarte.umrandungsFarbe = gewinnerHand.visualisierung(karte)
				break
			}
		}
	}
}

func (s *spielScreen) drawStapel(screen *ebiten.Image) {
	bild := assets.RequireImage("spielkarten/verdeckt.png")
	var drawOptions ebiten.DrawImageOptions
	scaleX := spielScreenKarteBreite / float64(bild.Bounds().Dx())
	scaleY := spielScreenKarteHöhe / float64(bild.Bounds().Dy())
	drawOptions.GeoM.Scale(scaleX, scaleY)
	drawOptions.GeoM.Translate(spielScreenStapelX, spielScreenStapelY)
	screen.DrawImage(bild, &drawOptions)
}

func (s *spielScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 71, G: 113, B: 72, A: 255})
	s.drawStapel(screen)
	for _, component := range s.components() {
		component.Draw(screen)
	}
}

type spielScreenBewegendeKarte struct {
	karte karte

	hatUmrandung    bool
	umrandungsFarbe color.Color

	currentRotation float64
	targetRotation  float64

	currentX, currentY float64
	targetX, targetY   float64

	currentAufgedecktProgress float64
	targetAufgedecktStatus    bool
	autoAufdecken             bool
}

var _ ui.Component = (*spielScreenBewegendeKarte)(nil)

func (s *spielScreenBewegendeKarte) Update() {
	if s == nil {
		return
	}

	s.currentX += math.Min(spielScreenKarteMaxSpeed, math.Max(-spielScreenKarteMaxSpeed, s.targetX-s.currentX))
	s.currentY += math.Min(spielScreenKarteMaxSpeed, math.Max(-spielScreenKarteMaxSpeed, s.targetY-s.currentY))
	s.currentRotation += math.Min(spielScreenKarteMaxRotationSpeed,
		math.Max(-spielScreenKarteMaxRotationSpeed, s.targetRotation-s.currentRotation))

	var targetAufgedecktProgress float64
	if s.targetAufgedecktStatus {
		targetAufgedecktProgress = 1
	}
	s.currentAufgedecktProgress += math.Min(spielScreenKarteMaxAufdeckenSpeed,
		math.Max(-spielScreenKarteMaxAufdeckenSpeed, targetAufgedecktProgress-s.currentAufgedecktProgress))

	if s.autoAufdecken && s.angekommen() {
		s.targetAufgedecktStatus = true
	}
}

func (s *spielScreenBewegendeKarte) angekommen() bool {
	const tolerance = 0.00001
	return math.Abs(s.currentX-s.targetX) <= tolerance &&
		math.Abs(s.currentY-s.targetY) <= tolerance &&
		math.Abs(s.currentRotation-s.targetRotation) <= tolerance
}

func (s *spielScreenBewegendeKarte) animationBeendet() bool {
	if s.autoAufdecken && !s.angekommen() {
		return false
	}

	const tolerance = 0.00001
	if s.targetAufgedecktStatus {
		return math.Abs(s.currentAufgedecktProgress-1) <= tolerance
	} else {
		return math.Abs(s.currentAufgedecktProgress-0) <= tolerance
	}
}

func (s *spielScreenBewegendeKarte) Draw(screen *ebiten.Image) {
	if s == nil {
		return
	}

	bild := assets.RequireImage("spielkarten/verdeckt.png")
	if s.currentAufgedecktProgress > 0.5 {
		bild = s.karte.image()
	}

	var geoM ebiten.GeoM

	scaleX := spielScreenKarteBreite / float64(bild.Bounds().Dx())
	scaleY := spielScreenKarteHöhe / float64(bild.Bounds().Dy())

	umdrehenScaleX := math.Abs(s.currentAufgedecktProgress-0.5) / 0.5

	geoM.Scale(umdrehenScaleX, 1)
	geoM.Scale(scaleX, scaleY)
	geoM.Translate(((1-umdrehenScaleX)/2)*spielScreenKarteBreite, 0)
	geoM.Rotate(s.currentRotation)
	geoM.Translate(s.currentX, s.currentY)

	screen.DrawImage(bild, &ebiten.DrawImageOptions{GeoM: geoM})
	if s.hatUmrandung {
		x, y := geoM.Apply(0, 0)
		untenRechtsX, untenRechtsY := geoM.Apply(float64(bild.Bounds().Dx()), float64(bild.Bounds().Dy()))
		width, height := untenRechtsX-x, untenRechtsY-y
		vector.StrokeRect(
			screen,
			float32(x),
			float32(y),
			float32(width),
			float32(height),
			12,
			s.umrandungsFarbe,
			true,
		)
	}
}
