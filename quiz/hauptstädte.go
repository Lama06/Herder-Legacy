package quiz

func NewHauptstädtDerNachbarländerDeutschlandsQuizConfig(zeitProFrage int) MultipleChoiceQuizConfig {
	return MultipleChoiceQuizConfig{
		Name:         "Hauptstädte der Nachbarländer Deutschlands",
		ZeitProFrage: zeitProFrage,
		Fragen:       hauptstädteDerNachbarländerDeutschlandsFragen,
	}
}

func NewHauptstädteEuropasQuizConfig(zeitProFrage int) MultipleChoiceQuizConfig {
	return MultipleChoiceQuizConfig{
		Name:         "Hauptstädte Europas",
		ZeitProFrage: zeitProFrage,
		Fragen:       hauptstädteEuropasFragen,
	}
}

func NewHauptstädteInternationalQuizConfig(zeitProFrage int) MultipleChoiceQuizConfig {
	return MultipleChoiceQuizConfig{
		Name:         "Internationale Hauptstädte",
		ZeitProFrage: zeitProFrage,
		Fragen:       hauptstädteInternationalFragen,
	}
}

var (
	hauptstädteDerNachbarländerDeutschlandsFragen = []MultipleChoiceQuizFrage{
		{
			Frage:            "Deutschland",
			Antwort:          "Berlin",
			FalscheAntworten: []string{"Europa"},
		},
		{
			Frage:   "Dänemark",
			Antwort: "Kopenhagen",
		},
		{
			Frage:            "Niederlande",
			Antwort:          "Amsterdam",
			FalscheAntworten: []string{"Käse"},
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

	hauptstädteEuropasFragen = append(
		hauptstädteDerNachbarländerDeutschlandsFragen,
		MultipleChoiceQuizFrage{
			Frage:   "Spanien",
			Antwort: "Madrid",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Portugal",
			Antwort: "Lisabon",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Andorra",
			Antwort: "Andorra la vella",
		},
		MultipleChoiceQuizFrage{
			Frage:            "Vereinigtes Königreich",
			Antwort:          "London",
			FalscheAntworten: []string{"Tee"},
		},
		MultipleChoiceQuizFrage{
			Frage:   "Irland",
			Antwort: "Dublin",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Norwegen",
			Antwort: "Oslow",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Schweden",
			Antwort: "Stockholm",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Finnland",
			Antwort: "Helsinki",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Estland",
			Antwort: "Tallinn",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Lettland",
			Antwort: "Riga",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Litauen",
			Antwort: "Vilnius",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Belarus",
			Antwort: "Minsk",
		},
		MultipleChoiceQuizFrage{
			Frage:            "Russland",
			Antwort:          "Moskau",
			FalscheAntworten: []string{"Vodka"},
		},
		MultipleChoiceQuizFrage{
			Frage:   "Ukraine",
			Antwort: "Kiew",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Slowakei",
			Antwort: "Bratislava",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Ungarn",
			Antwort: "Budapest",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Slowenien",
			Antwort: "Ljubljana",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Kroatien",
			Antwort: "Zagreb",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Bosnien und Herzegowina",
			Antwort: "Sarajevo",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Serbien",
			Antwort: "Belgrad",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Rumänien",
			Antwort: "Bukarest",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Bulgarien",
			Antwort: "Sofia",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Nordmazedonien",
			Antwort: "Skopje",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Griechenland",
			Antwort: "Athen",
		},
		MultipleChoiceQuizFrage{
			Frage:            "Italien",
			Antwort:          "Rom",
			FalscheAntworten: []string{"Pizza", "Spghetti"},
		},
		MultipleChoiceQuizFrage{
			Frage:   "Türkei",
			Antwort: "Ankara",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Albanien",
			Antwort: "Tirana",
		},
		MultipleChoiceQuizFrage{
			Frage:            "Lichtenstein",
			Antwort:          "Vaduz",
			FalscheAntworten: []string{"Kavier"},
		},
		MultipleChoiceQuizFrage{
			Frage:   "Monaco",
			Antwort: "Monaco",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Malta",
			Antwort: "Valetta",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Montenegro",
			Antwort: "Podgorica",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Kosovo",
			Antwort: "Pristina",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Moldawien",
			Antwort: "Chișinău",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Island",
			Antwort: "Reykjavik",
		},
	)

	hauptstädteInternationalFragen = append(
		hauptstädteEuropasFragen,
		MultipleChoiceQuizFrage{
			Frage:   "USA",
			Antwort: "Washington",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Kanada",
			Antwort: "Ottawa",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Mexiko",
			Antwort: "Mexiko City",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Nordkorea",
			Antwort: "Pjöngjang",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Südkorea",
			Antwort: "Seoul",
		},
		MultipleChoiceQuizFrage{
			Frage:   "China",
			Antwort: "Peeking",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Japan",
			Antwort: "Tokyo",
		},
		MultipleChoiceQuizFrage{
			Frage:   "Brasilien",
			Antwort: "Brasília",
		},
	)
)
