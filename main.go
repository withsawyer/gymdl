package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	
	"github.com/gin-gonic/gin"
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/cron"
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

func initGin(wg *sync.WaitGroup, c *config.WebConfig) {
	defer wg.Done()
	r := gin.Default()
	
	r.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "Hi,Gin!",
		})
	})
	err := r.Run(":8080")
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
	//初始化定时任务
	go cron.InitCron(c, wg)
	//启动http服务gin
	go initGin(wg, c.WebConfig)
	
	wg.Wait()
}
