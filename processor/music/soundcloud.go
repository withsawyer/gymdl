package music

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/processor"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type SoundCloudProcessor struct {
	cfg     *config.Config
	tempDir string
	songs   []*SongInfo
}

// Init  初始化
func (p *SoundCloudProcessor) Init(cfg *config.Config) {
	p.cfg = cfg
	p.songs = make([]*SongInfo, 0)
	p.tempDir = processor.BuildOutputDir(SoundcloudTempDir)
}

/* ---------------------- 基础接口实现 ---------------------- */

func (p *SoundCloudProcessor) Name() processor.LinkType {
	return processor.LinkSoundcloud
}

func (p *SoundCloudProcessor) Songs() []*SongInfo {
	return p.songs
}

/* ------------------------ 下载逻辑 ------------------------ */

func (p *SoundCloudProcessor) DownloadMusic(url string, callback func(string)) error {
	// TODO implement me
	panic("implement me")
}

func (p *SoundCloudProcessor) DownloadCommand(url string) *exec.Cmd {
	// TODO implement me
	panic("implement me")
}

func (p *SoundCloudProcessor) BeforeTidy() error {
	// TODO implement me
	panic("implement me")
}

func (p *SoundCloudProcessor) NeedRemoveDRM() bool {
	// TODO implement me
	panic("implement me")
}

func (p *SoundCloudProcessor) DRMRemove() error {
	// TODO implement me
	panic("implement me")
}

func (p *SoundCloudProcessor) TidyMusic() error {
	// TODO implement me
	panic("implement me")
}

func (p *SoundCloudProcessor) EncryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

func (p *SoundCloudProcessor) DecryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

/* ------------------------ 拓展方法 ------------------------ */
