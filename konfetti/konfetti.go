package konfetti

import (
	"image/color"
	"image/color/palette"
	"math/rand"

	"github.com/Lama06/Herder-Legacy/aabb"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	partikelAnzahlProSpawn = 100
	partikelSize           = 3
	maxSpeed               = 7
)

type KonfettiManager struct {
	nextKonfettiPartikelId int
	partikel               map[int]konfettiPartikel
}

func NewKonfettiManager() *KonfettiManager {
	return &KonfettiManager{
		nextKonfettiPartikelId: 0,
		partikel:               make(map[int]konfettiPartikel),
	}
}

func (k *KonfettiManager) SpawnKonfetti(x, y float64) {
	for i := 0; i < partikelAnzahlProSpawn; i++ {
		id := k.nextKonfettiPartikelId
		k.nextKonfettiPartikelId++
		k.partikel[id] = konfettiPartikel{
			x:      x,
			y:      y,
			speedX: -maxSpeed + 2*maxSpeed*rand.Float64(),
			speedY: -maxSpeed + 2*maxSpeed*rand.Float64(),
			farbe:  palette.Plan9[rand.Intn(len(palette.Plan9))],
		}
	}
}

func (k *KonfettiManager) Update() {
	for id, partikel := range k.partikel {
		if partikel.istUnsichtbar() {
			delete(k.partikel, id)
			continue
		}

		partikel.update()
		k.partikel[id] = partikel
	}
}

func (k *KonfettiManager) Draw(screen *ebiten.Image) {
	for _, partikel := range k.partikel {
		partikel.draw(screen)
	}
}

type konfettiPartikel struct {
	x, y           float64
	speedX, speedY float64
	farbe          color.Color
}

func (k *konfettiPartikel) istUnsichtbar() bool {
	return !aabb.Aabb{X: 0, Y: 0, Width: ui.Width, Height: ui.Height}.IsInside(k.x, k.y)
}

func (k *konfettiPartikel) draw(screen *ebiten.Image) {
	vector.DrawFilledRect(
		screen,
		float32(k.x),
		float32(k.y),
		partikelSize,
		partikelSize,
		k.farbe,
		false,
	)
}

func (k *konfettiPartikel) update() {
	k.x += k.speedX
	k.y += k.speedY
}
