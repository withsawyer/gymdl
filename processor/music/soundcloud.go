package music

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/processor"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type SoundCloudProcessor struct {
	cfg     *config.Config
	tempDir string
	songs   []*SongInfo
}

func NewSoundCloudProcessor(cfg *config.Config, baseTempDir string) processor.Processor {
	return &SoundCloudProcessor{cfg: cfg, tempDir: processor.BuildOutputDir(baseTempDir)}
}

/* ---------------------- 基础接口实现 ---------------------- */
func (sc *SoundCloudProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (sc *SoundCloudProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryMusic
}

func (sc *SoundCloudProcessor) Name() domain.LinkType {
	return domain.LinkSoundcloud
}

func (sc *SoundCloudProcessor) Songs() []*SongInfo {
	return sc.songs
}

/* ------------------------ 下载逻辑 ------------------------ */

func (sc *SoundCloudProcessor) DownloadMusic(url string) error {
	// TODO implement me
	panic("implement me")
}

func (sc *SoundCloudProcessor) DownloadCommand(url string) *exec.Cmd {
	// TODO implement me
	panic("implement me")
}

func (sc *SoundCloudProcessor) BeforeTidy() error {
	// TODO implement me
	panic("implement me")
}

func (sc *SoundCloudProcessor) NeedRemoveDRM() bool {
	// TODO implement me
	panic("implement me")
}

func (sc *SoundCloudProcessor) DRMRemove() error {
	// TODO implement me
	panic("implement me")
}

func (sc *SoundCloudProcessor) TidyMusic() error {
	// TODO implement me
	panic("implement me")
}

func (sc *SoundCloudProcessor) EncryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

func (sc *SoundCloudProcessor) DecryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

/* ------------------------ 拓展方法 ------------------------ */
