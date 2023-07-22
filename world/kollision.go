package world

import (
	"math"
	"sort"
)

type HitboxComponent struct {
	Width, Height float64
}

type RendererHitboxComponent struct{}

type RigidbodyComponent struct{}

func (w *World) rendererHitboxenAnwenden() {
	for entity := range w.Entites {
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
			entity.HitboxComponent = HitboxComponent{
				Width:  float64(entity.ImageRenderComponent.Image.Bounds().Dx()) * entity.ImageRenderComponent.Scale,
				Height: float64(entity.ImageRenderComponent.Image.Bounds().Dy()) * entity.ImageRenderComponent.Scale,
			}
		}
	}
}

type kollisionsReaktion struct {
	xVerschiebung, yVerschiebung         float64
	xVelocityAnpassen, yVelocityAnpassen bool
	xVelocityRichtung, yVelocityRichtung float64
}

func (k kollisionsReaktion) apply(entity *Entity) {
	entity.Position.X += k.xVerschiebung
	entity.Position.Y += k.yVerschiebung

	if entity.HatVelocityComponent {
		if k.xVelocityAnpassen {
			entity.VelocityComponent.VelocityX = k.xVelocityRichtung * math.Abs(entity.VelocityComponent.VelocityX)
		}
		if k.yVelocityAnpassen {
			entity.VelocityComponent.VelocityY = k.yVelocityRichtung * math.Abs(entity.VelocityComponent.VelocityY)
		}
	}
}

func getKollisionReaktion(entity, hindernis *Entity) kollisionsReaktion {
	entityAabb := entity.aabb()
	hindernisAabb := hindernis.aabb()

	möglicheReaktionen := make([]kollisionsReaktion, 0, 4)

	nachUntenAusweichenNeuesY := hindernisAabb.MaxY()
	nachUntenAusweichenVerschiebung := nachUntenAusweichenNeuesY - entityAabb.Y
	möglicheReaktionen = append(möglicheReaktionen, kollisionsReaktion{
		xVerschiebung:     0,
		yVerschiebung:     nachUntenAusweichenVerschiebung,
		xVelocityAnpassen: false,
		yVelocityAnpassen: true,
		yVelocityRichtung: 1,
	})

	nachObenAusweichenNeuesY := hindernisAabb.Y - entityAabb.Height
	nachObenAusweichenVerschiebung := nachObenAusweichenNeuesY - entityAabb.Y
	möglicheReaktionen = append(möglicheReaktionen, kollisionsReaktion{
		xVerschiebung:     0,
		yVerschiebung:     nachObenAusweichenVerschiebung,
		xVelocityAnpassen: false,
		yVelocityAnpassen: true,
		yVelocityRichtung: -1,
	})

	nachRechtsAusweichenNeuesX := hindernisAabb.MaxX()
	nachRechtsAusweichenVerschiebung := nachRechtsAusweichenNeuesX - entityAabb.X
	möglicheReaktionen = append(möglicheReaktionen, kollisionsReaktion{
		xVerschiebung:     nachRechtsAusweichenVerschiebung,
		yVerschiebung:     0,
		xVelocityAnpassen: true,
		xVelocityRichtung: 1,
		yVelocityAnpassen: false,
	})

	nachLinksAusweichenNeuesX := hindernisAabb.X - entityAabb.Width
	nachLinksAusweichenVerschiebung := nachLinksAusweichenNeuesX - entityAabb.X
	möglicheReaktionen = append(möglicheReaktionen, kollisionsReaktion{
		xVerschiebung:     nachLinksAusweichenVerschiebung,
		yVerschiebung:     0,
		xVelocityAnpassen: true,
		xVelocityRichtung: -1,
		yVelocityAnpassen: false,
	})

	sort.Slice(möglicheReaktionen, func(i int, j int) bool {
		iVerschiebungGesamt := math.Abs(möglicheReaktionen[i].xVerschiebung) + math.Abs(möglicheReaktionen[i].yVerschiebung)
		jVerschiebungGesamt := math.Abs(möglicheReaktionen[j].xVerschiebung) + math.Abs(möglicheReaktionen[j].yVerschiebung)
		return iVerschiebungGesamt < jVerschiebungGesamt
	})
	reaktionMitKürzesterVerschiebung := möglicheReaktionen[0]

	return reaktionMitKürzesterVerschiebung
}

func (w *World) kollisionenVerarbeiten() {
	var reaktionen []func()

	for entity := range w.Entites {
		if !entity.HatRigidbodyComponent {
			continue
		}

		entityAabb := entity.aabb()

		for hindernis := range w.Entites {
			if entity == hindernis {
				continue
			}

			if entity.Level != hindernis.Level {
				continue
			}

			if !hindernis.HatHitboxComponent {
				continue
			}

			hindernisAabb := hindernis.aabb()

			if !entityAabb.KollidiertMit(hindernisAabb) {
				continue
			}

			entity := entity
			reaktion := getKollisionReaktion(entity, hindernis)
			reaktionen = append(reaktionen, func() {
				reaktion.apply(entity)
			})
		}
	}

	for _, reaktion := range reaktionen {
		reaktion()
	}
}
