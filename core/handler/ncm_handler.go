package handler

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
)

// 网易云音乐处理器
type NCMHandler struct{}

func (ncm *NCMHandler) Platform() string {
	return "网易云音乐"
}

func (ncm *NCMHandler) DownloadMusic(url string, cfg *config.Config) error {
	return nil
}

// 构建下载命令
func (ncm *NCMHandler) DownloadCommand(cfg *config.Config, url string) *exec.Cmd {
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
func (ncm *NCMHandler) TidyMusic(cfg *config.Config, webdav *core.WebDAV) error {
	return nil
}

// 加密后缀
func (ncm *NCMHandler) EncryptedExts() []string {
	return []string{".ncm"}
}

// 非加密后缀
func (ncm *NCMHandler) DecryptedExts() []string {
	return []string{".flac", ".mp3", ".aac", "m4a", "ogg"}
}
