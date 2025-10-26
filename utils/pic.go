package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/nfnt/resize"
	"github.com/sirupsen/logrus"
)

// 缩放图片到 320x320px (黑底填充)
func ResizeImg(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	buffer, err := io.ReadAll(bufio.NewReader(file))
	if err != nil {
		return "", err
	}

	var img image.Image

	img, err = jpeg.Decode(bytes.NewReader(buffer))
	if err != nil {
		img, err = png.Decode(bytes.NewReader(buffer))
		if err != nil {
			return "", fmt.Errorf("Image decode error  %s", filePath)
		}
	}
	err = file.Close()
	if err != nil {
		return "", err
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	widthNew := 320
	heightNew := 320

	var m image.Image
	if width/height >= widthNew/heightNew {
		m = resize.Resize(uint(widthNew), uint(height)*uint(widthNew)/uint(width), img, resize.Lanczos3)
	} else {
		m = resize.Resize(uint(width*heightNew/height), uint(heightNew), img, resize.Lanczos3)
	}

	newImag := image.NewNRGBA(image.Rect(0, 0, 320, 320))
	if m.Bounds().Dx() > m.Bounds().Dy() {
		draw.Draw(newImag, image.Rectangle{
			Min: image.Point{Y: (320 - m.Bounds().Dy()) / 2},
			Max: image.Point{X: 320, Y: 320},
		}, m, m.Bounds().Min, draw.Src)
	} else {
		draw.Draw(newImag, image.Rectangle{
			Min: image.Point{X: (320 - m.Bounds().Dx()) / 2},
			Max: image.Point{X: 320, Y: 320},
		}, m, m.Bounds().Min, draw.Src)
	}

	out, err := os.Create(filePath + ".resize.jpg")
	if err != nil {
		return "", fmt.Errorf("Create image file error  %s", err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			logrus.Errorln(err)
		}
	}(out)

	err = jpeg.Encode(out, newImag, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	return filePath + ".resize.jpg", nil
}

// FetchImage 从网络下载图片并解码成 image.Image
func FetchImage(url string) (*image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("下载图片失败: %w", err)
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("图片解码失败: %w", err)
	}

	return &img, nil
}
