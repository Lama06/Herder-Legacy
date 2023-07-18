package breakout

import (
	"math"
	"math/rand"

	"github.com/Lama06/Herder-Legacy/ui"
	"golang.org/x/image/colornames"
)

type PongLevelConfig struct {
	AnzahlSpawnSteine            int
	SpawnDelay                   int
	GespawnterSteinSpeed         float64
	PlattformHoehe               float64
	PlattformSpeed               float64
	UpgradeWahrscheinlichkeit    float64
	BallMaxSpeedX, BallMaxSpeedY float64
}

var (
	PongLevelEinfachConfig = PongLevelConfig{
		AnzahlSpawnSteine:         6,
		SpawnDelay:                12 * 60,
		GespawnterSteinSpeed:      1.5,
		PlattformHoehe:            ui.Height / 3,
		PlattformSpeed:            40,
		UpgradeWahrscheinlichkeit: 0.25,
		BallMaxSpeedX:             9,
		BallMaxSpeedY:             6,
	}
	PongLevelNormalConfig = PongLevelConfig{
		AnzahlSpawnSteine:         10,
		SpawnDelay:                6 * 60,
		GespawnterSteinSpeed:      2,
		PlattformHoehe:            ui.Height / 4,
		PlattformSpeed:            50,
		UpgradeWahrscheinlichkeit: 0.20,
		BallMaxSpeedX:             9,
		BallMaxSpeedY:             6,
	}
	PongLevelHartConfig = PongLevelConfig{
		AnzahlSpawnSteine:         20,
		SpawnDelay:                4 * 60,
		GespawnterSteinSpeed:      4,
		PlattformHoehe:            ui.Height / 9,
		PlattformSpeed:            math.MaxInt,
		UpgradeWahrscheinlichkeit: 0.1,
		BallMaxSpeedX:             9,
		BallMaxSpeedY:             6,
	}
)

