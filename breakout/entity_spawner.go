package breakout

type spawnerComponent struct {
	creator            func(world *world, spawner *entity) *entity
	verbleibendeSpawns int
	spawnDelay         int
	nächsterSpawn      int
}

func (w *world) entitesSpawnen() {
	for spawner := range w.entities {
		if !spawner.hatSpawnerComponent {
			continue
		}

		if spawner.spawnerComponent.verbleibendeSpawns <= 0 {
			continue
		}

		if spawner.spawnerComponent.nächsterSpawn > 0 {
			spawner.spawnerComponent.nächsterSpawn--
			continue
		}

		spawner.spawnerComponent.verbleibendeSpawns--
		spawner.spawnerComponent.nächsterSpawn = spawner.spawnerComponent.spawnDelay
		w.entities[spawner.spawnerComponent.creator(w, spawner)] = struct{}{}
	}
}
