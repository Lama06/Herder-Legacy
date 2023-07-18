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
	}
}
