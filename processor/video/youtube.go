package video

import (
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/processor"
	"github.com/nichuanfang/gymdl/utils"
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

func (p *YoutubeProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (p *YoutubeProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryVideo
}

func (p *YoutubeProcessor) Name() domain.LinkType {
	return domain.LinkYoutube
}

func (p *YoutubeProcessor) Download(url string) error {
	utils.Debugf("开始下载视频:%s", url)
	return nil
}

func (p *YoutubeProcessor) Videos() []*VideoInfo {
	return p.videos
}

/* ------------------------ 下载逻辑 ------------------------ */

/* ------------------------ 拓展方法 ------------------------ */
