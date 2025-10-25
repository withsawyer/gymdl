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
	bot.Edit(msg, "链接分析成功,准备下载...")
	handler.HandlerMusic(link)
	//音乐下载成功
	time.Sleep(time.Second * 2)
	bot.Edit(msg, "音乐下载成功!")
	return nil
}
