// 通知模块
package bot

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	tb "gopkg.in/telebot.v4"
)

// NotifierInterface 定义发送通知的接口
type Notifier interface {
	Send(msg string)
}

type botNotifier struct {
	Bot    *tb.Bot
	ChatID int64
}

// Send 异步发送消息
func (b *botNotifier) Send(msg string) {
	go func() {
		_, err := b.Bot.Send(&tb.Chat{ID: b.ChatID}, msg)
		if err != nil {
			fmt.Errorf("[Telegram] 发送消息失败: %v\n", err)
		}
	}()
}

var (
	instance Notifier
	once     sync.Once
)

// InitBotNotifier 初始化全局单例，只能调用一次
func InitBotNotifier(bot *tb.Bot, chatID string) {
	chatId, _ := strconv.ParseInt(chatID, 10, 64)
	once.Do(func() {
		instance = &botNotifier{
			Bot:    bot,
			ChatID: chatId,
		}
	})
}

// GetNotifier 获取全局单例
func GetNotifier() Notifier {
	return instance
}

// SendMessage 方便调用
func SendMessage(msg string) {
	if instance == nil {
		log.Println("[Notifier] 尚未初始化，无法发送消息")
		return
	}
	instance.Send(msg)
}
