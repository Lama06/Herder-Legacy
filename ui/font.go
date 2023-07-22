package ui

import (
	_ "embed"
	"fmt"

	"github.com/Lama06/Herder-Legacy/assets"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	normalFontFace font.Face
	titleFontFace  font.Face
)

func init() {
	const dpi = 72

	fontData := assets.RequireFont("roboto.ttf")

	robotoFont, err := opentype.Parse(fontData)
	if err != nil {
		panic(fmt.Errorf("failed to parse roboto font: %w", err))
	}

	normalFontFace, err = opentype.NewFace(robotoFont, &opentype.FaceOptions{
		Size: 30,
		DPI:  dpi,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create font face: %w", err))
	}

	titleFontFace, err = opentype.NewFace(robotoFont, &opentype.FaceOptions{
		Size: 60,
		DPI:  dpi,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create font face: %w", err))
	}
}
