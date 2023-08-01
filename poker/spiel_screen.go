package poker

import (
	"image/color"
	"math"

	"github.com/Lama06/Herder-Legacy/assets"
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	status              spielScreenStatus
	info                *ui.Title
	stapel              kartenStapel
	eigeneKarten        [spielScreenAnzahlEigenerKarten]*spielScreenKarte
	linkerGegnerKarten  [spielScreenAnzahlGegnerKarten]*spielScreenKarte
	rechterGegnerKarten [spielScreenAnzahlGegnerKarten]*spielScreenKarte
	mittelkarten        [spielScreenAnzahlMittelkarten]*spielScreenKarte
}

var _ herderlegacy.Screen = (*spielScreen)(nil)

func NewSpielScreen() *spielScreen {
	return &spielScreen{
		status: spielStatusKartenAnfang,
		info: ui.NewTitle(ui.TitleConfig{
			Position:           ui.NewCenteredPosition(ui.Width/2, ui.Height/3),
			Text:               "Klicken um Karten zu ziehen",
			CustomColorPalette: false,
			ColorPalette:       ui.TitleColorPalette{},
		}),
		stapel: vollständigerKartenStapel.clone(),
	}
}

func (s *spielScreen) components() []ui.Component {
	components := []ui.Component{s.info}
	for _, eigeneKarte := range s.eigeneKarten {
		components = append(components, eigeneKarte)
	}
	for _, lingerGegnerKarte := range s.linkerGegnerKarten {
		components = append(components, lingerGegnerKarte)
	}
	for _, rechterGegnerKarte := range s.rechterGegnerKarten {
		components = append(components, rechterGegnerKarte)
	}
	for _, mittelkarte := range s.mittelkarten {
		components = append(components, mittelkarte)
	}
	return components
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
		for i := range s.eigeneKarten {
			s.eigeneKarten[i] = &spielScreenKarte{
				currentX:      spielScreenStapelX,
				currentY:      spielScreenStapelY,
				targetX:       spielScreenEigeneKarteX(i),
				targetY:       spielScreenEigeneKartenY,
				autoAufdecken: true,
				karte:         s.stapel.karteZiehen(),
			}
		}
		for i := range s.linkerGegnerKarten {
			s.linkerGegnerKarten[i] = &spielScreenKarte{
				targetRotation: spielScreenGegnerKartenRotation,
				currentX:       spielScreenStapelX,
				currentY:       spielScreenStapelY,
				targetX:        spielScreenLinkerGegnerKartenX,
				targetY:        spielScreenGegnerKarteY(i),
				karte:          s.stapel.karteZiehen(),
			}
		}
		for i := range s.rechterGegnerKarten {
			s.rechterGegnerKarten[i] = &spielScreenKarte{
				targetRotation: spielScreenGegnerKartenRotation,
				currentX:       spielScreenStapelX,
				currentY:       spielScreenStapelY,
				targetX:        spielScreenRechterGegnerKartenX,
				targetY:        spielScreenGegnerKarteY(i),
				karte:          s.stapel.karteZiehen(),
			}
		}
		for i := 0; i < 3; i++ {
			s.mittelkarten[i] = &spielScreenKarte{
				currentX: spielScreenStapelX,
				currentY: spielScreenStapelY,
				targetX:  spielScreenMittelKarteX(i),
				targetY:  spielScreenMittelKartenY,
				karte:    s.stapel.karteZiehen(),
			}
		}
	case spielStatusKartenWerdenGezogen:
		for _, eigeneKarte := range s.eigeneKarten {
			if !eigeneKarte.angekommen() {
				return
			}
		}

		for _, mittelKarte := range s.mittelkarten {
			if mittelKarte != nil && !mittelKarte.angekommen() {
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
			s.mittelkarten[i].targetAufgedecktStatus = true
		}
	case spielStatusMittelkartenWerdenAufgedeckt:
		for i := 0; i < 3; i++ {
			if !s.mittelkarten[i].fertigGewendet() {
				return
			}
		}

		s.status = spielStatus3AufgedeckteMittelKarten
	case spielStatus3AufgedeckteMittelKarten:
		if !istKlick() {
			return
		}

		s.status = spielStatusVierteKarteWirdGezogen
		s.mittelkarten[3] = &spielScreenKarte{
			currentX:      spielScreenStapelX,
			currentY:      spielScreenStapelY,
			targetX:       spielScreenMittelKarteX(3),
			targetY:       spielScreenMittelKartenY,
			karte:         s.stapel.karteZiehen(),
			autoAufdecken: true,
		}
	case spielStatusVierteKarteWirdGezogen:
		if !s.mittelkarten[3].fertigGewendet() {
			return
		}

		s.status = spielStatus4AufgedeckteMittelkarten
	case spielStatus4AufgedeckteMittelkarten:
		if !istKlick() {
			return
		}

		s.status = spielStatusFünfteKarteWirdGezogen
		s.mittelkarten[4] = &spielScreenKarte{
			currentX:      spielScreenStapelX,
			currentY:      spielScreenStapelY,
			targetX:       spielScreenMittelKarteX(4),
			targetY:       spielScreenMittelKartenY,
			karte:         s.stapel.karteZiehen(),
			autoAufdecken: true,
		}
	case spielStatusFünfteKarteWirdGezogen:
		if !s.mittelkarten[4].fertigGewendet() {
			return
		}

		s.status = spielStatus5AufgedeckteMittelkarten
	case spielStatus5AufgedeckteMittelkarten:

	case spielStatusSiegerermittlung:
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

type spielScreenKarte struct {
	karte karte

	currentRotation float64
	targetRotation  float64

	currentX, currentY float64
	targetX, targetY   float64

	currentAufgedecktProgress float64
	targetAufgedecktStatus    bool
	autoAufdecken             bool
}

var _ ui.Component = (*spielScreenKarte)(nil)

func (s *spielScreenKarte) Update() {
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

func (s *spielScreenKarte) angekommen() bool {
	const tolerance = 0.00001
	return math.Abs(s.currentX-s.targetX) <= tolerance &&
		math.Abs(s.currentY-s.targetY) <= tolerance &&
		math.Abs(s.currentRotation-s.targetRotation) <= tolerance
}

func (s *spielScreenKarte) fertigGewendet() bool {
	const tolerance = 0.00001
	if s.targetAufgedecktStatus {
		return math.Abs(s.currentAufgedecktProgress-1) <= tolerance
	} else {
		return math.Abs(s.currentAufgedecktProgress-0) <= tolerance
	}
}

func (s *spielScreenKarte) Draw(screen *ebiten.Image) {
	if s == nil {
		return
	}

	bild := assets.RequireImage("spielkarten/verdeckt.png")
	if s.currentAufgedecktProgress > 0.5 {
		bild = s.karte.image()
	}

	var drawOptions ebiten.DrawImageOptions

	scaleX := spielScreenKarteBreite / float64(bild.Bounds().Dx())
	scaleY := spielScreenKarteHöhe / float64(bild.Bounds().Dy())

	umdrehenScaleX := math.Abs(s.currentAufgedecktProgress-0.5) / 0.5

	drawOptions.GeoM.Rotate(s.currentRotation)
	drawOptions.GeoM.Scale(umdrehenScaleX, 1)
	drawOptions.GeoM.Scale(scaleX, scaleY)
	drawOptions.GeoM.Translate(((1-umdrehenScaleX)/2)*spielScreenKarteBreite, 0)
	drawOptions.GeoM.Translate(s.currentX, s.currentY)
	screen.DrawImage(bild, &drawOptions)
}
