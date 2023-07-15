package dame

import "strings"

type lehrer struct {
	name             string
	nameAkkusativ    string
	personalPronomen string

	info         string
	anfangsBrett brett
	regeln       regeln
	aiTiefe      int
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

var (
	alleLehrer = []lehrer{
		HerrPreuß,
		HerrMünch,
		HerrWeber,
		FrauDahmen,
		HerrJonaßon,
		HerrFahnenmüller,
		HerrSchwehner,
		FrauKüpper,
		FrauHannes,
	}
	HerrPreuß = lehrer{
		name:             "Herr Preuß",
		nameAkkusativ:    "Herrn Preuß",
		personalPronomen: "Er",
		info: `Herr Preuß spielt mit internationalen Dame Regeln.
Das heißt, dass auf einem 8x8 großen Spielbrett gespielt wird.
Normale Steine sowie Damen können sich diagonal bewegen und diagonal schlagen.`,
		anfangsBrett: mustParseBrett(
			"_l_l_l_l",
			"l_l_l_l_",
			"_l_l_l_l",
			"________",
			"________",
			"s_s_s_s_",
			"_s_s_s_s",
			"s_s_s_s_",
		),
		regeln:  internationaleRegeln,
		aiTiefe: 5,
	}
	HerrMünch = lehrer{
		name:             "Herr Münch",
		nameAkkusativ:    "Herrn Münch",
		personalPronomen: "Er",
		info: `Herr Münch spielt mit internationalen Dame Regeln, allerdings auf einem 8x16 Spielbrett, also auf zwei Spielbrettern übereinander.
Normale Steine sowie Damen können sich weiterhin diagonal bewegen und diagonal schlagen.`,
		anfangsBrett: mustParseBrett(
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
		regeln:  internationaleRegeln,
		aiTiefe: 4,
	}
	HerrWeber = lehrer{
		name:             "Herr Weber",
		nameAkkusativ:    "Herrn Weber",
		personalPronomen: "Er",
		info: `Herr Weber spielt mit altdeutschen Dameregeln.
Die Steine werden zwar diagonal bewegt, schlagen aber seitwärts (nach links und rechts), diagonal und vorwärts.
Damen können sich nach oben, unten, links und rechts bewegen.`,
		anfangsBrett: mustParseBrett(
			"llllllll",
			"llllllll",
			"________",
			"________",
			"________",
			"________",
			"ssssssss",
			"ssssssss",
		),
		regeln:  altdeutscheRegeln,
		aiTiefe: 5,
	}
	FrauDahmen = lehrer{
		name:             "Frau Dahmen",
		personalPronomen: "Sie",
		info:             "Frau Dahmen spielt nur mit Damen ;)",
		anfangsBrett: mustParseBrett(
			"_L_L_L_L",
			"L_L_L_L_",
			"_L_L_L_L",
			"________",
			"________",
			"S_S_S_S_",
			"_S_S_S_S",
			"S_S_S_S_",
		),
		regeln:  internationaleRegeln,
		aiTiefe: 4,
	}
	HerrJonaßon = lehrer{
		name:             "Herr Jonaßon",
		nameAkkusativ:    "Herrn Jonaßon",
		personalPronomen: "Er",
		info:             "Herr Jonaßon ist noch nicht sehr erfahren in Dame, deswegen spielt er mit einer Zeile mehr Steinen als du.",
		anfangsBrett: mustParseBrett(
			"_l_l_l_l",
			"l_l_l_l_",
			"_l_l_l_l",
			"l_l_l_l_",
			"________",
			"s_s_s_s_",
			"_s_s_s_s",
			"s_s_s_s_",
		),
		regeln:  internationaleRegeln,
		aiTiefe: 2,
	}
	HerrFahnenmüller = lehrer{
		name:             "Herr Fahnenmüller",
		nameAkkusativ:    "Herrn Fahnenmüller",
		personalPronomen: "Er",
		info:             "Herr Fahnenmüller hat viel Übung in Dame. Daher wird er versuchen, mit drei Steinen weniger zu spielen.",
		anfangsBrett: mustParseBrett(
			"_l_l_l_l",
			"l_l_l_l_",
			"___l____",
			"________",
			"________",
			"s_s_s_s_",
			"_s_s_s_s",
			"s_s_s_s_",
		),
		regeln:  internationaleRegeln,
		aiTiefe: 6,
	}
	HerrSchwehner = lehrer{
		name:             "Herr Schwehmer",
		nameAkkusativ:    "Herrn Schwehmer",
		personalPronomen: "Er",
		info:             "Herr Schwehmer spielt mit internationalen Regeln, aber sowohl auf schwarzen als auch weißen Feldern.",
		anfangsBrett: mustParseBrett(
			"llllllll",
			"llllllll",
			"llllllll",
			"________",
			"________",
			"ssssssss",
			"ssssssss",
			"ssssssss",
		),
		regeln:  internationaleRegeln,
		aiTiefe: 4,
	}
	FrauKüpper = lehrer{
		name:             "Frau Küpper",
		personalPronomen: "Sie",
		info:             "Frau Küpper spielt auf einem 5x5 großen Brett mit altdeutschen Dameregeln.",
		anfangsBrett: mustParseBrett(
			"lllll",
			"lllll",
			"_____",
			"sssss",
			"sssss",
		),
		regeln:  altdeutscheRegeln,
		aiTiefe: 4,
	}
	FrauHannes = lehrer{
		name:             "Frau Hannes",
		personalPronomen: "Sie",
		info:             "Frau Hannes spielt Dame mit normalen Regeln, aber im Querformat :)",
		anfangsBrett: mustParseBrett(
			"_l_l_l_l_l_l_l_l_l_l_l_l_l_l_l_l",
			"l_l_l_l_l_l_l_l_l_l_l_l_l_l_l_l_",
			"________________________________",
			"________________________________",
			"_s_s_s_s_s_s_s_s_s_s_s_s_s_s_s_s",
			"s_s_s_s_s_s_s_s_s_s_s_s_s_s_s_s_",
		),
		regeln:  internationaleRegeln,
		aiTiefe: 2,
	}
)
