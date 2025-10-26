package handler

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gcottom/audiometa/v3"
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/core/constants"
	"github.com/nichuanfang/gymdl/utils"
)

type AppleMusicHandler struct{}

/* ---------------------- 基础接口实现 ---------------------- */

func (am *AppleMusicHandler) Platform() string { return "AppleMusic" }

/* ---------------------- 下载逻辑 ---------------------- */

func (am *AppleMusicHandler) DownloadMusic(url string, cfg *config.Config) (*SongInfo, error) {
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

	var tidy string
	if cfg.MusicTidy.Mode == 1 {
		tidy = "default"
	} else {
		tidy = "webdav"
	}
	return &SongInfo{
		Tidy: tidy,
	}, nil
}

/* ---------------------- 构建下载命令 ---------------------- */

func (am *AppleMusicHandler) DownloadCommand(cfg *config.Config, url string) *exec.Cmd {
	cookiePath := filepath.Join(cfg.CookieCloud.CookieFilePath, cfg.CookieCloud.CookieFile)
	args := []string{
		"--cookies-path", cookiePath,
		"--download-mode", "nm3u8dlre",
		"--output-path", constants.BaseTempDir,
		"--temp-path", constants.BaseTempDir,
		"--album-folder-template", "AppleMusic",
		"--compilation-folder-template", "AppleMusic",
		"--no-album-folder-template", "AppleMusic",
		"--single-disc-folder-template", "{title}",
		"--multi-disc-folder-template", "{title}",
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

	// 只打开一次文件获取标签和文件大小
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	// 获取文件信息（大小）
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %w", err)
	}
	songInfo.MusicSize = int(info.Size())
	songInfo.FileExt = strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")

	// 读取音频标签
	tag, err := audiometa.OpenTag(f)
	if err != nil {
		utils.WarnWithFormat("[AppleMusic] ⚠️ 读取音频标签失败: %v", err)
	} else {
		songInfo.SongName = tag.GetTitle()
		songInfo.SongArtists = tag.GetArtist()
		songInfo.SongAlbum = tag.GetAlbum()
	}

	// 使用 ffprobe 获取比特率和时长
	// 将命令参数拆开，避免多余的 shell 执行
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration,bit_rate",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	)

	// 设置合理的超时，避免阻塞
	type ffprobeResult struct {
		out []byte
		err error
	}
	ch := make(chan ffprobeResult, 1)
	go func() {
		out, err := cmd.Output()
		ch <- ffprobeResult{out, err}
	}()

	select {
	case res := <-ch:
		if res.err != nil {
			utils.WarnWithFormat("[AppleMusic] ⚠️ ffprobe 获取音频信息失败: %v", res.err)
		} else {
			lines := strings.Split(strings.TrimSpace(string(res.out)), "\n")
			if len(lines) >= 2 {
				if duration, err := strconv.ParseFloat(lines[0], 64); err == nil {
					songInfo.Duration = int(duration)
				}
				if bitrate, err := strconv.Atoi(lines[1]); err == nil {
					songInfo.Bitrate = strconv.Itoa(bitrate / 1000)
				}
			}
		}
	case <-time.After(5 * time.Second):
		utils.WarnWithFormat("[AppleMusic] ⚠️ ffprobe 超时，跳过获取时长和比特率")
	}

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

	switch cfg.MusicTidy.Mode {
	case 1:
		return am.tidyToLocal(cfg, files)
	case 2:
		return am.tidyToWebDAV(cfg, files, webdav)
	default:
		return fmt.Errorf("未知整理模式: %d", cfg.MusicTidy.Mode)
	}
}

/* ---------------------- 本地整理 ---------------------- */

func (am *AppleMusicHandler) tidyToLocal(cfg *config.Config, files []os.DirEntry) error {
	dstDir := cfg.MusicTidy.DistDir
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
		return "", errors.New("未找到符合条件的解密文件")
	}

	return filepath.Join(constants.AppleMusicTempDir, latestFile.Name()), nil
}
