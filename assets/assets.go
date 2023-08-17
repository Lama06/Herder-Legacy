package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

//go:generate go run crop_images.go

//go:embed fonts/*
var fontsFS embed.FS

func LoadFont(filepath string) ([]byte, error) {
	return fontsFS.ReadFile("fonts/" + filepath)
}

func RequireFont(filepath string) []byte {
	font, err := LoadFont(filepath)
	if err != nil {
		panic(err)
	}
	return font
}

var (
	//go:embed sounds/*
	soundsFS   embed.FS
	soundCache = make(map[string]*audio.Player)
)

func LoadSound(filepath string) (*audio.Player, error) {
	cachedSound, ok := soundCache[filepath]
	if ok {
		return cachedSound, nil
	}

	soundData, err := soundsFS.ReadFile("sounds/" + filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read sound data: %w", err)
	}
	context := audio.CurrentContext()
	stream, err := mp3.DecodeWithSampleRate(context.SampleRate(), bytes.NewReader(soundData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode sound: %w", err)
	}
	player, err := context.NewPlayer(stream)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio player: %w", err)
	}
	soundCache[filepath] = player
	return player, nil
}

func RequireSound(filepath string) *audio.Player {
	sound, err := LoadSound(filepath)
	if err != nil {
		panic(err)
	}
	return sound
}

var (
	//go:embed images/*
	imagesFS   embed.FS
	imageCache = make(map[string]*ebiten.Image)
)

func LoadImage(filepath string) (*ebiten.Image, error) {
	cachedImage, ok := imageCache[filepath]
	if ok {
		return cachedImage, nil
	}

	imageData, err := imagesFS.ReadFile("images/" + filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	ebitenImg := ebiten.NewImageFromImage(img)
	imageCache[filepath] = ebitenImg
	return ebitenImg, nil
}

func RequireImage(filepath string) *ebiten.Image {
	img, err := LoadImage(filepath)
	if err != nil {
		panic(err)
	}
	return img
}
