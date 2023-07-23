package world

import (
	"math"

	"github.com/Lama06/Herder-Legacy/aabb"
	"github.com/Lama06/Herder-Legacy/astar"
)

type MoveToPositionComponent struct {
	Position Position
	Speed    float64
	Arrived  bool
}

type MoveToPositionsComponent struct {
	Positions       []Position
	Speed           float64
	CurrentPosition int
	Finished        bool
}

type PathfinderComponentState byte

const (
	PathfinderComponentStateNotStarted PathfinderComponentState = iota
	PathfinderComponentStateMovingToPortal
	PathfinderComponentStateMovingToDestination
	PathfinderComponentStateNoPath
	PathfinderComponentStateArrived
)

type PathfinderComponent struct {
	DestinationPosition Position
	DestinationLevel    Level
	Speed               float64
	State               PathfinderComponentState
}

const (
	pathfinderGridSize = 30
)

type pathfindingGridTile struct {
	x, y int
}

var _ astar.Node = pathfindingGridTile{}

func (p pathfindingGridTile) AstarEstimateCost(context any, other astar.Node) float64 {
	otherTile := other.(pathfindingGridTile)
	xDiff := math.Abs(float64(p.x - otherTile.x))
	yDiff := math.Abs(float64(p.y - otherTile.y))
	return math.Sqrt(xDiff*xDiff + yDiff*yDiff)
}

func (p pathfindingGridTile) AstarNeighbours(context any) []astar.Neighbour {
	blockedTiles := context.(map[pathfindingGridTile]struct{})

	var neighbours []astar.Neighbour
	for _, xOffset := range [...]int{-1, 0, 1} {
		for _, yOffset := range [...]int{-1, 0, 1} {
			if xOffset == 0 && yOffset == 0 {
				continue
			}

			neighbourTile := pathfindingGridTile{
				x: p.x + xOffset,
				y: p.y + yOffset,
			}

			if _, isBlocked := blockedTiles[neighbourTile]; isBlocked {
				continue
			}

			var cost float64
			if xOffset == 0 || yOffset == 0 {
				cost = 1
			} else {
				cost = math.Sqrt(2)
			}

			neighbours = append(neighbours, astar.Neighbour{
				Node: neighbourTile,
				Cost: cost,
			})
		}
	}
	return neighbours
}

func aabbInPathfindingGridZerlegen(hitbox aabb.Aabb) []pathfindingGridTile {
	xStart := int(hitbox.X / pathfinderGridSize)
	yStart := int(hitbox.Y / pathfinderGridSize)
	width := int(hitbox.Width/pathfinderGridSize) + 1
	height := int(hitbox.Height/pathfinderGridSize) + 1
	var tiles []pathfindingGridTile
	for x := xStart; x <= xStart+width; x++ {
		for y := yStart; y <= yStart+height; y++ {
			tileHitbox := aabb.Aabb{
				X:      float64(x) * pathfinderGridSize,
				Y:      float64(y) * pathfinderGridSize,
				Width:  pathfinderGridSize,
				Height: pathfinderGridSize,
			}
			if !tileHitbox.KollidiertMit(hitbox) {
				continue
			}
			tiles = append(tiles, pathfindingGridTile{x: x, y: y})
		}
	}
	return tiles
}

func (w *World) initBlockedPathfindingTiles() {
	if w.blockedPathfindingTiles != nil {
		return
	}
	w.blockedPathfindingTiles = make(map[Level]map[pathfindingGridTile]struct{})

	for entity := range w.Entities {
		if !entity.Static || !entity.HatHitboxComponent {
			continue
		}
		aabb := entity.aabb()
		blockedTiles := aabbInPathfindingGridZerlegen(aabb)

		levelBlockedTiles, ok := w.blockedPathfindingTiles[entity.Level]
		if !ok {
			levelBlockedTiles = make(map[pathfindingGridTile]struct{})
			w.blockedPathfindingTiles[entity.Level] = levelBlockedTiles
		}

		for _, blockedTile := range blockedTiles {
			levelBlockedTiles[blockedTile] = struct{}{}
		}
	}
}

func (w *World) moveEntitiesToPosition() {
	for entity := range w.Entities {
		if !entity.HatMoveToPositionComponent {
			continue
		}

		if entity.MoveToPositionComponent.Arrived {
			continue
		}

		xDistance := entity.MoveToPositionComponent.Position.X - entity.Position.X
		yDistance := entity.MoveToPositionComponent.Position.Y - entity.Position.Y

		if floatsRoughlyEqual(xDistance, 0) && floatsRoughlyEqual(yDistance, 0) {
			entity.Position = entity.MoveToPositionComponent.Position
			entity.MoveToPositionComponent.Arrived = true
			continue
		}

		maxSpeed := entity.MoveToPositionComponent.Speed

		xSpeed := math.Max(-maxSpeed, math.Min(maxSpeed, xDistance))
		ySpeed := math.Max(-maxSpeed, math.Min(maxSpeed, yDistance))

		entity.Position.X += xSpeed
		entity.Position.Y += ySpeed
	}
}

