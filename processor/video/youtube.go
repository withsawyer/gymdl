package video

import (
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/processor"
)

// youtube下载

/* ---------------------- 结构体与构造方法 ---------------------- */

type YoutubeProcessor struct {
	cfg     *config.Config
	tempDir string
	videos  []*VideoInfo
}

func NewYoutubeProcessor(cfg *config.Config, baseTempDir string) processor.Processor {
	return &YoutubeProcessor{cfg: cfg, tempDir: processor.BuildOutputDir(baseTempDir)}
}

/* ---------------------- 基础接口实现 ---------------------- */

func (yt *YoutubeProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (yt *YoutubeProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryVideo
}

func (yt *YoutubeProcessor) Name() domain.LinkType {
	return domain.LinkYoutube
}

func (yt *YoutubeProcessor) Videos() []*VideoInfo {
	return yt.videos
}

/* ------------------------ 下载逻辑 ------------------------ */

/* ------------------------ 拓展方法 ------------------------ */
