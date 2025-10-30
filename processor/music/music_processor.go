package music

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/processor"
	"gopkg.in/vansante/go-ffprobe.v2"
)

/* ---------------------- 音乐接口定义 ---------------------- */
type MusicProcessor interface {
	processor.Processor
	// 音乐处理器名称
	Name() domain.LinkType
	// 歌曲元信息列表
	Songs() []*SongInfo
	// 下载音乐
	DownloadMusic(url string) error
	// 构建下载命令
	DownloadCommand(url string) *exec.Cmd
	// 音乐整理之前的处理(如读取,嵌入元数据,刮削等)
	BeforeTidy() error
	// 是否需要移除DRM
	NeedRemoveDRM() bool
	// 移除DRM
	DRMRemove() error
	// 音乐整理
	TidyMusic() error
	// 加密后缀
	EncryptedExts() []string
	// 非加密后缀
	DecryptedExts() []string
}

/* ---------------------- 音乐结构体定义 ---------------------- */
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
	Lyric       string // 歌词
	Year        int    // 年份
	Tidy        string // 入库方式(默认/webdav)
}

/* ---------------------- 常量 ---------------------- */

var BaseTempDir = filepath.Join("data", "temp")

// 苹果音乐临时文件夹
var AppleMusicTempDir = filepath.Join(BaseTempDir, "AppleMusic")

// 网易云音乐临时文件夹
var NCMTempDir = filepath.Join(BaseTempDir, "NCM")

// QQ音乐临时文件夹
var QQTempDir = filepath.Join(BaseTempDir, "QQ")

// Youtube音乐临时文件夹
var YoutubeTempDir = filepath.Join(BaseTempDir, "Youtube")

// SoundCloud临时文件夹
var SoundcloudTempDir = filepath.Join(BaseTempDir, "Soundcloud")

// Spotify临时文件夹
var SpotifyTempDir = filepath.Join(BaseTempDir, "Spotify")

// 设定整理类型
func determineTidyType(cfg *config.Config) string {
	return map[int]string{1: "LOCAL", 2: "WEBDAV"}[cfg.MusicTidy.Mode]
}

/* ---------------------- 业务工具 ---------------------- */

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
