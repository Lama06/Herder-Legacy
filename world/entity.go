package world

import "github.com/Lama06/Herder-Legacy/aabb"

type Entity struct {
	Level    Level
	Position Position
	Static   bool

	HatVelocityComponent bool
	VelocityComponent    VelocityComponent

	HatHitboxComponent bool
	HitboxComponent    HitboxComponent

	HatRendererHitboxComponent bool
	RendererHitboxComponent    RendererHitboxComponent

	HatRigidbodyComponent bool
	RigidbodyComponent    RigidbodyComponent

	HatKeyboardMovementComponent bool
	KeyboardMovementComponent    KeyboardMovementComponent

	HatTouchInputComponent bool
	TouchInputComponent    TouchInputComponent

	HatMoveToPositionComponent bool
	MoveToPositionComponent    MoveToPositionComponent

	HatMoveToPositionsComponent bool
	MoveToPositionsComponent    MoveToPositionsComponent

	HatPathfinderComponent bool
	PathfinderComponent    PathfinderComponent

	HatPortalComponent bool
	PortalComponent    PortalComponent

	HatCameraComponent bool
	CameraComponent    CameraComponent

	HatRenderComponent bool
	RenderComponent    RenderComponent

	HatRectRenderComponent bool
	RectRenderComponent    RectRenderComponent

	HatKreisRenderComponent bool
	KreisRenderComponent    KreisRenderComponent

	HatImageRenderComponent bool
	ImageRenderComponent    ImageRenderComponent
}

func (e *Entity) aabb() aabb.Aabb {
	if !e.HatHitboxComponent {
		panic("missing hitbox")
	}
	return aabb.Aabb{
		X:      e.Position.X,
		Y:      e.Position.Y,
		Width:  e.HitboxComponent.Width,
		Height: e.HitboxComponent.Height,
	}
}

func (w *World) applyStaticToEntities() {
	for entity := range w.Entities {
		if !entity.Static {
			continue
		}

		entity.HatVelocityComponent = false
		entity.HatRigidbodyComponent = false
		entity.HatMoveToPositionComponent = false
		entity.HatMoveToPositionsComponent = false
		entity.HatPathfinderComponent = false
		entity.HatKeyboardMovementComponent = false
		entity.HatTouchInputComponent = false
	}
}
