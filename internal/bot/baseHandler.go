package bot

import (
	"fmt"
	"github.com/nichuanfang/gymdl/core/downloader"
	"github.com/nichuanfang/gymdl/core/downloader/music"
	"github.com/nichuanfang/gymdl/core/downloader/video"
	"strings"
	"time"

	"github.com/nichuanfang/gymdl/core"
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
	link, executor := downloader.ParseLink(text)
	utils.InfoWithFormat("[Telegram] 解析成功: %s", link)
	if link == "" || executor == nil {
		_, _ = bot.Edit(msg, "❌ 暂不支持该类型的链接")
		return nil
	}
	switch expr := executor.(type) {
	case music.Handler:
		_musicHandler(c, msg, link, expr)
	case video.Handler:
		_videoHandler(c, msg, link, expr)
	}
	return nil
}

func _musicHandler(c tb.Context, msg *tb.Message, link string, executor music.Handler) {
	user := c.Sender()
	bot := c.Bot()

	_, _ = bot.Edit(msg, fmt.Sprintf("✅ 已识别 **%s** 链接\n\n🎵 下载中,请稍候...", executor.Platform()), tb.ModeMarkdown)

	start := time.Now()

	// 3️⃣ 下载阶段
	utils.InfoWithFormat("[Telegram] 下载中...")
	songInfo, err := executor.Download(link, app.cfg)
	if err != nil {
		utils.ErrorWithFormat("[Telegram] 下载失败: %v", err)
		_, _ = bot.Edit(msg, fmt.Sprintf("❌ 下载失败：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return
	}

	// 4️⃣ 文件整理 & 处理
	utils.InfoWithFormat("[Telegram] 下载成功，整理中...")
	if err := executor.BeforeTidy(app.cfg, songInfo); err != nil {
		utils.ErrorWithFormat("[Telegram] 文件处理失败: %v", err)
		_, _ = bot.Edit(msg, fmt.Sprintf("⚠️ 文件处理阶段出错：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return
	}

	if err := executor.TidyMusic(app.cfg, core.GlobalWebDAV, songInfo); err != nil {
		utils.ErrorWithFormat("[Telegram] 文件整理失败: %v", err)
		_, _ = bot.Edit(msg, fmt.Sprintf("⚠️ 文件整理失败：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return
	}

	utils.InfoWithFormat("[Telegram] 整理成功，开始入库...")
	if app.cfg.ResourceTidy.Mode == 2 {
		_, _ = bot.Edit(msg, fmt.Sprintf("✅ 已识别 **%s** 链接\n\n🎵 开始入库...", executor.Platform()), tb.ModeMarkdown)
	}

	// 5️⃣ 成功反馈
	duration := time.Since(start)
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60

	// ✅ 构造详细入库成功提示
	fileSizeMB := float64(songInfo.MusicSize) / 1024.0 / 1024.0
	successMsg := fmt.Sprintf(
		`🎉 *入库成功！*

🎵 *歌曲:* %s
🎤 *艺术家:* %s
💿 *专辑:* %s
🎧 *格式:* %s
📊 *码率:* %s kbps
📦 *大小:* %.2f MB
☁️ *入库方式:* %s`,
		utils.TruncateString(songInfo.SongName, 80),
		utils.TruncateString(songInfo.SongArtists, 80),
		utils.TruncateString(songInfo.SongAlbum, 80),
		strings.ToUpper(songInfo.FileExt),
		songInfo.Bitrate,
		fileSizeMB,
		strings.ToUpper(songInfo.Tidy),
	)

	_, _ = bot.Edit(msg, successMsg, tb.ModeMarkdown)

	utils.InfoWithFormat("[Telegram] ✅ 用户 %s(%d) 下载成功 (%d分%d秒) -> %s", user.Username, user.ID, minutes, seconds, songInfo.SongName)
}

func _videoHandler(c tb.Context, msg *tb.Message, link string, executor video.Handler) {
	user := c.Sender()
	bot := c.Bot()

	_, _ = bot.Edit(msg, fmt.Sprintf("✅ 已识别 **%s** 链接\n\n🎵 下载中,请稍候...", executor.Platform()), tb.ModeMarkdown)

	start := time.Now()

	// 3️⃣ 下载阶段
	utils.InfoWithFormat("[Telegram] 视频下载中...")
	songInfo, err := executor.Download(link, app.cfg)
	if err != nil {
		utils.ErrorWithFormat("[Telegram] 下载失败: %v", err)
		_, _ = bot.Edit(msg, fmt.Sprintf("❌ 视频下载失败：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return
	}

	// 4️⃣ 文件整理 & 处理
	utils.InfoWithFormat("[Telegram] 下载成功，整理中...")
	if err := executor.BeforeTidy(app.cfg, songInfo); err != nil {
		utils.ErrorWithFormat("[Telegram] 文件处理失败: %v", err)
		_, _ = bot.Edit(msg, fmt.Sprintf("⚠️ 文件处理阶段出错：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return
	}

	utils.InfoWithFormat("[Telegram] 整理成功，开始入库...")
	if app.cfg.ResourceTidy.Mode == 2 {
		_, _ = bot.Edit(msg, fmt.Sprintf("✅ 已识别 **%s** 链接\n\n🎵 开始入库...", executor.Platform()), tb.ModeMarkdown)
	}

	// 5️⃣ 成功反馈
	duration := time.Since(start)
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60

	// ✅ 构造详细入库成功提示
	fileSizeMB := float64(songInfo.Size) / 1024.0 / 1024.0
	successMsg := fmt.Sprintf(
		`🎉 *入库成功！*

🎵 *名称:* %s
🎧 *格式:* %s
📊 *码率:* %s kbps
📦 *大小:* %.2f MB
☁️ *入库方式:* %s`,
		utils.TruncateString(songInfo.Title, 80),
		strings.ToUpper(songInfo.FileExt),
		songInfo.Bitrate,
		fileSizeMB,
		strings.ToUpper(songInfo.Tidy),
	)
	_, _ = bot.Edit(msg, successMsg, tb.ModeMarkdown)

	utils.InfoWithFormat("[Telegram] ✅ 用户 %s(%d) 下载成功 (%d分%d秒) -> %s", user.Username, user.ID, minutes, seconds, songInfo.Title)
}
