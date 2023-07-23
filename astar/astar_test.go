package astar_test

import (
	"math"
	"testing"

	"github.com/Lama06/Herder-Legacy/astar"
)

type position struct {
	x, y int
}

var _ astar.Node = position{}

func (p position) AstarEstimateCost(context any, to astar.Node) float64 {
	toPosition := to.(position)
	xDiff := math.Abs(float64(p.x) - float64(toPosition.x))
	yDiff := math.Abs(float64(p.y) - float64(toPosition.y))
	return math.Sqrt(xDiff*xDiff + yDiff*yDiff)
}

func (p position) AstarNeighbours(context any) []astar.Neighbour {
	blockedPositions := context.(map[position]struct{})

	var neighbours []astar.Neighbour
	for _, xOffset := range [...]int{-1, 0, 1} {
		for _, yOffset := range [...]int{-1, 0, 1} {
			if xOffset == 0 && yOffset == 0 {
				continue
			}

			neighbourPos := position{
				x: p.x + xOffset,
				y: p.y + yOffset,
			}

			if _, isBlocked := blockedPositions[neighbourPos]; isBlocked {
				continue
			}

			var cost float64
			if xOffset == 0 || yOffset == 0 {
				cost = 1
			} else {
				cost = math.Sqrt(2)
			}

			neighbours = append(neighbours, astar.Neighbour{
				Node: neighbourPos,
				Cost: cost,
			})
		}
	}
	return neighbours
}

func pathsEqual(path1 []astar.Node, path2 []position) bool {
	if len(path1) != len(path2) {
		return false
	}

	for i := range path1 {
		if path1[i] != path2[i] {
			return false
		}
	}

	return true
}

func TestAstar(t *testing.T) {
	testCases := map[string]struct {
		blockedPositions map[position]struct{}
		from             position
		to               position
		pathFound        bool
		path             []position
	}{
		"straight path": {
			blockedPositions: nil,
			from:             position{x: 0, y: 0},
			to:               position{x: 10, y: 0},
			pathFound:        true,
			path: []position{
				{x: 0, y: 0},
				{x: 1, y: 0},
				{x: 2, y: 0},
				{x: 3, y: 0},
				{x: 4, y: 0},
				{x: 5, y: 0},
				{x: 6, y: 0},
				{x: 7, y: 0},
				{x: 8, y: 0},
				{x: 9, y: 0},
				{x: 10, y: 0},
			},
		},
		"diagonal path": {
			blockedPositions: nil,
			from:             position{x: 0, y: 0},
			to:               position{x: -10, y: 10},
			pathFound:        true,
			path: []position{
				{x: 0, y: 0},
				{x: -1, y: 1},
				{x: -2, y: 2},
				{x: -3, y: 3},
				{x: -4, y: 4},
				{x: -5, y: 5},
				{x: -6, y: 6},
				{x: -7, y: 7},
				{x: -8, y: 8},
				{x: -9, y: 9},
				{x: -10, y: 10},
			},
		},
		"kein path möglich": {
			blockedPositions: map[position]struct{}{
				{x: -1, y: -1}: {},
				{x: -1, y: 0}:  {},
				{x: -1, y: 1}:  {},
				{x: 0, y: -1}:  {},
				{x: 0, y: 1}:   {},
				{x: 1, y: -1}:  {},
				{x: 1, y: 0}:   {},
				{x: 1, y: 1}:   {},
			},
			from:      position{x: 0, y: 0},
			to:        position{x: 13, y: 13},
			pathFound: false,
			path:      nil,
		},
		"path um Hindernis": {
			// X
			// XX
			// SXZ
			// XX
			blockedPositions: map[position]struct{}{
				{x: 1, y: -1}: {},
				{x: 1, y: 0}:  {},
				{x: 1, y: 1}:  {},
				{x: 0, y: -1}: {},
				{x: 0, y: -2}: {},
				{x: 0, y: 1}:  {},
			},
			from:      position{x: 0, y: 0},
			to:        position{x: 2, y: 0},
			pathFound: true,
			path: []position{
				{x: 0, y: 0},
				{x: -1, y: 1},
				{x: 0, y: 2},
				{x: 1, y: 2},
				{x: 2, y: 1},
				{x: 2, y: 0},
			},
		},
	}

	for name, testCase := range testCases {
		path, pathFound := astar.FindShortestPath(testCase.blockedPositions, testCase.from, testCase.to)
		if pathFound != testCase.pathFound || !pathsEqual(path, testCase.path) {
			t.Errorf(
				"%v, expected: %v, %v, got: %v, %v",
				name,
				testCase.path, testCase.pathFound,
				path, pathFound,
			)
		}
	}
}
