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
	"time"

	"github.com/disintegration/imaging"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
	},
}

// FetchImage 下载图片并返回字节数组
func FetchImage(url string) ([]byte, error) {
	if url == "" {
		return nil, fmt.Errorf("empty URL")
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image from %s: status code %d", url, resp.StatusCode)
	}

	img, format, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image from %s: %w", url, err)
	}

	// 直接缩放到目标尺寸
	resized := imaging.Fit(img, 320, 320, imaging.Lanczos)

	// 创建黑底并居中
	dst := imaging.New(320, 320, color.NRGBA{0, 0, 0, 255})
	dst = imaging.PasteCenter(dst, resized)

	var buf bytes.Buffer
	if strings.EqualFold(format, "png") {
		err = png.Encode(&buf, dst)
	} else {
		err = jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 90})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to encode image from %s: %w", url, err)
	}

	return buf.Bytes(), nil
}
