package cron

import (
	"fmt"
	"strings"
	"sync"
	"time"
	
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/utils"
	"go.uber.org/zap"
)

func InitScheduler(c *config.Config) gocron.Scheduler {
	platformInfo := core.PlatformInfo()
	utils.SugaredLogger().Infof("当前平台: %s", platformInfo.String())
	
	var startTimes sync.Map
	
	newScheduler, err := gocron.NewScheduler(
		gocron.WithLocation(time.Local),
		gocron.WithGlobalJobOptions(
			gocron.WithEventListeners(
				gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
					startTimes.Store(jobID, time.Now())
					logJobEvent("started", jobName, jobID, nil, nil)
				}),
				gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
					if v, ok := startTimes.LoadAndDelete(jobID); ok {
						start := v.(time.Time)
						logJobEvent("finished", jobName, jobID, &start, nil)
					} else {
						logJobEvent("finished", jobName, jobID, nil, nil)
					}
				}),
				gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
					if v, ok := startTimes.LoadAndDelete(jobID); ok {
						start := v.(time.Time)
						logJobEvent("failed", jobName, jobID, &start, err)
					} else {
						logJobEvent("failed", jobName, jobID, nil, err)
					}
				}),
			),
		),
	)
	if err != nil {
		utils.Logger().Fatal("Failed to initialize scheduler", zap.Error(err))
	}
	
	registerTasks(c, platformInfo, newScheduler)
	return newScheduler
}

func logJobEvent(event, jobName string, jobID uuid.UUID, start *time.Time, err error) {
	name := jobName
	if idx := strings.LastIndex(jobName, "."); idx != -1 {
		name = jobName[idx+1:]
	}
	
	fields := []zap.Field{
		zap.String("job_name", name),
		zap.String("job_id", jobID.String()),
	}
	
	if start != nil {
		fields = append(fields, zap.String("duration", fmt.Sprintf("%vms", time.Since(*start).Milliseconds())))
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
		utils.Logger().Error(fmt.Sprintf("[Job %s]", event), fields...)
	} else {
		utils.Logger().Info(fmt.Sprintf("[Job %s]", event), fields...)
	}
}

//registerTasks 注册定时任务
func registerTasks(c *config.Config, platform core.Platform, scheduler gocron.Scheduler) {
	//    - todo 核心服务(cookiecloud,webdav,ai)健康检查
	//    - todo 依赖的可执行文件或者pip更新检测
	//    - todo 定期从cookiecloud获取cookie数据并解密为cookie文件
	//       *  更新:  cookiecloud->处理成各个音乐平台的cookie数据->储存到本地(覆盖),以平台名称命名
	//       *  使用:  传入对应平台的cookie文件或者读取cookie文件加载cookie
	
	//执行一次依赖安装/更新
	scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartImmediately()),
		gocron.NewTask(installExecutableFile, c, platform))
	//todo 注册健康检查任务
	scheduler.NewJob(gocron.DurationJob(time.Minute), gocron.NewTask(healthCheck, c))
	//todo 注册依赖更新检测任务
	//todo 注册cookiecloud同步任务
}
