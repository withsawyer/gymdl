package video

import (
	"errors"
	"github.com/nichuanfang/gymdl/core"
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
)

// QQ音乐处理器
type DouYinVideoHandler struct{}

func (qm *DouYinVideoHandler) Platform() string {
	return "抖音"
}

func (qm *DouYinVideoHandler) Download(url string, cfg *config.Config) (*VideoInfo, error) {
	return &VideoInfo{}, errors.New("🚧 抖音正在开发中")
}

func (qm *DouYinVideoHandler) TidyVideo(cfg *config.Config, webdav *core.WebDAV, videoInfo *VideoInfo) error {
	return nil
}

// 构建下载命令
func (qm *DouYinVideoHandler) DownloadCommand(cfg *config.Config, url string) *exec.Cmd {
	return nil
}

// BeforeTidy 整理之前的处理
func (qm *DouYinVideoHandler) BeforeTidy(cfg *config.Config, videoInfo *VideoInfo) error {
	return nil
}

// 是否需要移除DRM
func (qm *DouYinVideoHandler) NeedRemoveDRM(cfg *config.Config) bool {
	return false
}

// 移除DRM
func (qm *DouYinVideoHandler) DRMRemove(cfg *config.Config, videoInfo *VideoInfo) error {
	return nil
}

// 加密后缀
func (qm *DouYinVideoHandler) EncryptedExts() []string {
	return []string{".mflac"}
}

// 非加密后缀
func (qm *DouYinVideoHandler) DecryptedExts() []string {
	return []string{".aac", ".m4a", ".flac", "mp3", "ogg"}
}
