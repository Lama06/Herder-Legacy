package dame

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

type feld byte

const (
	feldLeer feld = iota
	feldSteinLehrer
	feldSteinSchüler
	feldDameLehrer
	feldDameSchüler
)

func stein(von spieler) feld {
	switch von {
	case spielerLehrer:
		return feldSteinLehrer
	case spielerSchüler:
		return feldSteinSchüler
	default:
		panic("unreachable")
	}
}

func dame(von spieler) feld {
	switch von {
	case spielerLehrer:
		return feldDameLehrer
	case spielerSchüler:
		return feldDameSchüler
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
	case feldSteinSchüler.zeichen():
		return feldSteinSchüler, nil
	case feldDameLehrer.zeichen():
		return feldDameLehrer, nil
	case feldDameSchüler.zeichen():
		return feldDameSchüler, nil
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
	case feldSteinSchüler:
		return 's'
	case feldDameLehrer:
		return 'L'
	case feldDameSchüler:
		return 'S'
	default:
		panic("unreachable")
	}
}

func (f feld) eigentümer() (eigentümer spieler, ok bool) {
	switch f {
	case feldLeer:
		return spielerLehrer, false
	case feldSteinLehrer, feldDameLehrer:
		return spielerLehrer, true
	case feldSteinSchüler, feldDameSchüler:
		return spielerSchüler, true
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
	case feldSteinSchüler:
		return colornames.Orange, true
	case feldDameSchüler:
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
