package breakout

import (
	"image/color"
	"math/rand"

	"golang.org/x/image/colornames"
)

type kanonenUpgrade struct {
	anzahlKanonenKugeln int
	kanonenKugelDelay   int
	kanonenKugelSpeed   float64
	kanonenKugelStärke  int
}

func newRandomKanonenUpgrade() upgrade {
	return kanonenUpgrade{
		anzahlKanonenKugeln: 5 + rand.Intn(5),
		kanonenKugelDelay:   30 + rand.Intn(5*60),
		kanonenKugelSpeed:   0.5 + 7*rand.Float64(),
		kanonenKugelStärke:  1 + rand.Intn(2),
	}
}

var _ upgrade = kanonenUpgrade{}

func (f kanonenUpgrade) radius() float64 {
	if f.kanonenKugelStärke == 2 {
		return 25
	}
	return 15
}

func (f kanonenUpgrade) farbe() color.Color {
	return colornames.Royalblue
}

func (f kanonenUpgrade) collect(world *world, plattform *entity) {
	plattform.hatKanonenKugelSpawnerComponent = true
	plattform.kanonenKugelSpawnerComponent = kanonenKugelSpawnerComponent{
		kanonenKugelSpeedX:        float64(plattform.plattformComponent.schießrichtungX) * f.kanonenKugelSpeed,
		kanonenKugelSpeedY:        float64(plattform.plattformComponent.schießrichtungY) * f.kanonenKugelSpeed,
		verbleibendeKanonenKugeln: f.anzahlKanonenKugeln,
		kanoneKugelSchießDelay:    f.kanonenKugelDelay,
		nächsteKanonenKugel:       f.kanonenKugelDelay,
		kanonenKugelStärke:        f.kanonenKugelStärke,
	}
}

type kanonenKugelSpawnerComponent struct {
	kanonenKugelSpeedX, kanonenKugelSpeedY float64
	verbleibendeKanonenKugeln              int
	kanoneKugelSchießDelay                 int
	nächsteKanonenKugel                    int
	kanonenKugelStärke                     int
}

type kanonenKugelComponent struct {
	stärke int
}

func (w *world) kanonenKugelnSpawnen() {
	const kanonenKugelRadius = 20

	for kanonenKugelSpawner := range w.entities {
		if !kanonenKugelSpawner.hatKanonenKugelSpawnerComponent {
			continue
		}

		kanonenKugelSpawnerHitbox := kanonenKugelSpawner.hitbox()

		if kanonenKugelSpawner.kanonenKugelSpawnerComponent.verbleibendeKanonenKugeln == 0 {
			kanonenKugelSpawner.hatKanonenKugelSpawnerComponent = false
			continue
		}

		if kanonenKugelSpawner.kanonenKugelSpawnerComponent.nächsteKanonenKugel > 0 {
			kanonenKugelSpawner.kanonenKugelSpawnerComponent.nächsteKanonenKugel--
			continue
		}

		kanonenKugelSpawner.kanonenKugelSpawnerComponent.verbleibendeKanonenKugeln--
		kanonenKugelSpawner.kanonenKugelSpawnerComponent.nächsteKanonenKugel = kanonenKugelSpawner.kanonenKugelSpawnerComponent.kanoneKugelSchießDelay

		w.entities[&entity{
			position: position{
				x: kanonenKugelSpawnerHitbox.CenterX() - kanonenKugelRadius,
				y: kanonenKugelSpawnerHitbox.CenterY() - kanonenKugelRadius,
			},
			hatVelocityComponent: true,
			velocityComponent: velocityComponent{
				velocityX: kanonenKugelSpawner.kanonenKugelSpawnerComponent.kanonenKugelSpeedX,
				velocityY: kanonenKugelSpawner.kanonenKugelSpawnerComponent.kanonenKugelSpeedY,
			},
			hatRenderComponent: true,
			renderComponent: renderComponent{
				layer: renderLayerKanonenKugel,
			},
			hatCircleComponent: true,
			circleComponent: circleComponent{
				radius: kanonenKugelRadius,
				farbe:  colornames.Red,
			},
			hatKanonenKugelComponent: true,
			kanonenKugelComponent: kanonenKugelComponent{
				stärke: kanonenKugelSpawner.kanonenKugelSpawnerComponent.kanonenKugelStärke,
			},
			hatHitboxComponent: true,
			hitboxComponent: hitboxComponent{
				width:  kanonenKugelRadius * 2,
				height: kanonenKugelRadius * 2,
			},
			imAusEntfernen: true,
		}] = struct{}{}
	}
}

func (w *world) mitKanonenKugelnSteinZerstören() {
	for kanonenKugel := range w.entities {
		if !kanonenKugel.hatKanonenKugelComponent {
			continue
		}

		kanonenKugelHitbox := kanonenKugel.hitbox()

		for stein := range w.entities {
			if !stein.hatSteinComponent {
				continue
			}

			steinHitbox := stein.hitbox()

			if !kanonenKugelHitbox.KollidiertMit(steinHitbox) {
				continue
			}

			delete(w.entities, stein)
			fallendesUpgradeSpawnen(w, stein)

			kanonenKugel.kanonenKugelComponent.stärke--
			if kanonenKugel.kanonenKugelComponent.stärke == 0 {
				delete(w.entities, kanonenKugel)
			}
		}
	}
}
