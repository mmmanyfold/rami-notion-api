package notion

import (
	"image"
	_ "image/jpeg"
	"os"
)

func imageSize(imgPath string) (w int, h int, error error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return w, h, err
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return w, h, err
	}
	return img.Width, img.Height, nil
}
