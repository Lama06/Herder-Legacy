package world_test

import (
	"testing"

	"github.com/Lama06/Herder-Legacy/world"
)

func TestKollisionenVerhindern(t *testing.T) {
	type entityConfig struct {
		x, y float64

		hatHitbox     bool
		width, height float64

		hatVelocity          bool
		xVelocity, yVelocity float64

		hatRigidbody bool

		expectedX, expectedY float64
	}

	testCases := map[string]struct {
		ticks int

		entities []entityConfig
	}{
		"Entity prallt von statischem Hindernis ab": {
			ticks: 20,

			entities: []entityConfig{
				{
					x: -10,
					y: 0,

					hatHitbox: true,
					width:     2,
					height:    1,

					hatVelocity: true,
					xVelocity:   1,
					yVelocity:   0,

					hatRigidbody: true,

					expectedX: -13,
					expectedY: 0,
				},
				{
					x: 0,
					y: -20,

					hatHitbox: true,
					width:     5,
					height:    40,

					hatVelocity: false,

					hatRigidbody: false,

					expectedX: 0,
					expectedY: -20,
				},
			},
		},
		"Sich aufeinander zubewegende Entities prallen von einander ab": {
			ticks: 100,

			entities: []entityConfig{
				{
					x: 0,
					y: 0,

					hatHitbox: true,
					width:     3,
					height:    3,

					hatVelocity: true,
					xVelocity:   -0.1,
					yVelocity:   1,

					hatRigidbody: true,

					expectedX: -10,
					expectedY: -2.5,
				},
				{
					x: -100,
					y: 100,

					hatHitbox: true,
					width:     200,
					height:    10,

					hatVelocity: true,
					xVelocity:   -0.1,
					yVelocity:   -1,

					hatRigidbody: true,

					expectedX: -110,
					expectedY: 102.5,
				},
			},
		},
	}

	for name, testCase := range testCases {
		w := world.NewEmptyWorld()

		entities := make([]*world.Entity, len(testCase.entities))
		for i, entityConfig := range testCase.entities {
			entities[i] = w.SpawnEntity(&world.Entity{
				Position: world.Position{
					X: entityConfig.x,
					Y: entityConfig.y,
				},
				HatVelocityComponent: entityConfig.hatVelocity,
				VelocityComponent: world.VelocityComponent{
					VelocityX: entityConfig.xVelocity,
					VelocityY: entityConfig.yVelocity,
				},
				HatHitboxComponent: entityConfig.hatHitbox,
				HitboxComponent: world.HitboxComponent{
					Width:  entityConfig.width,
					Height: entityConfig.height,
				},
				HatRigidbodyComponent: entityConfig.hatRigidbody,
			})
		}

		for i := 0; i < testCase.ticks; i++ {
			w.Update()
		}

		for i, entity := range entities {
			if !floatsRoughlyEqual(entity.Position.X, testCase.entities[i].expectedX) ||
				!floatsRoughlyEqual(entity.Position.Y, testCase.entities[i].expectedY) {
				t.Errorf(
					"%v: expected position: %v %v, got position: %v %v",
					name,
					testCase.entities[i].expectedX, testCase.entities[i].expectedY,
					entity.Position.X, entity.Position.Y,
				)
			}
		}
	}
}
