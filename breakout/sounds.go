package breakout

import (
	"bytes"
	_ "embed"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

var (
	//go:embed sounds/stein.mp3
	steinSoundData []byte
	//go:embed sounds/plattform.mp3
	plattformSoundData []byte
	//go:embed sounds/upgrade.mp3
	upgradeSoundData []byte
	//go:embed sounds/wand.mp3
	wandSoundData []byte

	steinSound     *audio.Player
	plattformSound *audio.Player
	upgradeSound   *audio.Player
	wandSound      *audio.Player
)

func loadSound(herderLegacy herderlegacy.HerderLegacy, soundData []byte) *audio.Player {
	context := herderLegacy.AudioContext()
	stream, err := mp3.DecodeWithSampleRate(context.SampleRate(), bytes.NewReader(soundData))
	if err != nil {
		panic(err)
	}
	player, err := context.NewPlayer(stream)
	if err != nil {
		panic(err)
	}
	return player
}

func loadSounds(herderLegacy herderlegacy.HerderLegacy) {
	steinSound = loadSound(herderLegacy, steinSoundData)
	plattformSound = loadSound(herderLegacy, plattformSoundData)
	upgradeSound = loadSound(herderLegacy, upgradeSoundData)
	wandSound = loadSound(herderLegacy, wandSoundData)
}
