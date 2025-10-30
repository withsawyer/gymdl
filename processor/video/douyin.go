package video

import (
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
    "github.com/nichuanfang/gymdl/processor"
)

// 抖音下载

/* ---------------------- 结构体与构造方法 ---------------------- */

type DouyinProcessor struct {
	cfg     *config.Config
    tempDir string
	videos  []*VideoInfo
}

func NewDouyinProcessor(cfg *config.Config,baseTempDir string) processor.Processor {
    return &DouyinProcessor{cfg: cfg, tempDir: processor.BuildOutputDir(baseTempDir)}
}

/* ---------------------- 基础接口实现 ---------------------- */
func (dy *DouyinProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (dy *DouyinProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryVideo
}

func (dy *DouyinProcessor) Name() domain.LinkType {
	return domain.LinkDouyin
}

func (dy *DouyinProcessor) Videos() []*VideoInfo {
	return dy.videos
}

/* ------------------------ 下载逻辑 ------------------------ */

/* ------------------------ 拓展方法 ------------------------ */