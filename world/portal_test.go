package world_test

import (
	"testing"

	"github.com/Lama06/Herder-Legacy/world"
)

func TestPortalTeleportsEntity(t *testing.T) {
	w := world.NewEmptyWorld()
	entity := w.SpawnEntity(&world.Entity{
		Level: 0,
		Position: world.Position{
			X: 0,
			Y: 0,
		},
		HatVelocityComponent: true,
		VelocityComponent: world.VelocityComponent{
			VelocityX: 1.1,
			VelocityY: 0,
		},
		HatHitboxComponent: true,
		HitboxComponent: world.HitboxComponent{
			Width:  10,
			Height: 10,
		},
	})
	w.SpawnEntity(&world.Entity{
		Level: 0,
		Position: world.Position{
			X: 15,
			Y: 5,
		},
		HatHitboxComponent: true,
		HitboxComponent: world.HitboxComponent{
			Width:  1,
			Height: 1,
		},
		HatPortalComponent: true,
		PortalComponent: world.PortalComponent{
			DestinationLevel: 1,
			DestinationPosition: world.Position{
				X: 10,
				Y: 10,
			},
		},
	})

	for i := 0; i < 5; i++ {
		w.Update()
	}

	if entity.Level != 1 {
		t.Error("entity not teleported")
	}

	if !positionsAreRoughlyEqual(entity.Position, world.Position{X: 10, Y: 10}) {
		t.Error("wrong position")
	}
}
