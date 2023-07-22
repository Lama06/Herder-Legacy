package herder

import (
	"bytes"
	"embed"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

var imageCache = make(map[string]*ebiten.Image)

func loadImage(name string) (*ebiten.Image, error) {
	cachedImage, ok := imageCache[name]
	if ok {
		return cachedImage, nil
	}

	imageData, err := assets.ReadFile("assets/" + name + ".png")
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	ebitenImg := ebiten.NewImageFromImage(img)
	imageCache[name] = ebitenImg
	return ebitenImg, nil
}

func requireImage(name string) *ebiten.Image {
	img, err := loadImage(name)
	if err != nil {
		panic(err)
	}
	return img
}
