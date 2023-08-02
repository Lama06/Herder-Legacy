package poker

import "testing"

func TestParseHand(t *testing.T) {
	testCases := map[string]struct {
		karten [7]karte
		hand   hand
	}{
		"Höchste Karte": {
			karten: [7]karte{
				mustParseKarte("Herz König"),
				mustParseKarte("Pik Dame"),
				mustParseKarte("Kreuz Ass"),
				mustParseKarte("Herz 2"),
				mustParseKarte("Herz 3"),
				mustParseKarte("Herz 4"),
				mustParseKarte("Pik 5"),
			},
			hand: höchsteKarteHand{
				mustParseKarte("Kreuz Ass"),
				mustParseKarte("Herz König"),
				mustParseKarte("Pik Dame"),
				mustParseKarte("Pik 5"),
				mustParseKarte("Herz 4"),
			},
		},
		"Ein Paar": {
			karten: [7]karte{
				mustParseKarte("Herz König"),
				mustParseKarte("Kreuz Ass"),
				mustParseKarte("Herz 2"),
				mustParseKarte("Herz 3"),
				mustParseKarte("Herz 4"),
				mustParseKarte("Pik 5"),
				mustParseKarte("Pik König"),
			},
			hand: einPaarHand{
				paar: [2]karte{
					mustParseKarte("Herz König"),
					mustParseKarte("Pik König"),
				},
				beikarten: [3]karte{
					mustParseKarte("Kreuz Ass"),
					mustParseKarte("Pik 5"),
					mustParseKarte("Herz 4"),
				},
			},
		},
	}

	for name, testCase := range testCases {
		got := parseHand(testCase.karten)
		if got != testCase.hand {
			t.Errorf("%v: expected: %T %v, got: %T %v", name, testCase.hand, testCase.hand, got, got)
		}
	}
}

func BenchmarkParseHand(b *testing.B) {
	karten := [7]karte{
		mustParseKarte("Herz König"),
		mustParseKarte("Pik Dame"),
		mustParseKarte("Kreuz Ass"),
		mustParseKarte("Herz 2"),
		mustParseKarte("Herz 3"),
		mustParseKarte("Herz 4"),
		mustParseKarte("Pik 5"),
	}

	for i := 0; i < b.N; i++ {
		parseHand(karten)
	}
}
