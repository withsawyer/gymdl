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
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/core/constants"
	"github.com/nichuanfang/gymdl/utils"
)

// AppleMusicHandler 实现 MusicHandler 接口
type AppleMusicHandler struct{}

// DownloadMusic 使用 gamdl 下载 Apple Music 内容
func (am *AppleMusicHandler) DownloadMusic(url string, cfg *config.Config) error {
	start := time.Now()

	if err := os.MkdirAll(constants.AppleMusicTempDir, 0755); err != nil {
		utils.ErrorWithFormat("❌ 创建输出目录失败: %v", err)
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	cmd := am.DownloadCommand(cfg, url)
	utils.DebugWithFormat("[AppleMusicHandler] 🎵 开始下载: %s", url)
	utils.DebugWithFormat("[AppleMusicHandler] 执行命令: %s", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.ErrorWithFormat("❌ gamdl 下载失败: %v\n输出:\n%s", err, string(output))
		return fmt.Errorf("gamdl 下载失败: %w", err)
	}

	utils.DebugWithFormat("[AppleMusicHandler] 下载输出:\n%s", string(output))
	utils.InfoWithFormat("✅ 下载完成（耗时 %v）", time.Since(start).Truncate(time.Millisecond))
	return nil
}

// DownloadCommand 构建下载命令
func (am *AppleMusicHandler) DownloadCommand(cfg *config.Config, url string) *exec.Cmd {
	cookiePath := filepath.Join(cfg.CookieCloud.CookieFilePath, cfg.CookieCloud.CookieFile)
	return exec.Command(
		"gamdl",
		"--cookies-path", cookiePath,
		"--download-mode", "nm3u8dlre",
		"--output-path", constants.BaseTempDir,
		"--album-folder-template", "AppleMusic",
		"--compilation-folder-template", "AppleMusic",
		"--no-album-folder-template", "AppleMusic",
		"--single-disc-folder-template", "{title}",
		"--multi-disc-folder-template", "{title}",
		"--no-synced-lyrics",
		url,
	)
}

// BeforeTidy 清洗/解锁音乐文件
func (am *AppleMusicHandler) BeforeTidy(cfg *config.Config) error {
	if am.NeedRemoveDRM(cfg) {
		if err := am.DRMRemove(cfg); err != nil {
			utils.ErrorWithFormat("❌ DRM 移除失败: %v", err)
			return fmt.Errorf("DRM 移除失败: %w", err)
		}
		utils.InfoWithFormat("🔓 DRM 已移除")
	}
	return nil
}

// NeedRemoveDRM 判断是否需要去除 DRM
func (am *AppleMusicHandler) NeedRemoveDRM(cfg *config.Config) bool {
	return false
}

// DRMRemove 去除 Apple Music DRM（通过 ffmpeg 转码）
func (am *AppleMusicHandler) DRMRemove(cfg *config.Config) error {
	// TODO: 后续实现 DRM 移除逻辑
	return nil
}

// TidyMusic 将下载的 Apple Music 文件整理到最终输出目录
func (am *AppleMusicHandler) TidyMusic(cfg *config.Config, webdav *core.WebDAV) error {
	files, err := os.ReadDir(constants.AppleMusicTempDir)
	if err != nil {
		utils.ErrorWithFormat("❌ 读取临时目录失败: %v", err)
		return fmt.Errorf("读取临时目录失败: %w", err)
	}
	if len(files) == 0 {
		utils.WarnWithFormat("⚠️ 未找到待整理的音乐文件")
		return errors.New("未找到待整理的音乐文件")
	}

	switch cfg.MusicTidy.Mode {
	case 1:
		return am.tidyToLocal(cfg, files)
	case 2:
		return am.tidyToWebDAV(cfg, files, webdav)
	default:
		utils.ErrorWithFormat("❌ 未知整理模式: %d", cfg.MusicTidy.Mode)
		return fmt.Errorf("未知整理模式: %d", cfg.MusicTidy.Mode)
	}
}

// tidyToLocal 将音乐移动到本地目标目录
func (am *AppleMusicHandler) tidyToLocal(cfg *config.Config, files []os.DirEntry) error {
	dstDir := cfg.MusicTidy.DistDir
	if dstDir == "" {
		utils.WarnWithFormat("⚠️ 未配置输出目录")
		return errors.New("未配置输出目录")
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		utils.ErrorWithFormat("❌ 创建输出目录失败: %v", err)
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	for _, file := range files {
		if !utils.FilterMusicFile(file, am.EncryptedExts(), am.DecryptedExts()) {
			utils.DebugWithFormat("[AppleMusicHandler] 跳过非音乐文件: %s", file.Name())
			continue
		}
		srcPath := filepath.Join(constants.AppleMusicTempDir, file.Name())
		dstPath := filepath.Join(dstDir, utils.SanitizeFileName(file.Name()))
		if err := os.Rename(srcPath, dstPath); err != nil {
			utils.ErrorWithFormat("❌ 移动文件失败 %s → %s: %v", srcPath, dstPath, err)
			return fmt.Errorf("移动文件失败 %s → %s: %w", srcPath, dstPath, err)
		}
		utils.InfoWithFormat("📦 已整理: %s", dstPath)
	}
	return nil
}

// tidyToWebDAV 将音乐上传到 WebDAV
func (am *AppleMusicHandler) tidyToWebDAV(cfg *config.Config, files []os.DirEntry, webdav *core.WebDAV) error {
	if webdav == nil {
		utils.ErrorWithFormat("❌ WebDAV 未初始化")
		return errors.New("WebDAV 未初始化")
	}

	for _, file := range files {
		if !utils.FilterMusicFile(file, am.EncryptedExts(), am.DecryptedExts()) {
			utils.DebugWithFormat("[AppleMusicHandler] 跳过非音乐文件: %s", file.Name())
			continue
		}

		filePath := filepath.Join(constants.AppleMusicTempDir, file.Name())
		if err := webdav.Upload(filePath); err != nil {
			utils.ErrorWithFormat("❌ 上传文件失败 %s: %v", file.Name(), err)
			return fmt.Errorf("上传文件失败 %s: %w", file.Name(), err)
		}

		utils.InfoWithFormat("☁️ 已上传: %s", file.Name())

		ext := strings.ToLower(filepath.Ext(file.Name()))
		if utils.Contains(am.DecryptedExts(), ext) {
			if err := os.Remove(filePath); err != nil {
				utils.WarnWithFormat("⚠️ 删除临时文件失败: %s (%v)", filePath, err)
			} else {
				utils.DebugWithFormat("🧹 已删除临时文件: %s", filePath)
			}
		}
	}
	return nil
}

// EncryptedExts 返回加密后缀
func (am *AppleMusicHandler) EncryptedExts() []string {
	return []string{".m4p"}
}

// DecryptedExts 返回非加密后缀
func (am *AppleMusicHandler) DecryptedExts() []string {
	return []string{".aac", ".m4a", ".alac"}
}
