package music

import (
	"context"
	"fmt"
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"gopkg.in/vansante/go-ffprobe.v2"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// SongInfo 音乐信息
type SongInfo struct {
	SongName    string // 音乐名称
	SongArtists string // 艺术家
	SongAlbum   string // 专辑
	FileExt     string // 格式
	MusicSize   int    // 音乐大小
	Bitrate     string // 码率
	Duration    int    // 时长
	PicUrl      string // 封面图url
	Lyric       string //歌词
	Year        int    //年份
	Tidy        string // 入库方式(默认/webdav)
}

// Handler 音乐处理接口
type Handler interface {
	// 平台名称
	Platform() string
	// 下载音乐
	Download(url string, cfg *config.Config) (*SongInfo, error)
	// 构建下载命令
	DownloadCommand(cfg *config.Config, url string) *exec.Cmd
	// 音乐整理之前的处理(如嵌入元数据,刮削等)
	BeforeTidy(cfg *config.Config, songInfo *SongInfo) error
	// 是否需要移除DRM
	NeedRemoveDRM(cfg *config.Config) bool
	// 移除DRM
	DRMRemove(cfg *config.Config, songInfo *SongInfo) error
	// 音乐整理
	TidyMusic(cfg *config.Config, webdav *core.WebDAV, songInfo *SongInfo) error
	// 加密后缀
	EncryptedExts() []string
	// 非加密后缀
	DecryptedExts() []string
}

// 设定整理类型
func determineTidyType(cfg *config.Config) string {
	return map[int]string{1: "LOCAL", 2: "WEBDAV"}[cfg.ResourceTidy.Mode]
}

// ExtractSongInfo 通过ffprobe-go解析歌曲信息
func ExtractSongInfo(song *SongInfo, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	// 文件信息（大小和扩展名）
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 用 ffprobe 获取所有元信息
	data, err := ffprobe.ProbeURL(ctx, path)
	if err != nil {
		return fmt.Errorf("获取音频信息失败: %w", err)
	}
	song.MusicSize = int(info.Size())
	song.FileExt = strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")

	// 获取基础信息
	if data.Format != nil {
		if dur := data.Format.Duration(); dur > 0 {
			song.Duration = int(dur.Seconds())
		}
		if br, err := strconv.Atoi(data.Format.BitRate); err == nil {
			song.Bitrate = strconv.Itoa(br / 1000)
		}

		// 标签信息
		if tags := data.Format.TagList; tags != nil {
			song.SongName, _ = tags.GetString("title")
			song.SongArtists, _ = tags.GetString("artist")
			song.SongAlbum, _ = tags.GetString("album")
		}
	}

	return nil
}
