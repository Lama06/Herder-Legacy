package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type HorizontalerAnchor float64

const (
	HorizontalerAnchorLinks  HorizontalerAnchor = 0
	HorizontalerAnchorMitte  HorizontalerAnchor = 0.5
	HorizontalerAnchorRechts HorizontalerAnchor = 1
)

type VertikalerAnchor float64

const (
	VertikalerAnchorOben  VertikalerAnchor = 0
	VertikalerAnchorMitte VertikalerAnchor = 0.5
	VertikalerAnchorUnten VertikalerAnchor = 1
)

type Position struct {
	X, Y             float64
	AnchorHorizontal HorizontalerAnchor
	AnchorVertikal   VertikalerAnchor
}

func NewCenteredPosition(x, y float64) Position {
	return Position{
		X:                x,
		Y:                y,
		AnchorHorizontal: HorizontalerAnchorMitte,
		AnchorVertikal:   VertikalerAnchorMitte,
	}
}

func (p Position) eckeObenLinks(width, height float64) (x, y float64) {
	return p.X - float64(p.AnchorHorizontal)*width, p.Y - float64(p.AnchorVertikal)*height
}

func (p Position) isInside(width, height, x, y float64) bool {
	eckeObenLinksX, eckeObenLinksY := p.eckeObenLinks(width, height)

	if x < eckeObenLinksX || x > eckeObenLinksX+width {
		return false
	}

	if y < eckeObenLinksY || y > eckeObenLinksY+height {
		return false
	}

	return true
}

func (p Position) isHovered(width, height float64) bool {
	for _, touchId := range ebiten.AppendTouchIDs(nil) {
		touchX, touchY := ebiten.TouchPosition(touchId)
		if p.isInside(width, height, float64(touchX), float64(touchY)) {
			return true
		}
	}
	mausX, mausY := ebiten.CursorPosition()
	return p.isInside(width, height, float64(mausX), float64(mausY))
}

func (p Position) isClicked(width, height float64) bool {
	for _, touchId := range inpututil.AppendJustReleasedTouchIDs(nil) {
		touchX, touchY := inpututil.TouchPositionInPreviousTick(touchId)
		if p.isInside(width, height, float64(touchX), float64(touchY)) {
			return true
		}
	}
	mausX, mausY := ebiten.CursorPosition()
	return p.isInside(width, height, float64(mausX), float64(mausY)) &&
		inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}
