package bot

import (
	"fmt"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v4"
)

// StartCommand 响应 /start 命令
func StartCommand(c tb.Context) error {
	msg := `👋 欢迎来到 GymDL Bot!
我可以帮你管理健身相关的任务和信息 🏋️‍♂️

使用 /help 查看可用命令 📜`

	if err := c.Send(msg, &tb.SendOptions{ParseMode: tb.ModeMarkdown}); err != nil {
		logger.Error("Failed to send start message", zap.Error(err))
		return err
	}
	return nil
}

// HelpCommand 响应 /help 命令，自动生成已注册命令说明
func HelpCommand(c tb.Context) error {
	commands, err := c.Bot().Commands()
	if err != nil {
		logger.Error("Failed to get bot commands", zap.Error(err))
		return c.Send("⚠️ 无法获取命令列表")
	}

	if len(commands) == 0 {
		return c.Send("😅 当前没有可用命令")
	}

	helpMsg := "📖 *可用命令列表:*\n\n"
	for _, cmd := range commands {
		helpMsg += fmt.Sprintf("• `/%s` - %s\n", cmd.Text, cmd.Description)
	}

	helpMsg += "\n✨ *Tip:* 你可以直接在聊天中输入命令来体验功能!"

	if err := c.Send(helpMsg, &tb.SendOptions{ParseMode: tb.ModeMarkdown}); err != nil {
		logger.Error("Failed to send help message", zap.Error(err))
		return err
	}

	return nil
}

// WrapperCommand 注册wrapper并启动
func WrapperCommand(c tb.Context) error {
	//执行命令
	if !app.cfg.WrapperConfig.Enable {
		return c.Send("请先开启wrapper再使用该指令!")
	}

	return c.Send("准备颁发2FA,请稍后...")
}

// SetCommands 初始化 Telegram 命令列表
func SetCommands(c tb.Context) error {
	commands := []tb.Command{
		{Text: "start", Description: "启动 Bot 👋"},
		{Text: "help", Description: "获取帮助 📜"},
		{Text: "wrapper", Description: "am注册 🧩"},
	}

	if err := c.Bot().SetCommands(commands); err != nil {
		logger.Error("Failed to set commands", zap.Error(err))
		return err
	}

	const successMsg = "✅ 指令初始化成功，使用 /help 查看可用命令"
	if err := c.Send(successMsg); err != nil {
		logger.Error("Failed to send confirmation message", zap.Error(err))
		return err
	}

	return nil
}
