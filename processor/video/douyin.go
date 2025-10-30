package video

import (
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
)

// 抖音下载

/* ---------------------- 结构体与构造方法 ---------------------- */

type DouyinProcessor struct {
	cfg    *config.Config
	videos []*VideoInfo
}

func NewDouyinProcessor(cfg *config.Config) *DouyinProcessor {
	return &DouyinProcessor{cfg: cfg}
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
