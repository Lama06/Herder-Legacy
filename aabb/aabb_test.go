package aabb_test

import (
	"math"
	"testing"

	"github.com/Lama06/Herder-Legacy/aabb"
)

func floatsRoughlyEqual(a, b float64) bool {
	const tolerance = 0.00000001
	return math.Abs(a-b) <= tolerance
}

func aabbsRoughlyEqual(a, b aabb.Aabb) bool {
	return floatsRoughlyEqual(a.X, b.X) &&
		floatsRoughlyEqual(a.Y, b.Y) &&
		floatsRoughlyEqual(a.Width, b.Width) &&
		floatsRoughlyEqual(a.Height, b.Height)
}

func TestCenter(t *testing.T) {
	testCases := []struct {
		input            aabb.Aabb
		centerX, centerY float64
	}{
		{
			input: aabb.Aabb{
				X:      10,
				Y:      50,
				Width:  40,
				Height: 2,
			},
			centerX: 30,
			centerY: 51,
		},
		{
			input: aabb.Aabb{
				X:      0.5,
				Y:      1,
				Width:  0.5,
				Height: 3,
			},
			centerX: 0.75,
			centerY: 2.5,
		},
		{
			input: aabb.Aabb{
				X:      -20,
				Y:      -100,
				Width:  50,
				Height: 1,
			},
			centerX: 5,
			centerY: -99.5,
		},
	}

	for _, testCase := range testCases {
		if centerX := testCase.input.CenterX(); !floatsRoughlyEqual(centerX, testCase.centerX) {
			t.Errorf("expected: %v, got: %v", testCase.centerX, centerX)
		}

		if centerY := testCase.input.CenterY(); !floatsRoughlyEqual(centerY, testCase.centerY) {
			t.Errorf("expected: %v, got: %v", testCase.centerY, centerY)
		}
	}
}

func TestMax(t *testing.T) {
	testCases := []struct {
		input      aabb.Aabb
		maxX, maxY float64
	}{
		{
			input: aabb.Aabb{
				X:      10,
				Y:      50,
				Width:  40,
				Height: 2,
			},
			maxX: 50,
			maxY: 52,
		},
		{
			input: aabb.Aabb{
				X:      0.5,
				Y:      1,
				Width:  0.5,
				Height: 3,
			},
			maxX: 1,
			maxY: 4,
		},
		{
			input: aabb.Aabb{
				X:      -20,
				Y:      -100,
				Width:  50,
				Height: 1,
			},
			maxX: 30,
			maxY: -99,
		},
	}

	for _, testCase := range testCases {
		if maxX := testCase.input.MaxX(); !floatsRoughlyEqual(maxX, testCase.maxX) {
			t.Errorf("expected: %v, got: %v", testCase.maxX, maxX)
		}

		if maxY := testCase.input.MaxY(); !floatsRoughlyEqual(maxY, testCase.maxY) {
			t.Errorf("expected: %v, got: %v", testCase.maxY, maxY)
		}
	}
}

func TestArea(t *testing.T) {
	testCases := []struct {
		input aabb.Aabb
		area  float64
	}{
		{
			input: aabb.Aabb{
				X:      10,
				Y:      50,
				Width:  40,
				Height: 2,
			},
			area: 80,
		},
		{
			input: aabb.Aabb{
				X:      0.5,
				Y:      1,
				Width:  0.5,
				Height: 3,
			},
			area: 1.5,
		},
		{
			input: aabb.Aabb{
				X:      -20,
				Y:      -100,
				Width:  50,
				Height: 1,
			},
			area: 50,
		},
	}

	for _, testCase := range testCases {
		if area := testCase.input.Area(); !floatsRoughlyEqual(area, testCase.area) {
			t.Errorf("expected: %v, got: %v", testCase.area, area)
		}
	}
}

func TestKollission(t *testing.T) {
	testCases := []struct {
		first, second aabb.Aabb
		kollision     bool
	}{
		{
			first: aabb.Aabb{
				X:      0,
				Y:      0,
				Width:  100,
				Height: 100,
			},
			second: aabb.Aabb{
				X:      -50,
				Y:      -50,
				Width:  55,
				Height: 500,
			},
			kollision: true,
		},
		{
			first: aabb.Aabb{
				X:      0.5,
				Y:      0.5,
				Width:  0.5,
				Height: 0.5,
			},
			second: aabb.Aabb{
				X:      0.9,
				Y:      0.4,
				Width:  0.2,
				Height: 0.2,
			},
			kollision: true,
		},
		{
			first: aabb.Aabb{
				X:      0,
				Y:      0,
				Width:  13,
				Height: 7,
			},
			second: aabb.Aabb{
				X:      -3,
				Y:      -3,
				Width:  3,
				Height: 3,
			},
			kollision: false,
		},
		{
			first: aabb.Aabb{
				X:      13,
				Y:      13,
				Width:  1,
				Height: 1,
			},
			second: aabb.Aabb{
				X:      7,
				Y:      7,
				Width:  3,
				Height: 4.5,
			},
			kollision: false,
		},
	}

	for _, testCase := range testCases {
		if kollision := testCase.first.KollidiertMit(testCase.second); kollision != testCase.kollision {
			t.Errorf("expected: %v, got: %v", testCase.kollision, kollision)
		}

		if kollision := testCase.second.KollidiertMit(testCase.first); kollision != testCase.kollision {
			t.Errorf("expected: %v, got: %v", testCase.kollision, kollision)
		}
	}
}

