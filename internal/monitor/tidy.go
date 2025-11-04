package monitor

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/utils"
)

// tidyToLocal 将文件整理到目标目录，支持同盘移动和跨盘复制
func tidyToLocal(src string, cfg *config.Config) error {
	dst := filepath.Join(cfg.Tidy.DistDir, utils.SanitizeFileName(filepath.Base(src)))

	if sameDrive(src, dst) {
		// 同盘直接重命名
		if err := os.Rename(src, dst); err != nil {
			utils.WarnWithFormat("[Um] ⚠️ 移动失败 %s → %s: %v", src, dst, err)
			return err
		}
	} else {
		// 跨盘复制 + 删除
		if err := copyFile(src, dst); err != nil {
			utils.WarnWithFormat("[Um] ⚠️ 复制文件失败: %v", err)
			return err
		}

		// 删除源文件
		if err := os.Remove(src); err != nil {
			utils.WarnWithFormat("[Um] ⚠️ 删除源文件失败: %v", err)
			return err
		}
	}

	utils.InfoWithFormat("[Um] 📦 已整理: %s", dst)
	return nil
}

// copyFile 使用缓冲区复制文件，确保文件句柄关闭
func copyFile(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close() // 手动关闭源文件

	output, err := os.Create(dst)
	if err != nil {
		return err
	}

	// 使用 1MB 缓冲区
	buf := make([]byte, 1024*1024)
	_, err = io.CopyBuffer(output, input, buf)

	// 立即关闭目标文件
	if closeErr := output.Close(); closeErr != nil {
		utils.WarnWithFormat("[Um] ⚠️ 关闭目标文件失败: %v", closeErr)
	}

	return err
}

// sameDrive 判断两个路径是否在同一个磁盘或 UNC 网络共享
func sameDrive(path1, path2 string) bool {
	abs1, err1 := filepath.Abs(path1)
	abs2, err2 := filepath.Abs(path2)
	if err1 != nil || err2 != nil {
		return false
	}

	// 本地盘符比较 (C:, D:...)
	if len(abs1) >= 2 && len(abs2) >= 2 && abs1[1] == ':' && abs2[1] == ':' {
		return strings.EqualFold(abs1[:2], abs2[:2])
	}

	// UNC 网络路径比较 (\\NAS\share)
	if strings.HasPrefix(abs1, `\\`) && strings.HasPrefix(abs2, `\\`) {
		parts1 := strings.SplitN(abs1, `\`, 4)
		parts2 := strings.SplitN(abs2, `\`, 4)
		if len(parts1) >= 3 && len(parts2) >= 3 {
			return strings.EqualFold(parts1[1], parts2[1]) && strings.EqualFold(parts1[2], parts2[2])
		}
	}

	return false
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
