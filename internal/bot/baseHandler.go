package bot

import (
	"time"

	"github.com/nichuanfang/gymdl/core"
	tb "gopkg.in/telebot.v4"
)

// 普通文本处理器
func HandleText(c tb.Context) error {
	user := c.Sender()
	text := c.Text()

	bot := c.Bot()
	msg, _ := bot.Send(user, "正在分析链接")
	link, handler := core.ParseLink(text)
	if link == "" {
		bot.Edit(msg, "不支持的链接!")
		return nil
	}
	logger.Info("[Telegram] 处理链接:" + link)
	bot.Edit(msg, "链接分析成功,准备下载...")
	err := handler.DownloadMusic(link, app.cfg)
	//模拟音乐下载
	time.Sleep(time.Second * 2)
	if err != nil {
		bot.Edit(msg, "音乐下载失败:"+err.Error())
		return err
	}
	//音乐下载成功
	bot.Edit(msg, "音乐下载成功!")
	return nil
}
