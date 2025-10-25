package core

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
	"github.com/nichuanfang/gymdl/utils"
	"go.uber.org/zap"
)

type CookieCloud struct {
	Config *config.CookieCloudConfig
	Client *http.Client
}

type CookieCloudBody struct {
	Uuid      string `json:"uuid,omitempty"`
	Encrypted string `json:"encrypted,omitempty"`
}

type decryptResult struct {
	Key       string
	Decrypted []byte
}

var (
	// GlobalCookieCloud 全局单例对象
	GlobalCookieCloud *CookieCloud
	logger            *zap.Logger
)

// InitCookieCloud 初始化全局 CookieCloud，只会执行一次
func InitCookieCloud(cfg *config.CookieCloudConfig) {
	logger = utils.Logger()
	GlobalCookieCloud = &CookieCloud{
		Config: cfg,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

// CheckConnection 检查 CookieCloud 服务是否可用
func (cc *CookieCloud) CheckConnection() bool {
	if cc.Config == nil || cc.Config.CookieCloudUrl == "" || cc.Config.CookieCloudUUID == "" {
		logger.Error("⚠️CookieCloud config, URL or UUID is empty")
		return false
	}

	getUrl := fmt.Sprintf("%s/get/%s", strings.TrimSuffix(cc.Config.CookieCloudUrl, "/"), cc.Config.CookieCloudUUID)
	resp, err := cc.Client.Get(getUrl)
	if err != nil {
		logger.Warn(fmt.Sprintf("⚠️CookieCloud failed to request server: %v", err))
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Warn(fmt.Sprintf("⚠️CookieCloud server returned status %d", resp.StatusCode))
		return false
	}
	return true
}

// Sync 同步 CookieCloud 的 cookie 到本地
func (cc *CookieCloud) Sync() {
	if cc.Config == nil || len(cc.Config.CookieCloudKEY) == 0 {
		logger.Error("⚠️CookieCloud config or keys are empty")
		return
	}
	if cc.Config.CookieFile == "" {
		logger.Error("⚠️CookieCloud CookieFile is empty")
		return
	}

	dstDir := strings.TrimRight(cc.Config.CookieFilePath, "/")
	if dstDir == "" {
		dstDir = "."
	}
	dstPath := fmt.Sprintf("%s/%s", dstDir, cc.Config.CookieFile)

	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		logger.Error(fmt.Sprintf("⚠️CookieCloud failed to create dir %s: %v", dstDir, err))
		return
	}

	data, err := cc.fetchEncryptedData()
	if err != nil {
		logger.Warn(fmt.Sprintf("⚠️CookieCloud %v", err))
		return
	}

	decrypted := cc.decryptData(data)
	if decrypted == nil {
		logger.Error("⚠️CookieCloud all keys failed to decrypt cookie")
		return
	}

	netscapeData, err := ConvertToNetscapeFormat(decrypted)
	if err != nil {
		logger.Error(fmt.Sprintf("⚠️CookieCloud failed to convert to Netscape format: %v", err))
		return
	}

	if err := os.WriteFile(dstPath, netscapeData, 0o600); err != nil {
		logger.Error(fmt.Sprintf("⚠️CookieCloud failed to write cookie file %s: %v", dstPath, err))
		return
	}

	logger.Info("💡CookieCloud cookie updated successfully")
}

// fetchEncryptedData 获取加密的 cookie 数据
func (cc *CookieCloud) fetchEncryptedData() (string, error) {
	getUrl := fmt.Sprintf("%s/get/%s", strings.TrimSuffix(cc.Config.CookieCloudUrl, "/"), cc.Config.CookieCloudUUID)
	resp, err := cc.Client.Get(getUrl)
	if err != nil {
		return "", fmt.Errorf("failed to request server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read server response: %v", err)
	}

	var data CookieCloudBody
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("failed to parse json: %v", err)
	}
	if data.Encrypted == "" {
		return "", fmt.Errorf("encrypted data is empty")
	}
	return data.Encrypted, nil
}

// decryptData 使用配置的 key 解密数据
func (cc *CookieCloud) decryptData(encrypted string) []byte {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultCh := make(chan decryptResult, 1)
	var wg sync.WaitGroup
	var once sync.Once

	for _, key := range cc.Config.CookieCloudKEY {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}

			keyPassword := Md5String(cc.Config.CookieCloudUUID, "-", key)
			if len(keyPassword) < 16 {
				return
			}
			keyPassword = keyPassword[:16]

			dec, err := DecryptCryptoJsAesMsg(keyPassword, encrypted)
			if err == nil && len(dec) > 0 {
				once.Do(func() {
					resultCh <- decryptResult{Key: key, Decrypted: dec}
					cancel()
				})
			}
		}(key)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var result *decryptResult
	for r := range resultCh {
		result = &r
		break
	}

	if result == nil {
		return nil
	}
	return result.Decrypted
}

// ConvertToNetscapeFormat 将解密后的 cookie 转换为 Netscape 格式
func ConvertToNetscapeFormat(decryptedData []byte) ([]byte, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(decryptedData, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse decrypted data: %w", err)
	}

	cookieDataRaw, ok := raw["cookie_data"]
	if !ok {
		return nil, fmt.Errorf("missing 'cookie_data' key in decrypted data")
	}

	cookieData, ok := cookieDataRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'cookie_data' is not a map")
	}

	var b strings.Builder
	b.WriteString("# Netscape HTTP Cookie File\n")
	b.WriteString("# This file was generated by ConvertToNetscapeFormat\n")
	b.WriteString("# Fields: domain\tinclude_subdomains\tpath\tsecure\texpiry\tname\tvalue\n\n")

	for domain, v := range cookieData {
		cookies, ok := v.([]interface{})
		if !ok {
			continue
		}
		for _, c := range cookies {
			cm, ok := c.(map[string]interface{})
			if !ok || safeString(cm["name"]) == "" {
				continue
			}
			line := cookieToLine(domain, cm)
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}
	return []byte(b.String()), nil
}

func cookieToLine(domain string, c map[string]interface{}) string {
	includeSubdomains := "FALSE"
	if strings.HasPrefix(domain, ".") {
		includeSubdomains = "TRUE"
	}
	path := safeString(c["path"])
	secure := boolValue(c["secure"])
	expiry := int64Value(c["expiry"])
	name := safeString(c["name"])
	value := safeString(c["value"])

	return fmt.Sprintf("%s\t%s\t%s\t%t\t%d\t%s\t%s",
		domain,
		includeSubdomains,
		path,
		secure,
		expiry,
		name,
		value,
	)
}

func safeString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func boolValue(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

func int64Value(v interface{}) int64 {
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int:
		return int64(val)
	case int64:
		return val
	default:
		return 0
	}
}
