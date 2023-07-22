package openworld

import (
	"github.com/Lama06/Herder-Legacy/herder"
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/world"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type openWorldScreen struct {
	world *world.World
}

func NewOpenWorldScreen() herderlegacy.Screen {
	w := herder.CreateHerder()
	w.SpawnEntity(&world.Entity{
		Level: 0,
		Position: world.Position{
			X: 200,
			Y: 200,
		},
		HatRenderComponent: true,
		RenderComponent: world.RenderComponent{
			Layer: 1,
		},
		HatKreisRenderComponent: true,
		KreisRenderComponent: world.KreisRenderComponent{
			Farbe: colornames.Pink,
			Size:  50,
		},
		HatHitboxComponent: true,
		HitboxComponent: world.HitboxComponent{
			Width:  50,
			Height: 50,
		},
		HatRigidbodyComponent:        true,
		HatKeyboardMovementComponent: true,
		KeyboardMovementComponent: world.KeyboardMovementComponent{
			Unten:  ebiten.KeyS,
			Oben:   ebiten.KeyW,
			Links:  ebiten.KeyA,
			Rechts: ebiten.KeyD,
			Speed:  2,
		},
		HatTouchInputComponent: true,
		TouchInputComponent: world.TouchInputComponent{
			Speed:        2,
			Sensitivität: 0.01,
		},
		HatCameraComponent: true,
		CameraComponent: world.CameraComponent{
			OffsetX: 25,
			OffsetY: 25,
		},
	})

	return &openWorldScreen{
		world: w,
	}
}

func (o *openWorldScreen) Update() {
	o.world.Update()
}

func (o *openWorldScreen) Draw(screen *ebiten.Image) {
	o.world.Draw(screen)
}
