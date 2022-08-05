package image

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"image"
	_ "image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
)

func Size(imgPath string) (w int, h int, error error) {
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

func Download(url string) (path string, err error) {
	log.Println("downloading: ", url)
	uid := uuid.New()
	path = "/tmp/" + uid.String()

	response, err := http.Get(url)
	if err != nil {
		return path, err
	}
	defer response.Body.Close()

	// extension detection
	ext, err := Extension(response.Body)
	if err != nil {
		return path, err
	}

	// append extension
	path += ext

	file, err := os.Create(path)
	if err != nil {
		return path, err
	}
	defer file.Close()

	log.Println("saving to: ", path)
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return path, err
	}

	return path, nil
}

func Extension(f io.Reader) (extension string, err error) {
	mtype, err := mimetype.DetectReader(f)
	if err != nil {
		return extension, err
	}
	for _, format := range []string{"image/png", "image/jpeg"} {
		log.Printf("is format: %s = %t \n", format, mtype.Is(format))
	}

	return mtype.Extension(), err
}
