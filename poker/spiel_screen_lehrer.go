package poker

import (
	"math"
	"math/rand"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
)

const spielScreenLinkerGegnerKartenX = spielScreenKarteHöhe / 2

func spielScreenLehrerKartenPosition(
	position spielScreenLehrerPosition,
	index int,
) (
	karteX float64,
	karteY float64,
	karteRotation float64,
) {
	const (
		rotation        = math.Pi / 2
		abstandY        = spielScreenKarteBreite / 2
		höheGesamt      = abstandY + spielScreenKarteBreite
		abstandVomRandY = (ui.Height - höheGesamt) / 2

		rechterLehrerX = ui.Width + spielScreenKarteHöhe/2
	)

	switch position {
	case spielScreenLehrerPositionLinks:
		karteX = spielScreenLinkerGegnerKartenX
	case spielScreenLehrerPositionRechts:
		karteX = rechterLehrerX
	}

	karteY = abstandVomRandY + float64(index)*abstandY

	return karteX, karteY, rotation
}

type spielScreenLehrerPosition int

const (
	spielScreenLehrerPositionLinks spielScreenLehrerPosition = iota
	spielScreenLehrerPositionRechts
)

type spielScreenLehrer struct {
	position        spielScreenLehrerPosition
	aufgegeben      bool
	eigenerEinsatz  int
	karten          [2]karte
	bewegendeKarten [2]*spielScreenBewegendeKarte
}

func newSpielScreenLehrer(position spielScreenLehrerPosition) *spielScreenLehrer {
	return &spielScreenLehrer{
		position:       position,
		aufgegeben:     false,
		eigenerEinsatz: 0,
	}
}

func (s *spielScreenLehrer) kartenZiehen(karten [2]karte) {
	s.karten = karten
	for i := 0; i < 2; i++ {
		x, y, rotation := spielScreenLehrerKartenPosition(s.position, i)
		s.bewegendeKarten[i] = &spielScreenBewegendeKarte{
			karte:          karten[i],
			targetRotation: rotation,
			currentX:       spielScreenStapelX,
			currentY:       spielScreenStapelY,
			targetX:        x,
			targetY:        y,
		}
	}
}

func (s *spielScreenLehrer) getKarten() [2]karte {
	return s.karten
}

func (s *spielScreenLehrer) getBewegendeKarten() [2]*spielScreenBewegendeKarte {
	return s.bewegendeKarten
}

func (s *spielScreenLehrer) jettonSpawnPunkt() (float64, float64) {
	offset := rand.Float64() * 200
	switch s.position {
	case spielScreenLehrerPositionLinks:
		return -spielScreenJettonSize - offset, ui.Height / 2
	case spielScreenLehrerPositionRechts:
		return ui.Width + offset, ui.Height / 2
	default:
		panic("unreachable")
	}
}

func (s *spielScreenLehrer) setAufgegeben(aufgegeben bool) {
	s.aufgegeben = aufgegeben
}

func (s *spielScreenLehrer) hatAufgegeben() bool {
	return s.aufgegeben
}

func (s *spielScreenLehrer) einsatzErmitteln(
	herderLegacy herderlegacy.HerderLegacy,
	status spielScreenStatus,
	wirdGewinnen bool,
	callback func(einsatz int),
) {
	var einsatz int
	strategieZufallszahl := rand.Float64()
	switch {
	case status == spielStatusVerdeckteMittelkarten:
		einsatz = 1 // Einer geht immer :)
	case (strategieZufallszahl < 0.25) || (strategieZufallszahl >= 0.5 && !wirdGewinnen):
		einsatz = 1
	case strategieZufallszahl >= 0.5 && wirdGewinnen:
		switch status {
		case spielStatus3AufgedeckteMittelKarten:
			einsatz = 1 + rand.Intn(2)
		case spielStatus4AufgedeckteMittelkarten:
			einsatz = 2 + rand.Intn(3)
		case spielStatus5AufgedeckteMittelkarten:
			einsatz = 4 + rand.Intn(3)
		default:
			panic("unbekannter Spielstatus")
		}
	case strategieZufallszahl >= 0.25 && strategieZufallszahl < 0.5:
		einsatz = 8 + rand.Intn(7)
	default:
		panic("unreachable")
	}

	s.eigenerEinsatz += einsatz
	callback(einsatz)
}

func (s *spielScreenLehrer) gehtMit(
	herderLegacy herderlegacy.HerderLegacy,
	einsatz int,
	callback func(gehtMit bool),
) {
	schmerzgrenze := maxInt(3, s.eigenerEinsatz)
	entscheidung := schmerzgrenze >= einsatz
	if entscheidung {
		s.eigenerEinsatz += einsatz
	}
	callback(entscheidung)
}
