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

// =================================基础服务================================================

// verifyConfig 核心配置项检测
func verifyConfig(c *config.Config) error {
	fmt.Println(c.CookieCloud.CookieCloudUrl)
	return nil
}

// printBanner 打印启动横幅
func printBanner() {
	banner := `
 ________  _ _      ____  _    
/  __/\  \/// \__/|/  _ \/ \   
| |  _ \  / | |\/||| | \|| |   
| |_// / /  | |  ||| |_/|| |_/\
\____\/_/   \_/  \|\____/\____/
`
	fmt.Println(banner)
}

// initWebDAV 初始化webdav服务
func initWebDAV(c *config.WebDAVConfig) {
	utils.Logger().Info("已加载webdav服务:", c.WebDAVUrl)
	return
}

// initCookieCloud 初始化cookiecloud
func initCookieCloud(cookieCloudConfig *config.CookieCloudConfig) {
	utils.Successf("已加载cookiecloud服务: %s ", cookieCloudConfig.CookieCloudUrl)
}

// initAI 初始化AI服务
func initAI(c *config.AIConfig) {
	utils.Successf("已加载AI服务: %s ", c.BaseUrl)
}

//=================================后台服务================================================

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
	utils.Successf("Gin已成功启动")
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
	// banner
	printBanner()

	//加载配置
	c := config.LoadConfig(configFile)

	//核心配置项检测
	verify_err := verifyConfig(c)
	if verify_err != nil {
		return
	}

	//初始化日志模块
	err := utils.InitLogger(c.Log)
	if err != nil {
		return
	}
	defer utils.Sync()

	// 初始化webdav+连通性检测
	if c.MusicTidy.Mode == 2 {
		initWebDAV(c.WebDAV)
	}
	//初始化cookiecloud+连通性检测
	initCookieCloud(c.CookieCloud)

	//初始化AI服务+连通性检测
	initAI(c.AI)

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