func NewPongLevel(config PongLevelConfig) *world {
	const (
		ballRadius = 20
		ballStartX = ui.Width/2 - ballRadius
		ballStartY = ui.Height/2 - ballRadius

		plattformBreite = 40
		plattform1X     = plattformBreite
		plattform2X     = ui.Width - plattformBreite*2

		steinAbstand                  = 20
		steinBereichStartY            = steinAbstand
		steinBereichAbstandHorizontal = 120
		steinBereichStartX            = plattform1X + plattformBreite + steinBereichAbstandHorizontal + steinAbstand
		steinBereichBreite            = (plattform2X - steinBereichAbstandHorizontal - steinAbstand) - steinBereichStartX
	)

	var (
		plattformStartY = ui.Height/2 - config.PlattformHoehe/2

		steinBreite = (steinBereichBreite - (float64(config.AnzahlSpawnSteine)-1)*steinAbstand) / float64(config.AnzahlSpawnSteine)
		steinHoehe  = clamp(20, steinBreite/2, 50)
	)

	w := world{
		entities: map[*entity]struct{}{
			// Plattform 1
			{
				position: position{
					x: plattform1X,
					y: plattformStartY,
				},
				hatRenderComponent: true,
				renderComponent: renderComponent{
					layer: renderLayerPlattform,
				},
				hatRectComponent: true,
				rectComponent: rectComponent{
					width:  plattformBreite,
					height: config.PlattformHoehe,
					farbe:  colornames.Yellowgreen,
				},
				hatMovesWithInputComponent: true,
				moveWithInputComponent: moveWithInputComponent{
					x:         false,
					y:         true,
					maxSpeedY: config.PlattformSpeed,
					offsetY:   -config.PlattformHoehe / 2,
				},
				hatHitboxComponent: true,
				hitboxComponent: hitboxComponent{
					width:  plattformBreite,
					height: config.PlattformHoehe,
				},
				hatPlattformComponent: true,
				plattformComponent: plattformComponent{
					schießrichtungX: 1,
					schießrichtungY: 0,
				},
			}: {},
			// Plattform 2
			{
				position: position{
					x: plattform2X,
					y: plattformStartY,
				},
				hatRenderComponent: true,
				renderComponent: renderComponent{
					layer: renderLayerPlattform,
				},
				hatRectComponent: true,
				rectComponent: rectComponent{
					width:  plattformBreite,
					height: config.PlattformHoehe,
					farbe:  colornames.Yellowgreen,
				},
				hatMovesWithInputComponent: true,
				moveWithInputComponent: moveWithInputComponent{
					x:         false,
					y:         true,
					maxSpeedY: config.PlattformSpeed,
					offsetY:   -config.PlattformHoehe / 2,
				},
				hatHitboxComponent: true,
				hitboxComponent: hitboxComponent{
					width:  plattformBreite,
					height: config.PlattformHoehe,
				},
				hatPlattformComponent: true,
				plattformComponent: plattformComponent{
					schießrichtungX: -1,
					schießrichtungY: 0,
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
					velocityX: -config.BallMaxSpeedX,
					velocityY: 0,
				},
				hatRenderComponent: true,
				renderComponent: renderComponent{
					layer: renderLayerBall,
				},
				hatCircleComponent: true,
				circleComponent: circleComponent{
					radius: ballRadius,
					farbe:  colornames.Hotpink,
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
					minXSpeed: config.BallMaxSpeedY / 10,
				},
				hatAmRandAbprallenComponent: true,
				amRandAbprallenComponent: amRandAbprallenComponent{
					oben:   true,
					links:  false,
					rechts: false,
					unten:  true,
				},
				imAusEntfernen:        true,
				affectsAutomaticInput: true,
				istBall:               true,
			}: {},
		},
	}

	for spalte := 0; spalte < config.AnzahlSpawnSteine; spalte++ {
		spawnStein := entity{
			position: position{
				x: steinBereichStartX + float64(spalte)*(steinBreite+steinAbstand),
				y: steinBereichStartY,
			},
			hatRenderComponent: true,
			renderComponent: renderComponent{
				layer: renderLayerStein,
			},
			hatRectComponent: true,
			rectComponent: rectComponent{
				width:  steinBreite,
				height: steinHoehe,
				farbe:  colornames.Purple,
			},
			hatHitboxComponent: true,
			hitboxComponent: hitboxComponent{
				width:  steinBreite,
				height: steinHoehe,
			},
			hatSteinComponent: true,
			steinComponent: steinComponent{
				hatUpgrade: false,
			},
			hatSpawnerComponent: true,
			spawnerComponent: spawnerComponent{
				verbleibendeSpawns: math.MaxInt,
				spawnDelay:         config.SpawnDelay,
				nächsterSpawn:      config.SpawnDelay,
				creator: func(world *world, spawner *entity) *entity {
					stein := entity{
						position: position{
							x: spawner.position.x,
							y: spawner.position.y + steinHoehe + steinAbstand,
						},
						hatVelocityComponent: true,
						velocityComponent: velocityComponent{
							velocityX: 0,
							velocityY: config.GespawnterSteinSpeed,
						},
						imAusEntfernen:     true,
						hatRenderComponent: true,
						renderComponent: renderComponent{
							layer: renderLayerStein,
						},
						hatRectComponent: true,
						rectComponent: rectComponent{
							width:  steinBreite,
							height: steinHoehe,
							farbe:  colornames.Purple,
						},
						hatHitboxComponent: true,
						hitboxComponent: hitboxComponent{
							width:  steinBreite,
							height: steinHoehe,
						},
						hatSteinComponent: true,
						steinComponent: steinComponent{
							hatUpgrade: false,
						},
					}

					if rand.Float64() < config.UpgradeWahrscheinlichkeit {
						stein.steinComponent.hatUpgrade = true
						stein.steinComponent.upgrade = newRandomUpgrade()
						if rand.Float64() < 0.5 {
							stein.steinComponent.upgradeSpeedX = 1
						} else {
							stein.steinComponent.upgradeSpeedX = -1
						}

						stein.rectComponent.farbe = stein.steinComponent.upgrade.farbe()
					}

					return &stein
				},
			},
		}

		w.entities[&spawnStein] = struct{}{}
	}

	return &w
}
