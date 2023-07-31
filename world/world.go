package world

import (
	"image/color"

	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type World struct {
	Entities map[*Entity]struct{}

	Backgrounds map[Level]color.Color
	LevelNames  map[Level]string

	blockedPathfindingTiles map[Level]map[pathfindingGridTile]struct{}

	levelNameWidget *ui.Text

	debug         bool
	debugPasswort []rune
}

func NewEmptyWorld() *World {
	return &World{
		Entities: make(map[*Entity]struct{}),

		Backgrounds: make(map[Level]color.Color),
		LevelNames:  make(map[Level]string),

		levelNameWidget: ui.NewText(ui.TextConfig{
			Position: ui.Position{
				X:                ui.Width - 10,
				Y:                ui.Height - 10,
				AnchorHorizontal: ui.HorizontalerAnchorRechts,
				AnchorVertikal:   ui.VertikalerAnchorUnten,
			},
			Text:               "",
			CustomColorPalette: true,
			ColorPalette: ui.TextColorPalatte{
				Color:      colornames.Whitesmoke,
				HoverColor: colornames.White,
			},
		}),
	}
}

func (w *World) SpawnEntity(entity *Entity) *Entity {
	w.Entities[entity] = struct{}{}
	return entity
}

func (w *World) FindEntitiesWithTag(tag string) []*Entity {
	var result []*Entity
	for entity := range w.Entities {
		if entity.HasTag(tag) {
			result = append(result, entity)
		}
	}
	return result
}

func (w *World) FindEntityWithName(name string) *Entity {
	for entity := range w.Entities {
		if entity.Name == name {
			return entity
		}
	}
	return nil
}

func (w *World) Update() {
	w.applyStaticToEntities()
	w.rendererHitboxenAnwenden()
	w.initBlockedPathfindingTiles()
	w.applyVelocityToEntities()
	w.entitiesMitKeyboardSteuern()
	w.entitiesMitTouchSteuern()
	w.pathfind()
	w.moveEntitiesToPositions()
	w.moveEntitiesToPosition()
	w.teleportEntitiesTouchingPortal()
	w.kollisionenVerarbeiten()
	w.interaktionenHandeln()
	w.debugPasswortLesen()
	w.updateLevelName()
}

func (w *World) Draw(screen *ebiten.Image) {
	w.drawEntities(screen)
	w.pfadeVisualisierenDebug(screen)
	w.levelNameWidget.Draw(screen)
}
