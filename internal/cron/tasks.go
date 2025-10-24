package cron

import (
    "sync"
    "time"

    "github.com/nichuanfang/gymdl/config"
)

// installDependency 安装依赖项
func installDependency(c *config.Config) {
	logger.Info("开始初始化依赖项...")
	group := &sync.WaitGroup{}
	group.Add(3)
	go func() {
		defer group.Done()
		installFFmpeg(c)
	}()
	go func() {
		defer group.Done()
		installUm(c)
	}()
	go func() {
		defer group.Done()
		installPipDependency()
	}()
	group.Wait()
	logger.Info("依赖项更新完毕!")
}

// updateDependency 更新依赖项
func updateDependency(c *config.Config) {
    logger.Info("开始检测依赖项更新...")
    group := sync.WaitGroup{}
    group.Add(2)
    go func() {
        defer group.Done()
        updatePipDependency()
    }()
    go func() {
        defer group.Done()
        updateUm()
    }()
    group.Wait()
    logger.Info("依赖项更新检测完毕!")
}

// healthCheck 健康检查
func healthCheck(c *config.Config) {
	// todo 核心服务(cookiecloud,webdav,ai)健康检查
	return
}

// syncCookieCloud 同步cookie
func syncCookieCloud(c *config.Config) {
	// todo 定期从cookiecloud获取cookie数据并解密为cookie文件
	//   *  更新:  cookiecloud->处理成各个音乐平台的cookie数据->储存到本地(覆盖),以平台名称命名
	//   *  使用:  传入对应平台的cookie文件或者读取cookie文件加载cookie
	time.Sleep(time.Second * time.Duration(2))
	logger.Info("cookie更新成功")
}
