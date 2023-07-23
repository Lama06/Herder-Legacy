package world

import "github.com/hajimehoshi/ebiten/v2"

type World struct {
	Entities map[*Entity]struct{}

	blockedPathfindingTiles map[Level]map[pathfindingGridTile]struct{}
}

func NewEmptyWorld() *World {
	return &World{
		Entities: make(map[*Entity]struct{}),
	}
}

func (w *World) SpawnEntity(entity *Entity) *Entity {
	w.Entities[entity] = struct{}{}
	return entity
}

func (w *World) Update() {
	w.applyStaticToEntities()
	w.rendererHitboxenAnwenden()
	w.initBlockedPathfindingTiles()
	w.applyVelocityToEntities()
	w.entitiesMitKeyboardSteuern()
	w.entitiesMitTouchSteuern()
	w.pathfind()
	w.moveEntitiesToPositions()
	w.moveEntitiesToPosition()
	w.teleportEntitiesTouchingPortal()
	w.kollisionenVerarbeiten()
}

func (w *World) Draw(screen *ebiten.Image) {
	w.drawEntities(screen)
}
