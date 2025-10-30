package video

import (
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
)

// 小红书下载

/* ---------------------- 结构体与构造方法 ---------------------- */

type XiaohongshuProcessor struct {
	cfg    *config.Config
	videos []*VideoInfo
}

func NewXiaohongshuProcessor(cfg *config.Config) *XiaohongshuProcessor {
	return &XiaohongshuProcessor{cfg: cfg}
}

/* ---------------------- 基础接口实现 ---------------------- */
func (xhs *XiaohongshuProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (xhs *XiaohongshuProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryVideo
}

func (xhs *XiaohongshuProcessor) Name() domain.LinkType {
	return domain.LinkXiaohongshu
}

func (xhs *XiaohongshuProcessor) Videos() []*VideoInfo {
	return xhs.videos
}

/* ------------------------ 下载逻辑 ------------------------ */
