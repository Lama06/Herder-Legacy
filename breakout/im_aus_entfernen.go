package breakout

import (
	"github.com/Lama06/Herder-Legacy/aabb"
	"github.com/Lama06/Herder-Legacy/ui"
)

func (w *world) imAusEntfernen() {
	for entity := range w.entities {
		if !entity.imAusEntfernen {
			continue
		}

		if !entity.hitbox().KollidiertMit(aabb.Aabb{
			X: 0, Y: 0, Width: ui.Width, Height: ui.Height,
		}) {
			delete(w.entities, entity)
		}
	}
}
