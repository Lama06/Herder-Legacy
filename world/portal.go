package world

import "github.com/Lama06/Herder-Legacy/aabb"

type PortalComponent struct {
	Width, Height       float64
	DestinationLevel    Level
	DestinationPosition Position
}

func (w *World) teleportEntitiesTouchingPortal() {
	for portal := range w.Entites {
		if !portal.HatPortalComponent {
			continue
		}

		portalAabb := aabb.Aabb{
			X:      portal.Position.X,
			Y:      portal.Position.Y,
			Width:  portal.PortalComponent.Width,
			Height: portal.PortalComponent.Height,
		}

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
