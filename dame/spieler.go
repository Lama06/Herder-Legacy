package dame

import "github.com/Lama06/Herder-Legacy/minimax"

type spieler bool

var _ minimax.Spieler = spielerLehrer

const (
	spielerLehrer  spieler = false
	spielerSchüler spieler = true
)

func (s spieler) gegner() spieler {
	return !s
}

func (s spieler) MinimaxGegner() minimax.Spieler {
	return s.gegner()
}
