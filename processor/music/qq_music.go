package music

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/processor"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type QQMusicProcessor struct {
	cfg     *config.Config
	tempDir string
	songs   []*SongInfo
}

// Init  初始化
func (p *QQMusicProcessor) Init(cfg *config.Config) {
	p.cfg = cfg
	p.songs = make([]*SongInfo, 0)
	p.tempDir = processor.BuildOutputDir(QQTempDir)
}

/* ---------------------- 基础接口实现 ---------------------- */

func (p *QQMusicProcessor) Name() processor.LinkType {
	return processor.LinkQQMusic
}

func (p *QQMusicProcessor) Songs() []*SongInfo {
	return p.songs
}

/* ------------------------ 下载逻辑 ------------------------ */

func (p *QQMusicProcessor) DownloadMusic(url string, callback func(string)) error {
	// TODO implement me
	panic("implement me")
}

func (p *QQMusicProcessor) DownloadCommand(url string) *exec.Cmd {
	// TODO implement me
	panic("implement me")
}

func (p *QQMusicProcessor) BeforeTidy() error {
	// TODO implement me
	panic("implement me")
}

func (p *QQMusicProcessor) NeedRemoveDRM() bool {
	// TODO implement me
	panic("implement me")
}

func (p *QQMusicProcessor) DRMRemove() error {
	// TODO implement me
	panic("implement me")
}

func (p *QQMusicProcessor) TidyMusic() error {
	// TODO implement me
	panic("implement me")
}

func (p *QQMusicProcessor) EncryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

func (p *QQMusicProcessor) DecryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

/* ------------------------ 拓展方法 ------------------------ */
