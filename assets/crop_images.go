//go:build ignore

package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func cropImage(img image.Image) (result image.Image, cropped bool) {
	const transparentAlpha = 0

	var obenWegschneiden int
obenWegschneiden:
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			_, _, _, alpha := img.At(x, y).RGBA()
			if alpha != transparentAlpha {
				obenWegschneiden = y
				break obenWegschneiden
			}
		}
	}

	var untenWegschneiden int
untenWegschneiden:
	for y := img.Bounds().Dy() - 1; y >= 0; y-- {
		for x := 0; x < img.Bounds().Dx(); x++ {
			_, _, _, alpha := img.At(x, y).RGBA()
			if alpha != transparentAlpha {
				untenWegschneiden = img.Bounds().Dy() - y - 1
				break untenWegschneiden
			}
		}
	}

	var linksWegschneiden int
linksWegschneiden:
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			_, _, _, alpha := img.At(x, y).RGBA()
			if alpha != transparentAlpha {
				linksWegschneiden = x
				break linksWegschneiden
			}
		}
	}

	var rechtsWegschneiden int
rechtsWegschneiden:
	for x := img.Bounds().Dx() - 1; x >= 0; x-- {
		for y := 0; y < img.Bounds().Dy(); y++ {
			_, _, _, alpha := img.At(x, y).RGBA()
			if alpha != transparentAlpha {
				rechtsWegschneiden = img.Bounds().Dx() - x - 1
				break rechtsWegschneiden
			}
		}
	}

	if linksWegschneiden == 0 && rechtsWegschneiden == 0 && obenWegschneiden == 0 && untenWegschneiden == 0 {
		return nil, false
	}

	croppedImg := image.NewRGBA(image.Rect(
		0,
		0,
		img.Bounds().Dx()-linksWegschneiden-rechtsWegschneiden,
		img.Bounds().Dy()-untenWegschneiden-obenWegschneiden,
	))
	for originalY := obenWegschneiden; originalY < img.Bounds().Dy()-untenWegschneiden; originalY++ {
		croppedY := originalY - obenWegschneiden
		for originalX := linksWegschneiden; originalX < img.Bounds().Dx()-rechtsWegschneiden; originalX++ {
			croppedX := originalX - linksWegschneiden

			croppedImg.Set(croppedX, croppedY, img.At(originalX, originalY))
		}
	}

	return croppedImg, true
}

func cropDirectory(dirPath string) error {
	var errs []error

	fs.WalkDir(os.DirFS(dirPath), ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		path = filepath.Join(dirPath, path)

		data, err := os.ReadFile(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to read file: %w", err))
			return nil
		}

		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to decode image: %w", err))
			return nil
		}

		croppedImg, cropped := cropImage(img)
		if !cropped {
			return nil
		}

		buffer := new(bytes.Buffer)
		err = png.Encode(buffer, croppedImg)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to encode image: %w", err))
			return nil
		}
		croppedImgData := buffer.Bytes()

		log.Println(path)

		err = os.WriteFile(path, croppedImgData, 0)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to write to file: %w", err))
			return nil
		}

		return nil
	})

	return errors.Join(errs...)
}

func main() {
	err := cropDirectory("./images/")
	if err != nil {
		log.Println(err)
	}
}
