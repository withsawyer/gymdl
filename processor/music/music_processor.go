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

	"github.com/bogem/id3v2/v2"
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
	MusicPath   string //音乐文件路径
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
	song.MusicPath = path
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
			song.Lyric, _ = tags.GetString("lyrics")
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

// EmbedMetadata 为音频嵌入封面、元信息、歌词
func EmbedMetadata(song *SongInfo, filePath string) error {
	if song == nil {
		return fmt.Errorf("song info is nil")
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filePath), "."))
	utils.DebugWithFormat("🎧 Embedding metadata for [%s] (%s)", song.SongName, ext)

	tempFile := filePath + ".tmp" + filepath.Ext(filePath)
	coverPath, cleanup, err := prepareCover(song.PicUrl)
	if cleanup != nil {
		defer cleanup()
	}
	if err != nil {
		utils.DebugWithFormat("⚠️ No cover embedded: %v", err)
	}

	args := buildFFmpegArgs(ext, filePath, tempFile, song, coverPath)

	if err := runFFmpeg(args); err != nil {
		return fmt.Errorf("ffmpeg failed: %v", err)
	}

	if err := replaceFile(tempFile, filePath); err != nil {
		return fmt.Errorf("replace file failed: %v", err)
	}

	if ext == "mp3" && song.Lyric != "" {
		if err := writeID3Lyrics(filePath, song.Lyric); err != nil {
			utils.DebugWithFormat("❌ Failed to write lyrics: %v", err)
		} else {
			utils.DebugWithFormat("✅ Lyrics embedded via ID3v2 successfully")
		}
	}

	utils.DebugWithFormat("✨ Metadata embedding completed for [%s]", song.SongName)
	return nil
}

// EmbedLyricsOnly 仅为音频文件嵌入歌词
func EmbedLyricsOnly(filePath, lyrics string) error {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filePath), "."))

	// MP3 用 ID3v2 写入歌词
	if ext == "mp3" {
		return writeID3Lyrics(filePath, lyrics)
	}

	// 其他格式用 ffmpeg -metadata 写入歌词
	tempFile := filePath + ".tmp" + filepath.Ext(filePath)

	args := []string{
		"-y", "-i", filePath,
		"-c", "copy",
		"-metadata", fmt.Sprintf("lyrics=%s", lyrics),
		tempFile,
	}

	if err := runFFmpeg(args); err != nil {
		return fmt.Errorf("ffmpeg failed: %v", err)
	}

	if err := replaceFile(tempFile, filePath); err != nil {
		return fmt.Errorf("replace file failed: %v", err)
	}

	return nil
}

// prepareCover 下载或确认封面文件存在
func prepareCover(picURL string) (string, func(), error) {
	if picURL == "" {
		return "", nil, fmt.Errorf("no cover URL provided")
	}

	if strings.HasPrefix(picURL, "http") {
		tmpFile, err := os.CreateTemp("", "cover_*.jpg")
		if err != nil {
			return "", nil, fmt.Errorf("create temp cover failed: %v", err)
		}
		tmpFile.Close()

		if err := utils.DownloadFile(picURL, tmpFile.Name()); err != nil {
			os.Remove(tmpFile.Name())
			return "", nil, fmt.Errorf("download cover failed: %v", err)
		}
		return tmpFile.Name(), func() { _ = os.Remove(tmpFile.Name()) }, nil
	}

	if _, err := os.Stat(picURL); err != nil {
		return "", nil, fmt.Errorf("cover file not found: %v", err)
	}
	return picURL, nil, nil
}

// buildFFmpegArgs 根据格式生成对应参数
func buildFFmpegArgs(ext, input, output string, song *SongInfo, coverPath string) []string {
	args := []string{"-y", "-i", input}
	if coverPath != "" {
		args = append(args, "-i", coverPath)
	}

	args = append(args, metadataArgs(song)...)

	switch ext {
	case "mp3":
		args = append(args, "-c", "copy", "-id3v2_version", "3")
	case "flac", "m4a", "aac", "mp4", "ogg", "opus", "ape", "wv":
		args = append(args, "-c", "copy")
	default:
		args = append(args, "-c", "copy")
	}

	args = append(args, coverArgs(coverPath, ext)...)

	if song.Lyric != "" && ext != "mp3" {
		args = append(args, "-metadata", fmt.Sprintf("lyrics=%s", song.Lyric))
	}

	args = append(args, output)
	return args
}

// 通用元数据参数生成
func metadataArgs(song *SongInfo) []string {
	m := map[string]string{
		"title":        song.SongName,
		"artist":       song.SongArtists,
		"album":        song.SongAlbum,
		"album_artist": song.SongArtists,
	}
	if song.Year > 0 {
		m["date"] = fmt.Sprintf("%d", song.Year)
	}

	var args []string
	for k, v := range m {
		if v != "" {
			args = append(args, "-metadata", fmt.Sprintf("%s=%s", k, v))
		}
	}
	return args
}

// 封面参数生成
func coverArgs(coverPath, ext string) []string {
	if coverPath == "" {
		return nil
	}
	baseArgs := []string{
		"-map", "0:a?", "-map", "1:v?",
		"-metadata:s:v", "title=Cover",
	}
	if ext == "mp3" || ext == "flac" || ext == "ape" || ext == "wv" || strings.HasPrefix(ext, "m4") || ext == "aac" || ext == "mp4" {
		baseArgs = append(baseArgs,
			"-metadata:s:v", "comment=Cover (front)",
			"-disposition:v", "attached_pic",
		)
	}
	return baseArgs
}

// 执行 FFmpeg
func runFFmpeg(args []string) error {
	utils.DebugWithFormat("🚀 Running ffmpeg: ffmpeg %v", strings.Join(args, " "))
	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.DebugWithFormat("FFmpeg error:\n%s", string(output))
		return err
	}
	utils.DebugWithFormat("FFmpeg OK (%d bytes output)", len(output))
	return nil
}

// 写入 MP3 歌词（ID3v2）
func writeID3Lyrics(filePath, lyrics string) error {
	tag, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	if err != nil {
		return fmt.Errorf("open id3 tag failed: %v", err)
	}
	defer tag.Close()

	tag.AddUnsynchronisedLyricsFrame(id3v2.UnsynchronisedLyricsFrame{
		Encoding: id3v2.EncodingUTF8,
		Language: "chi",
		Lyrics:   lyrics,
	})

	return tag.Save()
}

// 文件替换工具
func replaceFile(src, dst string) error {
	if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.Rename(src, dst)
}
