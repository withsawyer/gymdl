package music

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type SpotifyProcessor struct {
	cfg       *config.Config
	SongInfos []*SongInfo
}

func NewSpotifyProcessor(cfg *config.Config) *SpotifyProcessor {
	return &SpotifyProcessor{cfg: cfg}
}

/* ---------------------- 基础接口实现 ---------------------- */
func (am *SpotifyProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (am *SpotifyProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryMusic
}

func (am *SpotifyProcessor) Name() domain.LinkType {
	return domain.LinkSpotify
}

/* ------------------------ 下载逻辑 ------------------------ */
func (sp *SpotifyProcessor) DownloadMusic(url string) error {
	//TODO implement me
	panic("implement me")
}

func (sp *SpotifyProcessor) DownloadCommand(url string) *exec.Cmd {
	//TODO implement me
	panic("implement me")
}

func (sp *SpotifyProcessor) BeforeTidy() ([]*SongInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (sp *SpotifyProcessor) NeedRemoveDRM() bool {
	//TODO implement me
	panic("implement me")
}

func (sp *SpotifyProcessor) DRMRemove() error {
	//TODO implement me
	panic("implement me")
}

func (sp *SpotifyProcessor) TidyMusic() error {
	//TODO implement me
	panic("implement me")
}

func (sp *SpotifyProcessor) EncryptedExts() []string {
	//TODO implement me
	panic("implement me")
}

func (sp *SpotifyProcessor) DecryptedExts() []string {
	//TODO implement me
	panic("implement me")
}
