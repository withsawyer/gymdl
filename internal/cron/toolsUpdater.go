package cron

import "time"

// installFFmpeg 安装ffmpeg
func installFFmpeg() {
	//1. 只需执行一次
	//2. 如果指定文件夹没有 执行安装
	time.Sleep(20 * time.Second)
	logger.Info("installing ffmpeg")
}

// installUm 安装um解锁器
func installUm() {
	//1. 	只需执行一次
	//2. 如果指定文件夹没有 执行安装; 如果有则判断是否需要更新
	time.Sleep(1 * time.Second)
	logger.Info("installing Um")
}

// updateUm 检查um版本更新
func updateUm() {
	//定时执行
}
