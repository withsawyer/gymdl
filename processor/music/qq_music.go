package music

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type QQMusicProcessor struct {
	cfg       *config.Config
	SongInfos []*SongInfo
}

func NewQQMusicProcessor(cfg *config.Config) *QQMusicProcessor {
	return &QQMusicProcessor{cfg: cfg}
}

/* ---------------------- 基础接口实现 ---------------------- */
func (am *QQMusicProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (am *QQMusicProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryMusic
}

func (am *QQMusicProcessor) Name() domain.LinkType {
	return domain.LinkQQMusic
}

/* ------------------------ 下载逻辑 ------------------------ */

func (qm *QQMusicProcessor) DownloadMusic(url string) error {
	//TODO implement me
	panic("implement me")
}

func (qm *QQMusicProcessor) DownloadCommand(url string) *exec.Cmd {
	//TODO implement me
	panic("implement me")
}

func (qm *QQMusicProcessor) BeforeTidy() ([]*SongInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (qm *QQMusicProcessor) NeedRemoveDRM() bool {
	//TODO implement me
	panic("implement me")
}

func (qm *QQMusicProcessor) DRMRemove() error {
	//TODO implement me
	panic("implement me")
}

func (qm *QQMusicProcessor) TidyMusic() error {
	//TODO implement me
	panic("implement me")
}

func (qm *QQMusicProcessor) EncryptedExts() []string {
	//TODO implement me
	panic("implement me")
}

func (qm *QQMusicProcessor) DecryptedExts() []string {
	//TODO implement me
	panic("implement me")
}
