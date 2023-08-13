package poker

import (
	"math"
	"math/rand"

	"github.com/Lama06/Herder-Legacy/assets"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

const spielScreenJettonSize = 50

type spielScreenBewegenderJetton struct {
	currentX, currentY float64
	targetX, targetY   float64
}

func newSpielScreenBewegenderJetton(spieler spielScreenSpieler) *spielScreenBewegenderJetton {
	startX, startY := spieler.jettonSpawnPunkt()
	targetX := spielScreenJettonAblageX + rand.Float64()*spielScreenJettonAblageBreite
	targetY := spielScreenJettonAblageY + rand.Float64()*spielScreenJettonAblageHöhe
	return &spielScreenBewegenderJetton{
		currentX: startX,
		currentY: startY,
		targetX:  targetX,
		targetY:  targetY,
	}
}

var _ ui.Component = (*spielScreenBewegenderJetton)(nil)

func (s *spielScreenBewegenderJetton) Update() {
	const maxSpeed = 10

	s.currentX += math.Min(maxSpeed, math.Max(-maxSpeed, s.targetX-s.currentX))
	s.currentY += math.Min(maxSpeed, math.Max(-maxSpeed, s.targetY-s.currentY))
}

func (s *spielScreenBewegenderJetton) Draw(screen *ebiten.Image) {
	img := assets.RequireImage("spielkarten/jeton.png")
	var options ebiten.DrawImageOptions
	options.GeoM.Scale(
		spielScreenJettonSize/float64(img.Bounds().Dx()),
		spielScreenJettonSize/float64(img.Bounds().Dy()),
	)
	options.GeoM.Translate(s.currentX, s.currentY)
	screen.DrawImage(img, &options)
}

func (s *spielScreenBewegenderJetton) angekommen() bool {
	const tolerance = 0.00001
	return math.Abs(s.currentX-s.targetX) <= tolerance &&
		math.Abs(s.currentY-s.targetY) <= tolerance
}
