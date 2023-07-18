package breakout

import "github.com/Lama06/Herder-Legacy/ui"

func (w *world) imAusEntfernen() {
	for entity := range w.entities {
		if !entity.imAusEntfernen {
			continue
		}

		var width, height float64
		if entity.hatHitboxComponent {
			width, height = entity.hitboxComponent.width, entity.hitboxComponent.height
		}

		if entity.position.x+width < 0 ||
			entity.position.y+height < 0 ||
			entity.position.x > ui.Width ||
			entity.position.y > ui.Height {
			delete(w.entities, entity)
		}
	}
}
