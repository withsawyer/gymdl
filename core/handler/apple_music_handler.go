package handler

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/constants"
	"github.com/nichuanfang/gymdl/utils"
)

// AppleMusicHandler 实现 MusicHandler 接口
type AppleMusicHandler struct{}

// DownloadMusic 使用 gamdl 下载 Apple Music 内容
func (am *AppleMusicHandler) DownloadMusic(url string, cfg *config.Config) error {
	start := time.Now()

	outputDir := constants.AppleMusicTempDir
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}
	//调用指令模块
	cmd := am.DownloadCommand(cfg)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	utils.InfoWithFormat("🎵 开始下载 Apple Music 内容: %s\n", url)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gamdl 下载失败: %w", err)
	}

	utils.InfoWithFormat("✅ 下载完成（耗时 %v）\n", time.Since(start).Truncate(time.Millisecond))
	return nil
}

// DownloadCommand 返回可复用的 exec.Cmd（便于调度器统一执行）
func (am *AppleMusicHandler) DownloadCommand(cfg *config.Config) *exec.Cmd {
	cookiePath := filepath.Join(cfg.CookieCloud.CookieFilePath, cfg.CookieCloud.CookieFile)
	return exec.Command(
		"gamdl",
		"--cookies-path", cookiePath,
		"--download-mode", "nm3u8dlre",
		"--output-path", constants.AppleMusicTempDir,
	)
}

// BeforeTidy 音乐整理之前的处理
func (am *AppleMusicHandler) BeforeTidy(cfg *config.Config) error {
	return nil
}

// NeedRemoveDRM 判断是否需要去除 DRM
func (am *AppleMusicHandler) NeedRemoveDRM(cfg *config.Config) bool {
	return false
}

// DRMRemove 去除 Apple Music DRM（通过 `ffmpeg` 转码）
func (am *AppleMusicHandler) DRMRemove(cfg *config.Config) error {
	srcDir := cfg.MusicTidy.DistDir
	files, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("读取输出目录失败: %w", err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(f.Name()), ".m4p") {
			continue
		}

		src := filepath.Join(srcDir, f.Name())
		dst := strings.TrimSuffix(src, filepath.Ext(src)) + "_drmfree.m4a"

		cmd := exec.Command("ffmpeg", "-y", "-i", src, "-c", "copy", dst)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Printf("🔓 移除 DRM: %s\n", f.Name())
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("移除 DRM 失败: %w", err)
		}

		// 删除原始加密文件
		_ = os.Remove(src)
	}
	return nil
}

// sanitizeFileName 处理非法字符，确保跨平台兼容性
func sanitizeFileName(name string) string {
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, c := range invalidChars {
		name = strings.ReplaceAll(name, c, "_")
	}
	return name
}

// TidyMusic 将下载的 Apple Music 文件整理到最终输出目录
func (am *AppleMusicHandler) TidyMusic(cfg *config.Config) error {
	srcDir := constants.AppleMusicTempDir
	if cfg.MusicTidy.Mode == 1 {
		dstDir := cfg.MusicTidy.DistDir
		if dstDir == "" {
			return errors.New("未配置输出目录")
		}

		files, err := os.ReadDir(srcDir)
		if err != nil {
			return fmt.Errorf("读取临时目录失败: %w", err)
		}

		if len(files) == 0 {
			return fmt.Errorf("未找到下载的音乐文件")
		}

		if err := os.MkdirAll(dstDir, 0755); err != nil {
			return fmt.Errorf("创建输出目录失败: %w", err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			src := filepath.Join(srcDir, file.Name())
			cleanName := sanitizeFileName(file.Name())
			dst := filepath.Join(dstDir, cleanName)

			if err := os.Rename(src, dst); err != nil {
				return fmt.Errorf("移动文件失败 %s → %s: %w", src, dst, err)
			}
			fmt.Printf("📦 已整理: %s\n", dst)
		}
	} else {
		//webdav上传
	}
	return nil
}
