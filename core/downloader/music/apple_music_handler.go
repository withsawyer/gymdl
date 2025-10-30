package music

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

type AppleMusicHandler struct{}

/* ---------------------- 基础接口实现 ---------------------- */

func (am *AppleMusicHandler) Platform() string { return "AppleMusic" }

/* ---------------------- 下载逻辑 ---------------------- */

func (am *AppleMusicHandler) Download(url string, cfg *config.Config) (*SongInfo, error) {
	if cfg.AdditionalConfig.EnableWrapper {
		utils.Logger().Debug("使用增强版am下载器")
		return am.wrapDownload(url, cfg)
	} else {
		utils.Logger().Debug("使用默认am下载器")
		return am.defaultDownload(url, cfg)
	}
}

// defaultDownload 默认下载器
func (am *AppleMusicHandler) defaultDownload(url string, cfg *config.Config) (*SongInfo, error) {
	start := time.Now()
	tempDir := constants.AppleMusicTempDir

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		utils.ErrorWithFormat("[AppleMusic] ❌ 创建输出目录失败: %v", err)
		return nil, fmt.Errorf("创建输出目录失败: %w", err)
	}

	cmd := am.DownloadCommand(cfg, url)
	utils.InfoWithFormat("[AppleMusic] 🎵 开始下载: %s", url)
	utils.DebugWithFormat("[AppleMusic] 执行命令: %s", strings.Join(cmd.Args, " "))

	output, err := cmd.CombinedOutput()
	logOut := strings.TrimSpace(string(output))
	if err != nil {
		utils.ErrorWithFormat("[AppleMusic] ❌ gamdl 下载失败: %v\n输出:\n%s", err, logOut)
		return nil, fmt.Errorf("gamdl 下载失败: %w", err)
	}

	if logOut != "" {
		utils.DebugWithFormat("[AppleMusic] 下载输出:\n%s", logOut)
	}
	utils.InfoWithFormat("[AppleMusic] ✅ 下载完成（耗时 %v）", time.Since(start).Truncate(time.Millisecond))

	return &SongInfo{}, nil
}

// wrapDownload 增强版下载器
func (am *AppleMusicHandler) wrapDownload(string, *config.Config) (*SongInfo, error) {
	return &SongInfo{}, nil
}

/* ---------------------- 构建下载命令 ---------------------- */

func (am *AppleMusicHandler) DownloadCommand(cfg *config.Config, url string) *exec.Cmd {
	cookiePath := filepath.Join(cfg.CookieCloud.CookieFilePath, cfg.CookieCloud.CookieFile)
	// https://github.com/glomatico/gamdl/commit/fdab6481ea246c2cf3415565c39da62a3b9dbd52 部分options改动
	args := []string{
		"--cookies-path", cookiePath,
		"--download-mode", "nm3u8dlre",
		"--output-path", constants.BaseTempDir,
		"--temp-path", constants.BaseTempDir,
		"--album-folder-template", "AppleMusic",
		"--compilation-folder-template", "AppleMusic",
		"--no-album-folder-template", "AppleMusic",
		"--single-disc-file-template", "{title}",
		"--multi-disc-file-template", "{title}",
		"--no-synced-lyrics",
		url,
	}
	return exec.Command("gamdl", args...)
}

/* ---------------------- DRM 处理 ---------------------- */

func (am *AppleMusicHandler) BeforeTidy(cfg *config.Config, songInfo *SongInfo) error {
	path, err := am.findLatestDecryptedFile()
	if err != nil {
		return err
	}
	err = ExtractSongInfo(songInfo, path)
	if err != nil {
		return err
	}
	// 设置整理类型
	songInfo.Tidy = determineTidyType(cfg)
	return nil
}

func (am *AppleMusicHandler) NeedRemoveDRM(cfg *config.Config) bool {
	// 后续可添加配置项，比如 cfg.AppleMusic.RemoveDRM
	return false
}

func (am *AppleMusicHandler) DRMRemove(cfg *config.Config, songInfo *SongInfo) error {
	// 预留 DRM 解锁逻辑，比如调用 ffmpeg 转码
	utils.DebugWithFormat("[AppleMusic] DRMRemove() 调用占位")
	return nil
}

/* ---------------------- 音乐整理 ---------------------- */