func TestIsInside(t *testing.T) {
	testCases := []struct {
		hitbox aabb.Aabb
		x, y   float64
		inside bool
	}{
		{
			hitbox: aabb.Aabb{
				X:      0,
				Y:      0,
				Width:  50,
				Height: 50,
			},
			x:      25,
			y:      25,
			inside: true,
		},
		{
			hitbox: aabb.Aabb{
				X:      -1,
				Y:      -2,
				Width:  2,
				Height: 3,
			},
			x:      0.2,
			y:      0.3,
			inside: true,
		},
		{
			hitbox: aabb.Aabb{
				X:      0,
				Y:      0,
				Width:  1,
				Height: 1,
			},
			x:      1,
			y:      1,
			inside: true,
		},
		{
			hitbox: aabb.Aabb{
				X:      0,
				Y:      0,
				Width:  1,
				Height: 1,
			},
			x:      1.1,
			y:      0.5,
			inside: false,
		},
	}

	for _, testCase := range testCases {
		if inside := testCase.hitbox.IsInside(testCase.x, testCase.y); inside != testCase.inside {
			t.Errorf("expected: %v, got: %v", testCase.inside, inside)
		}
	}
}

func TestIntersection(t *testing.T) {
	testCases := []struct {
		hitbox1, hitbox2 aabb.Aabb
		hasIntersection  bool
		intersection     aabb.Aabb
	}{
		{
			hitbox1: aabb.Aabb{
				X:      0,
				Y:      0,
				Width:  50,
				Height: 50,
			},
			hitbox2: aabb.Aabb{
				X:      45,
				Y:      45,
				Width:  100,
				Height: 13,
			},
			hasIntersection: true,
			intersection: aabb.Aabb{
				X:      45,
				Y:      45,
				Width:  5,
				Height: 5,
			},
		},
		{
			hitbox1: aabb.Aabb{
				X:      0,
				Y:      0,
				Width:  50,
				Height: 50,
			},
			hitbox2: aabb.Aabb{
				X:      51,
				Y:      50.01,
				Width:  100,
				Height: 13,
			},
			hasIntersection: false,
			intersection:    aabb.Aabb{},
		},
	}

	for _, testCase := range testCases {
		intersection, hasIntersection := testCase.hitbox1.Intersection(testCase.hitbox2)
		if hasIntersection != testCase.hasIntersection && intersection != testCase.intersection {
			t.Errorf("expected: %v, %v, got: %v, %v", testCase.intersection, testCase.hasIntersection, intersection, hasIntersection)
		}

		intersection, hasIntersection = testCase.hitbox2.Intersection(testCase.hitbox1)
		if hasIntersection != testCase.hasIntersection && intersection != testCase.intersection {
			t.Errorf("expected: %v, %v, got: %v, %v", testCase.intersection, testCase.hasIntersection, intersection, hasIntersection)
		}
	}
}

func TestVertikalZerschneiden(t *testing.T) {
	testCases := []struct {
		hitbox      aabb.Aabb
		prozentOben float64
		oben, unten aabb.Aabb
	}{
		{
			hitbox: aabb.Aabb{
				X:      0,
				Y:      0,
				Width:  77,
				Height: 50,
			},
			prozentOben: 0.1,
			oben: aabb.Aabb{
				X:      0,
				Y:      0,
				Width:  77,
				Height: 5,
			},
			unten: aabb.Aabb{
				X:      0,
				Y:      5,
				Width:  77,
				Height: 45,
			},
		},
		{
			hitbox: aabb.Aabb{
				X:      -13,
				Y:      -1,
				Width:  1,
				Height: 2,
			},
			prozentOben: 0.9,
			oben: aabb.Aabb{
				X:      -13,
				Y:      -1,
				Width:  1,
				Height: 1.8,
			},
			unten: aabb.Aabb{
				X:      -13,
				Y:      0.8,
				Width:  1,
				Height: 0.2,
			},
		},
	}

	for _, testCase := range testCases {
		oben, unten := testCase.hitbox.VertikalZerschneiden(testCase.prozentOben)
		if !aabbsRoughlyEqual(oben, testCase.oben) || !aabbsRoughlyEqual(unten, testCase.unten) {
			t.Errorf("expected: %v, %v, got: %v, %v", testCase.oben, testCase.unten, oben, unten)
		}
	}
}

func TestHorizontalZerschneiden(t *testing.T) {
	testCases := []struct {
		hitbox        aabb.Aabb
		prozentLinks  float64
		links, rechts aabb.Aabb
	}{
		{
			hitbox: aabb.Aabb{
				X:      0,
				Y:      0,
				Width:  70,
				Height: 50,
			},
			prozentLinks: 0.1,
			links: aabb.Aabb{
				X:      0,
				Y:      0,
				Width:  7,
				Height: 50,
			},
			rechts: aabb.Aabb{
				X:      7,
				Y:      0,
				Width:  63,
				Height: 50,
			},
		},
		{
			hitbox: aabb.Aabb{
				X:      -20,
				Y:      -1,
				Width:  1,
				Height: 2,
			},
			prozentLinks: 0.9,
			links: aabb.Aabb{
				X:      -20,
				Y:      -1,
				Width:  0.9,
				Height: 2,
			},
			rechts: aabb.Aabb{
				X:      -19.1,
				Y:      -1,
				Width:  0.1,
				Height: 2,
			},
		},
	}

	for _, testCase := range testCases {
		links, rechts := testCase.hitbox.HorizontalZerschneiden(testCase.prozentLinks)
		if !aabbsRoughlyEqual(links, testCase.links) || !aabbsRoughlyEqual(rechts, testCase.rechts) {
			t.Errorf("expected: %v, %v, got: %v, %v", testCase.links, testCase.rechts, links, rechts)
		}
	}
}
