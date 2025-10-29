package video

import (
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
)

//youtube下载

/* ---------------------- 结构体与构造方法 ---------------------- */

type YoutubeProcessor struct {
	cfg *config.Config
}

func NewYoutubeProcessor(cfg *config.Config) *YoutubeProcessor {
	return &YoutubeProcessor{cfg: cfg}
}

/* ---------------------- 基础接口实现 ---------------------- */

func (b *YoutubeProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (am *YoutubeProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryVideo
}

func (am *YoutubeProcessor) Name() domain.LinkType {
	return domain.LinkYoutube
}

/* ------------------------ 下载逻辑 ------------------------ */
