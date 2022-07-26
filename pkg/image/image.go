package image

import (
	"image"
	_ "image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func Size(imgPath string) (w uint64, h uint64, error error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return w, h, err
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return w, h, err
	}

	return uint64(img.Width), uint64(img.Height), nil
}

func Download(url string) (path string, err error) {
	log.Println("downloading: ", url)
	response, err := http.Get(url)
	if err != nil {
		return path, err
	}
	defer response.Body.Close()

	file, err := ioutil.TempFile("/tmp/", "img-*.png")
	if err != nil {
		return path, err
	}
	defer file.Close()

	log.Println("saving to: ", file.Name())
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return path, err
	}

	return file.Name(), nil
}
