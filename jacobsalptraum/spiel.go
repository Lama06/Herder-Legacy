package jacobsalptraum

import (
	"image/color"

	"github.com/Lama06/Herder-Legacy/minimax"
	"golang.org/x/image/colornames"
)

type Spieler bool

const (
	SpielerSchach      Spieler = true
	SpielerVierGewinnt Spieler = false
)

func (s Spieler) gegner() Spieler {
	return !s
}

func (s Spieler) MinimaxGegner() minimax.Spieler {
	return s.gegner()
}

type position struct {
	zeile, spalte int
}

func (p position) farbe() color.Color {
	if p.zeile%2 == p.spalte%2 {
		return colornames.Chocolate
	}
	return colornames.Beige
}

func (p position) valid(breite, höhe int) bool {
	return p.zeile >= 0 && p.zeile < höhe && p.spalte >= 0 && p.spalte < breite
}

type Regeln struct {
	MaximaleVierGewinntSteine int
}

var _ minimax.Regeln = Regeln{}

var StandardRegeln = Regeln{MaximaleVierGewinntSteine: 15}
