package music

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type NetEaseProcessor struct {
	cfg       *config.Config
	SongInfos []*SongInfo
}

func NewNetEaseProcessor(cfg *config.Config) *NetEaseProcessor {
	return &NetEaseProcessor{cfg: cfg}
}

/* ---------------------- 基础接口实现 ---------------------- */

func (am *NetEaseProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (am *NetEaseProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryMusic
}

func (am *NetEaseProcessor) Name() domain.LinkType {
	return domain.LinkNetEase
}

func (am *NetEaseProcessor) Songs() []*SongInfo {
	return am.SongInfos
}

/* ------------------------ 下载逻辑 ------------------------ */

func (ncm *NetEaseProcessor) DownloadMusic(url string) error {
	// TODO implement me
	panic("implement me")
}

func (ncm *NetEaseProcessor) DownloadCommand(url string) *exec.Cmd {
	// TODO implement me
	panic("implement me")
}

func (ncm *NetEaseProcessor) BeforeTidy() error {
	// TODO implement me
	panic("implement me")
}

func (ncm *NetEaseProcessor) NeedRemoveDRM() bool {
	// TODO implement me
	panic("implement me")
}

func (ncm *NetEaseProcessor) DRMRemove() error {
	// TODO implement me
	panic("implement me")
}

func (ncm *NetEaseProcessor) TidyMusic() error {
	// TODO implement me
	panic("implement me")
}

func (ncm *NetEaseProcessor) EncryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

func (ncm *NetEaseProcessor) DecryptedExts() []string {
	// TODO implement me
	panic("implement me")
}
