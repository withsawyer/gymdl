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
	if s.lastProgressTime == nil || currentTime.Sub(*s.lastProgressTime) >= 1*time.Second {
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
		// 构建结构化消息内容
		var parts []string
		parts = append(parts, "🎉 *入库成功！*")

		// 标题（必选字段）
		if title := strings.TrimSpace(videoInfo.Title); title != "" {
			parts = append(parts, fmt.Sprintf("📺 *标题:* %s", utils.TruncateString(title, 80)))
		}

		// 作者（可选字段）
		if author := strings.TrimSpace(videoInfo.Author); author != "" {
			parts = append(parts, fmt.Sprintf("🎤 *作者:* %s", utils.TruncateString(author, 40)))
		}

		// 分辨率（可选字段）
		if ratio := strings.TrimSpace(videoInfo.Ratio); ratio != "" {
			parts = append(parts, fmt.Sprintf("🎥 *分辨率:* %s", ratio))
		}

		// 创建时间（可选字段）
		if createTime := strings.TrimSpace(videoInfo.Time); createTime != "" {
			parts = append(parts, fmt.Sprintf("🕒 *发布时间:* %s", createTime))
		}

		// 封面（可选字段）
		if coverUrl := strings.TrimSpace(videoInfo.CoverUrl); coverUrl != "" {
			parts = append(parts, fmt.Sprintf("📷 *封面:* %s", coverUrl))
		}

		// 简介（可选字段）
		if desc := strings.TrimSpace(videoInfo.Desc); desc != "" {
			parts = append(parts, fmt.Sprintf("📝 *简介:* %s", utils.TruncateString(desc, 400)))
		}

		// 文件大小（可选字段）
		if fileSize := strings.TrimSpace(videoInfo.Size); fileSize != "" {
			parts = append(parts, fmt.Sprintf("📦 *大小:* %s", fileSize))
		}

		// 入库方式（必选字段）
		storageType := processor.DetermineTidyType(s.Cfg)
		parts = append(parts, fmt.Sprintf("☁️ *入库方式:* %s", storageType))

		// 合并所有非空部分
		successMsg := strings.Join(parts, "\n")
		_, _ = bot.Edit(msg, successMsg, tb.ModeMarkdown)
		return
	}

	// 🎶 多曲反馈
	var listBuilder strings.Builder
	for i, v := range videos {
		// 为每个视频创建结构化消息组件
		var videoParts []string

		// 标题（必选字段）
		if title := strings.TrimSpace(v.Title); title != "" {
			videoParts = append(videoParts, fmt.Sprintf("📺 *标题:* %s", utils.TruncateString(title, 60)))
		}

		// 作者（可选字段）
		if author := strings.TrimSpace(v.Author); author != "" {
			videoParts = append(videoParts, fmt.Sprintf("🎤 *作者:* %s", utils.TruncateString(author, 40)))
		}

		// 分辨率（可选字段）
		if ratio := strings.TrimSpace(v.Ratio); ratio != "" {
			videoParts = append(videoParts, fmt.Sprintf("🎥 *分辨率:* %s", ratio))
		}

		// 文件大小（可选字段）
		if fileSize := strings.TrimSpace(v.Size); fileSize != "" {
			videoParts = append(videoParts, fmt.Sprintf("📦 *大小:* %s", fileSize))
		}

		// 合并当前视频的非空字段
		if len(videoParts) > 0 {
			listBuilder.WriteString(strings.Join(videoParts, "\n"))

			// 添加分隔线（最后一个视频不添加）
			if i < count-1 {
				listBuilder.WriteString("\n──────────────────\n")
			}
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
