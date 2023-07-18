package breakout

import (
	"math/rand"

	"github.com/Lama06/Herder-Legacy/ui"
	"golang.org/x/image/colornames"
)

type StandardLevelConfig struct {
	SteinSpalten, SteinZeilen    int
	UpgradeWahrscheinlichkeit    float64
	PlattformBreite              float64
	PlattformMaxSpeed            float64
	BallMaxSpeedX, BallMaxSpeedY float64
}

var (
	StandardLevelSuperEinfachConfig = StandardLevelConfig{
		SteinSpalten:              7,
		SteinZeilen:               4,
		UpgradeWahrscheinlichkeit: 0.4,
		PlattformBreite:           ui.Width / 2,
		PlattformMaxSpeed:         5,
		BallMaxSpeedX:             4,
		BallMaxSpeedY:             4,
	}
	StandardLevelEinfachConfig = StandardLevelConfig{
		SteinSpalten:              10,
		SteinZeilen:               5,
		UpgradeWahrscheinlichkeit: 0.3,
		PlattformBreite:           ui.Width / 3,
		PlattformMaxSpeed:         6,
		BallMaxSpeedX:             6,
		BallMaxSpeedY:             4,
	}
	StandardLevelMediumConfig = StandardLevelConfig{
		SteinSpalten:              13,
		SteinZeilen:               7,
		UpgradeWahrscheinlichkeit: 0.1,
		PlattformBreite:           ui.Width / 4,
		PlattformMaxSpeed:         8,
		BallMaxSpeedX:             8,
		BallMaxSpeedY:             5,
	}
	StandardLevelHardConfig = StandardLevelConfig{
		SteinSpalten:              16,
		SteinZeilen:               12,
		UpgradeWahrscheinlichkeit: 0.03,
		PlattformBreite:           ui.Width / 8,
		PlattformMaxSpeed:         30,
		BallMaxSpeedX:             7,
		BallMaxSpeedY:             7,
	}
)

func NewStandardLevel(config StandardLevelConfig) *world {
	const (
		steineHoeheGesamt = ui.Height / 2
		steinAbstandX     = 20
		steinAbstandY     = 20

		ballRadius = 25
		ballStartX = ui.Width/2 - ballRadius
		ballStartY = ui.Height * (3.0 / 4.0)

		plattformHoehe = 50
		plattformY     = ui.Height - plattformHoehe*2
	)

	var (
		steinBreite = (ui.Width - float64(config.SteinSpalten+1)*steinAbstandX) / float64(config.SteinSpalten)
		steinHoehe  = (steineHoeheGesamt - float64(config.SteinZeilen+1)*steinAbstandY) / float64(config.SteinZeilen)

		platformStartX = ui.Width/2 - config.PlattformBreite/2
	)

	world := world{
		entities: map[*entity]struct{}{
			// Plattform
			{
				position: position{
					x: platformStartX,
					y: plattformY,
				},
				hatMovesWithInputComponent: true,
				moveWithInputComponent: moveWithInputComponent{
					x:         true,
					y:         false,
					offsetX:   -config.PlattformBreite / 2,
					maxSpeedX: config.PlattformMaxSpeed,
				},
				hatRenderComponent: true,
				renderComponent: renderComponent{
					layer: renderLayerPlattform,
				},
				hatRectComponent: true,
				rectComponent: rectComponent{
					width:  config.PlattformBreite,
					height: plattformHoehe,
					farbe:  colornames.Purple,
				},
				hatHitboxComponent: true,
				hitboxComponent: hitboxComponent{
					width:  config.PlattformBreite,
					height: plattformHoehe,
				},
				hatAmRandAbprallenComponent: true,
				amRandAbprallenComponent: amRandAbprallenComponent{
					links:  true,
					rechts: true,
				},
				hatPlattformComponent: true,
				plattformComponent: plattformComponent{
					schießrichtungX: 0,
					schießrichtungY: -1,
				},
			}: {},
			// Ball
			{
				position: position{
					x: ballStartX,
					y: ballStartY,
				},
				hatVelocityComponent: true,
				velocityComponent: velocityComponent{
					velocityX: 0,
					velocityY: config.BallMaxSpeedY,
				},
				hatRenderComponent: true,
				renderComponent: renderComponent{
					layer: renderLayerBall,
				},
				hatCircleComponent: true,
				circleComponent: circleComponent{
					radius: ballRadius,
					farbe:  colornames.Whitesmoke,
				},
				hatHitboxComponent: true,
				hitboxComponent: hitboxComponent{
					width:  ballRadius * 2,
					height: ballRadius * 2,
				},
				hatAnHitboxenAbprallenComponent: true,
				anHitboxenAbprallenComponent: anHitboxenAbprallenComponent{
					maxXSpeed:       config.BallMaxSpeedX,
					maxYSpeed:       config.BallMaxSpeedY,
					minYSpeed:       config.BallMaxSpeedY,
					steineZerstören: true,
				},
				hatAmRandAbprallenComponent: true,
				amRandAbprallenComponent: amRandAbprallenComponent{
					oben:   true,
					links:  true,
					rechts: true,
					unten:  false,
				},
				imAusEntfernen:        true,
				affectsAutomaticInput: true,
			}: {},
		},
	}

	for zeile := 0; zeile < config.SteinZeilen; zeile++ {
		for spalte := 0; spalte < config.SteinSpalten; spalte++ {
			stein := entity{
				position: position{
					x: float64(spalte)*(steinBreite+steinAbstandX) + steinAbstandX,
					y: float64(zeile)*(steinHoehe+steinAbstandY) + steinAbstandY,
				},
				hatRenderComponent: true,
				renderComponent: renderComponent{
					layer: renderLayerStein,
				},
				hatRectComponent: true,
				rectComponent: rectComponent{
					width:  steinBreite,
					height: steinHoehe,
					farbe:  colornames.Pink,
				},
				hatHitboxComponent: true,
				hitboxComponent: hitboxComponent{
					width:  steinBreite,
					height: steinHoehe,
				},
				hatSteinComponent: true,
				steinComponent: steinComponent{
					upgradeSpeedX: 0,
					upgradeSpeedY: 5,
				},
			}

			if rand.Float64() < config.UpgradeWahrscheinlichkeit {
				stein.steinComponent.hatUpgrade = true
				stein.steinComponent.upgrade = newRandomUpgrade()
				stein.rectComponent.farbe = stein.steinComponent.upgrade.farbe()
			}

			world.entities[&stein] = struct{}{}
		}
	}

	return &world
}
