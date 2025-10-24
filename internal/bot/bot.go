package bot

import (
	"time"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/utils"
	"go.uber.org/zap"
)

// tg机器人入口

var logger *zap.Logger

func InitBot(c *config.Config) {
	logger = utils.Logger()
	utils.Success("TelegramBot已成功启动")
	// 模拟阻塞
	time.Sleep(time.Hour)
}
