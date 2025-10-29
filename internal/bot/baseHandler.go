package bot

import (
	"fmt"
	"time"

	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/core/factory"
	"github.com/nichuanfang/gymdl/core/linkparser"
	"github.com/nichuanfang/gymdl/processor/music"
	"github.com/nichuanfang/gymdl/processor/video"
	"github.com/nichuanfang/gymdl/utils"
	tb "gopkg.in/telebot.v4"
)

// tg会话
type Session struct {
	text     string          //用户发送的消息
	context  tb.Context      //tg上下文
	user     *tb.User        //用户
	bot      tb.API          //机器人
	msg      *tb.Message     //初始化消息对象
	link     string          //有效链接
	linkType domain.LinkType //链接类型
	start    time.Time       //开始处理时间
}

// HandleText 精简版交互逻辑
func HandleText(c tb.Context) error {
	text := c.Text()
	user := c.Sender()
	b := c.Bot()

	//初始提示
	msg, _ := b.Send(user, "🔍 正在识别链接...")

	//解析链接link:有效链接 linkType:链接类型
	link, linkType := linkparser.ParseLink(text)

	if link == "" {
		_, _ = b.Edit(msg, "❌ 暂不支持该类型的链接")
		return nil
	}
	utils.InfoWithFormat("[Telegram] 解析成功: %s", link)
	proc := factory.GetProcessor(linkType, app.cfg)

	if proc == nil {
		return c.Send("未找到处理器")
	}
	//创建会话对象
	session := &Session{
		text:     text,
		context:  c,
		user:     user,
		bot:      b,
		msg:      msg,
		link:     link,
		linkType: linkType,
		start:    time.Now(),
	}
	// 根据不同类型的处理器执行不同逻辑
	switch proc.Category() {
	case domain.CategoryMusic:
		// 断言为 MusicProcessor
		if mp, ok := proc.(music.MusicProcessor); ok {
			return handleMusic(session, mp)
		}
		// 如果没有实现特定接口，也可以退回通用 Handle
		res, err := proc.Handle(text)
		if err != nil {
			return c.Send("处理失败：" + err.Error())
		}
		return c.Send(res)
	case domain.CategoryVideo:
		// 断言为 VideoProcessor
		if vp, ok := proc.(video.VideoProcessor); ok {
			return handleVideo(session, vp)
		}
		res, err := proc.Handle(text)
		if err != nil {
			return c.Send("处理失败：" + err.Error())
		}
		return c.Send(res)
	default:
		return c.Send("未知处理器类型")
	}
}

// ---------------------------
// 🎵 音乐处理逻辑
// ---------------------------
func handleMusic(session *Session, p music.MusicProcessor) error {
	bot := session.bot
	msg := session.msg
	//user := session.user
	//start := session.start

	_, _ = bot.Edit(msg, fmt.Sprintf("✅ 已识别 **%s** 链接\n\n🎵 下载中,请稍候...", p.Name()), tb.ModeMarkdown)

	//下载阶段
	utils.InfoWithFormat("[Telegram] 下载中...")
	err := p.DownloadMusic(session.link)
	if err != nil {
		utils.ErrorWithFormat("[Telegram] 下载失败: %v", err)
		bot.Edit(msg, fmt.Sprintf("❌ 下载失败：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return nil
	}

	// 4️⃣ 文件整理 & 处理
	utils.InfoWithFormat("[Telegram] 下载成功，整理中...")
	if _, err := p.BeforeTidy(); err != nil {
		utils.ErrorWithFormat("[Telegram] 文件处理失败: %v", err)
		bot.Edit(msg, fmt.Sprintf("⚠️ 文件处理阶段出错：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return nil
	}

	if err := p.TidyMusic(); err != nil {
		utils.ErrorWithFormat("[Telegram] 文件整理失败: %v", err)
		bot.Edit(msg, fmt.Sprintf("⚠️ 文件整理失败：\n```\n%s\n```", utils.TruncateString(err.Error(), 400)), tb.ModeMarkdown)
		return nil
	}

	utils.InfoWithFormat("[Telegram] 整理成功，开始入库...")
	if app.cfg.MusicTidy.Mode == 2 {
		_, _ = bot.Edit(msg, fmt.Sprintf("✅ 已识别 **%s** 链接\n\n🎵 开始入库...", p.Name()), tb.ModeMarkdown)
	}

	// 5️⃣ 成功反馈
	//duration := time.Since(start)
	//minutes := int(duration.Minutes())
	//seconds := int(duration.Seconds()) % 60

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

	//utils.InfoWithFormat("[Telegram] ✅ 用户 %s(%d) 下载成功 (%d分%d秒) -> %s", user.Username, user.ID, minutes, seconds, music.SongName)
	return nil
}

// ---------------------------
// 📺 视频处理逻辑
// ---------------------------
func handleVideo(session *Session, p video.VideoProcessor) error {
	// 假设视频处理器能返回缩略图 URL，可扩展为发送照片
	// c.Send(&tb.Photo{File: tb.FromURL(info.Thumbnail), Caption: text}, tb.ModeHTML)
	return session.context.Send("视频处理逻辑", tb.ModeHTML)
}
