package music

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/types"
	ncmutils "github.com/XiaoMengXinX/Music163Api-Go/utils"
	downloader "github.com/XiaoMengXinX/SimpleDownloader"
	"github.com/gcottom/audiometa/v3"
	"github.com/gcottom/flacmeta"
	"github.com/gcottom/mp3meta"
	"github.com/gcottom/mp4meta"
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/utils"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/processor"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type NetEaseProcessor struct {
	cfg     *config.Config
	songs   []*SongInfo
	tempDir string
}

// Init  初始化
func (ncm *NetEaseProcessor) Init(cfg *config.Config) {
	ncm.cfg = cfg
	ncm.songs = make([]*SongInfo, 0)
	ncm.tempDir = processor.BuildOutputDir(NCMTempDir)
}

/* ---------------------- 基础接口实现 ---------------------- */

func (ncm *NetEaseProcessor) Name() processor.LinkType {
	return processor.LinkNetEase
}

func (ncm *NetEaseProcessor) Songs() []*SongInfo {
	return ncm.songs
}

/* ------------------------ 下载逻辑 ------------------------ */

func (ncm *NetEaseProcessor) DownloadMusic(url string) error {
	start := time.Now()
	utils.InfoWithFormat("[NCM] 🎵 开始下载: %s", url)
	ncmType, musicID := utils.ParseMusicID(url)
	switch ncmType {
	case 1:
		//单曲下载
		return ncm.downloadSingle(musicID, start)
	case 2:
		//列表下载
		return ncm.downloadPlaylist(musicID, start)
	}
	return errors.New("不支持的下载类型")
}

func (ncm *NetEaseProcessor) DownloadCommand(url string) *exec.Cmd {
	return nil
}

func (ncm *NetEaseProcessor) BeforeTidy() error {
	songs := ncm.Songs()
	if len(songs) == 0 {
		return nil
	}

	// 并发控制（默认 4，可根据系统资源调整）
	const maxConcurrent = 4
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	processSong := func(song *SongInfo) {
		defer wg.Done()
		sem <- struct{}{}
		defer func() { <-sem }()

		rawPath := filepath.Join(ncm.tempDir, ncm.safeFileName(song))
		tempPath := filepath.Join(ncm.tempDir, ncm.safeTempFileName(song))

		// 打开原文件
		f, err := os.Open(rawPath)
		if err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("打开文件失败 [%s]: %w", rawPath, err))
			mu.Unlock()
			return
		}
		defer f.Close()

		// 读取音频标签
		tag, err := audiometa.OpenTag(f)
		if err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("读取音频标签失败 [%s]: %w", rawPath, err))
			mu.Unlock()
			return
		}

		// 设置元数据
		tag.SetArtist(song.SongArtists)
		tag.SetTitle(song.SongName)
		tag.SetAlbum(song.SongAlbum)
		tag.SetAlbumArtist(song.SongArtists)

		// 设置封面（无缓存）
		if img, err := utils.FetchImage(song.PicUrl); err == nil {
			tag.SetCoverArt(img)
		}

		// 年份处理
		if song.Year > 0 {
			switch t := tag.(type) {
			case *flacmeta.FLACTag:
				t.SetDate(strconv.Itoa(song.Year))
			case *mp3meta.MP3Tag:
				t.SetYear(song.Year)
			case *mp4meta.MP4Tag:
				t.SetYear(song.Year)
			}
		}

		// 设置歌词（仅 MP3 支持）
		if t, ok := tag.(*mp3meta.MP3Tag); ok {
			t.SetLyricist(song.Lyric)
		}

		// 创建临时文件
		f2, err := os.Create(tempPath)
		if err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("创建临时文件失败 [%s]: %w", tempPath, err))
			mu.Unlock()
			return
		}

		// 保存带标签的新文件
		if err = tag.Save(f2); err != nil {
			_ = f2.Close()
			_ = os.Remove(tempPath)
			mu.Lock()
			errs = append(errs, fmt.Errorf("保存元数据失败 [%s]: %w", rawPath, err))
			mu.Unlock()
			return
		}

		_ = f2.Close()
		_ = f.Close()

		// 用临时文件替换原文件
		if err = os.Rename(tempPath, rawPath); err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("替换文件失败 [%s]: %w", rawPath, err))
			mu.Unlock()
			return
		}

		utils.InfoWithFormat("[NCM] 🧩 已嵌入元数据: %s - %s", song.SongArtists, song.SongName)
	}

	for _, song := range songs {
		wg.Add(1)
		go processSong(song)
	}

	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("BeforeTidy 处理部分失败，共 %d 项: %v", len(errs), errs)
	}

	return nil
}

func (ncm *NetEaseProcessor) NeedRemoveDRM() bool {
	return false
}

