package bot

import (
	tb "gopkg.in/telebot.v4"
)

func (app *BotApp) registerHandlers() {
	app.bot.Use(app.LoggingMiddleware)
	app.bot.Use(app.AuthMiddleware)

	app.bot.Handle("/start", func(c tb.Context) error {
		return c.Send("Hello! You are authorized.")
	})

	app.bot.Handle("/echo", func(c tb.Context) error {
		return c.Send(c.Message().Text)
	})
}
