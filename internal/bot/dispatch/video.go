package dispatch

import (
	"fmt"

	"github.com/nichuanfang/gymdl/processor/video"
	"github.com/nichuanfang/gymdl/utils"
	tb "gopkg.in/telebot.v4"
)

// HandleVideo
// ---------------------------
// 📺 视频处理逻辑
// ---------------------------
func (s *Session) HandleVideo(p video.Processor) error {
	bot := s.Bot
	msg := s.Msg
	// user := s.User
	// start := s.Start

	_, _ = bot.Edit(msg, fmt.Sprintf("✅ 已识别【**%s**】链接\n\n🎵 下载中,请稍候...", p.Name()), tb.ModeMarkdown)

	// 下载阶段
	utils.InfoWithFormat("[Telegram] 下载中...")
	err := p.Download(s.Link)
	if err != nil {
		utils.ErrorWithFormat("[Telegram] 下载失败: %v", err)
		_, _ = bot.Edit(msg, fmt.Sprintf("❌ 下载失败：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return nil
	}

	// 4️⃣ 文件整理 & 处理
	utils.InfoWithFormat("[Telegram] 下载成功，整理中...")
	//if err := p.BeforeTidy(); err != nil {
	//	utils.ErrorWithFormat("[Telegram] 文件处理失败: %v", err)
	//	_, _ = bot.Edit(msg, fmt.Sprintf("⚠️ 文件处理阶段出错：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
	//	return nil
	//}

	//if err := p.TidyMusic(); err != nil {
	//	utils.ErrorWithFormat("[Telegram] 文件整理失败: %v", err)
	//	_, _ = bot.Edit(msg, fmt.Sprintf("⚠️ 文件整理失败：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
	//	return nil
	//}

	utils.InfoWithFormat("[Telegram] 整理成功，开始入库...")
	if s.Cfg.Tidy.Mode == 2 {
		_, _ = bot.Edit(msg, fmt.Sprintf("✅ 已识别 **%s** 链接\n\n🎵 开始入库...", p.Name()), tb.ModeMarkdown)
	}
	return nil
}
