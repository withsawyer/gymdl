package cron

import "os/exec"

// installPipDependency 安装pip依赖(执行一次)
func installPipDependency() {
	// 通过pip安装/更新
	exec.Command("pip", "install", "-U", "-r", "requirements.txt")
	logger.Info("【Pip Dependency】 has Installed.")
}

// updateByRequirements 检查pip依赖更新(定时执行)
func updatePipDependency() {
	// 通过pip更新
}
