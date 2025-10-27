package handler

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
	"time"

	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/types"
	ncmutils "github.com/XiaoMengXinX/Music163Api-Go/utils"
	downloader "github.com/XiaoMengXinX/SimpleDownloader"
	"github.com/gcottom/audiometa/v3"
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/core/constants"
	"github.com/nichuanfang/gymdl/utils"
)

type NCMHandler struct{}

/* ---------------------- 核心接口实现 ---------------------- */

func (ncm *NCMHandler) Platform() string { return "NetEaseCloudMusic" }

func (ncm *NCMHandler) DownloadMusic(url string, cfg *config.Config) (*SongInfo, error) {
	start := time.Now()
	utils.InfoWithFormat("[NCM] 🎵 开始下载: %s", url)

	// 创建缓存目录
	if err := os.MkdirAll(constants.NCMTempDir, 0755); err != nil {
		return nil, fmt.Errorf("创建临时目录失败: %w", err)
	}
	musicID := utils.ParseMusicID(url)
	detail, songURL, err := ncm.fetchSongData(musicID, cfg)
	if err != nil {
		return nil, err
	}

	if len(detail.Songs) == 0 || len(songURL.Data) == 0 || songURL.Data[0].Url == "" {
		return nil, errors.New("未获取到有效歌曲信息或歌曲无下载地址")
	}

	songInfo := ncm.buildSongInfo(cfg, detail, songURL)
	fileName := ncm.safeFileName(songInfo)

	if err := ncm.downloadFile(songURL.Data[0].Url, fileName, constants.NCMTempDir); err != nil {
		return nil, fmt.Errorf("下载失败: %w", err)
	}

	utils.InfoWithFormat("[NCM] ✅ 下载完成: %s （耗时 %v）", fileName, time.Since(start).Truncate(time.Millisecond))
	return songInfo, nil
}

/* ---------------------- 数据获取 ---------------------- */

func (ncm *NCMHandler) fetchSongData(musicID int, cfg *config.Config) (*types.SongsDetailData, *types.SongsURLData, error) {
	utils.DebugWithFormat("[NCM] 请求歌曲信息中... ID=%d", musicID)

	batch := api.NewBatch(
		api.BatchAPI{Key: api.SongDetailAPI, Json: api.CreateSongDetailReqJson([]int{musicID})},
		api.BatchAPI{Key: api.SongUrlAPI, Json: api.CreateSongURLJson(api.SongURLConfig{Ids: []int{musicID}})},
	)

	cookiePath := filepath.Join(cfg.CookieCloud.CookieFilePath, cfg.CookieCloud.CookieFile)
	musicU := utils.GetCookieValue(cookiePath, ".music.163.com", "MUSIC_U")

	req := ncmutils.RequestData{}
	if musicU != "" {
		req.Cookies = []*http.Cookie{{Name: "MUSIC_U", Value: musicU}}
	}

	result := batch.Do(req)
	if result.Error != nil {
		return nil, nil, fmt.Errorf("网易云API请求失败: %w", result.Error)
	}

	_, parsed := batch.Parse()

	var detail types.SongsDetailData
	var urls types.SongsURLData

	if err := json.Unmarshal([]byte(parsed[api.SongDetailAPI]), &detail); err != nil {
		return nil, nil, fmt.Errorf("解析歌曲详情失败: %w", err)
	}
	if err := json.Unmarshal([]byte(parsed[api.SongUrlAPI]), &urls); err != nil {
		return nil, nil, fmt.Errorf("解析歌曲URL失败: %w", err)
	}

	utils.DebugWithFormat("[NCM] 歌曲信息获取成功: %s", detail.Songs[0].Name)
	return &detail, &urls, nil
}

/* ---------------------- 数据构建 ---------------------- */

func (ncm *NCMHandler) buildSongInfo(cfg *config.Config, detail *types.SongsDetailData, urls *types.SongsURLData) *SongInfo {
	s := detail.Songs[0]
	u := urls.Data[0]
    // 整理方式
	tidy := determineTidyType(cfg)

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
	}
}

/* ---------------------- 下载逻辑 ---------------------- */

