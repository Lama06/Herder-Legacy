package world

import (
	"math"

	"github.com/Lama06/Herder-Legacy/aabb"
)

type HitboxComponent struct {
	Width, Height    float64
	OffsetX, OffsetY float64
}

type RendererHitboxComponent struct{}

type RigidbodyComponent struct{}

func (w *World) rendererHitboxenAnwenden() {
	for entity := range w.Entities {
		if !entity.HatRendererHitboxComponent {
			continue
		}

		switch {
		case entity.HatRectRenderComponent:
			entity.HatHitboxComponent = true
			entity.HitboxComponent = HitboxComponent{
				Width:  entity.RectRenderComponent.Width,
				Height: entity.RectRenderComponent.Height,
			}
		case entity.HatKreisRenderComponent:
			entity.HatHitboxComponent = true
			entity.HitboxComponent = HitboxComponent{
				Width:  entity.KreisRenderComponent.Size,
				Height: entity.KreisRenderComponent.Size,
			}
		case entity.HatImageRenderComponent:
			entity.HatHitboxComponent = true

			scale := entity.ImageRenderComponent.Scale
			if scale == 0 {
				scale = 1
			}

			hitbox := aabb.Aabb{
				X:      0,
				Y:      0,
				Width:  float64(entity.ImageRenderComponent.Image.Bounds().Dx()) * scale,
				Height: float64(entity.ImageRenderComponent.Image.Bounds().Dy()) * scale,
			}.Rotieren(entity.ImageRenderComponent.Rotation)

			entity.HitboxComponent = HitboxComponent{
				OffsetX: hitbox.X,
				OffsetY: hitbox.Y,
				Width:   hitbox.Width,
				Height:  hitbox.Height,
			}
		}
	}
}

func handleKollision(entity1, entity2 *Entity) {
	aabb1 := entity1.aabb()
	aabb2 := entity2.aabb()

	horizontalAusweichenDistanz := math.Min(math.Abs(aabb1.X-aabb2.MaxX()), math.Abs(aabb1.MaxX()-aabb2.X))
	vertikalAusweichenDistanz := math.Min(math.Abs(aabb1.Y-aabb2.MaxY()), math.Abs(aabb1.MaxY()-aabb2.Y))

	if horizontalAusweichenDistanz < vertikalAusweichenDistanz {
		if aabb1.X < aabb2.X {
			switch {
			case entity1.HatRigidbodyComponent && entity2.HatRigidbodyComponent:
				entity1.Position.X -= horizontalAusweichenDistanz / 2
				entity2.Position.X += horizontalAusweichenDistanz / 2
			case entity1.HatRigidbodyComponent:
				entity1.Position.X -= horizontalAusweichenDistanz
			case entity2.HatRigidbodyComponent:
				entity2.Position.X += horizontalAusweichenDistanz
			}
			if entity1.HatRigidbodyComponent && entity1.HatVelocityComponent {
				entity1.VelocityComponent.VelocityX = -math.Abs(entity1.VelocityComponent.VelocityX)
			}
			if entity2.HatRigidbodyComponent && entity2.HatVelocityComponent {
				entity2.VelocityComponent.VelocityX = math.Abs(entity2.VelocityComponent.VelocityX)
			}
		} else {
			switch {
			case entity1.HatRigidbodyComponent && entity2.HatRigidbodyComponent:
				entity2.Position.X -= horizontalAusweichenDistanz / 2
				entity1.Position.X += horizontalAusweichenDistanz / 2
			case entity1.HatRigidbodyComponent:
				entity1.Position.X += horizontalAusweichenDistanz
			case entity2.HatRigidbodyComponent:
				entity2.Position.X -= horizontalAusweichenDistanz
			}
			if entity1.HatRigidbodyComponent && entity1.HatVelocityComponent {
				entity1.VelocityComponent.VelocityX = math.Abs(entity1.VelocityComponent.VelocityX)
			}
			if entity2.HatRigidbodyComponent && entity2.HatVelocityComponent {
				entity2.VelocityComponent.VelocityX = -math.Abs(entity2.VelocityComponent.VelocityX)
			}
		}
	} else {
		if aabb1.Y < aabb2.Y {
			switch {
			case entity1.HatRigidbodyComponent && entity2.HatRigidbodyComponent:
				entity1.Position.Y -= vertikalAusweichenDistanz / 2
				entity2.Position.Y += vertikalAusweichenDistanz / 2
			case entity1.HatRigidbodyComponent:
				entity1.Position.Y -= vertikalAusweichenDistanz
			case entity2.HatRigidbodyComponent:
				entity2.Position.Y += vertikalAusweichenDistanz
			}
			if entity1.HatRigidbodyComponent && entity1.HatVelocityComponent {
				entity1.VelocityComponent.VelocityY = -math.Abs(entity1.VelocityComponent.VelocityY)
			}
			if entity2.HatRigidbodyComponent && entity2.HatVelocityComponent {
				entity2.VelocityComponent.VelocityY = math.Abs(entity2.VelocityComponent.VelocityY)
			}
		} else {
			switch {
			case entity1.HatRigidbodyComponent && entity2.HatRigidbodyComponent:
				entity2.Position.Y -= vertikalAusweichenDistanz / 2
				entity1.Position.Y += vertikalAusweichenDistanz / 2
			case entity1.HatRigidbodyComponent:
				entity1.Position.Y += vertikalAusweichenDistanz
			case entity2.HatRigidbodyComponent:
				entity2.Position.Y -= vertikalAusweichenDistanz
			}
			if entity1.HatRigidbodyComponent && entity1.HatVelocityComponent {
				entity1.VelocityComponent.VelocityY = math.Abs(entity1.VelocityComponent.VelocityY)
			}
			if entity2.HatRigidbodyComponent && entity2.HatVelocityComponent {
				entity2.VelocityComponent.VelocityY = -math.Abs(entity2.VelocityComponent.VelocityY)
			}
		}
	}
}

func (w *World) kollisionenVerarbeiten() {
	for entity1 := range w.Entities {
		if !entity1.HatRigidbodyComponent || !entity1.HatHitboxComponent {
			continue
		}

		aabb1 := entity1.aabb()

		for entity2 := range w.Entities {
			if entity1 == entity2 {
				continue
			}

			if entity1.Level != entity2.Level {
				continue
			}

			if !entity2.HatHitboxComponent {
				continue
			}

			aabb2 := entity2.aabb()

			if !aabb1.KollidiertMit(aabb2) {
				continue
			}

			handleKollision(entity1, entity2)
		}
	}
}
