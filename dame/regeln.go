package dame

type ZugRegeln struct {
	SteinBewegenRichtungen                map[Richtung]struct{}
	SteinSchlagenRichtungenAnfang         map[Richtung]struct{}
	SteinSchlagenRichtungenWeiterschlagen map[Richtung]struct{}

	DameBewegenRichtungen                map[Richtung]struct{}
	DameSchlagenRichtungenAnfang         map[Richtung]struct{}
	DameSchlagenRichtungenWeiterschlagen map[Richtung]struct{}

	SchlagZwang bool
}

func (z ZugRegeln) steinSchlagenRichtungen(weitschlagen bool) map[Richtung]struct{} {
	if weitschlagen {
		return z.SteinSchlagenRichtungenWeiterschlagen
	}
	return z.SteinSchlagenRichtungenAnfang
}

func (z ZugRegeln) dameSchlagenRichtungen(weiterschlagen bool) map[Richtung]struct{} {
	if weiterschlagen {
		return z.DameSchlagenRichtungenWeiterschlagen
	}
	return z.DameSchlagenRichtungenAnfang
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
