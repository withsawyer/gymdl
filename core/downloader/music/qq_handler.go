package music

import (
	"errors"
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
)

// QQ音乐处理器
type QQHandler struct{}

func (qm *QQHandler) Platform() string {
	return "QQ音乐"
}

func (qm *QQHandler) Download(url string, cfg *config.Config) (*SongInfo, error) {
	return &SongInfo{}, errors.New("🚧 开发中")
}

// 构建下载命令
func (qm *QQHandler) DownloadCommand(cfg *config.Config, url string) *exec.Cmd {
	return nil
}

// BeforeTidy 整理之前的处理
func (qm *QQHandler) BeforeTidy(cfg *config.Config, songInfo *SongInfo) error {
	return nil
}

// 是否需要移除DRM
func (qm *QQHandler) NeedRemoveDRM(cfg *config.Config) bool {
	return false
}

// 移除DRM
func (qm *QQHandler) DRMRemove(cfg *config.Config, songInfo *SongInfo) error {
	return nil
}

// 音乐整理
func (qm *QQHandler) TidyMusic(cfg *config.Config, webdav *core.WebDAV, songInfo *SongInfo) error {
	return nil
}

// 加密后缀
func (qm *QQHandler) EncryptedExts() []string {
	return []string{".mflac"}
}

// 非加密后缀
func (qm *QQHandler) DecryptedExts() []string {
	return []string{".aac", ".m4a", ".flac", "mp3", "ogg"}
}
