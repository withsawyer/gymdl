package cron

import (
	"os/exec"
)

// installPipDependency 安装pip依赖(执行一次)
func installPipDependency() {
	cmd := exec.Command("pip", "--disable-pip-version-check", "install", "-U", "-r", "requirements.txt")
	output, err := cmd.CombinedOutput()
	logger.Debug("\n" + string(output))
	if err != nil {
		logger.Error("⚠️Pip Dependency Install failed: " + err.Error())
		return
	}
	logger.Info("💡Pip Dependency Installed successfully")
}

// updatePipDependency 定时更新pip依赖
func updatePipDependency() {
	updateCmd := exec.Command("pip", "--disable-pip-version-check", "install", "-U", "-r", "requirements.txt")
	out, err := updateCmd.CombinedOutput()
	logger.Debug("\n" + string(out))
	if err != nil {
		logger.Error("⚠️Pip Dependency Update failed: " + err.Error())
		return
	}

	logger.Info("💡Pip Dependency has updated successfully")
}
