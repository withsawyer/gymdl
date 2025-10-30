package video

// bilibili下载

import (
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
    "github.com/nichuanfang/gymdl/processor"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type BiliBiliProcessor struct {
	cfg     *config.Config
    tempDir string
	videos  []*VideoInfo
}

func NewBiliBiliProcessor(cfg *config.Config,baseTempDir string) (*BiliBiliProcessor,error) {
    dir,err := processor.BuildOutputDir(baseTempDir)
    if err != nil {
        return nil,err
    }
	return &BiliBiliProcessor{cfg: cfg, tempDir: dir},nil
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

func (am *BiliBiliProcessor) Videos() []*VideoInfo {
	return am.videos
}

/* ------------------------ 下载逻辑 ------------------------ */

/* ------------------------ 拓展方法 ------------------------ */
