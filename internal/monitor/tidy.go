package monitor

import (
	"os"
	"path/filepath"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/utils"
)

// 整理到本地
func tidyToLocal(src string, cfg *config.Config) error {
	dst := filepath.Join(cfg.Tidy.DistDir, utils.SanitizeFileName(filepath.Base(src)))
	if err := os.Rename(src, dst); err != nil {
		utils.WarnWithFormat("[Um] ⚠️ 移动失败 %s → %s: %v", src, dst, err)
		return err
	}
	utils.InfoWithFormat("[Um] 📦 已整理: %s", dst)
	return nil
}

// 整理到webdav
func tidyToWebDAV(path string, webdav *core.WebDAV, cfg *config.Config) error {
	if err := webdav.Upload(path); err != nil {
		utils.WarnWithFormat("[Um] ☁️ 上传失败 %s: %v", utils.SanitizeFileName(filepath.Base(path)), err)
		return err
	}
	utils.InfoWithFormat("[Um] ☁️ 已上传: %s", filepath.Base(path))
	return nil
}
