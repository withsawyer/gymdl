package music

import (
	"os/exec"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/processor"
)

/* ---------------------- 结构体与构造方法 ---------------------- */

type YoutubeMusicProcessor struct {
	cfg     *config.Config
	tempDir string
	songs   []*SongInfo
}

func NewYoutubeMusicProcessor(cfg *config.Config, baseTempDir string) (processor.Processor, error) {
	dir, err := processor.BuildOutputDir(baseTempDir)
	if err != nil {
		return nil, err
	}
	return &YoutubeMusicProcessor{cfg: cfg, tempDir: dir}, nil
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
	return sp.songs
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

/* ------------------------ 拓展方法 ------------------------ */
