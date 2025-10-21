package bot

import (
	"time"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/utils"
)

// tg机器人入口

func InitBot(c *config.Config) {
	// 模拟阻塞
	utils.Logger.Info("tg机器人模块已加载")
	time.Sleep(time.Hour)
}
