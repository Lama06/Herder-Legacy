package dame

import "github.com/Lama06/Herder-Legacy/ai"

type spieler bool

var _ ai.Spieler = spielerLehrer

const (
	spielerLehrer   spieler = false
	spielerSchueler spieler = true
)

func (s spieler) gegner() spieler {
	return !s
}

func (s spieler) Gegner() ai.Spieler {
	return s.gegner()
}
