package cron

import (
	"time"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/utils"
)

// InitCron 初始化定时任务
func InitCron(config *config.Config) {
	//1.安装,检查更新可执行程序(yt-dlp gamdl um ffmpeg),
	//2.启动更新定时任务
	utils.Logger.Info("定时任务模块已加载")
	time.Sleep(time.Hour)
}
