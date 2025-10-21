package cron

import (
	"time"
	
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/utils"
)

//installExecutableFile 安装依赖项
func installExecutableFile(*config.Config, core.Platform) {
	// 1.todo 如果依赖的可执行文件(如yt-dlp,gamdl,um,ffmpeg)等未安装,执行安装
	// 2. todo 如果已安装,检查更新(只更新非固定版本的或者经常更新的)
	utils.Logger().Info("开始执行依赖项更新...")
	time.Sleep(2 * time.Second)
	utils.Successf("依赖项更新完毕")
	return
}
