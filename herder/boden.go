package herder

import (
	"fmt"

	"github.com/Lama06/Herder-Legacy/assets"
	"github.com/Lama06/Herder-Legacy/world"
)

type bodenArt string

const (
	bodenArtParkett bodenArt = "parkett"
)

func createBoden(
	w *world.World,
	level world.Level,
	boden bodenArt,
	x, y float64,
	zeilen, spalten int,
) {
	bodenImg := assets.RequireImage(fmt.Sprintf("boden/%v.png", boden))
	for zeile := 0; zeile < zeilen; zeile++ {
		for spalte := 0; spalte < spalten; spalte++ {
			w.SpawnEntity(&world.Entity{
				Level: level,
				Position: world.Position{
					X: float64(spalte) * float64(bodenImg.Bounds().Dx()) * 4,
					Y: float64(zeile) * float64(bodenImg.Bounds().Dy()) * 4,
				},
				HatRenderComponent: true,
				RenderComponent: world.RenderComponent{
					Layer: renderLayerBoden,
				},
				HatImageRenderComponent: true,
				ImageRenderComponent: world.ImageRenderComponent{
					Image: bodenImg,
					Scale: 4,
				},
			})
		}
	}
}
