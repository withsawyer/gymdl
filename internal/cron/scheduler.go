package cron

import (
	"time"
	
	"github.com/go-co-op/gocron/v2"
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/utils"
)

// InitScheduler 初始化定时任务
func InitScheduler(c *config.Config) gocron.Scheduler {
	platformInfo := core.PlatformInfo()
	utils.SugaredLogger().Infof("当前平台: %s", platformInfo.String())
	
	//初始化gocron
	newScheduler, _ := gocron.NewScheduler(gocron.WithLocation(time.Local))
	
	// 注册定时任务
	registerTasks(c, platformInfo, newScheduler)
	return newScheduler
}

//registerTasks 注册定时任务
func registerTasks(c *config.Config, platform core.Platform, scheduler gocron.Scheduler) {
	//    - todo 核心服务(cookiecloud,webdav,ai)健康检查
	//    - todo 依赖的可执行文件或者pip更新检测
	//    - todo 定期从cookiecloud获取cookie数据并解密为cookie文件
	//       *  更新:  cookiecloud->处理成各个音乐平台的cookie数据->储存到本地(覆盖),以平台名称命名
	//       *  使用:  传入对应平台的cookie文件或者读取cookie文件加载cookie
	
	//执行一次依赖安装/更新
	_, _ = scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartImmediately()),
		gocron.NewTask(installExecutableFile, c, platform))
	//todo 注册依赖更新检测任务
	//todo 注册cookiecloud同步任务
}
