package breakout

type timerComponent struct {
	verbleibendeZeit int
}

func (w *world) timersRunterzählen() {
	for entity := range w.entities {
		if !entity.hatTimerComponent {
			continue
		}

		entity.timerComponent.verbleibendeZeit--

		if entity.timerComponent.verbleibendeZeit <= 0 {
			delete(w.entities, entity)
		}
	}
}
