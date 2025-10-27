package bot

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/utils"
	"go.uber.org/zap"
	tb "gopkg.in/telebot.v4"
)

var (
	logger *zap.Logger
	app    *BotApp
)

type BotApp struct {
	bot *tb.Bot
	cfg *config.Config
}

// NewBotApp 创建机器人
func NewBotApp(cfg *config.Config) (*BotApp, error) {
	logger = utils.Logger()
	botSettings := tb.Settings{
		Token: cfg.Telegram.BotToken,
		// 默认使用 polling，后面可切换
		Poller: &tb.LongPoller{Timeout: 10},
	}

	if cfg.ProxyConfig.Scheme != "" && cfg.ProxyConfig.Host != "" && cfg.ProxyConfig.Port != 0 {
		botSettings.Client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(&url.URL{
					Scheme: cfg.ProxyConfig.Scheme,
					Host:   fmt.Sprintf("%s:%d", cfg.ProxyConfig.Host, cfg.ProxyConfig.Port),
				}),
			},
		}
	}

	if cfg.Telegram.Mode == 2 {
		botSettings.Poller = &tb.Webhook{
			Listen: ":" + fmt.Sprint(cfg.Telegram.WebhookPort),
			Endpoint: &tb.WebhookEndpoint{
				PublicURL: cfg.Telegram.WebhookURL,
			},
		}
	}

	bot, err := tb.NewBot(botSettings)
	if err != nil {
		return nil, err
	}

	app = &BotApp{
		bot: bot,
		cfg: cfg,
	}
	app.registerHandlers()
	return app, nil
}

// Start 启动机器人
func (app *BotApp) Start() {
	//移除webhook
	_ = app.bot.RemoveWebhook(true)
	app.bot.Start()
}

// Stop 关闭机器人
func (app *BotApp) Stop() {
	app.bot.Stop()
}

// registerCommands 注册命令 只需运行一次
func (app *BotApp) registerCommands() {

}
