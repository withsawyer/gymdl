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

func NewSpotifyProcessor(cfg *config.Config, baseTempDir string) (processor.Processor, error) {
	dir, err := processor.BuildOutputDir(baseTempDir)
	if err != nil {
		return nil, err
	}
	return &SpotifyProcessor{cfg: cfg, tempDir: dir}, nil
}

/* ---------------------- 基础接口实现 ---------------------- */
func (sp *SpotifyProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (sp *SpotifyProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryMusic
}

func (sp *SpotifyProcessor) Name() domain.LinkType {
	return domain.LinkSpotify
}

func (sp *SpotifyProcessor) Songs() []*SongInfo {
	return sp.songs
}

/* ------------------------ 下载逻辑 ------------------------ */
func (sp *SpotifyProcessor) DownloadMusic(url string) error {
	// TODO implement me
	panic("implement me")
}

func (sp *SpotifyProcessor) DownloadCommand(url string) *exec.Cmd {
	// TODO implement me
	panic("implement me")
}

func (sp *SpotifyProcessor) BeforeTidy() error {
	// TODO implement me
	panic("implement me")
}

func (sp *SpotifyProcessor) NeedRemoveDRM() bool {
	// TODO implement me
	panic("implement me")
}

func (sp *SpotifyProcessor) DRMRemove() error {
	// TODO implement me
	panic("implement me")
}

func (sp *SpotifyProcessor) TidyMusic() error {
	// TODO implement me
	panic("implement me")
}

func (sp *SpotifyProcessor) EncryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

func (sp *SpotifyProcessor) DecryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

/* ------------------------ 拓展方法 ------------------------ */
