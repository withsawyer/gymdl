package bot

import (
	"github.com/nichuanfang/gymdl/internal/bot/handlers"
	tb "gopkg.in/telebot.v4"
)

func (app *BotApp) registerHandlers() {
	app.bot.Use(app.LoggingMiddleware)
	app.bot.Use(app.AuthMiddleware)

	//欢迎语
	app.bot.Handle("/start", handlers.HandleStartCommand)

	//指令注册器
	app.bot.Handle("/setCommands", handlers.InitCommands)

	//普通文本
	app.bot.Handle(tb.OnText, handlers.HandleText)
}
