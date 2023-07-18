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

func (f kanonenUpgrade) collect(w *world, plattform *entity) {
	const kanonenKugelRadius = 20

	plattform.hatSpawnerComponent = true
	plattform.spawnerComponent = spawnerComponent{
		verbleibendeSpawns: f.anzahlKanonenKugeln,
		spawnDelay:         f.kanonenKugelDelay,
		nächsterSpawn:      f.kanonenKugelDelay,
		creator: func(w *world, spawner *entity) *entity {
			spawnerHitbox := spawner.hitbox()

			return &entity{
				position: position{
					x: spawnerHitbox.CenterX() - kanonenKugelRadius,
					y: spawnerHitbox.CenterY() - kanonenKugelRadius,
				},
				hatVelocityComponent: true,
				velocityComponent: velocityComponent{
					velocityX: float64(plattform.plattformComponent.schießrichtungX) * f.kanonenKugelSpeed,
					velocityY: float64(plattform.plattformComponent.schießrichtungY) * f.kanonenKugelSpeed,
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
					stärke: f.kanonenKugelStärke,
				},
				hatHitboxComponent: true,
				hitboxComponent: hitboxComponent{
					width:  kanonenKugelRadius * 2,
					height: kanonenKugelRadius * 2,
				},
				imAusEntfernen: true,
			}
		},
	}
}

type kanonenKugelComponent struct {
	stärke int
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
