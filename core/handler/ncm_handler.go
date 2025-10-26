package handler

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
)

// 网易云音乐处理器
type NCMHandler struct{}

func (ncm *NCMHandler) DownloadMusic(url string, cfg *config.Config) error {
	return nil
}

// 构建下载命令
func (ncm *NCMHandler) DownloadCommand(cfg *config.Config) *exec.Cmd {
	return nil
}

// 音乐整理之前的处理
func (ncm *NCMHandler) BeforeTidy(cfg *config.Config) error {
	return nil
}

// 是否需要移除DRM
func (ncm *NCMHandler) NeedRemoveDRM(cfg *config.Config) bool {
	return false
}

// 移除DRM
func (ncm *NCMHandler) DRMRemove(cfg *config.Config) error {
	return nil
}

// 音乐整理
func (ncm *NCMHandler) TidyMusic(cfg *config.Config) error {
	return nil
}
