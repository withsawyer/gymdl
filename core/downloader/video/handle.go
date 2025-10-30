package video

import (
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"os/exec"
)

// VideoInfo 音乐信息
type VideoInfo struct {
	Title    string // 视频标题
	FileExt  string // 格式
	Size     int    // 音乐大小
	Bitrate  string // 码率
	Duration int    // 时长
	PicUrl   string // 封面图url
	Tidy     string // 入库方式(默认/webdav)
}

// Handler 视频处理接口
type Handler interface {
	// 平台名称
	Platform() string
	// 下载音乐
	Download(url string, cfg *config.Config) (*VideoInfo, error)
	// 构建下载命令
	DownloadCommand(cfg *config.Config, url string) *exec.Cmd
	// 视频整理之前的处理(如刮削等)
	BeforeTidy(cfg *config.Config, videoInfo *VideoInfo) error
	// 是否需要移除DRM
	NeedRemoveDRM(cfg *config.Config) bool
	// 移除DRM
	DRMRemove(cfg *config.Config, videoInfo *VideoInfo) error
	// 音乐整理
	TidyVideo(cfg *config.Config, webdav *core.WebDAV, videoInfo *VideoInfo) error
	// 加密后缀
	EncryptedExts() []string
	// 非加密后缀
	DecryptedExts() []string
}
