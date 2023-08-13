package poker

import (
	"image/color"
	"math/rand"

	"github.com/Lama06/Herder-Legacy/assets"
	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/Lama06/Herder-Legacy/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func istKlick() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustReleasedTouchIDs(nil)) != 0
}

type spielScreenStatus int

const (
	spielStatusKartenAnfang spielScreenStatus = iota
	spielStatusKartenWerdenGezogen
	spielStatusVerdeckteMittelkarten
	spielStatusMittelkartenWerdenAufgedeckt
	spielStatus3AufgedeckteMittelKarten
	spielStatusVierteKarteWirdGezogen
	spielStatus4AufgedeckteMittelkarten
	spielStatusFünfteKarteWirdGezogen
	spielStatus5AufgedeckteMittelkarten
	spielStatusKartenWerdenAufgedeckt
	spielStatusSiegerermittlung
)

func (s spielScreenStatus) werSetztZuerst(spieler []spielScreenSpieler) spielScreenSpieler {
	var index int
	switch s {
	case spielStatusVerdeckteMittelkarten:
		index = 0
	case spielStatus3AufgedeckteMittelKarten:
		index = 1
	case spielStatus4AufgedeckteMittelkarten:
		index = 2
	case spielStatus5AufgedeckteMittelkarten:
		index = 3
	}
	return spieler[index%len(spieler)]
}

const (
	spielScreenJettonAblageX      = spielScreenLinkerGegnerKartenX + 10
	spielScreenJettonAblageY      = 10
	spielScreenJettonAblageBreite = spielScreenStapelX - spielScreenJettonAblageX - spielScreenJettonSize - 10
	spielScreenJettonAblageHöhe   = spielScreenMittelKartenY - spielScreenJettonAblageY - spielScreenJettonSize - 10

	spielScreenStapelY = 100.0
	spielScreenStapelX = (ui.Width - spielScreenKarteBreite) / 2.0

	spielScreenAnzahlMittelkarten          = 5
	spielScreenMittelKartenY               = spielScreenStapelY + spielScreenKarteHöhe*1.5
	spielScreenMittelKartenAbstandX        = spielScreenKarteBreite * 1.5
	spielScreenMittelKartenBreiteGesamt    = (spielScreenAnzahlMittelkarten-1)*spielScreenMittelKartenAbstandX + spielScreenKarteBreite
	spielScreenMittelKartenAbstandVomRandX = (ui.Width - spielScreenMittelKartenBreiteGesamt) / 2
)

func spielScreenMittelKarteX(karteIndex int) float64 {
	return spielScreenMittelKartenAbstandVomRandX + float64(karteIndex)*spielScreenMittelKartenAbstandX
}

type spielScreenSpieler interface {
	kartenZiehen(karte [2]karte)

	getKarten() [2]karte

	getBewegendeKarten() [2]*spielScreenBewegendeKarte

	jettonSpawnPunkt() (float64, float64)

	setAufgegeben(bool)

	hatAufgegeben() bool

	einsatzErmitteln(
		herderLegacy herderlegacy.HerderLegacy,
		status spielScreenStatus,
		wirdGewinnen bool,
		callback func(einsatz int),
	)

	gehtMit(
		herderLegacy herderlegacy.HerderLegacy,
		einsatz int,
		callback func(gehtMit bool),
	)
}

type spielScreen struct {
	herderLegacy herderlegacy.HerderLegacy

	status spielScreenStatus
	info   *ui.Title

	spieler []spielScreenSpieler

	stapel kartenStapel

	bewegendeJettons []*spielScreenBewegenderJetton

	mittelkarten          [5]karte
	bewegendeMittelkarten [5]*spielScreenBewegendeKarte
}

var _ herderlegacy.Screen = (*spielScreen)(nil)

func NewSpielScreen(herderLegacy herderlegacy.HerderLegacy) herderlegacy.Screen {
	stapel := vollständigerKartenStapel.clone()

	spieler := []spielScreenSpieler{
		newSpielScreenLehrer(spielScreenLehrerPositionLinks),
		newSpielScreenLehrer(spielScreenLehrerPositionRechts),
		newSpielScreenMensch(30),
	}
	rand.Shuffle(len(spieler), func(i, j int) {
		spieler[i], spieler[j] = spieler[j], spieler[i]
	})

	return &spielScreen{
		herderLegacy: herderLegacy,

		status: spielStatusKartenAnfang,
		info: ui.NewTitle(ui.TitleConfig{
			Position:           ui.NewCenteredPosition(ui.Width/2, ui.Height/3),
			Text:               "Klicken um Karten zu ziehen",
			CustomColorPalette: false,
			ColorPalette:       ui.TitleColorPalette{},
		}),
		spieler:      spieler,
		stapel:       stapel,
		mittelkarten: [5]karte{stapel.karteZiehen(), stapel.karteZiehen(), stapel.karteZiehen(), stapel.karteZiehen(), stapel.karteZiehen()},
	}
}

