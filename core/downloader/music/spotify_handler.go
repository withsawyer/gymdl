package music

import (
	"errors"
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
)

// 处理spotify music
type SpotifyHandler struct{}

func (sp *SpotifyHandler) Platform() string {
	return "Spotify"
}

func (sp *SpotifyHandler) Download(url string, cfg *config.Config) (*SongInfo, error) {
	return &SongInfo{}, errors.New("🚧 开发中")
}

// 构建下载命令
func (sp *SpotifyHandler) DownloadCommand(cfg *config.Config, url string) *exec.Cmd {
	return nil
}

// 音乐整理之前的逻辑
func (sp *SpotifyHandler) BeforeTidy(cfg *config.Config, songInfo *SongInfo) error {
	return nil
}

// 是否需要移除DRM
func (sp *SpotifyHandler) NeedRemoveDRM(cfg *config.Config) bool {
	return false
}

// 移除DRM
func (sp *SpotifyHandler) DRMRemove(cfg *config.Config, songInfo *SongInfo) error {
	return nil
}

// 音乐整理
func (sp *SpotifyHandler) TidyMusic(cfg *config.Config, webdav *core.WebDAV, songInfo *SongInfo) error {
	return nil
}

// 加密后缀
func (sp *SpotifyHandler) EncryptedExts() []string {
	return []string{""}
}

// 非加密后缀
func (sp *SpotifyHandler) DecryptedExts() []string {
	return []string{".aac", ".m4a", ".flac", "mp3", "ogg"}
}
