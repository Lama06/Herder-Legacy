package world

type PortalComponent struct {
	DestinationLevel    Level
	DestinationPosition Position
}

func (w *World) teleportEntitiesTouchingPortal() {
	for portal := range w.Entites {
		if !portal.HatPortalComponent {
			continue
		}

		portalAabb := portal.aabb()

		for entity := range w.Entites {
			if entity == portal {
				continue
			}

			if entity.Level != portal.Level {
				continue
			}

			entityAabb := entity.aabb()

			if !portalAabb.KollidiertMit(entityAabb) {
				continue
			}

			entity.Level = portal.PortalComponent.DestinationLevel
			entity.Position.X = portal.PortalComponent.DestinationPosition.X
			entity.Position.Y = portal.PortalComponent.DestinationPosition.Y
		}
	}
}