func (s *spielScreen) bewegendeKarten() []*spielScreenBewegendeKarte {
	bewegendeKarten := s.bewegendeMittelkarten[:]
	for _, spieler := range s.spieler {
		spielerKarten := spieler.getBewegendeKarten()
		bewegendeKarten = append(bewegendeKarten, spielerKarten[:]...)
	}
	return bewegendeKarten
}

func (s *spielScreen) components() []ui.Component {
	var components []ui.Component
	for _, bewegendeKarte := range s.bewegendeKarten() {
		components = append(components, bewegendeKarte)
	}
	for _, bewegenderJetton := range s.bewegendeJettons {
		components = append(components, bewegenderJetton)
	}
	components = append(components, s.info)
	return components
}

func (s *spielScreen) spielerHand(spieler spielScreenSpieler) hand {
	return parseHand([7]karte{
		spieler.getKarten()[0],
		spieler.getKarten()[1],
		s.mittelkarten[0],
		s.mittelkarten[1],
		s.mittelkarten[2],
		s.mittelkarten[3],
		s.mittelkarten[4],
	})
}

func (s *spielScreen) gewinner() spielScreenSpieler {
gewinnerKanidaten:
	for _, gewinnerKanidat := range s.spieler {
		gewinnerKanidatHand := s.spielerHand(gewinnerKanidat)

		for _, gegner := range s.spieler {
			if gewinnerKanidat == gegner {
				continue
			}

			gegnerHand := s.spielerHand(gegner)

			if compareHände(gewinnerKanidatHand, gegnerHand) != 1 {
				continue gewinnerKanidaten
			}
		}

		return gewinnerKanidat
	}

	// Kein eindeutiger Gewinner
	return s.spieler[rand.Intn(len(s.spieler))]
}

func (s *spielScreen) jettonsSetzen(spieler spielScreenSpieler, anzahl int) {
	for i := 0; i < anzahl; i++ {
		s.bewegendeJettons = append(s.bewegendeJettons, newSpielScreenBewegenderJetton(spieler))
	}
}

func (s *spielScreen) einsätzeErmitteln(status spielScreenStatus, callback func()) {
	var möglicheBeginner []spielScreenSpieler
	for _, spieler := range s.spieler {
		if !spieler.hatAufgegeben() {
			möglicheBeginner = append(möglicheBeginner, spieler)
		}
	}
	beginner := status.werSetztZuerst(möglicheBeginner)
	beginner.einsatzErmitteln(
		s.herderLegacy,
		status,
		s.gewinner() == beginner,
		func(einsatz int) {
			s.jettonsSetzen(beginner, einsatz)

			var spielerDieMitgehenKönnten []spielScreenSpieler
			for _, spieler := range s.spieler {
				if !spieler.hatAufgegeben() && spieler != beginner {
					spielerDieMitgehenKönnten = append(spielerDieMitgehenKönnten, spieler)
				}
			}

			var fragebObSpielerMitgehen func(int)
			fragebObSpielerMitgehen = func(spielerIndex int) {
				if spielerIndex == len(spielerDieMitgehenKönnten) {
					callback()
					return
				}

				gefragterSpieler := spielerDieMitgehenKönnten[spielerIndex]
				gefragterSpieler.gehtMit(
					s.herderLegacy,
					einsatz,
					func(gehtMit bool) {
						if !gehtMit {
							for _, karte := range gefragterSpieler.getBewegendeKarten() {
								karte.targetAufgedecktStatus = true
							}
							gefragterSpieler.setAufgegeben(true)
						} else {
							s.jettonsSetzen(gefragterSpieler, einsatz)
						}

						fragebObSpielerMitgehen(spielerIndex + 1)
					},
				)
			}
			fragebObSpielerMitgehen(0)
		},
	)
}

