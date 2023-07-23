package world_test

import (
	"testing"

	"github.com/Lama06/Herder-Legacy/world"
)

func TestMoveToPosition(t *testing.T) {
	w := world.NewEmptyWorld()
	entity := w.SpawnEntity(&world.Entity{
		Position: world.Position{
			X: -10,
			Y: -10,
		},
		HatMoveToPositionComponent: true,
		MoveToPositionComponent: world.MoveToPositionComponent{
			Position: world.Position{
				X: 13,
				Y: 13,
			},
			Speed: 1,
		},
	})

	for i := 0; i < 30; i++ {
		w.Update()
	}

	if !positionsAreRoughlyEqual(entity.Position, world.Position{X: 13, Y: 13}) {
		t.Errorf("expected position: (13, 13), got position: %v", entity.Position)
	}

	if !entity.MoveToPositionComponent.Arrived {
		t.Error("arrive flag false")
	}
}

func TestMoveToPositions(t *testing.T) {
	w := world.NewEmptyWorld()
	entity := w.SpawnEntity(&world.Entity{
		Position: world.Position{
			X: -10,
			Y: -10,
		},
		HatMoveToPositionsComponent: true,
		MoveToPositionsComponent: world.MoveToPositionsComponent{
			Speed: 1,
			Positions: []world.Position{
				{X: 0, Y: -10},
				{X: 0, Y: 0},
				{X: 10, Y: 10},
			},
		},
	})

	for i := 0; i < 50; i++ {
		w.Update()
	}

	if !positionsAreRoughlyEqual(entity.Position, world.Position{X: 10, Y: 10}) {
		t.Errorf("expected position: (0, 10), got position: %v", entity.Position)
	}

	if !entity.MoveToPositionsComponent.Finished {
		t.Error("finished flag false")
	}
}

func TestPathfindingInSameLevel(t *testing.T) {
	w := world.NewEmptyWorld()
	entity := w.SpawnEntity(&world.Entity{
		Position: world.Position{
			X: 0,
			Y: 0,
		},
		HatPathfinderComponent: true,
		PathfinderComponent: world.PathfinderComponent{
			DestinationPosition: world.Position{
				X: 1000,
				Y: 0,
			},
			DestinationLevel: 0,
			Speed:            10,
		},
	})

	w.SpawnEntity(&world.Entity{
		Position: world.Position{
			X: 100,
			Y: -2000,
		},
		Static:             true,
		HatHitboxComponent: true,
		HitboxComponent: world.HitboxComponent{
			Width:  100,
			Height: 4000,
		},
	})

	for i := 0; i < 200; i++ {
		w.Update()
	}

	if entity.PathfinderComponent.State == world.PathfinderComponentStateArrived {
		t.Error("entity arrived too early")
	}

	for i := 0; i < 500; i++ {
		w.Update()
	}

	if entity.PathfinderComponent.State != world.PathfinderComponentStateArrived {
		t.Error("entity did not arrive")
	}

	if !positionsAreRoughlyEqual(entity.Position, world.Position{X: 1000, Y: 0}) {
		t.Errorf("expected position: (1000, 0), got position: %v", entity.Position)
	}
}

func TestPathfindingToDifferentLevel(t *testing.T) {
	w := world.NewEmptyWorld()
	entity := w.SpawnEntity(&world.Entity{
		Level: 0,
		Position: world.Position{
			X: 0,
			Y: 0,
		},
		HatPathfinderComponent: true,
		PathfinderComponent: world.PathfinderComponent{
			DestinationPosition: world.Position{
				X: 0,
				Y: 0,
			},
			DestinationLevel: 1,
			Speed:            10,
		},
	})

	w.SpawnEntity(&world.Entity{
		Level: 0,
		Position: world.Position{
			X: 100,
			Y: 100,
		},
		Static:             true,
		HatPortalComponent: true,
		PortalComponent: world.PortalComponent{
			Width:            10,
			Height:           10,
			DestinationLevel: 1,
			DestinationPosition: world.Position{
				X: 0,
				Y: 0,
			},
		},
	})

	w.SpawnEntity(&world.Entity{
		Level: 1,
		Position: world.Position{
			X: 100,
			Y: 100,
		},
		Static:             true,
		HatPortalComponent: true,
		PortalComponent: world.PortalComponent{
			Width:            10,
			Height:           10,
			DestinationLevel: 0,
			DestinationPosition: world.Position{
				X: 0,
				Y: 0,
			},
		},
	})

	for i := 0; i < 50; i++ {
		w.Update()
	}

	if entity.PathfinderComponent.State != world.PathfinderComponentStateArrived {
		t.Error("entity did not arrive")
	}

	if !positionsAreRoughlyEqual(entity.Position, world.Position{X: 0, Y: 0}) {
		t.Errorf("expected position: (0, 0), got position: %v", entity.Position)
	}

	if entity.Level != 1 {
		t.Error("entity in wrong level")
	}
}
