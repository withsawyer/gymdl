package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/internal/bot"
	"github.com/nichuanfang/gymdl/internal/cron"
	"github.com/nichuanfang/gymdl/internal/gin/router"
	"github.com/nichuanfang/gymdl/utils"
)

var (
	configFile   string
	version      bool
	buildVersion = "dev-main"
)

func init() {
	flag.StringVar(&configFile, "c", "./config.json", "config file")
	flag.BoolVar(&version, "v", false, "display version")
	flag.Parse()
}

//=================================基础配置================================================

//=================================核心服务================================================

// initCron 启动定时任务
func initCron(wg *sync.WaitGroup, c *config.Config) {
	defer wg.Done()
	cron.InitCron(c)
}

// initGin 启动Web服务
func initGin(wg *sync.WaitGroup, c *config.Config) {
	defer wg.Done()
	//设置运行模式 debug/release/test
	gin.SetMode(c.WebConfig.GinMode)
	r := router.SetupRouter(c)
	err := r.Run(fmt.Sprintf("%s:%d", c.WebConfig.AppHost, c.WebConfig.AppPort))
	if err != nil {
		return
	}
}

// initGin 启动tg机器人
func initBot(wg *sync.WaitGroup, c *config.Config) {
	defer wg.Done()
	bot.InitBot(c)
}

//=====================================程序入口================================================

func main() {
	if version {
		fmt.Printf("version: %s, build with: %s\n", buildVersion, runtime.Version())
		return
	}
	//加载配置
	c := config.LoadConfig(configFile)
	//初始化日志模块
	err := utils.InitLogger(c.Log)
	if err != nil {
		return
	}
	defer utils.Sync()

	wg := &sync.WaitGroup{}
	wg.Add(3)

	//【协程1】 启动定时任务
	go initCron(wg, c)
	//【协程2】 启动web服务Gin
	go initGin(wg, c)
	//【协程3】 启动telegram机器人
	go initBot(wg, c)
	//阻塞主协程
	wg.Wait()
}
