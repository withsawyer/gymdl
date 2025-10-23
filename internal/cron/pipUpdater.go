package cron

import "time"

// installPipDependency 安装pip依赖(执行一次)
func installPipDependency() {
	//通过pip安装/更新
	time.Sleep(10 * time.Second)
	logger.Info("installing Pip")
}

// updateByRequirements 检查pip依赖更新(定时执行)
func updatePipDependency() {
	//通过pip更新
}
