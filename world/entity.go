package world

import "github.com/Lama06/Herder-Legacy/aabb"

type Entity struct {
	Level    Level
	Position Position

	HatVelocityComponent bool
	VelocityComponent    VelocityComponent

	HatHitboxComponent bool
	HitboxComponent    HitboxComponent

	HatRendererHitboxComponent bool
	RendererHitboxComponent    RendererHitboxComponent

	HatRigidbodyComponent bool
	RigidbodyComponent    RigidbodyComponent

	HatPortalComponent bool
	PortalComponent    PortalComponent

	HatKeyboardMovementComponent bool
	KeyboardMovementComponent    KeyboardMovementComponent

	HatTouchInputComponent bool
	TouchInputComponent    TouchInputComponent

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
