package music

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/processor"
	"github.com/nichuanfang/gymdl/utils"
)

/* ---------------------- 结构体与构造方法 ---------------------- */
type AppleMusicProcessor struct {
	cfg     *config.Config
	tempDir string
	songs   []*SongInfo
}

func NewAppleMusicProcessor(cfg *config.Config, baseTempDir string) processor.Processor {
	return &AppleMusicProcessor{cfg: cfg, tempDir: processor.BuildOutputDir(baseTempDir)}
}

/* ---------------------- 基础接口实现 ---------------------- */

func (am *AppleMusicProcessor) Handle(link string) (string, error) {
	panic("implement me")
}

func (am *AppleMusicProcessor) Category() domain.ProcessorCategory {
	return domain.CategoryMusic
}

func (am *AppleMusicProcessor) Name() domain.LinkType {
	return domain.LinkAppleMusic
}

func (am *AppleMusicProcessor) Songs() []*SongInfo {
	return am.songs
}

/* ------------------------ 下载逻辑 ------------------------ */

func (am *AppleMusicProcessor) DownloadMusic(url string) error {
	if am.cfg.AdditionalConfig.EnableWrapper {
		utils.Logger().Debug("使用增强版am下载器")
		return am.wrapDownload(url)
	} else {
		utils.Logger().Debug("使用默认am下载器")
		return am.defaultDownload(url)
	}
}

func (am *AppleMusicProcessor) DownloadCommand(url string) *exec.Cmd {
	cookiePath := filepath.Join(am.cfg.CookieCloud.CookieFilePath, am.cfg.CookieCloud.CookieFile)
	// https://github.com/glomatico/gamdl/commit/fdab6481ea246c2cf3415565c39da62a3b9dbd52 部分options改动
	rootDir := filepath.Dir(am.tempDir)
	baseDir := filepath.Base(am.tempDir)
	args := []string{
		"--cookies-path", cookiePath,
		"--download-mode", "nm3u8dlre",
		"--output-path", rootDir,
		"--temp-path", rootDir,
		"--album-folder-template", baseDir,
		"--compilation-folder-template", baseDir,
		"--no-album-folder-template", baseDir,
		"--single-disc-file-template", "{title}",
		"--multi-disc-file-template", "{title}",
		"--no-synced-lyrics",
		url,
	}
	return exec.Command("gamdl", args...)
}

func (am *AppleMusicProcessor) BeforeTidy() error {
    
	// TODO implement me
	return nil
}

func (am *AppleMusicProcessor) NeedRemoveDRM() bool {
	// TODO implement me
	return false
}

func (am *AppleMusicProcessor) DRMRemove() error {
	// TODO implement me
	return nil
}

func (am *AppleMusicProcessor) TidyMusic() error {
	// TODO implement me
	return nil
}

func (am *AppleMusicProcessor) EncryptedExts() []string {
	// TODO implement me
	return nil
}

func (am *AppleMusicProcessor) DecryptedExts() []string {
	// TODO implement me
	return nil
}

/* ------------------------ 拓展方法 ------------------------ */

// defaultDownload 默认下载器
func (am *AppleMusicProcessor) defaultDownload(url string) error {
	start := time.Now()
	cmd := am.DownloadCommand(url)
	utils.InfoWithFormat("[AppleMusic] 🎵 开始下载: %s", url)
	utils.DebugWithFormat("[AppleMusic] 执行命令: %s", strings.Join(cmd.Args, " "))
	err := processor.CreateOutputDir(am.tempDir)
	if err != nil {
		return err
	}
	output, err := cmd.CombinedOutput()
	logOut := strings.TrimSpace(string(output))
	if err != nil {
		utils.ErrorWithFormat("[AppleMusic] ❌ gamdl 下载失败: %v\n输出:\n%s", err, logOut)
		return fmt.Errorf("gamdl 下载失败: %w", err)
	}

	if logOut != "" {
		utils.DebugWithFormat("[AppleMusic] 下载输出:\n%s", logOut)
	}
	utils.InfoWithFormat("[AppleMusic] ✅ 下载完成（耗时 %v）", time.Since(start).Truncate(time.Millisecond))

	return nil
}

// wrapDownload todo 增强版下载器
func (am *AppleMusicProcessor) wrapDownload(string) error {
	panic("implement me")
}
