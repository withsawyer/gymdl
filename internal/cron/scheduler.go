package cron

import (
	"time"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/utils"
)

// InitCron 初始化定时任务
func InitCron(config *config.Config) {
	// 1.如果依赖的可执行文件(如yt-dlp,gamdl,um,ffmpeg)等未安装,执行安装
	// 2. 如果已安装,检查更新(只更新非固定版本的或者经常更新的)
	// 3. 开启定时任务
	//    - 核心服务(cookiecloud,webdav,ai)健康检查
	//    - 依赖的可执行文件或者pip更新检测
	//    - 定期从cookiecloud获取cookie数据并解密为cookie文件
	platformInfo := core.PlatformInfo()
	utils.Logger().Infof("当前平台:%s", platformInfo.String())
	time.Sleep(time.Hour)
}
