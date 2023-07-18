package breakout

import (
	"github.com/hajimehoshi/ebiten/v2"
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

	hatVelocityComponent bool
	velocityComponent    velocityComponent

	hatMovesWithInputComponent bool
	moveWithInputComponent     moveWithInputComponent

	hatRenderComponent bool
	renderComponent    renderComponent

	hatRectComponent bool
	rectComponent    rectComponent

	hatCircleComponent bool
	circleComponent    circleComponent

	hatAmRandAbprallenComponent bool
	amRandAbprallenComponent    amRandAbprallenComponent

	hatAnHitboxenAbprallenComponent bool
	anHitboxenAbprallenComponent    anHitboxenAbprallenComponent

	hatHitboxComponent bool
	hitboxComponent    hitboxComponent

	imAusEntfernen bool

	hatTimerComponent bool
	timerComponent    timerComponent

	hatSteinComponent bool
	steinComponent    steinComponent

	hatPlattformComponent bool
	plattformComponent    plattformComponent

	istBall bool

	hatFallendesUpgradeComponent bool
	fallendesUpgradeComponent    fallendesUpgradeComponent

	hatKanonenKugelComponent bool
	kanonenKugelComponent    kanonenKugelComponent

	hatKanonenKugelSpawnerComponent bool
	kanonenKugelSpawnerComponent    kanonenKugelSpawnerComponent

	affectsAutomaticInput bool
}

type world struct {
	entities map[*entity]struct{}

	rainbowUpgradeRemainingTime   int
	rainbowUpgradeChangeFrequency int

	fasterInputUpgradeRemainingTime int
	fasterInputUpgradeMultiplier    float64

	automaticInputUpgradeRemainingTime int
}

func (w *world) update() {
	w.moveWithVelocity()
	w.anHitboxenAbprallen()
	w.moveWithInput()
	w.imAusEntfernen()
	w.timersRunterzählen()
	w.fallendeUpgradesAufsammeln()
	w.kanonenKugelnSpawnen()
	w.mitKanonenKugelnSteinZerstören()
	w.changeColorsInRainbowMode()
	w.tickFasterInputUpgradeTimer()
	w.performAutomaticInput()
	w.amRandAbprallen()
}

func (w *world) draw(screen *ebiten.Image) {
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
