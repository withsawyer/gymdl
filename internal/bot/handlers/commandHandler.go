package handlers

import (
	"github.com/nichuanfang/gymdl/utils"
	"go.uber.org/zap"
	tb "gopkg.in/telebot.v4"
)

var logger *zap.Logger

func init() {
	logger = utils.Logger()
}

// start处理器
func HandleStartCommand(c tb.Context) error {
	return c.Send("Hello! You are authorized.")
}
