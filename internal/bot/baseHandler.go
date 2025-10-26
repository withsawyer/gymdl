package bot

import (
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/core/handler"
	tb "gopkg.in/telebot.v4"
)

// 普通文本处理器
func HandleText(c tb.Context) error {
	user := c.Sender()
	text := c.Text()

	bot := c.Bot()
	msg, _ := bot.Send(user, "正在分析链接")
	link, executor := handler.ParseLink(text)
	if link == "" {
		bot.Edit(msg, "不支持的链接!")
		return nil
	}
	logger.Info("[Telegram] 处理链接:" + link)
	bot.Edit(msg, "链接分析成功,准备下载...")
	err := executor.DownloadMusic(link, app.cfg)
	//模拟音乐下载
	if err != nil {
		bot.Edit(msg, err.Error())
		return err
	}
	//准备整理
	err = executor.BeforeTidy(app.cfg)
	if err != nil {
		return err
	}
	//整理
	err = executor.TidyMusic(app.cfg, core.GlobalWebDAV)
	if err != nil {
		return err
	}
	//音乐下载成功
	bot.Edit(msg, "音乐下载成功!")
	return nil
}
