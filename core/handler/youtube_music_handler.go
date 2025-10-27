package handler

import (
	"errors"
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
)

// 处理youtube music
type YoutubeMusicHandler struct{}

func (ytm *YoutubeMusicHandler) Platform() string {
	return "YoutubeMusic"
}

func (ytm *YoutubeMusicHandler) DownloadMusic(url string, cfg *config.Config) (*SongInfo, error) {
	return &SongInfo{}, errors.New("🚧 开发中")
}

// 构建下载命令
func (ytm *YoutubeMusicHandler) DownloadCommand(cfg *config.Config, url string) *exec.Cmd {
	return nil
}

// 音乐整理之前的逻辑
func (ytm *YoutubeMusicHandler) BeforeTidy(cfg *config.Config, songInfo *SongInfo) error {
	return nil
}

// 是否需要移除DRM
func (ytm *YoutubeMusicHandler) NeedRemoveDRM(cfg *config.Config) bool {
	return false
}

// 移除DRM
func (ytm *YoutubeMusicHandler) DRMRemove(cfg *config.Config, songInfo *SongInfo) error {
	return nil
}

// 音乐整理
func (ytm *YoutubeMusicHandler) TidyMusic(cfg *config.Config, webdav *core.WebDAV, songInfo *SongInfo) error {
	return nil
}

// 加密后缀
func (ytm *YoutubeMusicHandler) EncryptedExts() []string {
	return []string{""}
}

// 非加密后缀
func (ytm *YoutubeMusicHandler) DecryptedExts() []string {
	return []string{".aac", ".m4a", ".flac", "mp3", "ogg"}
}
