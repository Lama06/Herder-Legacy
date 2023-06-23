// Eine Implementierung des Minimax Algorithmuses: https://de.wikipedia.org/wiki/Minimax-Algorithmus
package ai

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Spieler interface {
	Gegner() Spieler
}

type Zug interface {
	Ergebnis() Brett
}

type Regeln any

type Brett interface {
	MoeglicheZuege(perspektive Spieler, regeln Regeln) []Zug

	Bewertung(perspektive Spieler, regeln Regeln) int
}

func rekursiveBrettBewertung(
	brett Brett,
	regeln Regeln,
	perspektive Spieler,
	amZug Spieler,
	maximaleTiefe int,
) int {
	if maximaleTiefe < 0 {
		panic("maximaleTiefe < 0")
	}

	if maximaleTiefe == 0 {
		return brett.Bewertung(perspektive, regeln)
	}

	folgendeZuege := brett.MoeglicheZuege(amZug, regeln)
	if len(folgendeZuege) == 0 {
		return brett.Bewertung(perspektive, regeln)
	}

	var besterFolgenderZugBewertung int
	for i, folgenderZug := range folgendeZuege {
		folgenderZugBewerung := rekursiveBrettBewertung(
			folgenderZug.Ergebnis(),
			regeln,
			perspektive,
			amZug.Gegner(),
			maximaleTiefe-1,
		)
		if i == 0 {
			besterFolgenderZugBewertung = folgenderZugBewerung
			continue
		}
		if amZug == perspektive {
			besterFolgenderZugBewertung = max(besterFolgenderZugBewertung, folgenderZugBewerung)
		} else {
			besterFolgenderZugBewertung = min(besterFolgenderZugBewertung, folgenderZugBewerung)
		}
	}
	return besterFolgenderZugBewertung
}

func BesterNaechsterZug(brett Brett, regeln Regeln, amZug Spieler, maximaleTiefe int) (zug Zug, ok bool) {
	if maximaleTiefe <= 0 {
		panic("maximaleTiefe <= 0")
	}

	moeglicheZuege := brett.MoeglicheZuege(amZug, regeln)
	if len(moeglicheZuege) == 0 {
		return nil, false
	}

	var (
		besterZug          Zug
		besterZugBewertung int
	)
	for _, moeglicherZug := range moeglicheZuege {
		zugErgebnisBewertung := rekursiveBrettBewertung(
			moeglicherZug.Ergebnis(),
			regeln,
			amZug,
			amZug.Gegner(),
			maximaleTiefe-1,
		)
		if besterZug == nil || zugErgebnisBewertung > besterZugBewertung {
			besterZug = moeglicherZug
			besterZugBewertung = zugErgebnisBewertung
		}
	}
	return besterZug, true
}
