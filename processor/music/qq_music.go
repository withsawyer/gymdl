package music

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/processor"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type QQMusicProcessor struct {
	cfg     *config.Config
	tempDir string
	songs   []*SongInfo
}

func NewQQMusicProcessor(cfg *config.Config, baseTempDir string) processor.Processor {
	return &QQMusicProcessor{cfg: cfg, tempDir: processor.BuildOutputDir(baseTempDir)}
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

func (am *QQMusicProcessor) Songs() []*SongInfo {
	return am.songs
}

/* ------------------------ 下载逻辑 ------------------------ */

func (qm *QQMusicProcessor) DownloadMusic(url string) error {
	// TODO implement me
	panic("implement me")
}

func (qm *QQMusicProcessor) DownloadCommand(url string) *exec.Cmd {
	// TODO implement me
	panic("implement me")
}

func (qm *QQMusicProcessor) BeforeTidy() error {
	// TODO implement me
	panic("implement me")
}

func (qm *QQMusicProcessor) NeedRemoveDRM() bool {
	// TODO implement me
	panic("implement me")
}

func (qm *QQMusicProcessor) DRMRemove() error {
	// TODO implement me
	panic("implement me")
}

func (qm *QQMusicProcessor) TidyMusic() error {
	// TODO implement me
	panic("implement me")
}

func (qm *QQMusicProcessor) EncryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

func (qm *QQMusicProcessor) DecryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

/* ------------------------ 拓展方法 ------------------------ */
