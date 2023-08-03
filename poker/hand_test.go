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
		"Straße": {
			karten: [7]karte{
				mustParseKarte("Karo 6"),
				mustParseKarte("Herz 7"),
				mustParseKarte("Pik 8"),
				mustParseKarte("Pik 10"),
				mustParseKarte("Karo 9"),
				mustParseKarte("Herz 10"),
				mustParseKarte("Kreuz 9"),
			},
			hand: straßeHand{
				mustParseKarte("Herz 10"),
				mustParseKarte("Kreuz 9"),
				mustParseKarte("Pik 8"),
				mustParseKarte("Herz 7"),
				mustParseKarte("Karo 6"),
			},
		},
		"Straight Flush": {
			karten: [7]karte{
				mustParseKarte("Pik Ass"),
				mustParseKarte("Pik König"),
				mustParseKarte("Herz König"),
				mustParseKarte("Herz Dame"),
				mustParseKarte("Herz Bube"),
				mustParseKarte("Herz 10"),
				mustParseKarte("Herz 9"),
			},
			hand: straightFlushHand{
				mustParseKarte("Herz König"),
				mustParseKarte("Herz Dame"),
				mustParseKarte("Herz Bube"),
				mustParseKarte("Herz 10"),
				mustParseKarte("Herz 9"),
			},
		},
		"Royal Flush": {
			karten: [7]karte{
				mustParseKarte("Herz Ass"),
				mustParseKarte("Herz König"),
				mustParseKarte("Herz Dame"),
				mustParseKarte("Herz Bube"),
				mustParseKarte("Herz 10"),
				mustParseKarte("Herz 9"),
				mustParseKarte("Pik 2"),
			},
			hand: royalFlush(symbolHerz),
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
