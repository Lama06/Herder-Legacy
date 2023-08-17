package flappyoinky

import (
	"math"
	"math/rand"
	"strconv"

	"github.com/Lama06/Herder-Legacy/aabb"
	"github.com/Lama06/Herder-Legacy/assets"
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
)

type flappyOinkyScreen struct {
	herderLegacy herderlegacy.HerderLegacy
	callback     func(punkte int) herderlegacy.Screen

	oinky *oinky

	hindernisCounter           int
	hindernisVerzögerung       float64
	hindernisSpeedX            float64
	nächstesHindernisCountdown int
	hindernisse                []*hindernis

	punkteText *ui.Title
}

func NewFlappyOinkyScreen(
	herderLegacy herderlegacy.HerderLegacy,
	callback func(punkte int) herderlegacy.Screen,
) herderlegacy.Screen {
	return &flappyOinkyScreen{
		herderLegacy: herderLegacy,
		callback:     callback,

		oinky: newOinky(),

		hindernisCounter:           0,
		hindernisVerzögerung:       hindernisVerzögerungStart,
		hindernisSpeedX:            hindernisSpeedXStart,
		nächstesHindernisCountdown: 0,

		punkteText: ui.NewTitle(ui.TitleConfig{
			Position: ui.Position{
				X:                ui.Width / 2,
				Y:                20,
				AnchorHorizontal: ui.HorizontalerAnchorMitte,
				AnchorVertikal:   ui.VertikalerAnchorOben,
			},
			Text: "0",
		}),
	}
}

func (f *flappyOinkyScreen) Update() {
	f.punkteText.Update()
	f.punkteText.SetText(strconv.Itoa(f.hindernisCounter))

	f.hindernisVerzögerung -= hindernisVerzögerungDecrease
	f.hindernisSpeedX += hindernisSpeedXIncrease

	f.nächstesHindernisCountdown--
	if f.nächstesHindernisCountdown <= 0 {
		f.nächstesHindernisCountdown = int(f.hindernisVerzögerung)
		f.hindernisCounter++
		f.hindernisse = append(f.hindernisse, newHindernis(f.hindernisSpeedX))
	}

	if f.oinky.istWeg() {
		f.herderLegacy.OpenScreen(f.callback(f.hindernisCounter))
		return
	}
	f.oinky.update()

	for i, hindernis := range f.hindernisse {
		if hindernis.istWeg() {
			f.hindernisse = append(f.hindernisse[:i], f.hindernisse[i+1:]...)
			continue
		}

		hindernisObenHitbox, hindernisUntenHitbox := hindernis.hitboxen()
		if f.oinky.hitbox().KollidiertMit(hindernisObenHitbox) || f.oinky.hitbox().KollidiertMit(hindernisUntenHitbox) {
			f.herderLegacy.OpenScreen(f.callback(f.hindernisCounter))
			return
		}

		hindernis.update()
	}
}

func (f *flappyOinkyScreen) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Skyblue)
	f.oinky.draw(screen)
	for _, hindernis := range f.hindernisse {
		hindernis.draw(screen)
	}
	f.punkteText.Draw(screen)
}

const (
	hindernisBreite              = 150
	hindernisFreierPlatz         = 350
	hindernisSpeedXStart         = 2
	hindernisSpeedXIncrease      = 0.001
	hindernisVerzögerungStart    = 5 * 60
	hindernisVerzögerungDecrease = 0.004
)

func createHindernisImage(breite, hoehe int) *ebiten.Image {
	hindernisTileImg := assets.RequireImage("flappyoinky/hindernis.png")
	hindernisTileScale := float64(breite) / float64(hindernisTileImg.Bounds().Dx())
	hindernisTileScaledHoehe := float64(hindernisTileImg.Bounds().Dy()) * hindernisTileScale

	hindernisImg := ebiten.NewImage(int(breite), int(hoehe))

	for y := 0.0; y < float64(hoehe); y += hindernisTileScaledHoehe {
		var tileDrawOptions ebiten.DrawImageOptions
		tileDrawOptions.GeoM.Scale(hindernisTileScale, hindernisTileScale)
		tileDrawOptions.GeoM.Translate(0, float64(y))
		hindernisImg.DrawImage(hindernisTileImg, &tileDrawOptions)
	}

	return hindernisImg
}

