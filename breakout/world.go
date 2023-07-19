package breakout

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type position struct {
	x, y float64
}

type plattformComponent struct {
	schießrichtungX, schießrichtungY int
}

type steinComponent struct {
	hatUpgrade                   bool
	upgrade                      upgrade
	upgradeSpeedX, upgradeSpeedY float64
}

type entity struct {
	position position

	hatRenderComponent bool
	renderComponent    renderComponent

	hatRectComponent bool
	rectComponent    rectComponent

	hatCircleComponent bool
	circleComponent    circleComponent

	hatVelocityComponent bool
	velocityComponent    velocityComponent

	hatMovesWithInputComponent bool
	moveWithInputComponent     moveWithInputComponent

	affectsAutomaticInput bool

	hatAmRandAbprallenComponent bool
	amRandAbprallenComponent    amRandAbprallenComponent

	hatHitboxComponent bool
	hitboxComponent    hitboxComponent

	hatAnHitboxenAbprallenComponent bool
	anHitboxenAbprallenComponent    anHitboxenAbprallenComponent

	imAusEntfernen bool

	hatSpawnerComponent bool
	spawnerComponent    spawnerComponent

	istBall bool

	hatSteinComponent bool
	steinComponent    steinComponent

	hatPlattformComponent bool
	plattformComponent    plattformComponent

	hatFallendesUpgradeComponent bool
	fallendesUpgradeComponent    fallendesUpgradeComponent

	hatKanonenKugelComponent bool
	kanonenKugelComponent    kanonenKugelComponent

	hatRainbowModeColorChangeComponent bool
	rainbowModeColorChangeComponent    rainbowModeColorChangeComponent
}

type world struct {
	entities map[*entity]struct{}

	rainbowUpgradeRemainingTime   int
	rainbowUpgradeChangeFrequency int

	fasterInputUpgradeRemainingTime int
	fasterInputUpgradeMultiplier    float64

	automaticInputUpgradeRemainingTime int

	zeitUpgradeRemainingTime int
	zeitUpgradeFaktor        float64
}

func (w *world) update() {
	w.moveWithVelocity()
	w.anHitboxenAbprallen()
	w.moveWithInput()
	w.imAusEntfernen()
	w.fallendeUpgradesAufsammeln()
	w.entitesSpawnen()
	w.mitKanonenKugelnSteinZerstören()
	w.changeColorsInRainbowMode()
	w.tickFasterInputUpgradeTimer()
	w.performAutomaticInput()
	w.amRandAbprallen()
	w.tickZeitUpgrade()
}

func (w *world) draw(screen *ebiten.Image) {
	screen.Fill(colornames.Black)
	w.renderObjects(screen)
}

func (w *world) gewonnen() bool {
	for stein := range w.entities {
		if stein.hatSteinComponent {
			return false
		}
	}
	return true
}

func (w *world) verloren() bool {
	for ball := range w.entities {
		if ball.istBall {
			return false
		}
	}

	return true
}
