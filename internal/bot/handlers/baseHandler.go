package handlers

import (
	"go.uber.org/zap"
	tb "gopkg.in/telebot.v4"
)

// InitCommands 初始化命令列表
func InitCommands(c tb.Context) error {
	// 定义命令列表
	commands := []tb.Command{
		{Text: "start", Description: "启动 bot"},
		{Text: "help", Description: "帮助说明"},
	}

	// 注册命令到 Telegram（全局默认 scope）
	err := c.Bot().SetCommands(commands)
	if err != nil {
		logger.Error("Failed to set commands", zap.Error(err))
		return err
	}

	return c.Send("指令初始化成功 ✅")
}

// 普通文本处理器
func HandleText(c tb.Context) error {
	return nil
}