func (ncm *NetEaseProcessor) DRMRemove() error {
	return nil
}

func (ncm *NetEaseProcessor) TidyMusic() error {
	files, err := os.ReadDir(ncm.tempDir)
	if err != nil {
		return fmt.Errorf("读取临时目录失败: %w", err)
	}
	if len(files) == 0 {
		utils.WarnWithFormat("[AppleMusic] ⚠️ 未找到待整理的音乐文件")
		return errors.New("未找到待整理的音乐文件")
	}

	switch ncm.cfg.Tidy.Mode {
	case 1:
		return ncm.tidyToLocal(files)
	case 2:
		return ncm.tidyToWebDAV(files, core.GlobalWebDAV)
	default:
		return fmt.Errorf("未知整理模式: %d", ncm.cfg.Tidy.Mode)
	}
}

func (ncm *NetEaseProcessor) EncryptedExts() []string {
	return []string{".ncm"}
}

func (ncm *NetEaseProcessor) DecryptedExts() []string {
	return []string{".flac", ".mp3", ".aac", ".m4a", ".ogg"}
}

/* ------------------------ 拓展方法 ------------------------ */

// downloadSingle 单曲下载
func (ncm *NetEaseProcessor) downloadSingle(musicID int, start time.Time) error {
	var err error

	utils.DebugWithFormat("[NCM] 获取单曲数据: ID=%d", musicID)
	detail, songURL, songLyric, err := ncm.FetchSongData(musicID, ncm.cfg)
	if err != nil {
		utils.ErrorWithFormat("[NCM] ❌ 获取歌曲数据失败: %v", err)
		return err
	}

	if len(detail.Songs) == 0 || len(songURL.Data) == 0 || songURL.Data[0].Url == "" {
		errMsg := "未获取到有效歌曲信息或歌曲无下载地址"
		utils.ErrorWithFormat("[NCM] ❌ %s: ID=%d", errMsg, musicID)
		return errors.New(errMsg)
	}

	// 构建歌曲元信息
	songInfo := ncm.buildSongInfo(ncm.cfg, detail, songURL, songLyric)
	fileName := ncm.safeFileName(songInfo)

	// 创建临时目录
	if err := processor.CreateOutputDir(ncm.tempDir); err != nil {
		utils.ErrorWithFormat("[NCM] ❌ 创建临时目录失败: %v", err)
		return err
	}

	// 下载文件
	utils.InfoWithFormat("[NCM] ⬇️ 开始下载: %s", fileName)
	if err := ncm.downloadFile(songURL.Data[0].Url, fileName, ncm.tempDir); err != nil {
		_ = processor.RemoveTempDir(ncm.tempDir)
		utils.ErrorWithFormat("[NCM] ❌ 下载失败: %v", err)
		return fmt.Errorf("下载失败: %w", err)
	}

	// 更新元信息列表
	ncm.songs = append(ncm.songs, songInfo)
	utils.InfoWithFormat("[NCM] ✅ 下载完成: %s （耗时 %v）", fileName, time.Since(start).Truncate(time.Millisecond))

	return nil
}

// downloadPlaylist 列表下载
func (ncm *NetEaseProcessor) downloadPlaylist(musicID int, start time.Time) error {
	utils.DebugWithFormat("[NCM] 获取歌单数据: ID=%d", musicID)
	detail, err := ncm.FetchPlaylistData(musicID, ncm.cfg)
	if err != nil {
		utils.ErrorWithFormat("[NCM] ❌ 获取歌单数据失败: %v", err)
		return err
	}

	if detail.Playlist.TrackCount == 0 {
		errMsg := "未获取到有效歌曲信息或歌曲无下载地址"
		utils.ErrorWithFormat("[NCM] ❌ %s: 歌单ID=%d", errMsg, musicID)
		return errors.New(errMsg)
	}

	utils.InfoWithFormat("[NCM] 开始下载歌单: %s (%d首)", detail.Playlist.Name, detail.Playlist.TrackCount)

	for index, track := range detail.Playlist.TrackIds {
		utils.InfoWithFormat("[NCM] 正在下载第%d首: ID=%d", index+1, track.Id)
		if err := ncm.downloadSingle(track.Id, start); err != nil {
			utils.ErrorWithFormat("[NCM] ❌ 歌单下载中断，第%d首下载失败: %v", index+1, err)
			return err
		}
	}

	utils.InfoWithFormat("[NCM] ✅ 歌单下载完成: %s （耗时 %v）", detail.Playlist.Name, time.Since(start).Truncate(time.Millisecond))
	return nil
}

