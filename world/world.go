package world

import "github.com/hajimehoshi/ebiten/v2"

type World struct {
	Entites map[*Entity]struct{}
}

func NewEmptyWorld() *World {
	return &World{
		Entites: make(map[*Entity]struct{}),
	}
}

func (w *World) SpawnEntity(entity *Entity) *Entity {
	w.Entites[entity] = struct{}{}
	return entity
}

func (w *World) Update() {
	w.applyVelocityToEntities()
	w.teleportEntitiesTouchingPortal()
	w.kollisionenVerhindern()
}

func (w *World) Draw(screen *ebiten.Image) {

}
