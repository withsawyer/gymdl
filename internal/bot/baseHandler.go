package bot

import (
	"fmt"
	"time"

	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/core/handler"
	"github.com/nichuanfang/gymdl/utils"
	tb "gopkg.in/telebot.v4"
)

// HandleText 精简版交互逻辑
func HandleText(c tb.Context) error {
	user := c.Sender()
	text := c.Text()
	bot := c.Bot()

	utils.InfoWithFormat("[Telegram] 用户 %s(%d) 提交内容: %s", user.Username, user.ID, text)

	// 1️⃣ 初始提示
	msg, _ := bot.Send(user, "🔍 正在识别链接...")

	// 2️⃣ 解析链接
	link, executor := handler.ParseLink(text)
	if link == "" || executor == nil {
		_, _ = bot.Edit(msg, "❌ 暂不支持该类型的链接")
		return nil
	}

	utils.InfoWithFormat("[Telegram] 解析成功: %s", link)
	_, _ = bot.Edit(msg, fmt.Sprintf("✅ 已识别 **%s** 链接：\n\n🎵 下载中,请稍候...", executor.Platform()), tb.ModeMarkdown)

	start := time.Now()

	// 3️⃣ 下载阶段
	if err := executor.DownloadMusic(link, app.cfg); err != nil {
		utils.ErrorWithFormat("[Telegram] 下载失败: %v", err)
		bot.Edit(msg, fmt.Sprintf("❌ 下载失败：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return nil
	}

	// 4️⃣ 文件整理 & 处理
	if err := executor.BeforeTidy(app.cfg); err != nil {
		utils.ErrorWithFormat("[Telegram] 文件处理失败: %v", err)
		bot.Edit(msg, fmt.Sprintf("⚠️ 文件处理阶段出错：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return nil
	}

	if err := executor.TidyMusic(app.cfg, core.GlobalWebDAV); err != nil {
		utils.ErrorWithFormat("[Telegram] 文件整理失败: %v", err)
		bot.Edit(msg, fmt.Sprintf("⚠️ 文件整理失败：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return nil
	}

	// 5️⃣ 成功反馈
	duration := time.Since(start)
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60
	successMsg := fmt.Sprintf(
		"🎉 入库成功！耗时 %d分%d秒",
		minutes, seconds,
	)

	_, _ = bot.Edit(msg, successMsg, tb.ModeMarkdown)
	utils.InfoWithFormat("[Telegram] ✅ 用户 %s(%d) 下载成功 (%d分%d秒)", user.Username, user.ID, minutes, seconds)

	return nil
}
