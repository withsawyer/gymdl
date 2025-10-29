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
	"github.com/nichuanfang/gymdl/internal/bot"
	"github.com/nichuanfang/gymdl/internal/cron"
	"github.com/nichuanfang/gymdl/internal/gin/router"
	"github.com/nichuanfang/gymdl/internal/monitor"
	"github.com/nichuanfang/gymdl/utils"
	"go.uber.org/zap"
)

var (
	configFile   string
	versionFlag  bool
	buildVersion = "dev-main"
)

func init() {
	flag.StringVar(&configFile, "c", "./config.json", "指定配置文件路径")
	flag.BoolVar(&versionFlag, "v", false, "显示版本信息")
	flag.Parse()
}

func printBanner() {
	green := "\033[32m"
	reset := "\033[0m"

	banner := `
   ________   _     _    ____  _    
  /  __/\  \/// \__/ |  /  _ \/ \   
  | |  _ \  / | |\/| |  | | \|| |   
  | |_// / /  | |  | |  | |_/|| |_/\
  \____\/_/   \_/  \_|  \____/\____/

=======================================
           🚀 服务启动中...
=======================================
`
	fmt.Println(green + banner + reset)
}

// 初始化 WebDAV 服务
func initWebDAV(c *config.WebDAVConfig) {
	core.InitWebDAV(c)
	if core.GlobalWebDAV.CheckConnection() {
		utils.ServiceIsOn("WebDAV 服务已加载")
	} else {
		utils.Warning("WebDAV 服务不可用，请检查配置或网络连接")
	}
}

// 初始化 Wrapper 服务
func initWrapper() {
	// 检测wrapper服务是否可用
	if utils.CheckWrapperConnection("wrapper") {
		utils.ServiceIsOn("Wrapper 服务已加载")
	} else {
		utils.Warning("Wrapper 服务不可用，请检查配置或网络连接")
	}
}

// 初始化 CookieCloud 服务
func initCookieCloud(cfg *config.CookieCloudConfig) {
	core.InitCookieCloud(cfg)
	if core.GlobalCookieCloud.CheckConnection() {
		utils.ServiceIsOn("CookieCloud 服务已加载")
	} else {
		utils.Warning("CookieCloud 服务不可用，请检查配置或网络连接")
	}
}

// 初始化 AI 服务
func initAI(c *config.AIConfig) {
	core.InitAI(c)
	if core.GlobalAI.CheckConnection() {
		utils.ServiceIsOn("AI 服务已加载")
	} else {
		utils.Warning("AI 服务不可用，请检查配置或网络连接")
	}
}

// 启动定时任务
func initCron(ctx context.Context, c *config.Config) {
	s := cron.InitScheduler(c)
	s.Start()
	utils.Success("定时任务调度器已启动")
	<-ctx.Done()
	_ = s.Shutdown()
	utils.Stop("定时任务调度器已关闭")
}

// 启动目录监控
func initMonitor(ctx context.Context, c *config.Config) {
	wm := monitor.NewWatchManager()

	for _, dir := range c.AdditionalConfig.MonitorDirs {
		// 监控主目录
		if err := wm.AddDir(dir); err != nil {
			utils.WarnWithFormat("[Monitor] 注册目录失败: %s (%v)", dir, err)
			continue
		}
		utils.InfoWithFormat("[Monitor] 注册目录: %s", dir)
	}

	go func() {
		utils.Success("目录监控已启动")
		wm.StartWorkerPool(runtime.NumCPU())
	}()

	<-ctx.Done()
	wm.Stop()
	utils.Stop("目录监控已关闭")
}

// 启动 Gin Web 服务
func initGin(ctx context.Context, c *config.Config) {
	gin.SetMode(c.WebConfig.GinMode)
	r := router.SetupRouter(c)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.WebConfig.AppPort),
		Handler: r,
	}

	go func() {
		protocol := "http"
		if c.WebConfig.Https {
			protocol = "https"
		}
		utils.Successf("Gin Web 服务已启动：%s://%s:%d", protocol, c.WebConfig.AppDomain, c.WebConfig.AppPort)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			utils.Logger().Error("Gin Web 服务运行出错", zap.Error(err))
		}
	}()

	<-ctx.Done()
	utils.Stop("Gin Web 服务正在关闭...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		utils.Logger().Error("Gin Web 服务关闭时发生错误", zap.Error(err))
	}
}

// 启动 Telegram Bot
func initBot(ctx context.Context, c *config.Config) {
	app, err := bot.NewBotApp(c)
	if err != nil {
		utils.Logger().Error("创建 Telegram Bot 应用失败", zap.Error(err))
		return
	}

	modeText := "轮询模式"
	if c.Telegram.Mode == 2 {
		modeText = fmt.Sprintf("Webhook 模式（端口: %d, URL: %s）", c.Telegram.WebhookPort, c.Telegram.WebhookURL)
	}

	go func() {
		utils.Success(fmt.Sprintf("Telegram Bot 已启动，运行模式：【%s】", modeText))
		app.Start()
	}()

	<-ctx.Done()
	app.Stop()
	utils.Stop("Telegram Bot 已停止运行")
}

// ===================== 程序入口 =====================
func main() {
	if versionFlag {
		fmt.Printf("版本号: %s，构建环境: %s\n", buildVersion, runtime.Version())
		return
	}

	printBanner()

	// 加载配置文件
	c := config.LoadConfig(configFile)

	// 初始化日志模块
	if err := utils.InitLogger(c.Log); err != nil {
		fmt.Println("日志模块初始化失败：", err)
		return
	}
	defer utils.Sync()

	// 初始化各服务
	if c.AdditionalConfig.EnableCron {
		initCookieCloud(c.CookieCloud)
	}

	if c.AI.Enable {
		initAI(c.AI)
	}

	if c.MusicTidy.Mode == 2 {
		initWebDAV(c.WebDAV)
	}

	if c.AdditionalConfig.EnableWrapper {
		initWrapper()
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// 用 map 映射模块启动逻辑，更优雅地管理协程
	services := map[string]func(context.Context, *config.Config){}

	// 是否启用定时任务
	if c.AdditionalConfig.EnableCron {
		services["cron"] = initCron
	}

	// 是否启用目录监控
	if c.AdditionalConfig.EnableDirMonitor {
		services["monitor"] = initMonitor
	}

	// 是否启用web服务
	if c.WebConfig.Enable {
		services["web"] = initGin
	}

	// 是否启用telegram
	if c.Telegram.Enable {
		services["telegram"] = initBot
	}

	// 启动所有服务
	for name, start := range services {
		wg.Add(1)
		go func(name string, start func(context.Context, *config.Config)) {
			defer wg.Done()
			start(ctx, c)
		}(name, start)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	utils.Logger().Info(fmt.Sprintf("收到退出信号 [%s]，正在关闭服务...", sig))

	cancel() // 通知协程退出
	wg.Wait()

	utils.Logger().Info("所有服务已正常退出，程序结束")
}
