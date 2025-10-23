package cron

import (
	"time"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/utils"
)

// 通过pip安装的 或者 可执行文件都在这个文件夹 并赋予+x权限
var binPath = "/usr/local/bin"

// installDependency 安装依赖项
func installDependency(*config.Config, core.Platform) {
	// todo 如果依赖的可执行文件(如yt-dlp,gamdl,um,ffmpeg)等未安装,执行安装

	// 依赖项名称:gamdl
	// 依赖项安装方式:pip
	// 安装指令:pip install -U gamdl
	logger.Info("开始执行依赖项更新...")
	time.Sleep(2 * time.Second)
	logger.Info("依赖项更新完毕")
	return
}

// updateDependency 更新依赖项
func updateDependency(c *config.Config, platform core.Platform) {
	// todo 依赖的可执行文件或者pip更新检测
	time.Sleep(time.Second * time.Duration(2))
	logger.Info("依赖更新成功")
}

// healthCheck 健康检查
func healthCheck(c *config.Config) {
	// todo 核心服务(cookiecloud,webdav,ai)健康检查
	utils.NetworkHealth("健康检查成功")
	return
}

// syncCookieCloud 同步cookie
func syncCookieCloud(c *config.Config) {
	// todo 定期从cookiecloud获取cookie数据并解密为cookie文件
	//   *  更新:  cookiecloud->处理成各个音乐平台的cookie数据->储存到本地(覆盖),以平台名称命名
	//   *  使用:  传入对应平台的cookie文件或者读取cookie文件加载cookie
	time.Sleep(time.Second * time.Duration(2))
	logger.Info("cookie更新成功")
}
