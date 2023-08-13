package poker

import (
	"image/color"
	"math"

	"github.com/Lama06/Herder-Legacy/assets"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	spielScreenKarteFormat = 489.0 / 338.0
	spielScreenKarteBreite = 200.0
	spielScreenKarteHöhe   = spielScreenKarteBreite * spielScreenKarteFormat
)

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
	const (
		maxSpeed          = 4
		maxRotationSpeed  = math.Pi / (3 * 60)
		maxAufdeckenSpeed = 1.0 / 60.0
	)

	if s == nil {
		return
	}

	s.currentX += math.Min(maxSpeed, math.Max(-maxSpeed, s.targetX-s.currentX))
	s.currentY += math.Min(maxSpeed, math.Max(-maxSpeed, s.targetY-s.currentY))
	s.currentRotation += math.Min(maxRotationSpeed,
		math.Max(-maxRotationSpeed, s.targetRotation-s.currentRotation))

	var targetAufgedecktProgress float64
	if s.targetAufgedecktStatus {
		targetAufgedecktProgress = 1
	}
	s.currentAufgedecktProgress += math.Min(maxAufdeckenSpeed,
		math.Max(-maxAufdeckenSpeed, targetAufgedecktProgress-s.currentAufgedecktProgress))

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

	geoM.Scale(scaleX*umdrehenScaleX, scaleY)
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
