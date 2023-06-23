package dame

type regeln struct {
	steinBewegenRichtungen                map[richtung]struct{}
	steinSchlagenRichtungenAnfang         map[richtung]struct{}
	steinSchlagenRichtungenWeiterschlagen map[richtung]struct{}

	dameBewegenRichtungen                map[richtung]struct{}
	dameSchlagenRichtungenAnfang         map[richtung]struct{}
	dameSchlagenRichtungenWeiterschlagen map[richtung]struct{}

	schlagZwang bool
}

func (r regeln) steinSchlagenRichtungen(weitschlagen bool) map[richtung]struct{} {
	if weitschlagen {
		return r.steinSchlagenRichtungenWeiterschlagen
	}
	return r.steinSchlagenRichtungenAnfang
}

func (r regeln) dameSchlagenRichtungen(weiterschlagen bool) map[richtung]struct{} {
	if weiterschlagen {
		return r.dameSchlagenRichtungenWeiterschlagen
	}
	return r.dameSchlagenRichtungenAnfang
}

var (
	internationaleRegeln = regeln{
		steinBewegenRichtungen:                richtungenDiagonalVorne,
		steinSchlagenRichtungenAnfang:         richtungenDiagonalVorne,
		steinSchlagenRichtungenWeiterschlagen: richtungenDiagonal,

		dameBewegenRichtungen:                richtungenDiagonal,
		dameSchlagenRichtungenAnfang:         richtungenDiagonal,
		dameSchlagenRichtungenWeiterschlagen: richtungenDiagonal,

		schlagZwang: true,
	}
	altdeutscheRegeln = regeln{
		steinBewegenRichtungen:                richtungenDiagonalVorne,
		steinSchlagenRichtungenAnfang:         richtungenSeiteDiagonalUndVorne,
		steinSchlagenRichtungenWeiterschlagen: richtungenSeiteDiagonalUndVorne,

		dameBewegenRichtungen:                richtungenAlle,
		dameSchlagenRichtungenAnfang:         richtungenAlle,
		dameSchlagenRichtungenWeiterschlagen: richtungenAlle,

		schlagZwang: true,
	}
)
