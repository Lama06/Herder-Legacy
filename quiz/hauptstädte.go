package quiz

var hauptstädteDerNachbarländerDeutschlandsFragen = []QuizFrage{
	{
		Frage:            "Hauptstadt von Deutschland",
		Antwort:          "Berlin",
		FalscheAntworten: []string{"Europa"},
	},
	{
		Frage:   "Dänemark",
		Antwort: "Kopenhagen",
	},
	{
		Frage:   "Niederlande",
		Antwort: "Amsterdam",
	},
	{
		Frage:   "Belgien",
		Antwort: "Brüssel",
	},
	{
		Frage:   "Luxemburg",
		Antwort: "Luxemburg",
	},
	{
		Frage:            "Frankreich",
		Antwort:          "Paris",
		FalscheAntworten: []string{"Baguette"},
	},
	{
		Frage:   "Schweiz",
		Antwort: "Bern",
	},
	{
		Frage:   "Österreich",
		Antwort: "Wien",
	},
	{
		Frage:   "Polen",
		Antwort: "Warschau",
	},
	{
		Frage:   "Tschechien",
		Antwort: "Prag",
	},
}

func NewHauptstädtDerNachbarländerDeutschlandsQuizConfig(zeitProFrage int) QuizConfig {
	return QuizConfig{
		Name:         "Hauptstädte der Nachbarländer Deutschlands",
		ZeitProFrage: zeitProFrage,
		Fragen:       hauptstädteDerNachbarländerDeutschlandsFragen,
	}
}

func NewHauptstädteEuropasQuizConfig(zeitProFrage int) QuizConfig {
	return QuizConfig{
		Name:         "Hauptstädte Europas",
		ZeitProFrage: zeitProFrage,
		Fragen: append(
			hauptstädteDerNachbarländerDeutschlandsFragen,
			QuizFrage{
				Frage:   "Spanien",
				Antwort: "Madrid",
			},
			QuizFrage{
				Frage:   "Portugal",
				Antwort: "Lisabon",
			},
			QuizFrage{
				Frage:   "Andorra",
				Antwort: "Andorra la vella",
			},
			QuizFrage{
				Frage:   "Vereinigtes Königreich",
				Antwort: "London",
			},
			QuizFrage{
				Frage:   "Irland",
				Antwort: "Dublin",
			},
			QuizFrage{
				Frage:   "Norwegen",
				Antwort: "Oslow",
			},
			QuizFrage{
				Frage:   "Schweden",
				Antwort: "Stockholm",
			},
			QuizFrage{
				Frage:   "Finnland",
				Antwort: "Helsinki",
			},
			QuizFrage{
				Frage:   "Estland",
				Antwort: "Tallinn",
			},
			QuizFrage{
				Frage:   "Lettland",
				Antwort: "Riga",
			},
			QuizFrage{
				Frage:   "Litauen",
				Antwort: "Vilnius",
			},
			QuizFrage{
				Frage:   "Belarus",
				Antwort: "Minsk",
			},
			QuizFrage{
				Frage:   "Russland",
				Antwort: "Moskau",
			},
			QuizFrage{
				Frage:   "Ukraine",
				Antwort: "Kiew",
			},
			QuizFrage{
				Frage:   "Slowakei",
				Antwort: "Bratislava",
			},
			QuizFrage{
				Frage:   "Ungarn",
				Antwort: "Budapest",
			},
			QuizFrage{
				Frage:   "Slowenien",
				Antwort: "Ljubljana",
			},
			QuizFrage{
				Frage:   "Kroatien",
				Antwort: "Zagreb",
			},
			QuizFrage{
				Frage:   "Bosnien und Herzegowina",
				Antwort: "Sarajevo",
			},
			QuizFrage{
				Frage:   "Serbien",
				Antwort: "Belgrad",
			},
			QuizFrage{
				Frage:   "Rumänien",
				Antwort: "Bukarest",
			},
			QuizFrage{
				Frage:   "Bulgarien",
				Antwort: "Sofia",
			},
			QuizFrage{
				Frage:   "Nordmazedonien",
				Antwort: "Skopje",
			},
			QuizFrage{
				Frage:   "Griechenland",
				Antwort: "Athen",
			},
			QuizFrage{
				Frage:   "Italien",
				Antwort: "Rom",
			},
			QuizFrage{
				Frage:   "Türkei",
				Antwort: "Ankara",
			},
			QuizFrage{
				Frage:   "Albanien",
				Antwort: "Tirana",
			},
		),
	}
}
