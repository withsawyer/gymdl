package video

import (
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
)

//小红书下载

/* ---------------------- 结构体与构造方法 ---------------------- */

type XiaohongshuProcessor struct {
	cfg *config.Config
}

func NewXiaohongshuProcessor(cfg *config.Config) *XiaohongshuProcessor {
	return &XiaohongshuProcessor{cfg: cfg}
}

/* ---------------------- 基础接口实现 ---------------------- */
func (b *XiaohongshuProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (am *XiaohongshuProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryVideo
}

func (am *XiaohongshuProcessor) Name() domain.LinkType {
	return domain.LinkXiaohongshu
}

/* ------------------------ 下载逻辑 ------------------------ */
