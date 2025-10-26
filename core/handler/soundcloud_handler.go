package handler

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
)

// SoundCloud音乐处理器
type SoundCloudHandler struct{}

func (sc *SoundCloudHandler) DownloadMusic(url string, cfg *config.Config) error {
	return nil
}

// 构建下载命令
func (sc *SoundCloudHandler) DownloadCommand(cfg *config.Config) *exec.Cmd {
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
func (sc *SoundCloudHandler) TidyMusic(cfg *config.Config) error {
	return nil
}
