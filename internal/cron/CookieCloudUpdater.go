package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nichuanfang/gymdl/config"
)

type CookieCloudBody struct {
	Uuid      string `json:"uuid,omitempty"`
	Encrypted string `json:"encrypted,omitempty"`
}

type decryptResult struct {
	Key       string
	Decrypted []byte
}

func syncCookieCloud(c *config.CookieCloudConfig) {
	if c == nil || len(c.CookieCloudKEY) == 0 {
		logger.Error("【CookieCloud】config or keys are empty")
		return
	}
	if c.CookieCloudUrl == "" || c.CookieCloudUUID == "" {
		logger.Error("【CookieCloud】CookieCloudUrl and CookieCloudUUID must be set")
		return
	}
	if c.CookieFile == "" {
		logger.Error("【CookieCloud】CookieFile is empty")
		return
	}

	dstDir := strings.TrimRight(c.CookieFilePath, "/")
	if dstDir == "" {
		dstDir = "."
	}
	dstPath := fmt.Sprintf("%s/%s", dstDir, c.CookieFile)

	// 创建目录（只需尝试一次）
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		logger.Error(fmt.Sprintf("【CookieCloud】failed to create dir %s: %v", dstDir, err))
		return
	}

	// HTTP 请求
	client := &http.Client{Timeout: 10 * time.Second}
	getUrl := fmt.Sprintf("%s/get/%s", strings.TrimSuffix(c.CookieCloudUrl, "/"), c.CookieCloudUUID)
	resp, err := client.Get(getUrl)
	if err != nil {
		logger.Warn(fmt.Sprintf("【CookieCloud】failed to request server: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Warn(fmt.Sprintf("【CookieCloud】server returned status %d", resp.StatusCode))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Warn(fmt.Sprintf("【CookieCloud】failed to read server response: %v", err))
		return
	}

	var data CookieCloudBody
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Warn(fmt.Sprintf("【CookieCloud】failed to parse json: %v", err))
		return
	}
	if data.Encrypted == "" {
		logger.Warn("【CookieCloud】encrypted data is empty")
		return
	}

	// 并发尝试所有 key 解密，只取第一个成功的结果
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultCh := make(chan decryptResult, 1)
	var wg sync.WaitGroup

	for _, key := range c.CookieCloudKEY {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				keyPassword := Md5String(c.CookieCloudUUID, "-", key)
				if len(keyPassword) < 16 {
					return
				}
				keyPassword = keyPassword[:16]

				dec, err := DecryptCryptoJsAesMsg(keyPassword, data.Encrypted)
				if err == nil && len(dec) > 0 {
					select {
					case resultCh <- decryptResult{Key: key, Decrypted: dec}:
						cancel() // 成功后取消其他 goroutine
					default:
					}
				}
			}
		}(key)
	}

	// 等待 goroutine 完成
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// 取第一个成功结果
	var result *decryptResult
	for r := range resultCh {
		result = &r
		break
	}

	if result == nil || len(result.Decrypted) == 0 {
		logger.Error("【CookieCloud】all keys failed to decrypt cookie")
		return
	}

	if err := os.WriteFile(dstPath, result.Decrypted, 0o600); err != nil {
		logger.Error(fmt.Sprintf("【CookieCloud】failed to write cookie file %s: %v", dstPath, err))
		return
	}

	logger.Info("【CookieCloud】cookie updated successfully")
}
