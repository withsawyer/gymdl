package music

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nichuanfang/gymdl/processor"
	"github.com/nichuanfang/gymdl/utils"
	"gopkg.in/vansante/go-ffprobe.v2"
)

/* ---------------------- 音乐接口定义 ---------------------- */
type Processor interface {
	processor.Processor
	// 歌曲元信息列表
	Songs() []*SongInfo
	// 下载音乐
	DownloadMusic(url string, callback func(string)) error
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
	Url         string //下载地址
	PicUrl      string // 封面图url
	Lyric       string // 歌词
	Year        int    // 年份
	Tidy        string // 入库方式(默认/webdav)
}

/* ---------------------- 常量 ---------------------- */

var BaseTempDir = filepath.Join("data", "temp", "music")

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

/* ---------------------- 音乐下载相关业务函数 ---------------------- */

// ExtractSongInfo 通过ffprobe-go解析歌曲信息
func ExtractSongInfo(path string) (*SongInfo, error) {
	song := &SongInfo{}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	// 文件信息（大小和扩展名）
	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 用 ffprobe 获取所有元信息
	data, err := ffprobe.ProbeURL(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("获取音频信息失败: %w", err)
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

	return song, nil
}

// 读取音乐目录 返回元信息列表
func ReadMusicDir(tempDir string, tidyType string, p Processor) ([]*SongInfo, error) {
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return nil, fmt.Errorf("读取临时目录失败: %w", err)
	}
	songs := make([]*SongInfo, 0, len(files))
	for _, f := range files {
		//目录跳过
		if f.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(f.Name()))
		if utils.Contains(p.DecryptedExts(), ext) {
			fullPath := filepath.Join(tempDir, f.Name())
			song, err := ExtractSongInfo(fullPath)
			if err != nil {
				return nil, fmt.Errorf("处理文件 %s 失败: %w", f.Name(), err)
			}
			song.Tidy = tidyType
			songs = append(songs, song)
		}
	}
	return songs, nil
}

// EmbedMetadata 嵌入音频元数据
func EmbedMetadata(song *SongInfo, filePath string) error {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filePath), ""))

	// 处理封面
	var coverPath string
	if song.PicUrl != "" {
		if strings.HasPrefix(song.PicUrl, "http://") || strings.HasPrefix(song.PicUrl, "https://") {
			tmpFile, err := os.CreateTemp("", "cover_*"+filepath.Ext(song.PicUrl))
			if err != nil {
				return fmt.Errorf("create temp file for cover failed: %v", err)
			}
			coverPath = tmpFile.Name()
			tmpFile.Close()

			if err := utils.DownloadFile(song.PicUrl, coverPath); err != nil {
				os.Remove(coverPath)
				return fmt.Errorf("download cover image failed: %v", err)
			}
			defer os.Remove(coverPath)
		} else {
			coverPath = song.PicUrl
		}
	}

	// 临时输出文件，保留原始扩展名，避免 FFmpeg 报错
	tempFile := filePath + ".tmp" + filepath.Ext(filePath)
	args := []string{"-y", "-i", filePath}

	// 封面输入
	if coverPath != "" {
		args = append(args, "-i", coverPath)
	}

	// 音频编码处理
	switch ext {
	case "mp3":
		args = append(args, "-c", "copy", "-id3v2_version", "3")
	case "flac", "m4a", "mp4", "aac", "ogg":
		args = append(args, "-c", "copy")
	default:
		args = append(args, "-c", "copy")
	}

	// 映射封面流并设置元数据
	if coverPath != "" {
		args = append(args,
			"-map", "0:a?", "-map", "1:v?",
			"-metadata:s:v", "title=Cover",
			"-metadata:s:v", "comment=Cover (front)",
			"-disposition:v", "attached_pic",
		)
	}

	// 内嵌元数据
	metadata := map[string]string{
		"title":  song.SongName,
		"artist": song.SongArtists,
		"album":  song.SongAlbum,
		"comment": func() string {
			if song.Lyric != "" {
				return song.Lyric
			}
			return ""
		}(),
	}

	if song.Year > 0 {
		metadata["date"] = fmt.Sprintf("%d", song.Year)
	}

	for k, v := range metadata {
		if v != "" {
			args = append(args, "-metadata", fmt.Sprintf("%s=%s", k, v))
		}
	}

	args = append(args, tempFile)

	// 调试输出
	utils.DebugWithFormat("Running ffmpeg with args: %v\n", args)

	// 执行 FFmpeg
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %v", err)
	}

	// 安全覆盖原文件
	if err := replaceFile(tempFile, filePath); err != nil {
		return fmt.Errorf("replace file failed: %v", err)
	}

	return nil
}

// replaceFile 安全覆盖文件，支持跨分区
func replaceFile(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// 跨分区，拷贝 + 删除
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer output.Close()

	if _, err := io.Copy(output, input); err != nil {
		return err
	}
	if err := output.Sync(); err != nil {
		return err
	}

	return os.Remove(src)
}
