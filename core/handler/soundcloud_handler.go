package handler

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
)

// SoundCloud音乐处理器
type SoundCloudHandler struct{}

func (sc *SoundCloudHandler) DownloadMusic(url string, cfg *config.Config) error {
	return nil
}

// 构建下载命令
func (sc *SoundCloudHandler) DownloadCommand(cfg *config.Config, url string) *exec.Cmd {
	return nil
}

// 音乐整理之前的处理
func (sc *SoundCloudHandler) BeforeTidy(cfg *config.Config) error {
	return nil
}

// 是否需要移除DRM
func (sc *SoundCloudHandler) NeedRemoveDRM(cfg *config.Config) bool {
	return false
}

// 移除DRM
func (sc *SoundCloudHandler) DRMRemove(cfg *config.Config) error {
	return nil
}

// 音乐整理
func (sc *SoundCloudHandler) TidyMusic(cfg *config.Config, webdav *core.WebDAV) error {
	return nil
}

// 加密后缀
func (sc *SoundCloudHandler) EncryptedExts() []string {
	return []string{""}
}

// 非加密后缀
func (sc *SoundCloudHandler) DecryptedExts() []string {
	return []string{".aac", ".m4a", ".flac", "mp3", "ogg"}
}
