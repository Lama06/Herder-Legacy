package breakout

import (
	"image/color"
	"math/rand"

	"golang.org/x/image/colornames"
)

type kanonenUpgrade struct {
	kugeln      int
	delay       int
	speed       float64
	kugelStärke int
}

func newRandomKanonenUpgrade() upgrade {
	return kanonenUpgrade{
		kugeln:      5 + rand.Intn(5),
		delay:       30 + rand.Intn(5*60),
		speed:       0.5 + 7*rand.Float64(),
		kugelStärke: 1 + rand.Intn(2),
	}
}

var _ upgrade = kanonenUpgrade{}

func (f kanonenUpgrade) radius() float64 {
	if f.kugelStärke == 2 {
		return 20
	}
	return 15
}

func (f kanonenUpgrade) farbe() color.Color {
	return colornames.Royalblue
}

func (f kanonenUpgrade) collect(world *world, plattform *entity) {
	plattform.hatKanonenKugelSpawnerComponent = true
	plattform.kanonenKugelSpawnerComponent = kanonenKugelSpawnerComponent{
		kanonenKugelSpeedX:        float64(plattform.plattformComponent.schießrichtungX) * f.speed,
		kanonenKugelSpeedY:        float64(plattform.plattformComponent.schießrichtungY) * f.speed,
		verbleibendeKanonenKugeln: f.kugeln,
		kanoneKugelSchießDelay:    f.delay,
		nächsteKanonenKugel:       f.delay,
		kanonenKugelStärke:        f.kugelStärke,
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
				x: kanonenKugelSpawnerHitbox.CenterX() - kanonenKugelRadius/2,
				y: kanonenKugelSpawnerHitbox.CenterY() - kanonenKugelRadius/2,
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
