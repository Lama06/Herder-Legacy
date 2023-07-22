package herder

import "github.com/Lama06/Herder-Legacy/world"

func CreateHerder() *world.World {
	w := world.NewEmptyWorld()
	bodenImg := requireImage("boden")
	for zeile := 0; zeile < 30; zeile++ {
		for spalte := 0; spalte < 30; spalte++ {
			w.SpawnEntity(&world.Entity{
				Level: 0,
				Position: world.Position{
					X: float64(spalte) * float64(bodenImg.Bounds().Dx()) * 2,
					Y: float64(zeile) * float64(bodenImg.Bounds().Dy()) * 2,
				},
				HatRenderComponent: true,
				RenderComponent: world.RenderComponent{
					Layer: -1,
				},
				HatImageRenderComponent: true,
				ImageRenderComponent: world.ImageRenderComponent{
					Image: bodenImg,
					Scale: 2,
				},
			})
		}
	}
	w.SpawnEntity(&world.Entity{
		Level: 0,
		Position: world.Position{
			X: 100,
			Y: 100,
		},
		HatRenderComponent: true,
		RenderComponent: world.RenderComponent{
			Layer: 0,
		},
		HatImageRenderComponent: true,
		ImageRenderComponent: world.ImageRenderComponent{
			Image: requireImage("tisch"),
			Scale: 1,
		},
		HatRendererHitboxComponent: true,
	})
	w.SpawnEntity(&world.Entity{
		Level: 0,
		Position: world.Position{
			X: 180,
			Y: 100,
		},
		HatRenderComponent: true,
		RenderComponent: world.RenderComponent{
			Layer: 0,
		},
		HatImageRenderComponent: true,
		ImageRenderComponent: world.ImageRenderComponent{
			Image: requireImage("tisch"),
			Scale: 2,
		},
		HatRendererHitboxComponent: true,
	})
	w.SpawnEntity(&world.Entity{
		Level: 0,
		Position: world.Position{
			X: 300,
			Y: 100,
		},
		HatRenderComponent: true,
		RenderComponent: world.RenderComponent{
			Layer: 0,
		},
		HatImageRenderComponent: true,
		ImageRenderComponent: world.ImageRenderComponent{
			Image: requireImage("tisch"),
			Scale: 1.798756,
		},
		HatRendererHitboxComponent: true,
	})
	return w
}
