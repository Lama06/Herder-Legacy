package dame

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image/color"
)

type feld byte

const (
	feldLeer feld = iota
	feldSteinLehrer
	feldSteinSchueler
	feldDameLehrer
	feldDameSchueler
)

func stein(von spieler) feld {
	switch von {
	case spielerLehrer:
		return feldSteinLehrer
	case spielerSchueler:
		return feldSteinSchueler
	default:
		panic("unreachable")
	}
}

func dame(von spieler) feld {
	switch von {
	case spielerLehrer:
		return feldDameLehrer
	case spielerSchueler:
		return feldDameSchueler
	default:
		panic("unreachable")
	}
}

func parseFeld(zeichen rune) (feld, error) {
	switch zeichen {
	case feldLeer.zeichen():
		return feldLeer, nil
	case feldSteinLehrer.zeichen():
		return feldSteinLehrer, nil
	case feldSteinSchueler.zeichen():
		return feldSteinSchueler, nil
	case feldDameLehrer.zeichen():
		return feldDameLehrer, nil
	case feldDameSchueler.zeichen():
		return feldDameSchueler, nil
	default:
		return 0, fmt.Errorf("ungültiges Zeichen für Feld: %c", zeichen)
	}
}

func (f feld) zeichen() rune {
	switch f {
	case feldLeer:
		return '_'
	case feldSteinLehrer:
		return 'l'
	case feldSteinSchueler:
		return 's'
	case feldDameLehrer:
		return 'L'
	case feldDameSchueler:
		return 'S'
	default:
		panic("unreachable")
	}
}

func (f feld) eigentuemer() (eigentuemer spieler, ok bool) {
	switch f {
	case feldLeer:
		return spielerLehrer, false
	case feldSteinLehrer, feldDameLehrer:
		return spielerLehrer, true
	case feldSteinSchueler, feldDameSchueler:
		return spielerSchueler, true
	default:
		panic("unreachable")
	}
}

func (f feld) farbe() (clr color.Color, ok bool) {
	switch f {
	case feldLeer:
		return nil, false
	case feldSteinLehrer:
		return colornames.Lightblue, true
	case feldDameLehrer:
		return colornames.Darkblue, true
	case feldSteinSchueler:
		return colornames.Orange, true
	case feldDameSchueler:
		return colornames.Darkorange, true
	default:
		panic("unreachable")
	}
}

func (f feld) draw(screen *ebiten.Image, x, y, size float64) {
	clr, ok := f.farbe()
	if !ok {
		return
	}
	centerX, centerY := float32(x)+float32(size)/2, float32(y)+float32(size)/2
	radius := float32(size) * 0.35
	vector.DrawFilledCircle(screen, centerX, centerY, radius, clr, true)
}
