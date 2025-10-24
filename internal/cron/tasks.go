package cron

import (
    "sync"

    "github.com/nichuanfang/gymdl/config"
)

// installDependency 安装依赖项
func installDependency(c *config.Config) {
	group := sync.WaitGroup{}
	group.Add(3)
	go func() {
		defer group.Done()
		installPipDependency()
	}()
	go func() {
		defer group.Done()
		installUm()
	}()
    go func() {
        defer group.Done()
        syncCookieCloud(c.CookieCloud)
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