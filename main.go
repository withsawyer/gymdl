package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/cron"
	"github.com/nichuanfang/gymdl/internal/router"
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

// initGin 启动Web服务
func initGin(wg *sync.WaitGroup, c *config.Config) {
	defer wg.Done()
	r := router.SetupRouter(c)
	err := r.Run(fmt.Sprintf("%s:%d", c.WebConfig.AppHost, c.WebConfig.AppPort))
	if err != nil {
		return
	}
}

func main() {
	if version {
		fmt.Printf("version: %s, build with: %s\n", buildVersion, runtime.Version())
		return
	}
	//加载配置
	c := config.LoadConfig(configFile)
	
	wg := &sync.WaitGroup{}
	wg.Add(2)
	
	//【协程1】 初始化定时任务
	go cron.InitCron(c, wg)
	//【协程2】 启动http服务gin
	go initGin(wg, c)
	
	//阻塞主协程
	wg.Wait()
}
