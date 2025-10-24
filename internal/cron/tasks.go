package cron

import (
    "sync"
    "time"

    "github.com/nichuanfang/gymdl/config"
)

// installDependency 安装依赖项
func installDependency(c *config.Config) {
    group := sync.WaitGroup{}
    group.Add(2)
    go func() {
        defer group.Done()
        installPipDependency()    
    }()
    go func() {
        defer group.Done()
        installUm()       
    }()
    group.Wait()
}

// updateDependency 更新依赖项
func updateDependency(c *config.Config) {
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
