package dame

import "strings"

type lehrer struct {
	name             string
	nameAkkusativ    string
	personalPronomen string
	info             string
	spielOptionen    SpielOptionen
}

func (l lehrer) nameAkkusativOrDefault() string {
	if l.nameAkkusativ == "" {
		return l.name
	}
	return l.nameAkkusativ
}

func (l lehrer) personalPronomenSatzanfang() string {
	return strings.ToUpper(l.personalPronomen[:1]) + l.personalPronomen[1:]
}

func (l lehrer) personalPronomenSatzmitte() string {
	return strings.ToLower(l.personalPronomen[:1]) + l.personalPronomen[1:]
}

var alleLehrer = []lehrer{
	{
		name:             "Herr Preuß",
		nameAkkusativ:    "Herrn Preuß",
		personalPronomen: "Er",
		info: `Herr Preuß spielt mit internationalen Dame Regeln.
		Das heißt, dass auf einem 8x8 großen Spielbrett gespielt wird.
		Normale Steine sowie Damen können sich diagonal bewegen und diagonal schlagen.`,
		spielOptionen: SpielOptionen{
			StartBrett: MustParseBrett(
				"_l_l_l_l",
				"l_l_l_l_",
				"_l_l_l_l",
				"________",
				"________",
				"s_s_s_s_",
				"_s_s_s_s",
				"s_s_s_s_",
			),
			ZugRegeln: InternationaleZugRegeln,
			AiTiefe:   5,
		},
	},
	{
		name:             "Herr Münch",
		nameAkkusativ:    "Herrn Münch",
		personalPronomen: "Er",
		info: `Herr Münch spielt mit internationalen Dame Regeln, allerdings auf einem 8x16 Spielbrett, also auf zwei Spielbrettern übereinander.
		Normale Steine sowie Damen können sich weiterhin diagonal bewegen und diagonal schlagen.`,
		spielOptionen: SpielOptionen{
			StartBrett: MustParseBrett(
				"_l_l_l_l",
				"l_l_l_l_",
				"_l_l_l_l",
				"l_l_l_l_",
				"________",
				"________",
				"________",
				"________",
				"________",
				"________",
				"________",
				"________",
				"_s_s_s_s",
				"s_s_s_s_",
				"_s_s_s_s",
				"s_s_s_s_",
			),
			ZugRegeln: InternationaleZugRegeln,
			AiTiefe:   4,
		},
	},
	{
		name:             "Herr Weber",
		nameAkkusativ:    "Herrn Weber",
		personalPronomen: "Er",
		info: `Herr Weber spielt mit altdeutschen Dameregeln.
		Die Steine werden zwar diagonal bewegt, schlagen aber seitwärts (nach links und rechts), diagonal und vorwärts.
		Damen können sich nach oben, unten, links und rechts bewegen.`,
		spielOptionen: SpielOptionen{
			StartBrett: MustParseBrett(
				"llllllll",
				"llllllll",
				"________",
				"________",
				"________",
				"________",
				"ssssssss",
				"ssssssss",
			),
			ZugRegeln: AltdeutscheZugRegeln,
			AiTiefe:   5,
		},
	},
	{
		name:             "Frau Dahmen",
		personalPronomen: "Sie",
		info:             "Frau Dahmen spielt nur mit Damen ;)",
		spielOptionen: SpielOptionen{
			StartBrett: MustParseBrett(
				"_L_L_L_L",
				"L_L_L_L_",
				"_L_L_L_L",
				"________",
				"________",
				"S_S_S_S_",
				"_S_S_S_S",
				"S_S_S_S_",
			),
			ZugRegeln: InternationaleZugRegeln,
			AiTiefe:   4,
		},
	},
	{
		name:             "Herr Jonaßon",
		nameAkkusativ:    "Herrn Jonaßon",
		personalPronomen: "Er",
		info:             "Herr Jonaßon ist noch nicht sehr erfahren in Dame, deswegen spielt er mit einer Zeile mehr Steinen als du.",
		spielOptionen: SpielOptionen{
			StartBrett: MustParseBrett(
				"_l_l_l_l",
				"l_l_l_l_",
				"_l_l_l_l",
				"l_l_l_l_",
				"________",
				"s_s_s_s_",
				"_s_s_s_s",
				"s_s_s_s_",
			),
			ZugRegeln: InternationaleZugRegeln,
			AiTiefe:   2,
		},
	},
	{
		name:             "Herr Fahnenmüller",
		nameAkkusativ:    "Herrn Fahnenmüller",
		personalPronomen: "Er",
		info:             "Herr Fahnenmüller hat viel Übung in Dame. Daher wird er versuchen, mit drei Steinen weniger zu spielen.",
		spielOptionen: SpielOptionen{
			StartBrett: MustParseBrett(
				"_l_l_l_l",
				"l_l_l_l_",
				"___l____",
				"________",
				"________",
				"s_s_s_s_",
				"_s_s_s_s",
				"s_s_s_s_",
			),
			ZugRegeln: InternationaleZugRegeln,
			AiTiefe:   6,
		},
	},
	{
		name:             "Herr Schwehmer",
		nameAkkusativ:    "Herrn Schwehmer",
		personalPronomen: "Er",
		info:             "Herr Schwehmer spielt mit internationalen Regeln, aber sowohl auf schwarzen als auch weißen Feldern.",
		spielOptionen: SpielOptionen{
			StartBrett: MustParseBrett(
				"llllllll",
				"llllllll",
				"llllllll",
				"________",
				"________",
				"ssssssss",
				"ssssssss",
				"ssssssss",
			),
			ZugRegeln: InternationaleZugRegeln,
			AiTiefe:   4,
		},
	},
	{
		name:             "Frau Küpper",
		personalPronomen: "Sie",
		info:             "Frau Küpper spielt auf einem 5x5 großen Brett mit altdeutschen Dameregeln.",
		spielOptionen: SpielOptionen{
			StartBrett: MustParseBrett(
				"lllll",
				"lllll",
				"_____",
				"sssss",
				"sssss",
			),
			ZugRegeln: AltdeutscheZugRegeln,
			AiTiefe:   4,
		},
	},
	{
		name:             "Frau Hannes",
		personalPronomen: "Sie",
		info:             "Frau Hannes spielt Dame mit normalen Regeln, aber im Querformat :)",
		spielOptionen: SpielOptionen{
			StartBrett: MustParseBrett(
				"_l_l_l_l_l_l_l_l_l_l_l_l_l_l_l_l",
				"l_l_l_l_l_l_l_l_l_l_l_l_l_l_l_l_",
				"________________________________",
				"________________________________",
				"_s_s_s_s_s_s_s_s_s_s_s_s_s_s_s_s",
				"s_s_s_s_s_s_s_s_s_s_s_s_s_s_s_s_",
			),
			ZugRegeln: InternationaleZugRegeln,
			AiTiefe:   2,
		},
	},
}
