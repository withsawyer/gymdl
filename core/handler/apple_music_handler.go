package handler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/constants"
)

type AppleMusicHandler struct{}

// 苹果音乐处理器
func (mh *AppleMusicHandler) HandlerMusic(url string, cfg *config.Config) {
	// 输出目录
	var outputDir = constants.AppleMusicTempDir
	// 构建 gamdl 命令
	cookieFileFullPath := filepath.Join(cfg.CookieCloud.CookieFilePath, cfg.CookieCloud.CookieFile)
	cmd := exec.Command(
		"gamdl",
		"--cookies-path", cookieFileFullPath,
		"--download-mode", "nm3u8dlre",
		"--output-path", outputDir,
		url,
	)

	// 运行命令并捕获输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("下载失败: %v\n", err)
		return
	}

	// 查找下载后的文件（假设只有一个音乐文件）
	files, err := os.ReadDir(outputDir)
	if err != nil {
		fmt.Printf("读取输出目录失败: %v\n", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			fmt.Printf("音乐文件已下载: %s\n", filepath.Join(outputDir, file.Name()))
			break
		}
	}
}
