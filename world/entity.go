package world

import "github.com/Lama06/Herder-Legacy/aabb"

type Entity struct {
	Level    Level
	Position Position

	HatVelocityComponent bool
	VelocityComponent    VelocityComponent

	HatHitboxComponent bool
	HitboxComponent    HitboxComponent

	HatImageHitboxComponent bool
	ImageHitboxComponent    ImageHitboxComponent

	HatKollisionenVerhindernComponent bool
	KollisionenVerhindernComponent    KollisionenVerhindernComponent

	HatPortalComponent bool
	PortalComponent    PortalComponent
}

func (e *Entity) aabb() aabb.Aabb {
	if e.HatHitboxComponent {
		return aabb.Aabb{
			X:      e.Position.X,
			Y:      e.Position.Y,
			Width:  e.HitboxComponent.Width,
			Height: e.HitboxComponent.Height,
		}
	}
	return aabb.Aabb{
		X:      e.Position.X,
		Y:      e.Position.Y,
		Width:  0,
		Height: 0,
	}
}