func (s *spielScreen) Update() {
	for _, component := range s.components() {
		component.Update()
	}

	switch s.status {
	case spielStatusKartenAnfang:
		if !istKlick() {
			return
		}

		s.status = spielStatusKartenWerdenGezogen
		s.info.SetText("")
		for _, spieler := range s.spieler {
			spieler.kartenZiehen([2]karte{s.stapel.karteZiehen(), s.stapel.karteZiehen()})
		}
		for i := 0; i < 3; i++ {
			s.bewegendeMittelkarten[i] = &spielScreenBewegendeKarte{
				currentX: spielScreenStapelX,
				currentY: spielScreenStapelY,
				targetX:  spielScreenMittelKarteX(i),
				targetY:  spielScreenMittelKartenY,
				karte:    s.mittelkarten[i],
			}
		}
	case spielStatusKartenWerdenGezogen:
		for _, karte := range s.bewegendeKarten() {
			if karte != nil && !karte.animationBeendet() {
				return
			}
		}

		s.status = spielStatusVerdeckteMittelkarten
	case spielStatusVerdeckteMittelkarten:
		if !istKlick() {
			return
		}

		s.einsätzeErmitteln(spielStatusVerdeckteMittelkarten, func() {
			s.status = spielStatusMittelkartenWerdenAufgedeckt
			for i := 0; i < 3; i++ {
				s.bewegendeMittelkarten[i].targetAufgedecktStatus = true
			}
		})
	case spielStatusMittelkartenWerdenAufgedeckt:
		for i := 0; i < 3; i++ {
			if !s.bewegendeMittelkarten[i].animationBeendet() {
				return
			}
		}

		s.status = spielStatus3AufgedeckteMittelKarten
	case spielStatus3AufgedeckteMittelKarten:
		if !istKlick() {
			return
		}

		s.einsätzeErmitteln(spielStatus3AufgedeckteMittelKarten, func() {
			s.status = spielStatusVierteKarteWirdGezogen
			s.bewegendeMittelkarten[3] = &spielScreenBewegendeKarte{
				currentX:      spielScreenStapelX,
				currentY:      spielScreenStapelY,
				targetX:       spielScreenMittelKarteX(3),
				targetY:       spielScreenMittelKartenY,
				karte:         s.mittelkarten[3],
				autoAufdecken: true,
			}
		})
	case spielStatusVierteKarteWirdGezogen:
		if !s.bewegendeMittelkarten[3].animationBeendet() {
			return
		}

		s.status = spielStatus4AufgedeckteMittelkarten
	case spielStatus4AufgedeckteMittelkarten:
		if !istKlick() {
			return
		}

		s.einsätzeErmitteln(spielStatus4AufgedeckteMittelkarten, func() {
			s.status = spielStatusFünfteKarteWirdGezogen
			s.bewegendeMittelkarten[4] = &spielScreenBewegendeKarte{
				currentX:      spielScreenStapelX,
				currentY:      spielScreenStapelY,
				targetX:       spielScreenMittelKarteX(4),
				targetY:       spielScreenMittelKartenY,
				karte:         s.mittelkarten[4],
				autoAufdecken: true,
			}
		})
	case spielStatusFünfteKarteWirdGezogen:
		if !s.bewegendeMittelkarten[4].animationBeendet() {
			return
		}

		s.status = spielStatus5AufgedeckteMittelkarten
	case spielStatus5AufgedeckteMittelkarten:
		if !istKlick() {
			return
		}

		s.einsätzeErmitteln(spielStatus5AufgedeckteMittelkarten, func() {
			for _, spieler := range s.spieler {
				for _, karte := range spieler.getBewegendeKarten() {
					karte.targetAufgedecktStatus = true
				}
			}
			s.status = spielStatusKartenWerdenAufgedeckt
		})
	case spielStatusKartenWerdenAufgedeckt:
		for _, karte := range s.bewegendeKarten() {
			if !karte.animationBeendet() {
				return
			}
		}

		s.status = spielStatusSiegerermittlung
	case spielStatusSiegerermittlung:
		gewinner := s.gewinner()
		gewinnerHand := s.spielerHand(gewinner)

		s.info.SetText(gewinnerHand.displayName())

		for _, karte := range gewinnerHand.karten() {
			if gewinnerHand.visualisierung(karte) == nil {
				continue
			}

			for _, bewegendeKarte := range s.bewegendeKarten() {
				if bewegendeKarte.karte != karte {
					continue
				}

				bewegendeKarte.hatUmrandung = true
				bewegendeKarte.umrandungsFarbe = gewinnerHand.visualisierung(karte)
				break
			}
		}
	}
}

func (s *spielScreen) drawStapel(screen *ebiten.Image) {
	bild := assets.RequireImage("spielkarten/verdeckt.png")
	var drawOptions ebiten.DrawImageOptions
	scaleX := spielScreenKarteBreite / float64(bild.Bounds().Dx())
	scaleY := spielScreenKarteHöhe / float64(bild.Bounds().Dy())
	drawOptions.GeoM.Scale(scaleX, scaleY)
	drawOptions.GeoM.Translate(spielScreenStapelX, spielScreenStapelY)
	screen.DrawImage(bild, &drawOptions)
}

func (s *spielScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 71, G: 113, B: 72, A: 255})
	s.drawStapel(screen)
	for _, component := range s.components() {
		component.Draw(screen)
	}
}
