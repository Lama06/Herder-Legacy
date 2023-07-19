package dame

type ZugRegeln struct {
	SteinBewegenRichtungen                Richtungen
	SteinSchlagenRichtungenAnfang         Richtungen
	SteinSchlagenRichtungenWeiterschlagen Richtungen

	DameBewegenRichtungen                Richtungen
	DameSchlagenRichtungenAnfang         Richtungen
	DameSchlagenRichtungenWeiterschlagen Richtungen

	SchlagZwang bool
}

func (z ZugRegeln) steinSchlagenRichtungen(weitschlagen bool) Richtungen {
	if weitschlagen {
		return z.SteinSchlagenRichtungenWeiterschlagen
	}
	return z.SteinSchlagenRichtungenAnfang
}

func (z ZugRegeln) dameSchlagenRichtungen(weiterschlagen bool) Richtungen {
	if weiterschlagen {
		return z.DameSchlagenRichtungenWeiterschlagen
	}
	return z.DameSchlagenRichtungenAnfang
}

func (z ZugRegeln) clone() ZugRegeln {
	return ZugRegeln{
		SteinBewegenRichtungen:                z.SteinBewegenRichtungen.clone(),
		SteinSchlagenRichtungenAnfang:         z.SteinSchlagenRichtungenAnfang.clone(),
		SteinSchlagenRichtungenWeiterschlagen: z.SteinSchlagenRichtungenWeiterschlagen.clone(),

		DameBewegenRichtungen:                z.DameBewegenRichtungen.clone(),
		DameSchlagenRichtungenAnfang:         z.DameSchlagenRichtungenAnfang.clone(),
		DameSchlagenRichtungenWeiterschlagen: z.DameSchlagenRichtungenWeiterschlagen.clone(),

		SchlagZwang: z.SchlagZwang,
	}
}

var (
	InternationaleZugRegeln = ZugRegeln{
		SteinBewegenRichtungen:                RichtungenDiagonalVorne,
		SteinSchlagenRichtungenAnfang:         RichtungenDiagonalVorne,
		SteinSchlagenRichtungenWeiterschlagen: RichtungenDiagonal,

		DameBewegenRichtungen:                RichtungenDiagonal,
		DameSchlagenRichtungenAnfang:         RichtungenDiagonal,
		DameSchlagenRichtungenWeiterschlagen: RichtungenDiagonal,

		SchlagZwang: true,
	}
	AltdeutscheZugRegeln = ZugRegeln{
		SteinBewegenRichtungen:                RichtungenDiagonalVorne,
		SteinSchlagenRichtungenAnfang:         RichtungenSeiteDiagonalUndVorne,
		SteinSchlagenRichtungenWeiterschlagen: RichtungenSeiteDiagonalUndVorne,

		DameBewegenRichtungen:                RichtungenAlle,
		DameSchlagenRichtungenAnfang:         RichtungenAlle,
		DameSchlagenRichtungenWeiterschlagen: RichtungenAlle,

		SchlagZwang: true,
	}
)
