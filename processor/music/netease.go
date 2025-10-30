package music

import (
	"github.com/nichuanfang/gymdl/utils"
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/processor"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type NetEaseProcessor struct {
	cfg     *config.Config
	songs   []*SongInfo
	tempDir string
}

// Init  初始化
func (p *NetEaseProcessor) Init(cfg *config.Config) {
	p.cfg = cfg
	p.songs = make([]*SongInfo, 0)
	p.tempDir = processor.BuildOutputDir(NCMTempDir)
}

/* ---------------------- 基础接口实现 ---------------------- */

func (p *NetEaseProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (p *NetEaseProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryMusic
}

func (p *NetEaseProcessor) Name() domain.LinkType {
	return domain.LinkNetEase
}

func (p *NetEaseProcessor) Songs() []*SongInfo {
	return p.songs
}

/* ------------------------ 下载逻辑 ------------------------ */

func (p *NetEaseProcessor) DownloadMusic(url string) error {

	utils.InfoWithFormat("55.69 MiB / 55.69 MiB [==============] 598.77 MiB p/s 100.00% 32s：%s", url)

	return nil
}

func (p *NetEaseProcessor) DownloadCommand(url string) *exec.Cmd {
	// TODO implement me
	panic("implement me")
}

func (p *NetEaseProcessor) BeforeTidy() error {
	// TODO implement me
	panic("implement me")
}

func (p *NetEaseProcessor) NeedRemoveDRM() bool {
	// TODO implement me
	panic("implement me")
}

func (p *NetEaseProcessor) DRMRemove() error {
	// TODO implement me
	panic("implement me")
}

func (p *NetEaseProcessor) TidyMusic() error {
	// TODO implement me
	panic("implement me")
}

func (p *NetEaseProcessor) EncryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

func (p *NetEaseProcessor) DecryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

/* ------------------------ 拓展方法 ------------------------ */
