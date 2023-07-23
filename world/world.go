package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	Entities    map[*Entity]struct{}
	Backgrounds map[Level]color.Color

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

func (w *World) FindEntitiesWithTag(tag string) []*Entity {
	var result []*Entity
	for entity := range w.Entities {
		if entity.HasTag(tag) {
			result = append(result, entity)
		}
	}
	return result
}

func (w *World) FindEntityWithName(name string) *Entity {
	for entity := range w.Entities {
		if entity.Name == name {
			return entity
		}
	}
	return nil
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
