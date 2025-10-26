package handler

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
)

// 处理youtube music
type YoutubeMusicHandler struct{}

func (ytm *YoutubeMusicHandler) DownloadMusic(url string, cfg *config.Config) error {
	return nil
}

// 构建下载命令
func (ytm *YoutubeMusicHandler) DownloadCommand(cfg *config.Config) *exec.Cmd {
	return nil
}

// 音乐整理之前的逻辑
func (ytm *YoutubeMusicHandler) BeforeTidy(cfg *config.Config) error {
	return nil
}

// 是否需要移除DRM
func (ytm *YoutubeMusicHandler) NeedRemoveDRM(cfg *config.Config) bool {
	return false
}

// 移除DRM
func (ytm *YoutubeMusicHandler) DRMRemove(cfg *config.Config) error {
	return nil
}

// 音乐整理
func (ytm *YoutubeMusicHandler) TidyMusic(cfg *config.Config) error {
	return nil
}
