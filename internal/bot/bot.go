package bot

import (
	"time"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/utils"
)

// tg机器人入口

func InitBot(c *config.Config) {
	utils.Success("TelegramBot已成功启动")
	// 模拟阻塞
	time.Sleep(time.Hour)
}
