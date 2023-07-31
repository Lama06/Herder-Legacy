package herder

import (
	"math"

	"github.com/Lama06/Herder-Legacy/assets"
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/Lama06/Herder-Legacy/world"
)

func CreateHerder(herderLegacy herderlegacy.HerderLegacy) *world.World {
	w := world.NewEmptyWorld()
	w.LevelNames[0] = "Weber"
	w.LevelNames[1] = "E202"
	bodenImg := assets.RequireImage("boden.png")
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
			X: 0,
			Y: 0,
		},
		HatRenderComponent: true,
		RenderComponent: world.RenderComponent{
			Layer: 0,
		},
		HatImageRenderComponent: true,
		ImageRenderComponent: world.ImageRenderComponent{
			Image: assets.RequireImage("tisch.png"),
			Scale: 1,
		},
		HatRendererHitboxComponent: true,
		HatPathfinderComponent:     true,
		PathfinderComponent: world.PathfinderComponent{
			DestinationPosition: world.Position{
				X: 500,
				Y: 500,
			},
			DestinationLevel: 1,
			Speed:            1,
		},
	})
	w.SpawnEntity(&world.Entity{
		Level: 0,
		Position: world.Position{
			X: 180,
			Y: 100,
		},
		Static:             true,
		HatRenderComponent: true,
		RenderComponent: world.RenderComponent{
			Layer: 0,
		},
		HatImageRenderComponent: true,
		ImageRenderComponent: world.ImageRenderComponent{
			Image:    assets.RequireImage("tisch.png"),
			Scale:    2,
			Rotation: math.Pi / 4,
		},
		HatRendererHitboxComponent: true,
		HatInteraktionComponent:    true,
		InteraktionComponent: world.InteraktionComponent{
			OffsetX: 0,
			OffsetY: 0,
			Width:   50,
			Height:  50,
			Callback: func() {
				thisScreen := herderLegacy.CurrentScreen()

				herderLegacy.OpenScreen(ui.NewCountdown(ui.CountdownConfig{
					Position:  ui.NewCenteredPosition(ui.Width/2, ui.Height/2),
					StartZeit: 10 * 60,
					AbgelaufenCallback: func() {
						herderLegacy.OpenScreen(thisScreen)
					},
				}))
			},
		},
	})
	w.SpawnEntity(&world.Entity{
		Level: 0,
		Position: world.Position{
			X: 500,
			Y: 200,
		},
		Static:             true,
		HatRenderComponent: true,
		RenderComponent: world.RenderComponent{
			Layer: 0,
		},
		HatImageRenderComponent: true,
		ImageRenderComponent: world.ImageRenderComponent{
			Image: assets.RequireImage("tisch.png"),
			Scale: 1,
		},
		HatPortalComponent: true,
		PortalComponent: world.PortalComponent{
			Width:               100,
			Height:              100,
			DestinationLevel:    1,
			DestinationPosition: world.Position{X: 0, Y: 0},
		},
	})
	w.SpawnEntity(&world.Entity{
		Level: 1,
		Position: world.Position{
			X: 200,
			Y: 200,
		},
		Static:             true,
		HatRenderComponent: true,
		RenderComponent: world.RenderComponent{
			Layer: 0,
		},
		HatImageRenderComponent: true,
		ImageRenderComponent: world.ImageRenderComponent{
			Image: assets.RequireImage("tisch.png"),
			Scale: 1,
		},
		HatPortalComponent: true,
		PortalComponent: world.PortalComponent{
			Width:               100,
			Height:              100,
			DestinationLevel:    0,
			DestinationPosition: world.Position{X: 0, Y: 0},
		},
	})
	return w
}
