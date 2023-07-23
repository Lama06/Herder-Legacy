package world

import (
	"image/color"
	"sort"

	"github.com/Lama06/Herder-Legacy/aabb"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CameraComponent struct {
	OffsetX, OffsetY float64
}

type RenderLayer int

type RenderComponent struct {
	Layer RenderLayer
}

type RectRenderComponent struct {
	Width, Height float64
	Farbe         color.Color
}

type KreisRenderComponent struct {
	Size  float64
	Farbe color.Color
}

type ImageRenderComponent struct {
	Image *ebiten.Image
	Scale float64
}

func (w *World) findCamera() *Entity {
	for entity := range w.Entities {
		if !entity.HatCameraComponent {
			continue
		}
		return entity
	}
	return nil
}

func (w *World) drawEntities(screen *ebiten.Image) {
	camera := w.findCamera()
	if camera == nil {
		return
	}

	var renderableEntities []*Entity
	for entity := range w.Entities {
		if !entity.HatRenderComponent {
			continue
		}

		if entity.Level != camera.Level {
			continue
		}

		renderableEntities = append(renderableEntities, entity)
	}

	sort.Slice(renderableEntities, func(i, j int) bool {
		return renderableEntities[i].RenderComponent.Layer < renderableEntities[j].RenderComponent.Layer
	})

	for _, entity := range renderableEntities {
		screenX := ui.Width/2 + entity.Position.X - camera.Position.X + camera.CameraComponent.OffsetX
		screenY := ui.Height/2 + entity.Position.Y - camera.Position.Y + camera.CameraComponent.OffsetY

		screenAabb := aabb.Aabb{X: 0, Y: 0, Width: ui.Width, Height: ui.Height}

		switch {
		case entity.HatRectRenderComponent:
			rectAabb := aabb.Aabb{
				X:      screenX,
				Y:      screenY,
				Width:  entity.RectRenderComponent.Width,
				Height: entity.RectRenderComponent.Height,
			}
			if !rectAabb.KollidiertMit(screenAabb) {
				continue
			}

			vector.DrawFilledRect(
				screen,
				float32(screenX),
				float32(screenY),
				float32(entity.RectRenderComponent.Width),
				float32(entity.RectRenderComponent.Height),
				entity.KreisRenderComponent.Farbe,
				true,
			)
		case entity.HatKreisRenderComponent:
			circleAabb := aabb.Aabb{
				X:      screenX,
				Y:      screenY,
				Width:  entity.KreisRenderComponent.Size,
				Height: entity.KreisRenderComponent.Size,
			}
			if !circleAabb.KollidiertMit(screenAabb) {
				continue
			}

			centerX := screenX + entity.KreisRenderComponent.Size/2
			centerY := screenY + entity.KreisRenderComponent.Size/2
			vector.DrawFilledCircle(
				screen,
				float32(centerX),
				float32(centerY),
				float32(entity.KreisRenderComponent.Size)/2,
				entity.KreisRenderComponent.Farbe,
				true,
			)
		case entity.HatImageRenderComponent:
			scale := entity.ImageRenderComponent.Scale
			if scale == 0 {
				scale = 1
			}

			imageAabb := aabb.Aabb{
				X:      screenX,
				Y:      screenY,
				Width:  float64(entity.ImageRenderComponent.Image.Bounds().Dx()) * scale,
				Height: float64(entity.ImageRenderComponent.Image.Bounds().Dy()) * scale,
			}
			if !imageAabb.KollidiertMit(screenAabb) {
				continue
			}

			var drawOptions ebiten.DrawImageOptions
			drawOptions.GeoM.Scale(scale, scale)
			drawOptions.GeoM.Translate(screenX, screenY)
			screen.DrawImage(entity.ImageRenderComponent.Image, &drawOptions)
		}
	}
}