func (ncm *NCMHandler) downloadFile(url, fileName, saveDir string) error {
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

func (ncm *NCMHandler) fixHost(host string) string {
	replacer := strings.NewReplacer("m8.", "m7.", "m801.", "m701.", "m804.", "m701.", "m704.", "m701.")
	return replacer.Replace(host)
}

/* ---------------------- 工具函数 ---------------------- */

func (ncm *NCMHandler) detectExt(url string) string {
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

func (ncm *NCMHandler) safeFileName(info *SongInfo) string {
	replacer := strings.NewReplacer("/", " ", "?", " ", "*", " ", ":", " ",
		"|", " ", "\\", " ", "<", " ", ">", " ", "\"", " ")
	return replacer.Replace(fmt.Sprintf("%s - %s.%s",
		strings.ReplaceAll(info.SongArtists, "/", ","),
		info.SongName,
		info.FileExt))
}

func (ncm *NCMHandler) safeTempFileName(info *SongInfo) string {
	replacer := strings.NewReplacer("/", " ", "?", " ", "*", " ", ":", " ",
		"|", " ", "\\", " ", "<", " ", ">", " ", "\"", " ")
	return replacer.Replace(fmt.Sprintf("%s - %s.%s",
		strings.ReplaceAll(info.SongArtists, "/", ","),
		fmt.Sprintf("%s_temp", info.SongName),
		info.FileExt))
}

/* ---------------------- 整理逻辑 ---------------------- */

func (ncm *NCMHandler) TidyMusic(cfg *config.Config, webdav *core.WebDAV, songInfo *SongInfo) error {
	files, err := os.ReadDir(constants.NCMTempDir)
	if err != nil {
		return fmt.Errorf("读取临时目录失败: %w", err)
	}
	if len(files) == 0 {
		return errors.New("未找到待整理的音乐文件")
	}

	switch cfg.MusicTidy.Mode {
	case 1:
		return ncm.tidyToLocal(cfg, files)
	case 2:
		return ncm.tidyToWebDAV(cfg, files, webdav)
	default:
		return fmt.Errorf("未知整理模式: %d", cfg.MusicTidy.Mode)
	}
}

func (ncm *NCMHandler) tidyToLocal(cfg *config.Config, files []os.DirEntry) error {
	if cfg.MusicTidy.DistDir == "" {
		return errors.New("未配置输出目录")
	}
	if err := os.MkdirAll(cfg.MusicTidy.DistDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	for _, file := range files {
		if !utils.FilterMusicFile(file, ncm.EncryptedExts(), ncm.DecryptedExts()) {
			continue
		}
		src := filepath.Join(constants.NCMTempDir, file.Name())
		dst := filepath.Join(cfg.MusicTidy.DistDir, utils.SanitizeFileName(file.Name()))
		if err := utils.MoveFile(src, dst); err != nil {
			utils.WarnWithFormat("[NCM] ⚠️ 移动失败 %s → %s: %v", src, dst, err)
			continue
		}
		utils.InfoWithFormat("[NCM] 📦 已整理: %s", dst)
	}
	return nil
}

func (ncm *NCMHandler) tidyToWebDAV(cfg *config.Config, files []os.DirEntry, webdav *core.WebDAV) error {
	if webdav == nil {
		return errors.New("WebDAV 未初始化")
	}
	for _, file := range files {
		if !utils.FilterMusicFile(file, ncm.EncryptedExts(), ncm.DecryptedExts()) {
			continue
		}
		filePath := filepath.Join(constants.NCMTempDir, file.Name())
		if err := webdav.Upload(filePath); err != nil {
			utils.WarnWithFormat("[NCM] ☁️ 上传失败: %s (%v)", file.Name(), err)
			continue
		}
		utils.InfoWithFormat("[NCM] ☁️ 已上传: %s", file.Name())
		_ = os.Remove(filePath)
	}
	return nil
}

/* ---------------------- 元数据嵌入 ---------------------- */

func (ncm *NCMHandler) BeforeTidy(cfg *config.Config, info *SongInfo) error {
	// 原文件路径
	rawPath := filepath.Join(constants.NCMTempDir, ncm.safeFileName(info))
	// 临时文件路径
	tempPath := filepath.Join(constants.NCMTempDir, ncm.safeTempFileName(info))

	// 打开原文件
	f, err := os.Open(rawPath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	// 读取音频标签
	tag, err := audiometa.OpenTag(f)
	if err != nil {
		return fmt.Errorf("读取音频标签失败: %w", err)
	}

	// 设置元数据
	tag.SetArtist(info.SongArtists)
	tag.SetTitle(info.SongName)
	tag.SetAlbum(info.SongAlbum)

	// 设置封面图（如果获取成功）
	if image, err := utils.FetchImage(info.PicUrl); err == nil {
		tag.SetCoverArt(image)
	}

	// 创建临时文件（用于保存修改后的数据）
	f2, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
	}

	// 将标签保存到临时文件
	if err = tag.Save(f2); err != nil {
		return fmt.Errorf("保存元数据失败: %w", err)
	}
	_ = f2.Close()
	_ = f.Close()
	_ = os.Remove(rawPath)
	err = os.Rename(tempPath, rawPath)
	if err != nil {
		return err
	}

	utils.InfoWithFormat("[NCM] 🧩 已嵌入元数据: %s - %s", info.SongArtists, info.SongName)
	return nil
}

/* ---------------------- 其他接口实现 ---------------------- */

func (ncm *NCMHandler) DownloadCommand(cfg *config.Config, url string) *exec.Cmd { return nil }
func (ncm *NCMHandler) NeedRemoveDRM(cfg *config.Config) bool                    { return false }
func (ncm *NCMHandler) DRMRemove(cfg *config.Config, songInfo *SongInfo) error   { return nil }
func (ncm *NCMHandler) EncryptedExts() []string                                  { return []string{".ncm"} }
func (ncm *NCMHandler) DecryptedExts() []string {
	return []string{".flac", ".mp3", ".aac", ".m4a", ".ogg"}
}
