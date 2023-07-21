package world_test

import (
	"testing"

	"github.com/Lama06/Herder-Legacy/world"
)

func TestVelocity(t *testing.T) {
	w := world.NewEmptyWorld()

	entity := w.SpawnEntity(&world.Entity{
		Position: world.Position{
			X: 10,
			Y: 10,
		},
		HatVelocityComponent: true,
		VelocityComponent: world.VelocityComponent{
			VelocityX: 5,
			VelocityY: -5,
		},
	})

	w.Update()
	w.Update()
	w.Update()

	if !floatsRoughlyEqual(entity.Position.X, 25) || !floatsRoughlyEqual(entity.Position.Y, -5) {
		t.Errorf("entity did not move correctly")
	}
}