type hindernis struct {
	x                             float64
	xGeschwindigkeit              float64
	obereSchranke, untereSchranke float64
	obenBild, untenBild           *ebiten.Image
}

func newHindernis(xGeschwindigkeit float64) *hindernis {
	obereSchranke := (ui.Height - hindernisFreierPlatz) * rand.Float64()
	untereSchranke := obereSchranke + hindernisFreierPlatz

	return &hindernis{
		x:                ui.Width,
		xGeschwindigkeit: xGeschwindigkeit,
		obereSchranke:    obereSchranke,
		untereSchranke:   untereSchranke,
		obenBild:         createHindernisImage(hindernisBreite, int(obereSchranke)),
		untenBild:        createHindernisImage(hindernisBreite, ui.Height-int(untereSchranke)),
	}
}

func (h *hindernis) istWeg() bool {
	return h.x+hindernisBreite < 0
}

func (h *hindernis) hitboxen() (oben, unten aabb.Aabb) {
	return aabb.Aabb{
			X:      h.x,
			Y:      0,
			Width:  hindernisBreite,
			Height: h.obereSchranke,
		}, aabb.Aabb{
			X:      h.x,
			Y:      h.untereSchranke,
			Width:  hindernisBreite,
			Height: ui.Height - h.untereSchranke,
		}
}

func (h *hindernis) update() {
	h.x -= h.xGeschwindigkeit
}

func (h *hindernis) draw(screen *ebiten.Image) {
	var obenDrawOptions ebiten.DrawImageOptions
	obenDrawOptions.GeoM.Translate(h.x, 0)
	screen.DrawImage(h.obenBild, &obenDrawOptions)

	var untenDrawOptions ebiten.DrawImageOptions
	untenDrawOptions.GeoM.Translate(h.x, h.untereSchranke)
	screen.DrawImage(h.untenBild, &untenDrawOptions)
}

const (
	oinkySize                       = 100
	oinkyX                          = ui.Width/2 - oinkySize/2
	oinkyYStart                     = ui.Height/2 - oinkySize/2
	oinkyYGeschwindigkeitNachSprung = -10
	oinkyYGeschwindigkeitStart      = oinkyYGeschwindigkeitNachSprung * 2
	oinkyBeschleunigungY            = 0.5
	oinkyMaxRotation                = math.Pi / 4
	oinkyMaxRotationÄnderung        = 0.1
)

type oinky struct {
	y                float64
	yGeschwindigkeit float64
	rotation         float64
}

func newOinky() *oinky {
	return &oinky{
		y:                oinkyYStart,
		yGeschwindigkeit: oinkyYGeschwindigkeitStart,
	}
}

func (o *oinky) istWeg() bool {
	return o.y > ui.Height || o.y+oinkySize < 0
}

func (o *oinky) hitbox() aabb.Aabb {
	return aabb.Aabb{
		X:      oinkyX,
		Y:      o.y,
		Width:  oinkySize,
		Height: oinkySize,
	}
}

func (o *oinky) update() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustPressedTouchIDs(nil)) != 0 {
		o.yGeschwindigkeit = oinkyYGeschwindigkeitNachSprung
	}

	o.y += o.yGeschwindigkeit
	o.yGeschwindigkeit += oinkyBeschleunigungY

	rotationÄnderung := o.targetRotation() - o.rotation
	rotationÄnderung = math.Max(-oinkyMaxRotationÄnderung, math.Min(oinkyMaxRotationÄnderung, rotationÄnderung))
	o.rotation += rotationÄnderung
}

func (o *oinky) targetRotation() float64 {
	rotation := (o.yGeschwindigkeit / oinkyYGeschwindigkeitNachSprung) * -oinkyMaxRotation
	return math.Min(oinkyMaxRotation, math.Max(-oinkyMaxRotation, rotation))
}

func (o *oinky) draw(screen *ebiten.Image) {
	oinkyImg := assets.RequireImage("flappyoinky/oinky.png")
	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Scale(oinkySize/float64(oinkyImg.Bounds().Dx()), oinkySize/float64(oinkyImg.Bounds().Dy()))
	drawOptions.GeoM.Rotate(o.rotation)
	drawOptions.GeoM.Translate(oinkyX, o.y)
	screen.DrawImage(oinkyImg, &drawOptions)
}
