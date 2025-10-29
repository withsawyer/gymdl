package video

//bilibili下载

import (
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type BiliBiliProcessor struct {
	cfg *config.Config
}

func NewBiliBiliProcessor(cfg *config.Config) *BiliBiliProcessor {
	return &BiliBiliProcessor{cfg: cfg}
}

/* ---------------------- 基础接口实现 ---------------------- */
func (am *BiliBiliProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (am *BiliBiliProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryVideo
}

func (am *BiliBiliProcessor) Name() domain.LinkType {
	return domain.LinkBilibili
}

/* ------------------------ 下载逻辑 ------------------------ */
