package bot

import (
	"github.com/nichuanfang/gymdl/internal/bot/handlers"
	tb "gopkg.in/telebot.v4"
)

func (app *BotApp) registerHandlers() {
	app.bot.Use(app.LoggingMiddleware)
	app.bot.Use(app.AuthMiddleware)

	//欢迎语
	app.bot.Handle("/start", handlers.StartCommand)

	//帮助信息
	app.bot.Handle("/help", handlers.HelpCommand)

	//指令注册器
	app.bot.Handle("/setCommands", handlers.SetCommands)

	//普通文本
	app.bot.Handle(tb.OnText, handlers.HandleText)
}
