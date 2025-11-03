package main

import (
	"fmt"
	"github.com/nichuanfang/gymdl/utils"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/nichuanfang/gymdl/config"
)

// 简单初始化日志系统，避免空指针异常
func initLogger() {
	// 确保日志目录存在
	logDir := filepath.Join("./data", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("创建日志目录失败: %v\n", err)
	}

	// 创建临时日志配置
	logConfig := &config.LogConfig{
		Level: 2, // Info级别
		Mode:  1, // 只输出到控制台
		File:  filepath.Join(logDir, "test.log"),
	}

	// 初始化日志
	if err := utils.InitLogger(logConfig); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		// 如果初始化失败，继续执行但可能会有日志相关的错误
	}
}

func main() {
	// 初始化日志
	initLogger()

	url := "https://aweme.snssdk.com/aweme/v1/play/?line=0&logo_name=aweme_diversion_search&ratio=720p&video_id=v0d00fg10000d27evvfog65kbjdbm9hg"
	options := utils.NewDefaultDownloadOptions()
	options.SavePath = "./data/temp"
	options.FileName = "test_video.mp4"
	options.IgnoreSSL = true // 忽略SSL验证，避免证书问题

	// 设置进度回调
	options.ProgressFunc = func(progress *utils.DownloadProgress) {
		fmt.Printf("进度: %.1f%%, 速度: %s, 已下载: %s/%s\r",
			progress.Progress,
			progress.FormattedSpeed,
			progress.FormattedDownloaded,
			progress.FormattedSize,
		)
	}

	fmt.Println("开始下载测试...")
	fmt.Printf("URL: %s\n", url)

	manager, err := utils.DownloadFile(url, options)
	if err != nil {
		log.Fatalf("创建下载管理器失败: %v", err)
	}

	// 开始下载
	err = manager.Start()
	if err != nil {
		log.Fatalf("开始下载失败: %v", err)
	}

	// 等待下载完成或失败
	for {
		progress := manager.GetProgress()
		if progress.Status == utils.StatusCompleted {
			fmt.Println("\n下载完成!")
			break
		} else if progress.Status == utils.StatusFailed {
			fmt.Printf("\n下载失败: %s\n", progress.ErrorMessage)
			break
		}

		// 每秒检查一次状态
		// 注意：实际应用中应该使用更优雅的方式等待完成
		// 这里为了简化测试使用轮询
		time.Sleep(1 * time.Second)
	}
}
