package music

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nichuanfang/gymdl/processor"
	"github.com/nichuanfang/gymdl/utils"
	"go.senan.xyz/taglib"
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
	MusicSize   int64  // 音乐大小
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
			song, err := ReadTags(fullPath)
			//嵌入默认标签
			FillDefaultTags(fullPath, song)
			if err != nil {
				return nil, fmt.Errorf("处理文件 %s 失败: %w", f.Name(), err)
			}
			song.Tidy = tidyType
			songs = append(songs, song)
		}
	}
	return songs, nil
}

// ReadTags 读取音乐元数据
func ReadTags(path string) (*SongInfo, error) {
	tags, err := taglib.ReadTags(path)
	if err != nil {
		return nil, err
	}

	props, err := taglib.ReadProperties(path)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	songInfo := &SongInfo{
		FileExt:   strings.TrimPrefix(filepath.Ext(path), "."),
		MusicSize: fileInfo.Size(),
		Bitrate:   strconv.Itoa(int(props.Bitrate)),
		Duration:  int(props.Length),
	}

	if t, ok := tags[taglib.Title]; ok && len(t) > 0 {
		songInfo.SongName = t[0]
	} else {
		songInfo.SongName = filepath.Base(path)
	}

	if a, ok := tags[taglib.Artist]; ok && len(a) > 0 {
		songInfo.SongArtists = a[0]
	}

	if al, ok := tags[taglib.Album]; ok && len(al) > 0 {
		songInfo.SongAlbum = al[0]
	}

	if d, ok := tags[taglib.Date]; ok && len(d) > 0 {
		songInfo.Year, _ = strconv.Atoi(d[0])
	}

	if l, ok := tags[taglib.Lyrics]; ok && len(l) > 0 {
		songInfo.Lyric = l[0]
	}

	return songInfo, nil
}

// FillDefaultTags 标签写入默认值
func FillDefaultTags(path string, info *SongInfo) {
	updates := make(map[string][]string)

	if info.Year == 0 {
		info.Year = 2020
		updates[taglib.Date] = []string{strconv.Itoa(info.Year)}
	}

	if info.Lyric == "" {
		info.Lyric = "[00:00:00]此歌曲为没有填词的纯音乐，请您欣赏"
		updates[taglib.Lyrics] = []string{info.Lyric}
	}

	if len(updates) > 0 {
		if err := taglib.WriteTags(path, updates, 0); err != nil {
			utils.WarnWithFormat("write default tags failed: %w", err)
		}
	}
}

// WriteTags 写入标签
func WriteTags(song *SongInfo, filePath string) error {
	//写入标签
	err := taglib.WriteTags(filePath, map[string][]string{
		taglib.Title:       {song.SongName},
		taglib.Artist:      {song.SongArtists},
		taglib.Album:       {song.SongAlbum},
		taglib.AlbumArtist: {song.SongAlbum},
		//年份
		taglib.Date: {strconv.Itoa(song.Year)},
		//歌词
		taglib.Lyrics: {song.Lyric},
		//流派
		//taglib.Genre:   {song.SongName},
	}, 0)
	if err != nil {
		utils.WarnWithFormat("write metadate failed: %w", err)
		return nil
	}
	// 写入封面图片
	if song.PicUrl != "" {
		// 下载图片
		imageData, errImage := utils.FetchImage(song.PicUrl)
		if errImage != nil {
			utils.WarnWithFormat("fetch image failed: %w", err)
			return nil
		}
		// 写入图片到音频文件
		err = taglib.WriteImage(filePath, imageData)
		if err != nil {
			utils.WarnWithFormat("write image failed: %w", err)
			return nil
		}
	}
	return nil
}
