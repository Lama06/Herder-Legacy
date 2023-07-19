package breakout

import (
	"image/color"
	"math/rand"
)

func newRandomUpgrade() upgrade {
	possibleUpgradeCreators := []func() upgrade{
		newRandomKanonenUpgrade,
		newRandomRainbowUpgrade,
		newRandomFasterInputUpgrade,
		newRandomAutomaticInputUpgrade,
		newRandomHiddenUpgrade,
		newRandomZeitUpgrade,
		newRandomUpgradeShuffleUpgrade,
	}
	creator := possibleUpgradeCreators[rand.Intn(len(possibleUpgradeCreators))]
	return creator()
}

type upgrade interface {
	farbe() color.Color

	radius() float64

	collect(world *world, plattform *entity)
}

type fallendesUpgradeComponent struct {
	upgrade upgrade
}

func fallendesUpgradeSpawnen(world *world, stein *entity) {
	if !stein.steinComponent.hatUpgrade {
		return
	}

	world.entities[&entity{
		position:             stein.position,
		hatVelocityComponent: true,
		velocityComponent: velocityComponent{
			velocityX: stein.steinComponent.upgradeSpeedX,
			velocityY: stein.steinComponent.upgradeSpeedY,
		},
		hatRenderComponent: true,
		renderComponent: renderComponent{
			layer: renderLayerFallendesUpgrade,
		},
		hatCircleComponent: true,
		circleComponent: circleComponent{
			radius: stein.steinComponent.upgrade.radius(),
			farbe:  stein.steinComponent.upgrade.farbe(),
		},
		hatHitboxComponent: true,
		hitboxComponent: hitboxComponent{
			width:  stein.steinComponent.upgrade.radius() * 2,
			height: stein.steinComponent.upgrade.radius() * 2,
		},
		hatFallendesUpgradeComponent: true,
		fallendesUpgradeComponent: fallendesUpgradeComponent{
			upgrade: stein.steinComponent.upgrade,
		},
		imAusEntfernen:        true,
		affectsAutomaticInput: true,
	}] = struct{}{}
}

func (w *world) fallendeUpgradesAufsammeln() {
	for plattform := range w.entities {
		if !plattform.hatPlattformComponent {
			continue
		}

		plattformHitbox := plattform.hitbox()

		for fallendesUpgrade := range w.entities {
			if !fallendesUpgrade.hatFallendesUpgradeComponent {
				continue
			}

			fallendesUpgradeHitbox := fallendesUpgrade.hitbox()

			if !plattformHitbox.KollidiertMit(fallendesUpgradeHitbox) {
				continue
			}

			delete(w.entities, fallendesUpgrade)
			fallendesUpgrade.fallendesUpgradeComponent.upgrade.collect(w, plattform)
		}
	}
}
