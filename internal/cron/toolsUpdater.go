package cron

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const destDir = "./data/bin"

// installUm 安装或更新 Um
func installUm() {
	version := getLatestVersion()

	if !shouldUpdate(version) {
		logger.Info("【Um】已安装，版本：" + version)
		return
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		logger.Error("创建目录失败: " + err.Error())
		return
	}

	url := buildUmURL(version)
	logger.Info("Downloading: " + url)

	if err := downloadAndExtract(url, destDir); err != nil {
		logger.Error("安装失败: " + err.Error())
		return
	}

	logger.Info("【Um】安装完成，版本：" + version)
}

// updateUm 检查并更新 Um
func updateUm() {
	version := getLatestVersion()
	if !shouldUpdate(version) {
		logger.Debug(fmt.Sprintf("Um 已是最新版本: %s", version))
		return
	}

	logger.Info(fmt.Sprintf("检测到新版本: %s，正在更新 Um...", version))
	installUm()
}

// shouldUpdate 判断是否需要更新
func shouldUpdate(latestVersion string) bool {
	binaryPath := filepath.Join(destDir, umBinaryName())

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return true
	}

	currentVersion, err := getLocalVersion(binaryPath)
	if err != nil {
		return true
	}

	return currentVersion != latestVersion
}

// getLocalVersion 获取本地 Um 版本
func getLocalVersion(binaryPath string) (string, error) {
	out, err := exec.Command(binaryPath, "--version").Output()
	if err != nil {
		return "", err
	}
	return parseVersion(string(out)), nil
}

// getLatestVersion 返回最新版本
func getLatestVersion() string {
	// TODO: 可以改为从 GitLab API 获取最新 release
	return "v0.2.17"
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
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
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

		if strings.HasSuffix(header.Name, target) {
			return saveFile(tr, filepath.Join(destDir, target))
		}
	}

	return fmt.Errorf("未找到二进制文件: %s", target)
}

// umBinaryName 返回系统对应的 Um 文件名
func umBinaryName() string {
	if runtime.GOOS == "windows" {
		return "um.exe"
	}
	return "um"
}

// saveFile 保存文件并设置可执行权限
func saveFile(r io.Reader, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, r); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
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
	// 输出示例: "Unlock Music CLI version v0.2.17 (go1.25.1,darwin/amd64)"
	fields := strings.Fields(output)
	for i, f := range fields {
		if f == "version" && i+1 < len(fields) {
			return strings.TrimSpace(fields[i+1])
		}
	}
	return ""
}
