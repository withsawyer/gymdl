package music

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/processor"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type SpotifyProcessor struct {
	cfg     *config.Config
	tempDir string
	songs   []*SongInfo
}

// Init  初始化
func (p *SpotifyProcessor) Init(cfg *config.Config) {
	p.cfg = cfg
	p.songs = make([]*SongInfo, 0)
	p.tempDir = processor.BuildOutputDir(SpotifyTempDir)
}

/* ---------------------- 基础接口实现 ---------------------- */
func (p *SpotifyProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (p *SpotifyProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryMusic
}

func (p *SpotifyProcessor) Name() domain.LinkType {
	return domain.LinkSpotify
}

func (p *SpotifyProcessor) Songs() []*SongInfo {
	return p.songs
}

/* ------------------------ 下载逻辑 ------------------------ */
func (p *SpotifyProcessor) DownloadMusic(url string) error {
	// TODO implement me
	panic("implement me")
}

func (p *SpotifyProcessor) DownloadCommand(url string) *exec.Cmd {
	// TODO implement me
	panic("implement me")
}

func (p *SpotifyProcessor) BeforeTidy() error {
	// TODO implement me
	panic("implement me")
}

func (p *SpotifyProcessor) NeedRemoveDRM() bool {
	// TODO implement me
	panic("implement me")
}

func (p *SpotifyProcessor) DRMRemove() error {
	// TODO implement me
	panic("implement me")
}

func (p *SpotifyProcessor) TidyMusic() error {
	// TODO implement me
	panic("implement me")
}

func (p *SpotifyProcessor) EncryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

func (p *SpotifyProcessor) DecryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

/* ------------------------ 拓展方法 ------------------------ */
