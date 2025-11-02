package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"time"

	"github.com/disintegration/imaging"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
	},
}

// FetchImage 下载图片
func FetchImage(url string) ([]byte, error) {
	// 发起请求
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// 限制最大读取大小，防止内存攻击
	const maxImageSize = 10 << 20
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxImageSize))
	if err != nil {
		return nil, err
	}

	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	//缩放到目标尺寸 320x320
	resized := imaging.Resize(img, 320, 320, imaging.Lanczos)

	// 编码为 JPEG 格式，压缩质量设为85以平衡性能与质量
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resized, &jpeg.Options{Quality: 85})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
