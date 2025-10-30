package music

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/processor"
)

/* ---------------------- 结构体与构造方法 ---------------------- */
type AppleMusicProcessor struct {
	cfg       *config.Config
	SongInfos []*SongInfo
}
  
func NewAppleMusicProcessor(cfg *config.Config) processor.Processor {
	return &AppleMusicProcessor{cfg: cfg}
}

/* ---------------------- 基础接口实现 ---------------------- */

func (am *AppleMusicProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (am *AppleMusicProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryMusic
}

func (am *AppleMusicProcessor) Name() domain.LinkType {
	return domain.LinkAppleMusic
}

func (am *AppleMusicProcessor) Songs() []*SongInfo {
	return am.SongInfos
}

/* ------------------------ 下载逻辑 ------------------------ */

func (am *AppleMusicProcessor) DownloadMusic(url string) error {
	// TODO implement me
	panic("implement me")
}

func (am *AppleMusicProcessor) DownloadCommand(url string) *exec.Cmd {
	// TODO implement me
	panic("implement me")
}

func (am *AppleMusicProcessor) BeforeTidy() error {
	// TODO implement me
	panic("implement me")
}

func (am *AppleMusicProcessor) NeedRemoveDRM() bool {
	// TODO implement me
	panic("implement me")
}

func (am *AppleMusicProcessor) DRMRemove() error {
	// TODO implement me
	panic("implement me")
}

func (am *AppleMusicProcessor) TidyMusic() error {
	// TODO implement me
	panic("implement me")
}

func (am *AppleMusicProcessor) EncryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

func (am *AppleMusicProcessor) DecryptedExts() []string {
	// TODO implement me
	panic("implement me")
}
