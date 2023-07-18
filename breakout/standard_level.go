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
	Störkugel                    bool
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
		Störkugel:                 false,
	}
	StandardLevelEinfachConfig = StandardLevelConfig{
		SteinSpalten:              10,
		SteinZeilen:               5,
		UpgradeWahrscheinlichkeit: 0.3,
		PlattformBreite:           ui.Width / 3,
		PlattformMaxSpeed:         6,
		BallMaxSpeedX:             6,
		BallMaxSpeedY:             4,
		Störkugel:                 false,
	}
	StandardLevelMediumConfig = StandardLevelConfig{
		SteinSpalten:              13,
		SteinZeilen:               7,
		UpgradeWahrscheinlichkeit: 0.2,
		PlattformBreite:           ui.Width / 4,
		PlattformMaxSpeed:         8,
		BallMaxSpeedX:             8,
		BallMaxSpeedY:             5,
		Störkugel:                 false,
	}
	StandardLevelHardConfig = StandardLevelConfig{
		SteinSpalten:              16,
		SteinZeilen:               12,
		UpgradeWahrscheinlichkeit: 0.2,
		PlattformBreite:           ui.Width / 8,
		PlattformMaxSpeed:         30,
		BallMaxSpeedX:             12,
		BallMaxSpeedY:             7,
		Störkugel:                 true,
	}
	StandardLevelExpertConfig = StandardLevelConfig{
		SteinSpalten:              18,
		SteinZeilen:               14,
		UpgradeWahrscheinlichkeit: 0.2,
		PlattformBreite:           ui.Width / 16,
		PlattformMaxSpeed:         100,
		BallMaxSpeedX:             20,
		BallMaxSpeedY:             9,
		Störkugel:                 true,
	}
)

func NewStandardLevel(config StandardLevelConfig) *world {
	const (
		steineHoeheGesamt = ui.Height / 2
		steinAbstandX     = 20
		steinAbstandY     = 20

		ballRadius = 20
		ballStartX = ui.Width/2 - ballRadius
		ballStartY = ui.Height * (3.0 / 4.0)

		plattformHoehe = 40
		plattformY     = ui.Height - plattformHoehe*2

		störkugelRadius    = 13 // Das kann nur Unglück bringen
		störkugelStartX    = ui.Width/2 - störkugelRadius
		störkugelStartY    = ballStartY + ballRadius*4
		störkugelMaxSpeedX = 2
		störkugelMaySpeedY = 2
	)

	var (
		steinBreite = (ui.Width - float64(config.SteinSpalten+1)*steinAbstandX) / float64(config.SteinSpalten)
		steinHoehe  = (steineHoeheGesamt - float64(config.SteinZeilen+1)*steinAbstandY) / float64(config.SteinZeilen)

		platformStartX = ui.Width/2 - config.PlattformBreite/2
	)

	w := world{
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
					maxXSpeed: config.BallMaxSpeedX,
					maxYSpeed: config.BallMaxSpeedY,
					minYSpeed: config.BallMaxSpeedY / 10,
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
				istBall:               true,
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

			w.entities[&stein] = struct{}{}
		}
	}

	if config.Störkugel {
		störkugel := entity{
			position: position{
				x: störkugelStartX,
				y: störkugelStartY,
			},
			hatVelocityComponent: true,
			velocityComponent: velocityComponent{
				velocityX: 0,
				velocityY: -störkugelMaySpeedY,
			},
			hatRenderComponent: true,
			renderComponent: renderComponent{
				layer: renderLayerStörkugel,
			},
			hatCircleComponent: true,
			circleComponent: circleComponent{
				radius: störkugelRadius,
				farbe:  colornames.Brown,
			},
			hatHitboxComponent: true,
			hitboxComponent: hitboxComponent{
				width:  störkugelRadius * 2,
				height: störkugelRadius * 2,
			},
			hatAmRandAbprallenComponent: true,
			amRandAbprallenComponent: amRandAbprallenComponent{
				oben:   true,
				unten:  true,
				links:  true,
				rechts: true,
			},
			hatAnHitboxenAbprallenComponent: true,
			anHitboxenAbprallenComponent: anHitboxenAbprallenComponent{
				maxXSpeed: störkugelMaxSpeedX,
				maxYSpeed: störkugelMaySpeedY,
			},
		}

		w.entities[&störkugel] = struct{}{}
	}

	return &w
}
