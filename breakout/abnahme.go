package breakout

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func (w *world) abnahmeModusPasswortLesen() {
	const passwort = "Abnahme"

	w.abnahmeModusPasswort = ebiten.AppendInputChars(w.abnahmeModusPasswort)
	if len(w.abnahmeModusPasswort) > len(passwort) {
		w.abnahmeModusPasswort = w.abnahmeModusPasswort[len(w.abnahmeModusPasswort)-len(passwort) : len(w.abnahmeModusPasswort)]
	}

	if string(w.abnahmeModusPasswort) != passwort {
		return
	}

	w.automaticInputUpgradeRemainingTime = math.MaxInt
}
