package world_test

import (
	"testing"

	"github.com/Lama06/Herder-Legacy/world"
)

func TestKollisionenVerhindern(t *testing.T) {
	testCases := []struct {
		entityX, entityY                 float64
		entityWidth, entityHeight        float64
		entityXVelocity, entityYVelocity float64

		hindernisX, hindernisY                 float64
		hindernisWidth, hindernisHeight        float64
		hindernisXVelocity, hindernisYVelocity float64

		ticks int

		expectedEntityX, expectedEntityY float64
	}{
		{
			entityX:         -10,
			entityY:         0,
			entityWidth:     2,
			entityHeight:    1,
			entityXVelocity: 1,
			entityYVelocity: 0,

			hindernisX:         0,
			hindernisY:         -20,
			hindernisWidth:     5,
			hindernisHeight:    40,
			hindernisXVelocity: 0,
			hindernisYVelocity: 0,

			ticks: 20,

			expectedEntityX: -2,
			expectedEntityY: 0,
		},
		{
			entityX:         0,
			entityY:         0,
			entityWidth:     3,
			entityHeight:    3,
			entityXVelocity: -0.1,
			entityYVelocity: 1,

			hindernisX:         -100,
			hindernisY:         100,
			hindernisWidth:     200,
			hindernisHeight:    10,
			hindernisXVelocity: -0.1,
			hindernisYVelocity: -1,

			ticks: 100,

			expectedEntityX: -10,
			expectedEntityY: -3,
		},
	}

	for _, testCase := range testCases {
		w := world.NewEmptyWorld()
		entity := w.SpawnEntity(&world.Entity{
			Position: world.Position{
				X: testCase.entityX,
				Y: testCase.entityY,
			},
			HatHitboxComponent: true,
			HitboxComponent: world.HitboxComponent{
				Width:  testCase.entityWidth,
				Height: testCase.entityHeight,
			},
			HatVelocityComponent: true,
			VelocityComponent: world.VelocityComponent{
				VelocityX: testCase.entityXVelocity,
				VelocityY: testCase.entityYVelocity,
			},
			HatKollisionenVerhindernComponent: true,
		})
		w.SpawnEntity(&world.Entity{
			Position: world.Position{
				X: testCase.hindernisX,
				Y: testCase.hindernisY,
			},
			HatHitboxComponent: true,
			HitboxComponent: world.HitboxComponent{
				Width:  testCase.hindernisWidth,
				Height: testCase.hindernisHeight,
			},
			HatVelocityComponent: true,
			VelocityComponent: world.VelocityComponent{
				VelocityX: testCase.hindernisXVelocity,
				VelocityY: testCase.hindernisYVelocity,
			},
		})

		for i := 0; i < testCase.ticks; i++ {
			w.Update()
		}

		if !floatsRoughlyEqual(entity.Position.X, testCase.expectedEntityX) ||
			!floatsRoughlyEqual(entity.Position.Y, testCase.expectedEntityY) {
			t.Errorf("expected position: %v %v, got position: %v %v",
				testCase.expectedEntityX, testCase.expectedEntityY,
				entity.Position.X, entity.Position.Y)
		}
	}
}
