package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
)

// FetchImage 下载图片并返回字节数组
func FetchAndResizeImage(url string) ([]byte, error) {
	if url == "" {
		return nil, fmt.Errorf("empty URL")
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: status code %d", resp.StatusCode)
	}

	img, format, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// 使用高性能的缩放算法
	resized := imaging.Fit(img, 320, 320, imaging.Lanczos)

	// 创建黑底 320x320
	dst := imaging.New(320, 320, color.NRGBA{0, 0, 0, 255})

	// 居中绘制缩放后的图片
	dst = imaging.PasteCenter(dst, resized)

	// 编码输出为 JPEG
	var buf bytes.Buffer
	switch strings.ToLower(format) {
	case "png":
		err = png.Encode(&buf, dst)
	default:
		err = jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 90})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to encode resized image: %w", err)
	}

	return buf.Bytes(), nil
}
