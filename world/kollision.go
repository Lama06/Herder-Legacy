package world

import (
	"math"
	"sort"
)

type HitboxComponent struct {
	Width, Height float64
}

type ImageHitboxComponent struct{}

type KollisionenVerhindernComponent struct{}

func kollisionVerhindern(entity, hindernis *Entity) {
	type möglichkeit struct {
		xVerschiebung, yVerschiebung float64
	}

	entityAabb := entity.aabb()
	hindernisAabb := hindernis.aabb()

	möglichkeiten := make([]möglichkeit, 0, 4)

	nachUntenAusweichenNeuesY := hindernisAabb.MaxY()
	nachUntenAusweichenVerschiebung := nachUntenAusweichenNeuesY - entityAabb.Y
	möglichkeiten = append(möglichkeiten, möglichkeit{
		xVerschiebung: 0,
		yVerschiebung: nachUntenAusweichenVerschiebung,
	})

	nachObenAusweichenNeuesY := hindernisAabb.Y - entityAabb.Height
	nachObenAusweichenVerschiebung := nachObenAusweichenNeuesY - entityAabb.Y
	möglichkeiten = append(möglichkeiten, möglichkeit{
		xVerschiebung: 0,
		yVerschiebung: nachObenAusweichenVerschiebung,
	})

	nachRechtsAusweichenNeuesX := hindernisAabb.MaxX()
	nachRechtsAusweichenVerschiebung := nachRechtsAusweichenNeuesX - entityAabb.X
	möglichkeiten = append(möglichkeiten, möglichkeit{
		xVerschiebung: nachRechtsAusweichenVerschiebung,
		yVerschiebung: 0,
	})

	nachLinksAusweichenNeuesX := hindernisAabb.X - entityAabb.Width
	nachLinksAusweichenVerschiebung := nachLinksAusweichenNeuesX - entityAabb.X
	möglichkeiten = append(möglichkeiten, möglichkeit{
		xVerschiebung: nachLinksAusweichenVerschiebung,
		yVerschiebung: 0,
	})

	sort.Slice(möglichkeiten, func(i int, j int) bool {
		iVerschiebungGesamt := math.Abs(möglichkeiten[i].xVerschiebung) + math.Abs(möglichkeiten[i].yVerschiebung)
		jVerschiebungGesamt := math.Abs(möglichkeiten[j].xVerschiebung) + math.Abs(möglichkeiten[j].yVerschiebung)
		return iVerschiebungGesamt < jVerschiebungGesamt
	})
	kürzesteVerschiebung := möglichkeiten[0]

	entity.Position.X += kürzesteVerschiebung.xVerschiebung
	entity.Position.Y += kürzesteVerschiebung.yVerschiebung
}

func (w *World) kollisionenVerhindern() {
	for entity := range w.Entites {
		if !entity.HatKollisionenVerhindernComponent {
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

			hindernisAabb := hindernis.aabb()

			if !entityAabb.KollidiertMit(hindernisAabb) {
				continue
			}

			kollisionVerhindern(entity, hindernis)
		}
	}
}
