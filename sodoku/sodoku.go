package sodoku

import (
	"image/color"
	"math/rand"
	"strconv"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type position struct {
	zeile, spalte int
}

func (p position) nächste() (position, bool) {
	if p.zeile == 8 && p.spalte == 8 {
		return position{}, false
	}
	if p.spalte == 8 {
		return position{spalte: 0, zeile: p.zeile + 1}, true
	}
	return position{spalte: p.spalte + 1, zeile: p.zeile}, true
}

type sodoku [9][9]byte

func (s sodoku) möglicheZahlen() [9][9]map[byte]struct{} {
	var möglicheZahlen [9][9]map[byte]struct{}
	for zeile := 0; zeile < 9; zeile++ {
		for spalte := 0; spalte < 9; spalte++ {
			if s[zeile][spalte] != 0 {
				continue
			}

			möglicheZahlen[zeile][spalte] = make(map[byte]struct{}, 9)
			for zahl := 1; zahl <= 9; zahl++ {
				möglicheZahlen[zeile][spalte][byte(zahl)] = struct{}{}
			}
		}
	}

	for startZeile := 0; startZeile < 9; startZeile++ {
		for startSpalte := 0; startSpalte < 9; startSpalte++ {
			zahl := s[startZeile][startSpalte]
			if zahl == 0 {
				continue
			}

			for spalte := 0; spalte < 9; spalte++ {
				if s[startZeile][spalte] != 0 {
					continue
				}
				delete(möglicheZahlen[startZeile][spalte], zahl)
			}

			for zeile := 0; zeile < 9; zeile++ {
				if s[zeile][startSpalte] != 0 {
					continue
				}
				delete(möglicheZahlen[zeile][startSpalte], zahl)
			}

			quadratX := (startSpalte / 3) * 3
			quadratY := (startZeile / 3) * 3
			for zeile := quadratY; zeile < quadratY+3; zeile++ {
				for spalte := quadratX; spalte < quadratX+3; spalte++ {
					if s[zeile][spalte] != 0 {
						continue
					}
					delete(möglicheZahlen[zeile][spalte], zahl)
				}
			}
		}
	}

	return möglicheZahlen
}

func (s sodoku) teilweiseVereinfachen() (sodoku, bool) {
	möglicheZahlen := s.möglicheZahlen()

	var vereinfacht bool
	for zeile := 0; zeile < 9; zeile++ {
		for spalte := 0; spalte < 9; spalte++ {
			if s[zeile][spalte] != 0 {
				continue
			}

			if len(möglicheZahlen[zeile][spalte]) == 1 {
				for zahl := range möglicheZahlen[zeile][spalte] {
					s[zeile][spalte] = zahl
					vereinfacht = true
				}
			}
		}
	}

	return s, vereinfacht
}

func (s sodoku) vereinfachen() sodoku {
	for {
		var vereinfacht bool
		if s, vereinfacht = s.teilweiseVereinfachen(); !vereinfacht {
			return s
		}
	}
}

func (s sodoku) hatFehlerInSpalten() bool {
	zahlen := make(map[byte]struct{}, 9)
	for spalte := 0; spalte < 9; spalte++ {
		for zeile := 0; zeile < 9; zeile++ {
			zahl := s[zeile][spalte]
			if zahl == 0 {
				continue
			}
			if _, doppelt := zahlen[zahl]; doppelt {
				return true
			}
			zahlen[zahl] = struct{}{}
		}
		clear(zahlen)
	}
	return false
}

func (s sodoku) hatFehlerInZeilen() bool {
	zahlen := make(map[byte]struct{}, 9)
	for zeile := 0; zeile < 9; zeile++ {
		for spalte := 0; spalte < 9; spalte++ {
			zahl := s[zeile][spalte]
			if zahl == 0 {
				continue
			}
			if _, doppelt := zahlen[zahl]; doppelt {
				return true
			}
			zahlen[zahl] = struct{}{}
		}
		clear(zahlen)
	}
	return false
}

func (s sodoku) hatFelherInQuadrat() bool {
	zahlen := make(map[byte]struct{}, 9)
	for quadratX := 0; quadratX < 3; quadratX++ {
		for quadratY := 0; quadratY < 3; quadratY++ {
			for spalte := quadratX * 3; spalte < quadratX*3+3; spalte++ {
				for zeile := quadratY * 3; zeile < quadratY*3+3; zeile++ {
					zahl := s[zeile][spalte]
					if zahl == 0 {
						continue
					}
					if _, doppelt := zahlen[zahl]; doppelt {
						return true
					}
					zahlen[zahl] = struct{}{}
				}
			}
			clear(zahlen)
		}
	}
	return false
}

func (s sodoku) hatFehler() bool {
	return s.hatFehlerInZeilen() || s.hatFehlerInSpalten() || s.hatFelherInQuadrat()
}

func (s sodoku) nächstesFreiesFeld(pos position) (position, bool) {
	for s[pos.zeile][pos.spalte] != 0 {
		var ok bool
		pos, ok = pos.nächste()
		if !ok {
			return position{}, false
		}
	}
	return pos, true
}

func (s sodoku) rekursivLösen(pos position, zufall, alle bool) []sodoku {
	zahlen := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if zufall {
		rand.Shuffle(len(zahlen), func(i, j int) {
			zahlen[i], zahlen[j] = zahlen[j], zahlen[i]
		})
	}

	var lösungen []sodoku
	for _, zahl := range zahlen {
		s[pos.zeile][pos.spalte] = byte(zahl)
		if s.hatFehler() {
			continue
		}
		nächstesFeld, ok := pos.nächste()
		if !ok {
			if alle {
				lösungen = append(lösungen, s)
				continue
			} else {
				return []sodoku{s}
			}
		}
		nächstesFreiesFeld, ok := s.nächstesFreiesFeld(nächstesFeld)
		if !ok {
			if alle {
				lösungen = append(lösungen, s)
				continue
			} else {
				return []sodoku{s}
			}
		}
		lösungen = append(lösungen, s.rekursivLösen(nächstesFreiesFeld, zufall, alle)...)
		if !alle && len(lösungen) == 1 {
			return lösungen
		}
	}
	return lösungen
}

func (s sodoku) lösen(zufall, alle bool) []sodoku {
	ersteFreiePosition, ok := s.nächstesFreiesFeld(position{0, 0})
	if !ok {
		if !s.hatFehler() {
			return []sodoku{s}
		}
		return nil
	}
	return s.rekursivLösen(ersteFreiePosition, zufall, alle)
}

func sodokuGenerieren() sodoku {
	s := sodoku{}.lösen(true, false)[0]
	for {
		legalesSodoku := s
		s[rand.Intn(9)][rand.Intn(9)] = 0
		if len(s.lösen(false, true)) != 1 {
			return legalesSodoku
		}
	}
}

type sodokuDrawer struct {
	x, y, width, height float64
	zahlen              [9][9]*ui.Text
}

func newSodokuDrawer(x, y, width, height float64) *sodokuDrawer {
	var zahlen [9][9]*ui.Text
	for zeile := 0; zeile < 9; zeile++ {
		for spalte := 0; spalte < 9; spalte++ {
			zahlen[zeile][spalte] = ui.NewText(ui.TextConfig{})
		}
	}
	return &sodokuDrawer{
		x: x, y: y, width: width, height: height,
		zahlen: zahlen,
	}
}

func (s *sodokuDrawer) draw(screen *ebiten.Image) {
	size, feldSize, abstandX, abstandY := s.calculatePosition()

	for zeile := 0; zeile <= 9; zeile++ {
		strokeWidth := 2
		if zeile%3 == 0 {
			strokeWidth = 4
		}
		vector.StrokeLine(
			screen,
			float32(s.x+abstandX), float32(s.y+abstandY+float64(zeile)*feldSize),
			float32(s.x+abstandX+size), float32(s.y+abstandY+float64(zeile)*feldSize),
			float32(strokeWidth),
			color.Black,
			true,
		)
	}

	for spalte := 0; spalte <= 9; spalte++ {
		strokeWidth := 2
		if spalte%3 == 0 {
			strokeWidth = 4
		}
		vector.StrokeLine(
			screen,
			float32(s.x+abstandX+float64(spalte)*feldSize), float32(s.y+abstandY),
			float32(s.x+abstandX+float64(spalte)*feldSize), float32(s.y+abstandY+size),
			float32(strokeWidth),
			color.Black,
			true,
		)
	}

	for zeile := 0; zeile < 9; zeile++ {
		for spalte := 0; spalte < 9; spalte++ {
			text := s.zahlen[zeile][spalte]
			text.Draw(screen)
		}
	}
}

func (s *sodokuDrawer) update(sodoku sodoku) {
	_, feldSize, abstandX, abstandY := s.calculatePosition()

	for zeile := 0; zeile < 9; zeile++ {
		for spalte := 0; spalte < 9; spalte++ {
			text := s.zahlen[zeile][spalte]
			if sodoku[zeile][spalte] == 0 {
				text.SetText("")
			} else {
				text.SetText(strconv.Itoa(int(sodoku[zeile][spalte])))
			}
			text.SetPosition(ui.NewCenteredPosition(
				s.x+abstandX+float64(spalte)*feldSize+feldSize/2,
				s.y+abstandY+float64(zeile)*feldSize+feldSize/2,
			))
			s.zahlen[zeile][spalte].Update()
		}
	}
}

func (s *sodokuDrawer) calculatePosition() (size, feldSize, abstandX, abstandY float64) {
	size = min(s.width, s.height)
	feldSize = size / 9
	abstandX = (s.width - size) / 2
	abstandY = (s.height - size) / 2
	return size, feldSize, abstandX, abstandY
}

func (s *sodokuDrawer) mausZuKoordinaten(mouseX, mouseY float64) (position, bool) {
	_, feldSize, abstandX, abstandY := s.calculatePosition()
	pos := position{
		zeile:  int((mouseY - (s.y + abstandY)) / feldSize),
		spalte: int((mouseX - (s.x + abstandX)) / feldSize),
	}
	if pos.zeile < 0 || pos.spalte < 0 || pos.zeile >= 9 || pos.spalte >= 9 {
		return position{}, false
	}
	return pos, true
}

type sodokuScreen struct {
	herderLegacy herderlegacy.HerderLegacy

	drawer *sodokuDrawer
	sodoku sodoku
}

func NewSodokuScreen(herderLegacy herderlegacy.HerderLegacy) herderlegacy.Screen {
	return &sodokuScreen{
		herderLegacy: herderLegacy,
		drawer:       newSodokuDrawer(0, 0, ui.Width, ui.Height),
		sodoku:       sodokuGenerieren(),
	}
}

var _ herderlegacy.Screen = (*sodokuScreen)(nil)

func (s *sodokuScreen) onClick(mausX, mausY int) {
	pos, ok := s.drawer.mausZuKoordinaten(float64(mausX), float64(mausY))
	if !ok {
		return
	}

	widgets := make([]ui.ListScreenWidget, 9)
	for i := range widgets {
		i := i
		widgets[i] = ui.ListScreenButtonWidget{
			Text: strconv.Itoa(i + 1),
			Callback: func() {
				s.herderLegacy.OpenScreen(s)
				neuesSodoku := s.sodoku
				neuesSodoku[pos.zeile][pos.spalte] = byte(i + 1)
				if !neuesSodoku.hatFehler() {
					s.sodoku = neuesSodoku
				}
			},
		}
	}
	s.herderLegacy.OpenScreen(ui.NewListScreen(s.herderLegacy, ui.ListScreenConfig{
		Title:   "Zahl setzen",
		Widgets: widgets,
		CancelAction: func() herderlegacy.Screen {
			return s
		},
	}))
}

func (s *sodokuScreen) Update() {
	s.drawer.update(s.sodoku)
	if inpututil.IsKeyJustReleased(ebiten.KeyW) && ebiten.IsKeyPressed(ebiten.KeyU) {
		s.sodoku = s.sodoku.lösen(false, false)[0]
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		s.onClick(ebiten.CursorPosition())
	}
}

func (s *sodokuScreen) Draw(screen *ebiten.Image) {
	screen.Fill(ui.BackgroundColor)
	s.drawer.draw(screen)
}
