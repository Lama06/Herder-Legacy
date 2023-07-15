package dame

import "testing"

var zügeTests = map[string]struct {
	ausgangsSituation brett
	amZug             spieler
	regeln            regeln
	erwarteteZüge     züge
}{
	"internationale Regeln: Stein bewegen": {
		ausgangsSituation: mustParseBrett(
			"____l___",
			"_____L__",
			"________",
			"________",
			"________",
			"________",
			"________",
			"_s______",
		),
		amZug:  spielerSchüler,
		regeln: internationaleRegeln,
		erwarteteZüge: züge{
			zug{
				zugSchritt{
					von:                    position{zeile: 7, spalte: 1},
					zu:                     position{zeile: 6, spalte: 2},
					hatGeschlagenePosition: false,
					ergebnis: mustParseBrett(
						"____l___",
						"_____L__",
						"________",
						"________",
						"________",
						"________",
						"__s_____",
						"________",
					),
				},
			},
			zug{
				zugSchritt{
					von:                    position{zeile: 7, spalte: 1},
					zu:                     position{zeile: 6, spalte: 0},
					hatGeschlagenePosition: false,
					ergebnis: mustParseBrett(
						"____l___",
						"_____L__",
						"________",
						"________",
						"________",
						"________",
						"s_______",
						"________",
					),
				},
			},
		},
	},

	"internationale Regeln: mit Stein schlagen": {
		ausgangsSituation: mustParseBrett(
			"________",
			"________",
			"________",
			"________",
			"____L_l_",
			"___s_l__",
			"______l_",
			"________",
		),
		amZug:  spielerSchüler,
		regeln: internationaleRegeln,
		erwarteteZüge: züge{
			zug{
				zugSchritt{
					von:                    position{zeile: 5, spalte: 3},
					zu:                     position{zeile: 3, spalte: 5},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 4, spalte: 4},
					ergebnis: mustParseBrett(
						"________",
						"________",
						"________",
						"_____s__",
						"______l_",
						"_____l__",
						"______l_",
						"________",
					),
				},
				zugSchritt{
					von:                    position{zeile: 3, spalte: 5},
					zu:                     position{zeile: 5, spalte: 7},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 4, spalte: 6},
					ergebnis: mustParseBrett(
						"________",
						"________",
						"________",
						"________",
						"________",
						"_____l_s",
						"______l_",
						"________",
					),
				},
				zugSchritt{
					von:                    position{zeile: 5, spalte: 7},
					zu:                     position{zeile: 7, spalte: 5},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 6, spalte: 6},
					ergebnis: mustParseBrett(
						"________",
						"________",
						"________",
						"________",
						"________",
						"_____l__",
						"________",
						"_____s__",
					),
				},
			},
		},
	},

	"internationale Regeln: mit Stein rückwärts schlagen": {
		ausgangsSituation: mustParseBrett(
			"________",
			"________",
			"____l___",
			"_s_s_s__",
			"________",
			"________",
			"________",
			"________",
		),
		amZug:  spielerLehrer,
		regeln: internationaleRegeln,
		erwarteteZüge: züge{
			zug{
				zugSchritt{
					von:                    position{zeile: 2, spalte: 4},
					zu:                     position{zeile: 4, spalte: 6},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 3, spalte: 5},
					ergebnis: mustParseBrett(
						"________",
						"________",
						"________",
						"_s_s____",
						"______l_",
						"________",
						"________",
						"________",
					),
				},
			},
			zug{
				zugSchritt{
					von:                    position{zeile: 2, spalte: 4},
					zu:                     position{zeile: 4, spalte: 2},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 3, spalte: 3},
					ergebnis: mustParseBrett(
						"________",
						"________",
						"________",
						"_s___s__",
						"__l_____",
						"________",
						"________",
						"________",
					),
				},
				zugSchritt{
					von:                    position{zeile: 4, spalte: 2},
					zu:                     position{zeile: 2, spalte: 0},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 3, spalte: 1},
					ergebnis: mustParseBrett(
						"________",
						"________",
						"l_______",
						"_____s__",
						"________",
						"________",
						"________",
						"________",
					),
				},
			},
		},
	},

	"internationale Regeln: mit Stein schlagen und Stein zur Dame umwandeln": {
		ausgangsSituation: mustParseBrett(
			"________",
			"___L_l__",
			"__s_____",
			"________",
			"________",
			"________",
			"________",
			"________",
		),
		amZug:  spielerSchüler,
		regeln: internationaleRegeln,
		erwarteteZüge: züge{
			zug{
				zugSchritt{
					von:                    position{zeile: 2, spalte: 2},
					zu:                     position{zeile: 0, spalte: 4},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 1, spalte: 3},
					ergebnis: mustParseBrett(
						"____S___",
						"_____l__",
						"________",
						"________",
						"________",
						"________",
						"________",
						"________",
					),
				},
			},
		},
	},

	"internationale Regeln: Dame bewegen": {
		ausgangsSituation: mustParseBrett(
			"s_______",
			"_______s",
			"________",
			"________",
			"____L___",
			"________",
			"________",
			"_s_____S",
		),
		amZug:  spielerLehrer,
		regeln: internationaleRegeln,
		erwarteteZüge: züge{
			// Nach unten rechts
			zug{
				zugSchritt{
					von:                    position{zeile: 4, spalte: 4},
					zu:                     position{zeile: 5, spalte: 5},
					hatGeschlagenePosition: false,
					ergebnis: mustParseBrett(
						"s_______",
						"_______s",
						"________",
						"________",
						"________",
						"_____L__",
						"________",
						"_s_____S",
					),
				},
			},
			zug{
				zugSchritt{
					von:                    position{zeile: 4, spalte: 4},
					zu:                     position{zeile: 6, spalte: 6},
					hatGeschlagenePosition: false,
					ergebnis: mustParseBrett(
						"s_______",
						"_______s",
						"________",
						"________",
						"________",
						"________",
						"______L_",
						"_s_____S",
					),
				},
			},

			// Nach unten links
			zug{
				zugSchritt{
					von:                    position{zeile: 4, spalte: 4},
					zu:                     position{zeile: 5, spalte: 3},
					hatGeschlagenePosition: false,
					ergebnis: mustParseBrett(
						"s_______",
						"_______s",
						"________",
						"________",
						"________",
						"___L____",
						"________",
						"_s_____S",
					),
				},
			},
			zug{
				zugSchritt{
					von:                    position{zeile: 4, spalte: 4},
					zu:                     position{zeile: 6, spalte: 2},
					hatGeschlagenePosition: false,
					ergebnis: mustParseBrett(
						"s_______",
						"_______s",
						"________",
						"________",
						"________",
						"________",
						"__L_____",
						"_s_____S",
					),
				},
			},

			// Nach oben rechts
			zug{
				zugSchritt{
					von:                    position{zeile: 4, spalte: 4},
					zu:                     position{zeile: 3, spalte: 5},
					hatGeschlagenePosition: false,
					ergebnis: mustParseBrett(
						"s_______",
						"_______s",
						"________",
						"_____L__",
						"________",
						"________",
						"________",
						"_s_____S",
					),
				},
			},
			zug{
				zugSchritt{
					von:                    position{zeile: 4, spalte: 4},
					zu:                     position{zeile: 2, spalte: 6},
					hatGeschlagenePosition: false,
					ergebnis: mustParseBrett(
						"s_______",
						"_______s",
						"______L_",
						"________",
						"________",
						"________",
						"________",
						"_s_____S",
					),
				},
			},

			// Nach oben links
			zug{
				zugSchritt{
					von:                    position{zeile: 4, spalte: 4},
					zu:                     position{zeile: 3, spalte: 3},
					hatGeschlagenePosition: false,
					ergebnis: mustParseBrett(
						"s_______",
						"_______s",
						"________",
						"___L____",
						"________",
						"________",
						"________",
						"_s_____S",
					),
				},
			},
			zug{
				zugSchritt{
					von:                    position{zeile: 4, spalte: 4},
					zu:                     position{zeile: 2, spalte: 2},
					hatGeschlagenePosition: false,
					ergebnis: mustParseBrett(
						"s_______",
						"_______s",
						"__L_____",
						"________",
						"________",
						"________",
						"________",
						"_s_____S",
					),
				},
			},
			zug{
				zugSchritt{
					von:                    position{zeile: 4, spalte: 4},
					zu:                     position{zeile: 1, spalte: 1},
					hatGeschlagenePosition: false,
					ergebnis: mustParseBrett(
						"s_______",
						"_L_____s",
						"________",
						"________",
						"________",
						"________",
						"________",
						"_s_____S",
					),
				},
			},
		},
	},

	"internationale Regeln: mit Dame schlagen": {
		ausgangsSituation: mustParseBrett(
			"________",
			"_____s__",
			"________",
			"___L____",
			"________",
			"________",
			"________",
			"________",
		),
		amZug:  spielerLehrer,
		regeln: internationaleRegeln,
		erwarteteZüge: züge{
			zug{
				zugSchritt{
					von:                    position{zeile: 3, spalte: 3},
					zu:                     position{zeile: 0, spalte: 6},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 1, spalte: 5},
					ergebnis: mustParseBrett(
						"______L_",
						"________",
						"________",
						"________",
						"________",
						"________",
						"________",
						"________",
					),
				},
			},
		},
	},

	"internationale Regeln: mit Dame mehrmals schlagen": {
		ausgangsSituation: mustParseBrett(
			"________",
			"___s____",
			"________",
			"________",
			"____S___",
			"________",
			"________",
			"_L______",
		),
		amZug:  spielerLehrer,
		regeln: internationaleRegeln,
		erwarteteZüge: züge{
			zug{
				zugSchritt{
					von:                    position{zeile: 7, spalte: 1},
					zu:                     position{zeile: 3, spalte: 5},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 4, spalte: 4},
					ergebnis: mustParseBrett(
						"________",
						"___s____",
						"________",
						"_____L__",
						"________",
						"________",
						"________",
						"________",
					),
				},
				zugSchritt{
					von:                    position{zeile: 3, spalte: 5},
					zu:                     position{zeile: 0, spalte: 2},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 1, spalte: 3},
					ergebnis: mustParseBrett(
						"__L_____",
						"________",
						"________",
						"________",
						"________",
						"________",
						"________",
						"________",
					),
				},
			},
		},
	},

	"internationale Regeln: mit Dame mehrmals in einer Richtung schlagen": {
		ausgangsSituation: mustParseBrett(
			"L_______",
			"________",
			"__s_____",
			"________",
			"________",
			"_____s__",
			"________",
			"________",
		),
		amZug:  spielerLehrer,
		regeln: internationaleRegeln,
		erwarteteZüge: züge{
			zug{
				zugSchritt{
					von:                    position{zeile: 0, spalte: 0},
					zu:                     position{zeile: 3, spalte: 3},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 2, spalte: 2},
					ergebnis: mustParseBrett(
						"________",
						"________",
						"________",
						"___L____",
						"________",
						"_____s__",
						"________",
						"________",
					),
				},
				zugSchritt{
					von:                    position{zeile: 3, spalte: 3},
					zu:                     position{zeile: 6, spalte: 6},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 5, spalte: 5},
					ergebnis: mustParseBrett(
						"________",
						"________",
						"________",
						"________",
						"________",
						"________",
						"______L_",
						"________",
					),
				},
			},
		},
	},

	"internationale Regeln: mit Dame mehrmals schlagen mit verschiedenen Möglichkeiten": {
		ausgangsSituation: mustParseBrett(
			"l_______",
			"________",
			"______s_",
			"_s___S__",
			"________",
			"________",
			"__s_S___",
			"_L______",
		),
		amZug:  spielerLehrer,
		regeln: internationaleRegeln,
		erwarteteZüge: züge{
			zug{
				zugSchritt{
					von:                    position{zeile: 7, spalte: 1},
					zu:                     position{zeile: 5, spalte: 3},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 6, spalte: 2},
					ergebnis: mustParseBrett(
						"l_______",
						"________",
						"______s_",
						"_s___S__",
						"________",
						"___L____",
						"____S___",
						"________",
					),
				},
				zugSchritt{
					von:                    position{zeile: 5, spalte: 3},
					zu:                     position{zeile: 7, spalte: 5},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 6, spalte: 4},
					ergebnis: mustParseBrett(
						"l_______",
						"________",
						"______s_",
						"_s___S__",
						"________",
						"________",
						"________",
						"_____L__",
					),
				},
				zugSchritt{
					von:                    position{zeile: 7, spalte: 5},
					zu:                     position{zeile: 2, spalte: 0},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 3, spalte: 1},
					ergebnis: mustParseBrett(
						"l_______",
						"________",
						"L_____s_",
						"_____S__",
						"________",
						"________",
						"________",
						"________",
					),
				},
			},
			zug{
				zugSchritt{
					von:                    position{zeile: 7, spalte: 1},
					zu:                     position{zeile: 5, spalte: 3},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 6, spalte: 2},
					ergebnis: mustParseBrett(
						"l_______",
						"________",
						"______s_",
						"_s___S__",
						"________",
						"___L____",
						"____S___",
						"________",
					),
				},
				zugSchritt{
					von:                    position{zeile: 5, spalte: 3},
					zu:                     position{zeile: 2, spalte: 0},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 3, spalte: 1},
					ergebnis: mustParseBrett(
						"l_______",
						"________",
						"L_____s_",
						"_____S__",
						"________",
						"________",
						"____S___",
						"________",
					),
				},
				zugSchritt{
					von:                    position{zeile: 2, spalte: 0},
					zu:                     position{zeile: 7, spalte: 5},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 6, spalte: 4},
					ergebnis: mustParseBrett(
						"l_______",
						"________",
						"______s_",
						"_____S__",
						"________",
						"________",
						"________",
						"_____L__",
					),
				},
			},
		},
	},

	"altdeutsche Regeln: mit Stein schlagen": {
		ausgangsSituation: mustParseBrett(
			"________",
			"________",
			"________",
			"________",
			"___l____",
			"_l______",
			"l_______",
			"s_______",
		),
		amZug:  spielerSchüler,
		regeln: altdeutscheRegeln,
		erwarteteZüge: züge{
			zug{
				zugSchritt{
					von:                    position{zeile: 7, spalte: 0},
					zu:                     position{zeile: 5, spalte: 0},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 6, spalte: 0},
					ergebnis: mustParseBrett(
						"________",
						"________",
						"________",
						"________",
						"___l____",
						"sl______",
						"________",
						"________",
					),
				},
				zugSchritt{
					von:                    position{zeile: 5, spalte: 0},
					zu:                     position{zeile: 5, spalte: 2},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 5, spalte: 1},
					ergebnis: mustParseBrett(
						"________",
						"________",
						"________",
						"________",
						"___l____",
						"__s_____",
						"________",
						"________",
					),
				},
				zugSchritt{
					von:                    position{zeile: 5, spalte: 2},
					zu:                     position{zeile: 3, spalte: 4},
					hatGeschlagenePosition: true,
					geschlagenePosition:    position{zeile: 4, spalte: 3},
					ergebnis: mustParseBrett(
						"________",
						"________",
						"________",
						"____s___",
						"________",
						"________",
						"________",
						"________",
					),
				},
			},
		},
	},

	"kein Zug möglich wenn gewonnen": {
		ausgangsSituation: mustParseBrett(
			"________",
			"________",
			"________",
			"________",
			"________",
			"________",
			"________",
			"s_s_s_s_",
		),
		amZug:         spielerSchüler,
		regeln:        internationaleRegeln,
		erwarteteZüge: nil,
	},

	"4x3 Brett": {
		ausgangsSituation: mustParseBrett(
			"l___",
			"____",
			"s___",
		),
		amZug:  spielerSchüler,
		regeln: internationaleRegeln,
		erwarteteZüge: züge{
			zug{
				zugSchritt{
					von:                    position{zeile: 2, spalte: 0},
					zu:                     position{zeile: 1, spalte: 1},
					hatGeschlagenePosition: false,
					ergebnis: mustParseBrett(
						"l___",
						"_s__",
						"____",
					),
				},
			},
		},
	},
}

func TestZüge(t *testing.T) {
	for name, zugTest := range zügeTests {
		erhalteneZüge := zugTest.ausgangsSituation.möglicheZüge(zugTest.amZug, zugTest.regeln, true)
		if !erhalteneZüge.equals(zugTest.erwarteteZüge) {
			t.Errorf(`Zugtest fehlgeschlagen: %v

Erhalten:
%v

Erwartet:
%v`, name, erhalteneZüge, zugTest.erwarteteZüge)
		}
	}
}
