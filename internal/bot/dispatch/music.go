package dispatch

import (
	"fmt"
	"github.com/nichuanfang/gymdl/processor/music"
	"github.com/nichuanfang/gymdl/utils"
	tb "gopkg.in/telebot.v4"
)

// HandleMusic
// ---------------------------
// 🎵 音乐处理逻辑
// ---------------------------
func (s *Session) HandleMusic(p music.Processor) error {
	bot := s.Bot
	msg := s.Msg
	// user := s.User
	// start := s.Start

	_, _ = bot.Edit(msg, fmt.Sprintf("✅ 已识别 **%s** 链接\n\n🎵 下载中,请稍候...", p.Name()), tb.ModeMarkdown)

	// 下载阶段
	utils.InfoWithFormat("[Telegram] 下载中...")
	err := p.DownloadMusic(s.Link)
	if err != nil {
		utils.ErrorWithFormat("[Telegram] 下载失败: %v", err)
		_, _ = bot.Edit(msg, fmt.Sprintf("❌ 下载失败：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return nil
	}

	// 4️⃣ 文件整理 & 处理
	utils.InfoWithFormat("[Telegram] 下载成功，整理中...")
	if err := p.BeforeTidy(); err != nil {
		utils.ErrorWithFormat("[Telegram] 文件处理失败: %v", err)
		_, _ = bot.Edit(msg, fmt.Sprintf("⚠️ 文件处理阶段出错：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return nil
	}

	if err := p.TidyMusic(); err != nil {
		utils.ErrorWithFormat("[Telegram] 文件整理失败: %v", err)
		_, _ = bot.Edit(msg, fmt.Sprintf("⚠️ 文件整理失败：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return nil
	}

	utils.InfoWithFormat("[Telegram] 整理成功，开始入库...")
	if s.Cfg.MusicTidy.Mode == 2 {
		_, _ = bot.Edit(msg, fmt.Sprintf("✅ 已识别 **%s** 链接\n\n🎵 开始入库...", p.Name()), tb.ModeMarkdown)
	}

	// 5️⃣ 成功反馈
	// duration := time.Since(start)
	// minutes := int(duration.Minutes())
	// seconds := int(duration.Seconds()) % 60

	// ✅ 构造详细入库成功提示
	/*fileSizeMB := float64(music.MusicSize) / 1024.0 / 1024.0
	  	successMsg := fmt.Sprintf(
	  		`🎉 *入库成功！*

	  🎵 *歌曲:* %s
	  🎤 *艺术家:* %s
	  💿 *专辑:* %s
	  🎧 *格式:* %s
	  📊 *码率:* %s kbps
	  📦 *大小:* %.2f MB
	  ☁️ *入库方式:* %s`,
	  		utils.TruncateString(music.SongName, 80),
	  		utils.TruncateString(music.SongArtists, 80),
	  		utils.TruncateString(music.SongAlbum, 80),
	  		strings.ToUpper(music.FileExt),
	  		music.Bitrate,
	  		fileSizeMB,
	  		strings.ToUpper(music.Tidy),
	  	)

	  	_, _ = bot.Edit(msg, successMsg, tb.ModeMarkdown)*/

	// utils.InfoWithFormat("[Telegram] ✅ 用户 %s(%d) 下载成功 (%d分%d秒) -> %s", user.Username, user.ID, minutes, seconds, music.SongName)
	return nil
}