func (am *AppleMusicHandler) TidyMusic(cfg *config.Config, webdav *core.WebDAV, songInfo *SongInfo) error {
	files, err := os.ReadDir(constants.AppleMusicTempDir)
	if err != nil {
		return fmt.Errorf("读取临时目录失败: %w", err)
	}
	if len(files) == 0 {
		utils.WarnWithFormat("[AppleMusic] ⚠️ 未找到待整理的音乐文件")
		return errors.New("未找到待整理的音乐文件")
	}

	switch cfg.ResourceTidy.Mode {
	case 1:
		return am.tidyToLocal(cfg, files)
	case 2:
		return am.tidyToWebDAV(cfg, files, webdav)
	default:
		return fmt.Errorf("未知整理模式: %d", cfg.ResourceTidy.Mode)
	}
}

/* ---------------------- 本地整理 ---------------------- */

func (am *AppleMusicHandler) tidyToLocal(cfg *config.Config, files []os.DirEntry) error {
	dstDir := cfg.ResourceTidy.DistDir
	if dstDir == "" {
		return errors.New("未配置输出目录")
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	for _, f := range files {
		if !utils.FilterMusicFile(f, am.EncryptedExts(), am.DecryptedExts()) {
			utils.DebugWithFormat("[AppleMusic] 跳过非音乐文件: %s", f.Name())
			continue
		}
		src := filepath.Join(constants.AppleMusicTempDir, f.Name())
		dst := filepath.Join(dstDir, utils.SanitizeFileName(f.Name()))
		if err := os.Rename(src, dst); err != nil {
			utils.WarnWithFormat("[AppleMusic] ⚠️ 移动失败 %s → %s: %v", src, dst, err)
			continue
		}
		utils.InfoWithFormat("[AppleMusic] 📦 已整理: %s", dst)
	}
	return nil
}

/* ---------------------- WebDAV 整理 ---------------------- */

func (am *AppleMusicHandler) tidyToWebDAV(cfg *config.Config, files []os.DirEntry, webdav *core.WebDAV) error {
	if webdav == nil {
		return errors.New("WebDAV 未初始化")
	}

	for _, f := range files {
		if !utils.FilterMusicFile(f, am.EncryptedExts(), am.DecryptedExts()) {
			utils.DebugWithFormat("[AppleMusic] 跳过非音乐文件: %s", f.Name())
			continue
		}

		filePath := filepath.Join(constants.AppleMusicTempDir, f.Name())
		if err := webdav.Upload(filePath); err != nil {
			utils.WarnWithFormat("[AppleMusic] ☁️ 上传失败 %s: %v", f.Name(), err)
			continue
		}
		utils.InfoWithFormat("[AppleMusic] ☁️ 已上传: %s", f.Name())

		ext := strings.ToLower(filepath.Ext(f.Name()))
		if utils.Contains(am.DecryptedExts(), ext) {
			if err := os.Remove(filePath); err == nil {
				utils.DebugWithFormat("[AppleMusic] 🧹 已删除临时文件: %s", filePath)
			} else {
				utils.WarnWithFormat("[AppleMusic] ⚠️ 删除临时文件失败: %s (%v)", filePath, err)
			}
		}
	}
	return nil
}

/* ---------------------- 扩展定义 ---------------------- */

func (am *AppleMusicHandler) EncryptedExts() []string { return []string{".m4p"} }
func (am *AppleMusicHandler) DecryptedExts() []string { return []string{".aac", ".m4a", ".alac"} }

// 获取最新入库文件
func (am *AppleMusicHandler) findLatestDecryptedFile() (string, error) {
	files, err := os.ReadDir(constants.AppleMusicTempDir)
	if err != nil {
		return "", fmt.Errorf("读取临时目录失败: %w", err)
	}

	var latestFile os.DirEntry
	var latestModTime time.Time

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(f.Name()))
		if !utils.Contains(am.DecryptedExts(), ext) {
			continue
		}

		info, err := f.Info()
		if err != nil {
			continue // 无法读取信息则跳过
		}

		modTime := info.ModTime()
		if latestFile == nil || modTime.After(latestModTime) {
			latestFile = f
			latestModTime = modTime
		}
	}

	if latestFile == nil {
		return "", errors.New("未找到符合条件的音乐文件,检测是否下载失败")
	}

	return filepath.Join(constants.AppleMusicTempDir, latestFile.Name()), nil
}
