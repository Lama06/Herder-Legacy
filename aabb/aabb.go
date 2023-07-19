package aabb

import "math"

type Aabb struct {
	X, Y, Width, Height float64
}

func (aabb Aabb) CenterX() float64 {
	return aabb.X + aabb.Width/2
}

func (aabb Aabb) CenterY() float64 {
	return aabb.Y + aabb.Height/2
}

func (aabb Aabb) MaxX() float64 {
	return aabb.X + aabb.Width
}

func (aabb Aabb) MaxY() float64 {
	return aabb.Y + aabb.Height
}

func (aabb Aabb) Area() float64 {
	return aabb.Width * aabb.Height
}

func (aabb1 Aabb) KollidiertMit(aabb2 Aabb) bool {
	if aabb1.X >= aabb2.MaxX() {
		return false
	}

	if aabb1.MaxX() <= aabb2.X {
		return false
	}

	if aabb1.Y >= aabb2.MaxY() {
		return false
	}

	if aabb1.MaxY() <= aabb2.Y {
		return false
	}

	return true
}

func (aabb Aabb) IsInside(x, y float64) bool {
	if x < aabb.X || x > aabb.MaxX() {
		return false
	}

	if y < aabb.Y || y > aabb.MaxY() {
		return false
	}

	return true
}

func (aabb1 Aabb) Intersection(aabb2 Aabb) (intersection Aabb, ok bool) {
	if !aabb1.KollidiertMit(aabb2) {
		return Aabb{}, false
	}

	x := math.Max(aabb1.X, aabb2.X)
	y := math.Max(aabb1.Y, aabb2.Y)

	return Aabb{
		X:      x,
		Y:      y,
		Width:  math.Min(aabb1.MaxX(), aabb2.MaxX()) - x,
		Height: math.Min(aabb1.MaxY(), aabb2.MaxY()) - y,
	}, true
}

func (aabb Aabb) VertikalZerschneiden(prozentOben float64) (oben, unten Aabb) {
	oben = Aabb{
		X:      aabb.X,
		Y:      aabb.Y,
		Width:  aabb.Width,
		Height: aabb.Height * prozentOben,
	}
	unten = Aabb{
		X:      aabb.X,
		Y:      aabb.Y + aabb.Height*prozentOben,
		Width:  aabb.Width,
		Height: aabb.Height - aabb.Height*prozentOben,
	}
	return oben, unten
}

func (aabb Aabb) VonObenZerschneiden(prozentOben float64) (oben Aabb) {
	oben, _ = aabb.VertikalZerschneiden(prozentOben)
	return oben
}

func (aabb Aabb) VonUntenZerschneiden(prozentUnten float64) (unten Aabb) {
	_, unten = aabb.VertikalZerschneiden(1 - prozentUnten)
	return unten
}

func (aabb Aabb) HorizontalZerschneiden(prozentLinks float64) (links, rechts Aabb) {
	links = Aabb{
		X:      aabb.X,
		Y:      aabb.Y,
		Width:  aabb.Width * prozentLinks,
		Height: aabb.Height,
	}
	rechts = Aabb{
		X:      aabb.X + aabb.Width*prozentLinks,
		Y:      aabb.Y,
		Width:  aabb.Width - aabb.Width*prozentLinks,
		Height: aabb.Height,
	}
	return links, rechts
}

func (aabb Aabb) VonLinksZerschneiden(prozent float64) (links Aabb) {
	links, _ = aabb.HorizontalZerschneiden(prozent)
	return links
}

func (aabb Aabb) VonRechtsZerschneiden(prozent float64) (rechts Aabb) {
	_, rechts = aabb.HorizontalZerschneiden(1 - prozent)
	return rechts
}
