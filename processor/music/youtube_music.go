package music

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/processor"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type YoutubeMusicProcessor struct {
	cfg     *config.Config
	tempDir string
	songs   []*SongInfo
}

// Init  初始化
func (p *YoutubeMusicProcessor) Init(cfg *config.Config) {
	p.cfg = cfg
	p.songs = make([]*SongInfo, 0)
	p.tempDir = processor.BuildOutputDir(YoutubeTempDir)
}

/* ---------------------- 基础接口实现 ---------------------- */

func (p *YoutubeMusicProcessor) Name() processor.LinkType {
	return processor.LinkYoutubeMusic
}

func (p *YoutubeMusicProcessor) Songs() []*SongInfo {
	return p.songs
}

/* ------------------------ 下载逻辑 ------------------------ */
func (p *YoutubeMusicProcessor) DownloadMusic(url string) error {
	// TODO implement me
	panic("implement me")
}

func (p *YoutubeMusicProcessor) DownloadCommand(url string) *exec.Cmd {
	// TODO implement me
	panic("implement me")
}

func (p *YoutubeMusicProcessor) BeforeTidy() error {
	// TODO implement me
	panic("implement me")
}

func (p *YoutubeMusicProcessor) NeedRemoveDRM() bool {
	// TODO implement me
	panic("implement me")
}

func (p *YoutubeMusicProcessor) DRMRemove() error {
	// TODO implement me
	panic("implement me")
}

func (p *YoutubeMusicProcessor) TidyMusic() error {
	// TODO implement me
	panic("implement me")
}

func (p *YoutubeMusicProcessor) EncryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

func (p *YoutubeMusicProcessor) DecryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

/* ------------------------ 拓展方法 ------------------------ */
