package cron

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/nichuanfang/gymdl/config"
)

// installFFmpeg 安装ffmpeg
func installFFmpeg(c *config.Config) {
	switch runtime.GOOS {
	case "darwin":
        logger.Info("【ffmpeg】Macos skip installation,you need install it manually")
        output, _ := exec.Command("which", "ffmpeg").CombinedOutput()
        logger.Info(string(output))
		return
	case "windows":
        logger.Info("【ffmpeg】Windows skip installation,you need install it manually")
        output, _ := exec.Command("where", "ffmpeg").CombinedOutput()
        logger.Info(string(output))
		return
	}
	var ffmpegPath string
	ffmpegPath = filepath.Join(c.Ffmpeg.FfmpegPath, "ffmpeg")
	// 检查 ffmpeg 是否存在
	if _, err := os.Stat(ffmpegPath); err == nil {
		logger.Info("【ffmpeg】Already installed, skip installation.")
		return
	}

	logger.Info("【ffmpeg】Not found, start installing...")

	cmd := exec.Command("/bin/sh", "/app/install_ffmpeg.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.Error("【ffmpeg】Install failed: " + err.Error())
		return
	}

	logger.Info("【ffmpeg】Installed successfully.")
}

// installUm 安装um解锁器
func installUm(*config.Config) {
	// 1. 	只需执行一次
	// 2. 如果指定文件夹没有 执行安装; 如果有则判断是否需要更新
	time.Sleep(1 * time.Second)
	logger.Info("【Um】 has Installed.")
}

// updateUm 检查um版本更新
func updateUm() {
	// 定时执行
}
