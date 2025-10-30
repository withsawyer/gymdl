package music

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type YoutubeMusicProcessor struct {
	cfg       *config.Config
	SongInfos []*SongInfo
}

func (y *YoutubeMusicProcessor) IsPlaylist() bool {
	// TODO implement me
	panic("implement me")
}

func NewYoutubeMusicProcessor(cfg *config.Config) *YoutubeMusicProcessor {
	return &YoutubeMusicProcessor{cfg: cfg}
}

/* ---------------------- 基础接口实现 ---------------------- */
func (am *YoutubeMusicProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (am *YoutubeMusicProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryMusic
}

func (am *YoutubeMusicProcessor) Name() domain.LinkType {
	return domain.LinkYoutubeMusic
}

func (sp *YoutubeMusicProcessor) Songs() []*SongInfo {
	return sp.SongInfos
}

/* ------------------------ 下载逻辑 ------------------------ */
func (ytm *YoutubeMusicProcessor) DownloadMusic(url string) error {
	// TODO implement me
	panic("implement me")
}

func (ytm *YoutubeMusicProcessor) DownloadCommand(url string) *exec.Cmd {
	// TODO implement me
	panic("implement me")
}

func (ytm *YoutubeMusicProcessor) BeforeTidy() error {
	// TODO implement me
	panic("implement me")
}

func (ytm *YoutubeMusicProcessor) NeedRemoveDRM() bool {
	// TODO implement me
	panic("implement me")
}

func (ytm *YoutubeMusicProcessor) DRMRemove() error {
	// TODO implement me
	panic("implement me")
}

func (ytm *YoutubeMusicProcessor) TidyMusic() error {
	// TODO implement me
	panic("implement me")
}

func (ytm *YoutubeMusicProcessor) EncryptedExts() []string {
	// TODO implement me
	panic("implement me")
}

func (ytm *YoutubeMusicProcessor) DecryptedExts() []string {
	// TODO implement me
	panic("implement me")
}
