package sodoku

import (
	"testing"
)

func TestLösen(t *testing.T) {
	testCases := []struct {
		input sodoku

		hatLösung bool
		lösung    sodoku
	}{
		{
			input: sodoku{
				{4, 1, 0 /**/, 0, 6, 5 /**/, 0, 0, 7},
				{0, 0, 6 /**/, 0, 0, 7 /**/, 4, 8, 0},
				{2, 0, 7 /**/, 4, 9, 0 /**/, 0, 0, 6},
				/***********************************/
				{0, 6, 0 /**/, 0, 7, 0 /**/, 1, 0, 0},
				{3, 0, 1 /**/, 5, 0, 0 /**/, 0, 7, 2},
				{0, 9, 0 /**/, 0, 4, 2 /**/, 3, 0, 8},
				/***********************************/
				{1, 0, 8 /**/, 6, 0, 0 /**/, 0, 2, 9},
				{0, 2, 0 /**/, 0, 1, 8 /**/, 6, 4, 0},
				{6, 0, 0 /**/, 3, 0, 0 /**/, 0, 1, 0},
			},

			hatLösung: true,
			lösung: sodoku{
				{4, 1, 3 /**/, 8, 6, 5 /**/, 2, 9, 7},
				{9, 5, 6 /**/, 2, 3, 7 /**/, 4, 8, 1},
				{2, 8, 7 /**/, 4, 9, 1 /**/, 5, 3, 6},
				/***********************************/
				{8, 6, 2 /**/, 9, 7, 3 /**/, 1, 5, 4},
				{3, 4, 1 /**/, 5, 8, 6 /**/, 9, 7, 2},
				{7, 9, 5 /**/, 1, 4, 2 /**/, 3, 6, 8},
				/***********************************/
				{1, 3, 8 /**/, 6, 5, 4 /**/, 7, 2, 9},
				{5, 2, 9 /**/, 7, 1, 8 /**/, 6, 4, 3},
				{6, 7, 4 /**/, 3, 2, 9 /**/, 8, 1, 5},
			},
		},
	}

	for _, testCase := range testCases {
		lösung, ok := testCase.input.lösen()
		if ok != testCase.hatLösung || lösung != testCase.lösung {
			t.Fail()
		}
	}
}