package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/internal/cron"
	"github.com/nichuanfang/gymdl/internal/gin/router"
	"github.com/nichuanfang/gymdl/utils"
	"go.uber.org/zap"
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

// printBanner 打印启动横幅
func printBanner() {
	green := "\033[32m"
	reset := "\033[0m"

	banner := `
   ________   _ _       ____   _    
  /  __/\  \/// \__/|  /  _ \ / \   
  | |  _ \  / | |\/||  | | \|| |   
  | |_// / /  | |  ||  | |_/|| |_/\
  \____\/_/   \_/  \|  \____/\____/
==========================================
          🚀 Service Starting...
==========================================
`
	fmt.Println(green + banner + reset)
}

// initWebDAV 初始化webdav服务
func initWebDAV(c *config.WebDAVConfig) {
	utils.ServiceIsOnf("已加载webdav服务")
	return
}

// initCookieCloud 初始化cookiecloud
func initCookieCloud(cookieCloudConfig *config.CookieCloudConfig) {
	core.InitCookieCloud(cookieCloudConfig) // 初始化全局 CookieCloud
	if core.GlobalCookieCloud.CheckConnection() {
		utils.ServiceIsOnf("已加载cookiecloud服务")
	} else {
		utils.Warning("CookieCloud service is not available")
	}
}

// initAI 初始化AI服务
func initAI(c *config.AIConfig) {
	core.InitAI(c)
	if core.GlobalAI.CheckConnection() {
		utils.ServiceIsOnf("已加载AI服务")
	} else {
		utils.Warning("AI service is not available")
	}

}

// =================================后台服务================================================

// initCron 启动定时任务
func initCron(ctx context.Context, wg *sync.WaitGroup, c *config.Config) {
	defer wg.Done()
	s := cron.InitScheduler(c)
	s.Start()
	utils.Success("Scheduler is started")
	<-ctx.Done()
	utils.Stop("定时任务已关闭")
}

// initGin 启动Web服务
func initGin(ctx context.Context, wg *sync.WaitGroup, c *config.Config) {
	defer wg.Done()
	// 设置运行模式 debug/release/test
	gin.SetMode(c.WebConfig.GinMode)
	r := router.SetupRouter(c)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.WebConfig.AppPort),
		Handler: r,
	}

	go func() {
		var httpFlag string
		if c.WebConfig.Https {
			httpFlag = "https"
		} else {
			httpFlag = "http"
		}
		utils.Successf(fmt.Sprintf("Gin server is starting on %s://%s:%d", httpFlag, c.WebConfig.AppDomain, c.WebConfig.AppPort))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			utils.Logger().Error("Gin server error", zap.Any("error", err))
		}
	}()
	<-ctx.Done()
	utils.Stop("Gin服务已关闭")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		utils.Logger().Error("Gin服务关闭错误", zap.Any("error", err))
	}
}

// initGin 启动tg机器人
func initBot(ctx context.Context, wg *sync.WaitGroup, c *config.Config) {
	defer wg.Done()
	// todo 执行启动操作
	utils.Success("Telegram Bot is started")
	<-ctx.Done()
	utils.Stop("Telegram Bot 退出")
}

// =====================================程序入口================================================

func main() {
	if version {
		fmt.Printf("version: %s, build with: %s\n", buildVersion, runtime.Version())
		return
	}
	// banner
	printBanner()

	// 加载配置
	c := config.LoadConfig(configFile)

	// 初始化日志模块
	err := utils.InitLogger(c.Log)
	if err != nil {
		return
	}
	defer utils.Sync()

	// 初始化webdav+连通性检测
	if c.MusicTidy.Mode == 2 {
		initWebDAV(c.WebDAV)
	}
	// 初始化cookiecloud+连通性检测
	initCookieCloud(c.CookieCloud)

	// 初始化AI服务+连通性检测 暂时停用以节省api-key
	//initAI(c.AI)

	// 创建可取消上下文
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(3)

	// 【协程1】 启动定时任务
	go initCron(ctx, wg, c)
	// 【协程2】 启动web服务Gin
	go initGin(ctx, wg, c)
	// 【协程3】 启动telegram机器人
	go initBot(ctx, wg, c)

	// 捕捉系统信号，优雅退出
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// 文丑丑
	<-sig
	utils.Logger().Info("收到退出信号，开始关闭服务...")
	cancel() // 通知所有协程退出

	// 阻塞主协程
	wg.Wait()
	utils.Logger().Info("所有服务已退出，程序结束")
}
