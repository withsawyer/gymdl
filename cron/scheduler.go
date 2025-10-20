package cron

import (
	"sync"
	"time"
	
	"github.com/nichuanfang/gymdl/config"
)

// InitCron 初始化定时任务
func InitCron(config *config.Config, wg *sync.WaitGroup) {
	//1.安装,检查更新可执行程序(yt-dlp gamdl um ffmpeg),
	//2.启动更新定时任务
	defer wg.Done()
	//暂时模拟业务
	time.Sleep(time.Hour)
}
