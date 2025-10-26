package handler

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/bytedance/gopkg/util/logger"
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
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	cmd := am.DownloadCommand(cfg, url)
	utils.InfoWithFormat("🎵 开始下载 Apple Music 内容: %s\n", url)
	utils.DebugWithFormat("DownloadCommand： %s ", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gamdl 下载失败: %w\n\n输出:\n\n %s", err, string(output))
	}

	logger.Debug("\n\n" + string(output))
	utils.InfoWithFormat("✅ 下载完成（耗时 %v）\n", time.Since(start).Truncate(time.Millisecond))
	return nil
}

// DownloadCommand 构建下载命令
func (am *AppleMusicHandler) DownloadCommand(cfg *config.Config, url string) *exec.Cmd {
	cookiePath := filepath.Join(cfg.CookieCloud.CookieFilePath, cfg.CookieCloud.CookieFile)
	return exec.Command(
		"gamdl",
		"--cookies-path", cookiePath,
		"--no-config-file", "true",
		"--download-mode", "nm3u8dlre",
		"--overwrite", "true",
		"--output-path", constants.BaseTempDir,
		"--album-folder-template", "AppleMusic",
		"--compilation-folder-template", "AppleMusic",
		"--no-album-folder-template", "AppleMusic",
		"--no-synced-lyrics",
		url,
	)
}

// BeforeTidy 清洗/解锁音乐文件
func (am *AppleMusicHandler) BeforeTidy(cfg *config.Config) error {
	if am.NeedRemoveDRM(cfg) {
		if err := am.DRMRemove(cfg); err != nil {
			return fmt.Errorf("DRM 移除失败: %w", err)
		}
		logger.Info("🔓 DRM 已移除")
	}
	return nil
}

// NeedRemoveDRM 判断是否需要去除 DRM
func (am *AppleMusicHandler) NeedRemoveDRM(cfg *config.Config) bool {
	// 当前默认不去 DRM，可根据 cfg 配置动态调整
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
		return fmt.Errorf("读取临时目录失败: %w", err)
	}
	if len(files) == 0 {
		return errors.New("未找到待整理的音乐文件")
	}

	switch cfg.MusicTidy.Mode {
	case 1:
		return am.tidyToLocal(cfg, files)
	case 2:
		return am.tidyToWebDAV(cfg, files, webdav)
	default:
		return fmt.Errorf("未知整理模式: %d", cfg.MusicTidy.Mode)
	}
}

// tidyToLocal 将音乐移动到本地目标目录
func (am *AppleMusicHandler) tidyToLocal(cfg *config.Config, files []os.DirEntry) error {
	dstDir := cfg.MusicTidy.DistDir
	if dstDir == "" {
		return errors.New("未配置输出目录")
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	for _, file := range files {
		if !utils.FilterMusicFile(file, am.EncryptedExts(), am.DecryptedExts()) {
			continue
		}
		srcPath := filepath.Join(constants.AppleMusicTempDir, file.Name())
		dstPath := filepath.Join(dstDir, utils.SanitizeFileName(file.Name()))
		if err := os.Rename(srcPath, dstPath); err != nil {
			return fmt.Errorf("移动文件失败 %s → %s: %w", srcPath, dstPath, err)
		}
		utils.InfoWithFormat("📦 已整理: %s\n", dstPath)
	}
	return nil
}

// tidyToWebDAV 将音乐上传到 WebDAV
func (am *AppleMusicHandler) tidyToWebDAV(cfg *config.Config, files []os.DirEntry, webdav *core.WebDAV) error {
	if webdav == nil {
		return errors.New("WebDAV 未初始化")
	}

	for _, file := range files {
		if !utils.FilterMusicFile(file, am.EncryptedExts(), am.DecryptedExts()) {
			continue
		}
		filePath := filepath.Join(constants.AppleMusicTempDir, file.Name())
		if err := webdav.Upload(filePath); err != nil {
			return fmt.Errorf("上传文件失败 %s: %w", file.Name(), err)
		}
		utils.InfoWithFormat("☁️ 已上传: %s\n", file.Name())
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
