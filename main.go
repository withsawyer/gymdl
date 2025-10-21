package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"

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

// initCron 启动定时任务
func initCron(wg *sync.WaitGroup, c *config.Config) {
	defer wg.Done()
	cron.InitCron(c)
}

// initGin 启动Web服务
func initGin(wg *sync.WaitGroup, c *config.Config) {
	defer wg.Done()
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

func main() {
	if version {
		fmt.Printf("version: %s, build with: %s\n", buildVersion, runtime.Version())
		return
	}
	//加载配置
	c := config.LoadConfig(configFile)
	//初始化日志模块
	utils.InitLogger(c.Log.Mode, c.Log.Level, c.Log.File)

	wg := &sync.WaitGroup{}
	wg.Add(3)

	//【协程1】 初始化定时任务
	go initCron(wg, c)
	//【协程2】 启动http服务gin
	go initGin(wg, c)
	//【协程3】 启动telegram机器人
	go initBot(wg, c)
	//阻塞主协程
	wg.Wait()
}
