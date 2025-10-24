package cron

import (
    "os/exec"
)

// installPipDependency 安装pip依赖(执行一次)
func installPipDependency() {
	// 通过pip安装/更新
	command := exec.Command("pip", "install", "-U", "-r", "requirements.txt")
	if output,err := command.CombinedOutput(); err != nil {
        logger.Info("\n"+string(output))
		logger.Error("【Pip Dependency】Installed failed")
	} else {
        logger.Info("\n"+string(output))
		logger.Info("【Pip Dependency】has Installed.")
	}
}

// updateByRequirements 检查pip依赖更新(定时执行)
func updatePipDependency() {
	// 通过pip更新
}
