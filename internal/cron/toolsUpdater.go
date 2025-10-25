package cron

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const destDir = "./data/bin"

var (
	httpClient = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: false,
		},
	}

	latestVersionCache struct {
		version   string
		timestamp time.Time
		mu        sync.Mutex
	}

	localVersionCache struct {
		version string
		mu      sync.Mutex
	}
)

// installUm 安装或更新 Um
func installUm() {
	version := getLatestVersion()
	if version == "" {
		logger.Error("⚠️未获取到最新版本，安装失败")
		return
	}

	if !shouldUpdate(version) {
		logger.Info("💡Um Installed successfully")
		return
	}

	url := buildUmURL(version)
	logger.Info("Downloading: " + url)

	if err := downloadAndExtract(url, destDir); err != nil {
		logger.Error("Installed failed: " + err.Error())
		return
	}

	// 更新本地版本缓存
	localVersionCache.mu.Lock()
	localVersionCache.version = version
	localVersionCache.mu.Unlock()

	logger.Info("💡Um Installed successfully")
}

// updateUm 检查并更新 Um
func updateUm() {
	version := getLatestVersion()
	if version == "" {
		return
	}

	if !shouldUpdate(version) {
		logger.Debug("💡Um 已是最新版本: " + version)
		return
	}

	logger.Info("💡检测到新版本: " + version + "，正在更新 Um...")
	installUm()
}

// shouldUpdate 判断是否需要更新
func shouldUpdate(latestVersion string) bool {
	binaryPath := filepath.Join(destDir, umBinaryName())

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return true
	}

	currentVersion := getLocalVersionCached(binaryPath)
	if currentVersion == "" {
		return true
	}

	return currentVersion != latestVersion
}

// getLocalVersionCached 获取本地 Um 版本，带缓存
func getLocalVersionCached(binaryPath string) string {
	localVersionCache.mu.Lock()
	defer localVersionCache.mu.Unlock()

	if localVersionCache.version != "" {
		return localVersionCache.version
	}

	out, err := exec.Command(binaryPath, "--version").Output()
	if err != nil {
		return ""
	}

	v := parseVersion(string(out))
	localVersionCache.version = v
	return v
}

// getLatestVersion 从 Releases 页面获取最新版本号，带缓存
func getLatestVersion() string {
	latestVersionCache.mu.Lock()
	defer latestVersionCache.mu.Unlock()

	if time.Since(latestVersionCache.timestamp) < 360*time.Minute && latestVersionCache.version != "" {
		return latestVersionCache.version
	}

	url := "https://git.um-react.app/um/cli/releases/"
	resp, err := httpClient.Get(url)
	if err != nil {
		logger.Error("⚠️获取 Releases 页面失败: " + err.Error())
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("⚠️获取 Releases 页面失败，状态码: " + fmt.Sprint(resp.StatusCode))
		return ""
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Error("⚠️解析页面 HTML 失败: " + err.Error())
		return ""
	}

	selection := doc.Find("#release-list > li").First()
	text := selection.Text()
	re := regexp.MustCompile(`v\d+\.\d+\.\d+`)
	match := re.FindString(text)
	if match == "" {
		logger.Error("⚠️未匹配到版本号")
		return ""
	}

	latestVersionCache.version = match
	latestVersionCache.timestamp = time.Now()
	return match
}

// buildUmURL 构建下载 URL
func buildUmURL(version string) string {
	return fmt.Sprintf(
		"https://git.um-react.app/um/cli/releases/download/%s/um-%s-amd64-%s.tar.gz",
		version, runtime.GOOS, version,
	)
}

// downloadAndExtract 下载并解压 Um
func downloadAndExtract(url, destDir string) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	return extractUmBinary(resp.Body, destDir)
}

// extractUmBinary 从 tar.gz 流中提取 Um
func extractUmBinary(r io.Reader, destDir string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("gzip 解压失败: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	target := umBinaryName()

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar 读取失败: %w", err)
		}

		if filepath.Base(header.Name) == target {
			return saveFile(tr, filepath.Join(destDir, target))
		}
	}

	return fmt.Errorf("未找到二进制文件: " + target)
}

// umBinaryName 返回系统对应的 Um 文件名
func umBinaryName() string {
	if runtime.GOOS == "windows" {
		return "um.exe"
	}
	return "um"
}

// saveFile 保存文件并设置可执行权限，使用缓冲写入
func saveFile(r io.Reader, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer out.Close()

	buf := bufio.NewWriter(out)
	if _, err := io.Copy(buf, r); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}
	if err := buf.Flush(); err != nil {
		return fmt.Errorf("刷新缓冲失败: %w", err)
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(path, 0755); err != nil {
			return fmt.Errorf("设置可执行权限失败: %w", err)
		}
	}
	return nil
}

// parseVersion 从 Um --version 输出中提取版本号
func parseVersion(output string) string {
	re := regexp.MustCompile(`v\d+\.\d+\.\d+`)
	return re.FindString(output)
}
