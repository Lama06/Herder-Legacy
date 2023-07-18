package breakout

type breakoutLevel struct {
	name         string
	worldCreator func() *world
}

var breakoutLevelListe = []breakoutLevel{
	{
		name: "Standard - Supereinfach",
		worldCreator: func() *world {
			return NewStandardLevel(StandardLevelSuperEinfachConfig)
		},
	},
	{
		name: "Standard - Einfach",
		worldCreator: func() *world {
			return NewStandardLevel(StandardLevelEinfachConfig)
		},
	},
	{
		name: "Standard - Medium",
		worldCreator: func() *world {
			return NewStandardLevel(StandardLevelMediumConfig)
		},
	},
	{
		name: "Standard - Hart",
		worldCreator: func() *world {
			return NewStandardLevel(StandardLevelHardConfig)
		},
	},
	{
		name: "Standard - Herr Hammdorf Modus",
		worldCreator: func() *world {
			return NewStandardLevel(StandardLevelExpertConfig)
		},
	},
	{
		name: "Pong - Einfach",
		worldCreator: func() *world {
			return NewPongLevel(PongLevelEinfachConfig)
		},
	},
	{
		name: "Pong - Is ja Normal ne",
		worldCreator: func() *world {
			return NewPongLevel(PongLevelNormalConfig)
		},
	},
	{
		name: "Pong - Ist das möglich?",
		worldCreator: func() *world {
			return NewPongLevel(PongLevelHartConfig)
		},
	},
}
