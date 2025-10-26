package handler

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
)

// 处理spotify music
type SpotifyHandler struct{}

func (sp *SpotifyHandler) DownloadMusic(url string, cfg *config.Config) error {
	return nil
}

// 构建下载命令
func (sp *SpotifyHandler) DownloadCommand(cfg *config.Config) *exec.Cmd {
	return nil
}

// 音乐整理之前的逻辑
func (sp *SpotifyHandler) BeforeTidy(cfg *config.Config) error {
	return nil
}

// 是否需要移除DRM
func (sp *SpotifyHandler) NeedRemoveDRM(cfg *config.Config) bool {
	return false
}

// 移除DRM
func (sp *SpotifyHandler) DRMRemove(cfg *config.Config) error {
	return nil
}

// 音乐整理
func (sp *SpotifyHandler) TidyMusic(cfg *config.Config) error {
	return nil
}
