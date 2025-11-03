package dispatch

import (
	"fmt"
	"strings"
	"time"

	"github.com/nichuanfang/gymdl/processor"

	"github.com/nichuanfang/gymdl/processor/video"
	"github.com/nichuanfang/gymdl/utils"
	tb "gopkg.in/telebot.v4"
)

// ---------------------------
// 📺 视频处理逻辑
// ---------------------------

// ReportProgress 实现video.ProgressReporter接口，限制发送频率为2秒一次
func (s *Session) ReportProgress(progress string) {
	// 检查距离上次发送进度条的时间间隔
	currentTime := time.Now()
	// 如果是第一次发送或者时间间隔大于等于2秒，则发送进度条
	if s.lastProgressTime == nil || currentTime.Sub(*s.lastProgressTime) >= 2*time.Second {
		utils.DebugWithFormat("[Telegram] 发送进度条: %s", progress)
		s._sendVideoProgress(progress)
		// 更新上次发送时间，创建新的时间实例
		s.lastProgressTime = &currentTime
	} else {
		utils.DebugWithFormat("[Telegram] 进度条发送频率限制，距离上次发送间隔: %v", currentTime.Sub(*s.lastProgressTime))
	}
}

func (s *Session) HandleVideo(p video.Processor) error {
	bot := s.Bot
	msg := s.Msg
	// user := s.User
	// start := s.Start

	_, _ = bot.Edit(msg, fmt.Sprintf("✅ 已识别【**%s**】链接\n\n🎵 开始分析资源,请稍候...", p.Name()), tb.ModeMarkdown)

	// 下载阶段
	utils.InfoWithFormat("[Telegram] 正在解析下载资源,请稍候...")
	err := p.Download(s.Link, s)
	if err != nil {
		utils.ErrorWithFormat("[Telegram] 下载失败: %v", err)
		_, _ = bot.Edit(msg, fmt.Sprintf("❌ 下载失败：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return nil
	}
	// 文件整理 & 处理
	utils.InfoWithFormat("[Telegram] 下载成功，整理中...")
	if err := p.Tidy(); err != nil {
		utils.ErrorWithFormat("[Telegram] 文件整理失败: %v", err)
		_, _ = bot.Edit(msg, fmt.Sprintf("⚠️ 文件整理失败：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return nil
	}
	utils.InfoWithFormat("[Telegram] 整理成功，开始入库...")
	if s.Cfg.Tidy.Mode == 2 {
		_, _ = bot.Edit(msg, fmt.Sprintf("✅ 已识别 **%s** 链接\n\n🎵 开始入库...", p.Name()), tb.ModeMarkdown)
	}
	// 成功反馈
	s.sendVideoFeedback(p)
	utils.InfoWithFormat("[Telegram] 入库成功!")
	return nil
}

func (s *Session) sendVideoFeedback(p video.Processor) {
	bot := s.Bot
	msg := s.Msg

	videos := p.Videos()
	count := len(videos)

	if count == 0 {
		_, _ = bot.Edit(msg, "⚠️ 没有检测到有效视频", tb.ModeMarkdown)
		return
	}

	// 🎵 单曲反馈
	if count == 1 {
		videoInfo := videos[0]
		fileSize := videoInfo.Size

		successMsg := fmt.Sprintf(
			`🎉 *入库成功！*
📺 *标题:* %s  
🎤 *作者:* %s  
🎥 *分辨率:* %s  
🕒 *创建时间:* %s
📷 *封面:* %s
📝 *简介:* %s
📦 *大小:* %s
☁️ *入库方式:* %s`,
			utils.TruncateString(videoInfo.Title, 80),
			utils.TruncateString(videoInfo.Author, 40),
			videoInfo.Ratio,
			videoInfo.Time,
			videoInfo.CoverUrl,
			utils.TruncateString(videoInfo.Desc, 400),
			fileSize,
			processor.DetermineTidyType(s.Cfg),
		)
		_, _ = bot.Edit(msg, successMsg, tb.ModeMarkdown)
		return
	}

	// 🎶 多曲反馈
	var listBuilder strings.Builder
	for i, v := range videos {
		fileSize := v.Size
		listBuilder.WriteString(fmt.Sprintf(
			"📺 《%s》\n🎤 作者：%s\n🎥 分辨率:：%s\n📦 大小：%s",
			utils.TruncateString(v.Desc, 60),
			utils.TruncateString(v.Author, 40),
			utils.TruncateString(v.Ratio, 40),
			fileSize,
		))

		// 如果不是最后一首，添加长横线分隔
		if i < count-1 {
			listBuilder.WriteString("\n──────────────────\n")
		} else {
			listBuilder.WriteString("\n")
		}
	}

	successMsg := fmt.Sprintf(
		`🎉 *入库成功！*

已成功添加 *%d* 视频至影库：
──────────────────
%s──────────────────
☁️ *入库方式:* %s
`, count, listBuilder.String(), processor.DetermineTidyType(s.Cfg))

	_, _ = bot.Edit(msg, successMsg, tb.ModeMarkdown)
}

func (s *Session) _sendVideoProgress(progress string) {
	bot := s.Bot
	msg := s.Msg
	_, _ = bot.Edit(msg, progress, tb.ModeMarkdown)
}
