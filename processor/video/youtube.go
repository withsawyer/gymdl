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

func NewYoutubeProcessor(cfg *config.Config, baseTempDir string) (*YoutubeProcessor, error) {
	dir, err := processor.BuildOutputDir(baseTempDir)
	if err != nil {
		return nil, err
	}
	return &YoutubeProcessor{cfg: cfg, tempDir: dir}, nil
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
