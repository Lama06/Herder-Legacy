package quiz

import (
	"bytes"
	_ "embed"

	"github.com/Lama06/Herder-Legacy/herderlegacy"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

var (
	//go:embed assets/richtig.mp3
	richtigSoundData []byte
	//go:embed assets/falsch.mp3
	falschSoundData []byte

	richtigSound *audio.Player
	falschSound  *audio.Player
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
	richtigSound = loadSound(herderLegacy, richtigSoundData)
	falschSound = loadSound(herderLegacy, falschSoundData)
}
