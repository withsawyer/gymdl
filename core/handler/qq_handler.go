package handler

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
)

// QQ音乐处理器
type QQHandler struct{}

func (qm *QQHandler) DownloadMusic(url string, cfg *config.Config) error {
	return nil
}

// 构建下载命令
func (qm *QQHandler) DownloadCommand(cfg *config.Config) *exec.Cmd {
	return nil
}

// BeforeTidy 整理之前的处理
func (qm *QQHandler) BeforeTidy(cfg *config.Config) error {
	return nil
}

// 是否需要移除DRM
func (qm *QQHandler) NeedRemoveDRM(cfg *config.Config) bool {
	return false
}

// 移除DRM
func (qm *QQHandler) DRMRemove(cfg *config.Config) error {
	return nil
}

// 音乐整理
func (qm *QQHandler) TidyMusic(cfg *config.Config) error {
	return nil
}