func (w *World) moveEntitiesToPositions() {
	for entity := range w.Entities {
		if !entity.HatMoveToPositionsComponent {
			continue
		}

		if entity.MoveToPositionsComponent.Finished {
			continue
		}

		currentPosition := entity.MoveToPositionsComponent.Positions[entity.MoveToPositionsComponent.CurrentPosition]

		if !entity.HatMoveToPositionComponent || entity.MoveToPositionComponent.Position != currentPosition {
			entity.HatMoveToPositionComponent = true
			entity.MoveToPositionComponent = MoveToPositionComponent{
				Position: currentPosition,
				Speed:    entity.MoveToPositionsComponent.Speed,
			}
			continue
		}

		if !entity.MoveToPositionComponent.Arrived {
			continue
		}

		entity.MoveToPositionsComponent.CurrentPosition++
		if entity.MoveToPositionsComponent.CurrentPosition == len(entity.MoveToPositionsComponent.Positions) {
			entity.MoveToPositionsComponent.Finished = true
		}
	}
}

func findShortestPath(w *World, level Level, from Position, to Position) ([]Position, bool) {
	fromTile := pathfindingGridTile{
		x: int(math.Round(from.X / pathfinderGridSize)),
		y: int(math.Round(from.Y / pathfinderGridSize)),
	}
	toTile := pathfindingGridTile{
		x: int(math.Round(to.X / pathfinderGridSize)),
		y: int(math.Round(to.Y / pathfinderGridSize)),
	}
	nodePath, found := astar.FindShortestPath(w.blockedPathfindingTiles[level], fromTile, toTile)
	if !found {
		return nil, false
	}
	positionPath := make([]Position, len(nodePath))
	for i, node := range nodePath {
		positionPath[i] = Position{
			X: float64(node.(pathfindingGridTile).x) * pathfinderGridSize,
			Y: float64(node.(pathfindingGridTile).y) * pathfinderGridSize,
		}
	}
	positionPath[0] = from
	if len(positionPath) >= 1 {
		positionPath[len(positionPath)-1] = to
	}
	return positionPath, true
}

func findShortestPathInSlice(paths [][]Position) []Position {
	if len(paths) == 0 {
		panic("paths empty")
	}

	var shortest []Position
	for i, path := range paths {
		if i == 0 || len(path) < len(shortest) {
			shortest = path
		}
	}
	return shortest
}

func findShortestPathToPortal(w *World, start Position, level Level) (path []Position, found bool) {
	var possiblePaths [][]Position

	for portal := range w.Entities {
		if !portal.HatPortalComponent || !portal.Static || portal.Level != level {
			continue
		}

		pathToPortal, found := findShortestPath(w, level, start, portal.Position)
		if !found {
			continue
		}
		possiblePaths = append(possiblePaths, pathToPortal)
	}

	if len(possiblePaths) == 0 {
		return nil, false
	}
	return findShortestPathInSlice(possiblePaths), true
}

func findShortestPathFromPortal(w *World, destination Position, level Level) (path []Position, found bool) {
	var possiblePaths [][]Position

	for portal := range w.Entities {
		if !portal.HatPortalComponent || portal.Level != level {
			continue
		}

		pathToDestination, found := findShortestPath(w, level, portal.Position, destination)
		if !found {
			continue
		}
		possiblePaths = append(possiblePaths, pathToDestination)
	}

	if len(possiblePaths) == 0 {
		return nil, false
	}
	return findShortestPathInSlice(possiblePaths), true
}

func (w *World) pathfind() {
	for entity := range w.Entities {
		if !entity.HatPathfinderComponent {
			continue
		}

		switch entity.PathfinderComponent.State {
		case PathfinderComponentStateNoPath, PathfinderComponentStateArrived:
			continue
		case PathfinderComponentStateNotStarted:
			if entity.PathfinderComponent.DestinationLevel != entity.Level {
				path, found := findShortestPathToPortal(w, entity.Position, entity.Level)
				if !found {
					entity.PathfinderComponent.State = PathfinderComponentStateNoPath
					continue
				}
				entity.PathfinderComponent.State = PathfinderComponentStateMovingToPortal
				entity.HatMoveToPositionsComponent = true
				entity.MoveToPositionsComponent = MoveToPositionsComponent{
					Positions: path,
					Speed:     entity.PathfinderComponent.Speed,
				}
				continue
			}

			path, found := findShortestPath(w, entity.Level, entity.Position, entity.PathfinderComponent.DestinationPosition)
			if !found {
				entity.PathfinderComponent.State = PathfinderComponentStateNoPath
				continue
			}
			entity.PathfinderComponent.State = PathfinderComponentStateMovingToDestination
			entity.HatMoveToPositionsComponent = true
			entity.MoveToPositionsComponent = MoveToPositionsComponent{
				Positions: path,
				Speed:     entity.PathfinderComponent.Speed,
			}
		case PathfinderComponentStateMovingToPortal:
			if !entity.MoveToPositionsComponent.Finished {
				continue
			}

			path, found := findShortestPathFromPortal(w, entity.PathfinderComponent.DestinationPosition, entity.PathfinderComponent.DestinationLevel)
			if !found {
				entity.PathfinderComponent.State = PathfinderComponentStateNoPath
				continue
			}
			entity.PathfinderComponent.State = PathfinderComponentStateMovingToDestination
			entity.Level = entity.PathfinderComponent.DestinationLevel
			entity.Position = path[0]
			entity.HatMoveToPositionsComponent = true
			entity.MoveToPositionsComponent = MoveToPositionsComponent{
				Positions: path,
				Speed:     entity.PathfinderComponent.Speed,
			}
		case PathfinderComponentStateMovingToDestination:
			if !entity.MoveToPositionsComponent.Finished {
				continue
			}

			entity.PathfinderComponent.State = PathfinderComponentStateArrived
		}
	}
}
