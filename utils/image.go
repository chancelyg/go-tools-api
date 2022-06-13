package utils

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
)

func imageClipCompliant(img image.Image, width int, height int) bool {
	return width > img.Bounds().Max.X || height > img.Bounds().Max.Y || width < 0 || height < 0
}

func ClipImage(filePath string, width int, height int) (image.Image, error) {
	fileExtension := strings.Split(filePath, ".")[len(strings.Split(filePath, "."))-1]
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	// decode图片
	switch fileExtension {
	case "png":
		img, _ := png.Decode(f)
		if imageClipCompliant(img, width, height) {
			return nil, errors.New("not support size")
		}
		return img.(*image.RGBA).SubImage(image.Rect(0, 0, width, height)).(*image.RGBA), nil
	case "jpeg":
	case "jpg":
		img, _ := jpeg.Decode(f)
		if imageClipCompliant(img, width, height) {
			return nil, errors.New("not support size")
		}
		return img.(*image.YCbCr).SubImage(image.Rect(0, 0, width, height)).(*image.YCbCr), nil
	}
	panic("Not support format")

}

func ClipImageWithTmpFile(filePath string, width int, height int) (string, error) {
	image, err := ClipImage(filePath, width, height)
	if err != nil {
		return "", err
	}

	file, err := ioutil.TempFile("", "IMAGE_*_.png")
	if err != nil {
		return "", err
	}

	err = png.Encode(file, image)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}
