package video

import (
	"errors"
	"fmt"
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/processor"
	"github.com/nichuanfang/gymdl/utils"
	"os"
	"path/filepath"
)

// youtube下载

/* ---------------------- 结构体与构造方法 ---------------------- */

type YoutubeProcessor struct {
	cfg     *config.Config
	tempDir string
	videos  []*VideoInfo
}

// Init  初始化
func (p *YoutubeProcessor) Init(cfg *config.Config) {
	p.cfg = cfg
	p.videos = make([]*VideoInfo, 0)
	p.tempDir = processor.BuildOutputDir(YoutubeTempDir)
}

/* ---------------------- 基础接口实现 ---------------------- */

func (p *YoutubeProcessor) Name() processor.LinkType {
	return processor.LinkYoutube
}

func (p *YoutubeProcessor) Videos() []*VideoInfo {
	return p.videos
}

/* ------------------------ 下载逻辑 ------------------------ */

func (p *YoutubeProcessor) Download(url string, callback func(progress string)) error {
	utils.Debugf("开始下载视频:%s", url)
	return nil
}

/* ------------------------ 拓展方法 ------------------------ */

func (p *YoutubeProcessor) Tidy() error {
	files, err := os.ReadDir(p.tempDir)
	if err != nil {
		return fmt.Errorf("读取临时目录失败: %w", err)
	}
	if len(files) == 0 {
		utils.WarnWithFormat("[DouYinVideo] ⚠️ 未找到待整理的资源文件")
		return errors.New("未找到待整理的资源文件")
	}

	switch p.cfg.Tidy.Mode {
	case 1:
		return p.tidyToLocal(files)
	case 2:
		return p.tidyToWebDAV(files, core.GlobalWebDAV)
	default:
		return fmt.Errorf("未知整理模式: %d", p.cfg.Tidy.Mode)
	}
}

// 整理到本地
func (p *YoutubeProcessor) tidyToLocal(files []os.DirEntry) error {
	dstDir := p.cfg.Tidy.DistDir
	if dstDir == "" {
		_ = processor.RemoveTempDir(p.tempDir)
		return errors.New("未配置输出目录")
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		_ = processor.RemoveTempDir(p.tempDir)
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	for _, f := range files {
		src := filepath.Join(p.tempDir, f.Name())
		dst := filepath.Join(dstDir, utils.SanitizeFileName(f.Name()))
		if err := os.Rename(src, dst); err != nil {
			utils.WarnWithFormat("[DouYinVideo] ⚠️ 移动失败 %s → %s: %v", src, dst, err)
			continue
		}
		utils.InfoWithFormat("[DouYinVideo] 📦 已整理: %s", dst)
	}
	//清除临时目录
	err := processor.RemoveTempDir(p.tempDir)
	if err != nil {
		utils.WarnWithFormat("[DouYinVideo] ⚠️ 删除临时目录失败: %s (%v)", p.tempDir, err)
		return err
	}
	utils.DebugWithFormat("[DouYinVideo] 🧹 已删除临时目录: %s", p.tempDir)
	return nil
}

// 整理到webdav
func (p *YoutubeProcessor) tidyToWebDAV(files []os.DirEntry, webdav *core.WebDAV) error {
	if webdav == nil {
		_ = processor.RemoveTempDir(p.tempDir)
		return errors.New("WebDAV 未初始化")
	}

	for _, f := range files {
		filePath := filepath.Join(p.tempDir, f.Name())
		if err := webdav.Upload(filePath); err != nil {
			utils.WarnWithFormat("[DouYinVideo] ☁️ 上传失败 %s: %v", f.Name(), err)
			continue
		}
		utils.InfoWithFormat("[DouYinVideo] ☁️ 已上传: %s", f.Name())
	}
	//清除临时目录
	err := processor.RemoveTempDir(p.tempDir)
	if err != nil {
		utils.WarnWithFormat("[DouYinVideo] ⚠️ 删除临时目录失败: %s (%v)", p.tempDir, err)
		return err
	}
	utils.DebugWithFormat("[DouYinVideo] 🧹 已删除临时目录: %s", p.tempDir)
	return nil
}
