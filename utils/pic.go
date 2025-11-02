package utils

import (
	"fmt"
	"io"
	"net/http"
)

// FetchImage 下载图片并返回字节数组
func FetchImage(url string) ([]byte, error) {
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

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image body: %w", err)
	}

	return imageData, nil
}
