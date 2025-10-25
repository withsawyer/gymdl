package handlers

import (
	"github.com/nichuanfang/gymdl/core"
	tb "gopkg.in/telebot.v4"
)

// 普通文本处理器
func HandleText(c tb.Context) error {
	ask, _ := core.GlobalAI.Ask(c.Text())
	return c.Send(ask)
}