// FetchSongData 获取单曲信息
func (ncm *NetEaseProcessor) FetchSongData(musicID int, cfg *config.Config) (*types.SongsDetailData, *types.SongsURLData, *types.SongLyricData, error) {
	utils.DebugWithFormat("[NCM] 请求歌曲信息中... ID=%d", musicID)

	batch := api.NewBatch(
		api.BatchAPI{Key: api.SongDetailAPI, Json: api.CreateSongDetailReqJson([]int{musicID})},
		api.BatchAPI{Key: api.SongUrlAPI, Json: api.CreateSongURLJson(api.SongURLConfig{Ids: []int{musicID}})},
		api.BatchAPI{Key: api.SongLyricAPI, Json: api.CreateSongLyricReqJson(musicID)},
	)

	cookiePath := filepath.Join(cfg.CookieCloud.CookieFilePath, cfg.CookieCloud.CookieFile)
	musicU := utils.GetCookieValue(cookiePath, ".music.163.com", "MUSIC_U")

	req := ncmutils.RequestData{}
	if musicU != "" {
		req.Cookies = []*http.Cookie{{Name: "MUSIC_U", Value: musicU}}
	}

	result := batch.Do(req)
	if result.Error != nil {
		return nil, nil, nil, fmt.Errorf("网易云API请求失败: %w", result.Error)
	}

	_, parsed := batch.Parse()

	var detail types.SongsDetailData
	var urls types.SongsURLData
	var lyrics types.SongLyricData

	if err := json.Unmarshal([]byte(parsed[api.SongDetailAPI]), &detail); err != nil {
		return nil, nil, nil, fmt.Errorf("解析歌曲详情失败: %w", err)
	}
	if err := json.Unmarshal([]byte(parsed[api.SongUrlAPI]), &urls); err != nil {
		return nil, nil, nil, fmt.Errorf("解析歌曲URL失败: %w", err)
	}
	if err := json.Unmarshal([]byte(parsed[api.SongLyricAPI]), &lyrics); err != nil {
		return nil, nil, nil, fmt.Errorf("解析歌曲歌词失败: %w", err)
	}

	utils.DebugWithFormat("[NCM] 歌曲信息获取成功: %s", detail.Songs[0].Name)
	return &detail, &urls, &lyrics, nil
}

// FetchPlaylistData 获取播放列表信息
func (ncm *NetEaseProcessor) FetchPlaylistData(musicID int, cfg *config.Config) (*types.PlaylistDetailData, error) {
	utils.DebugWithFormat("[NCM] 请求歌单信息中... ID=%d", musicID)

	batch := api.NewBatch(
		api.BatchAPI{Key: api.PlaylistDetailAPI, Json: api.CreatePlaylistDetailReqJson(musicID)},
	)

	cookiePath := filepath.Join(cfg.CookieCloud.CookieFilePath, cfg.CookieCloud.CookieFile)
	musicU := utils.GetCookieValue(cookiePath, ".music.163.com", "MUSIC_U")

	req := ncmutils.RequestData{}
	if musicU != "" {
		req.Cookies = []*http.Cookie{{Name: "MUSIC_U", Value: musicU}}
	}

	result := batch.Do(req)
	if result.Error != nil {
		return nil, fmt.Errorf("网易云API请求失败: %w", result.Error)
	}

	_, parsed := batch.Parse()

	var detail types.PlaylistDetailData

	if err := json.Unmarshal([]byte(parsed[api.PlaylistDetailAPI]), &detail); err != nil {
		return nil, fmt.Errorf("解析歌单详情失败: %w", err)
	}
	utils.DebugWithFormat("[NCM] 歌单信息获取成功: %s", detail.Playlist.Name)
	return &detail, nil
}

// downloadFile 下载文件
func (ncm *NetEaseProcessor) downloadFile(url, fileName, saveDir string) error {
	utils.DebugWithFormat("[NCM] 开始下载文件: %s", fileName)

	d := downloader.NewDownloader().
		SetSavePath(saveDir).
		SetBreakPoint(true).
		SetTimeOut(300 * time.Second)

	task, _ := d.NewDownloadTask(url)
	task.CleanTempFiles()
	task.ReplaceHostName(ncm.fixHost(task.GetHostName())).
		ForceHttps().
		ForceMultiThread()

	return task.SetFileName(fileName).Download()
}

// fixHost 主机名修正
func (ncm *NetEaseProcessor) fixHost(host string) string {
	replacer := strings.NewReplacer("m8.", "m7.", "m801.", "m701.", "m804.", "m701.", "m704.", "m701.")
	return replacer.Replace(host)
}

