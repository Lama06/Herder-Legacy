package jacobsalptraum

import (
	"image"

	"github.com/Lama06/Herder-Legacy/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

type feld byte

const (
	feldLeer feld = iota

	feldBauer
	feldSpringer
	feldTurm
	feldLäufer
	feldDame
	feldKönig

	feldVierGewinntStein
)

var feldBuchstaben = map[feld]rune{
	feldLeer:     '_',
	feldBauer:    'B',
	feldSpringer: 'S',
	feldTurm:     'T',
	feldLäufer:   'L',
	feldDame:     'D',
	feldKönig:    'K',
}

func parseFeld(buchstabe rune) (feld, bool) {
	for f := range feldBuchstaben {
		if feldBuchstaben[f] == buchstabe {
			return f, true
		}
	}
	return 0, false
}

func (f feld) image() *ebiten.Image {
	switch f {
	case feldLeer:
		return nil
	case feldBauer, feldSpringer, feldTurm, feldLäufer, feldDame, feldKönig:
		offset := int(f - feldBauer)
		return assets.RequireImage("brettspiel/schach_weiß.png").SubImage(image.Rect(offset*16, 0, (offset+1)*16, 16)).(*ebiten.Image)
	case feldVierGewinntStein:
		return assets.RequireImage("brettspiel/dame.png").SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image)
	default:
		panic("unreachable")
	}
}