// buildSongInfo 构建歌曲信息
func (ncm *NetEaseProcessor) buildSongInfo(cfg *config.Config, detail *types.SongsDetailData, urls *types.SongsURLData, lyric *types.SongLyricData) *SongInfo {
	s := detail.Songs[0]
	u := urls.Data[0]
	// 整理方式
	tidy := processor.DetermineTidyType(cfg)

	ncmLyric := utils.ParseNCMLyric(lyric)
	if ncmLyric == "" {
		ncmLyric = "[00:00:00]此歌曲为没有填词的纯音乐，请您欣赏"
	}
	year := utils.ParseNCMYear(detail)

	return &SongInfo{
		SongName:    s.Name,
		SongArtists: utils.ParseArtist(s),
		SongAlbum:   s.Al.Name,
		FileExt:     ncm.detectExt(u.Url),
		MusicSize:   u.Size,
		Bitrate:     strconv.Itoa((8 * u.Size / (s.Dt / 1000)) / 1000),
		Duration:    s.Dt / 1000,
		PicUrl:      s.Al.PicUrl,
		Tidy:        tidy,
		Lyric:       ncmLyric,
		Year:        year,
	}
}

// detectExt 检测扩展名
func (ncm *NetEaseProcessor) detectExt(url string) string {
	if idx := strings.Index(url, "?"); idx != -1 {
		url = url[:idx]
	}
	ext := strings.ToLower(path.Ext(url))
	switch ext {
	case ".mp3", ".aac", ".m4a", ".ogg", ".flac":
		return strings.TrimPrefix(ext, ".")
	default:
		return "mp3"
	}
}

// safeFileName 合法的文件名
func (ncm *NetEaseProcessor) safeFileName(info *SongInfo) string {
	replacer := strings.NewReplacer("/", " ", "?", " ", "*", " ", ":", " ",
		"|", " ", "\\", " ", "<", " ", ">", " ", "\"", " ")
	return replacer.Replace(fmt.Sprintf("%s - %s.%s",
		strings.ReplaceAll(info.SongArtists, "/", ","),
		info.SongName,
		info.FileExt))
}

// safeTempFileName 合法临时文件路径
func (ncm *NetEaseProcessor) safeTempFileName(info *SongInfo) string {
	replacer := strings.NewReplacer("/", " ", "?", " ", "*", " ", ":", " ",
		"|", " ", "\\", " ", "<", " ", ">", " ", "\"", " ")
	return replacer.Replace(fmt.Sprintf("%s - %s.%s",
		strings.ReplaceAll(info.SongArtists, "/", ","),
		fmt.Sprintf("%s_temp", info.SongName),
		info.FileExt))
}

// 整理到本地
func (ncm *NetEaseProcessor) tidyToLocal(files []os.DirEntry) error {
	dstDir := ncm.cfg.Tidy.DistDir
	if dstDir == "" {
		_ = processor.RemoveTempDir(ncm.tempDir)
		return errors.New("未配置输出目录")
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		_ = processor.RemoveTempDir(ncm.tempDir)
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	for _, f := range files {
		if !utils.FilterMusicFile(f, ncm.EncryptedExts(), ncm.DecryptedExts()) {
			utils.DebugWithFormat("[AppleMusic] 跳过非音乐文件: %s", f.Name())
			continue
		}
		src := filepath.Join(ncm.tempDir, f.Name())
		dst := filepath.Join(dstDir, utils.SanitizeFileName(f.Name()))
		if err := os.Rename(src, dst); err != nil {
			utils.WarnWithFormat("[AppleMusic] ⚠️ 移动失败 %s → %s: %v", src, dst, err)
			continue
		}
		utils.InfoWithFormat("[AppleMusic] 📦 已整理: %s", dst)
	}
	// 清除临时目录
	err := processor.RemoveTempDir(ncm.tempDir)
	if err != nil {
		return err
	}
	return nil
}

// 整理到webdav
func (ncm *NetEaseProcessor) tidyToWebDAV(files []os.DirEntry, webdav *core.WebDAV) error {
	if webdav == nil {
		_ = processor.RemoveTempDir(ncm.tempDir)
		return errors.New("WebDAV 未初始化")
	}

	for _, f := range files {
		if !utils.FilterMusicFile(f, ncm.EncryptedExts(), ncm.DecryptedExts()) {
			utils.DebugWithFormat("[AppleMusic] 跳过非音乐文件: %s", f.Name())
			continue
		}

		filePath := filepath.Join(ncm.tempDir, f.Name())
		if err := webdav.Upload(filePath); err != nil {
			utils.WarnWithFormat("[AppleMusic] ☁️ 上传失败 %s: %v", f.Name(), err)
			continue
		}
		utils.InfoWithFormat("[AppleMusic] ☁️ 已上传: %s", f.Name())
	}
	// 清除临时目录
	err := processor.RemoveTempDir(ncm.tempDir)
	if err != nil {
		return err
	}
	return nil
}
