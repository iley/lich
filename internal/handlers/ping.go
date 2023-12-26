package handlers

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/iley/lich/internal/telegram"
)

func MakePingHandler() telegram.Handler {
	return func(bot *telegram.Bot, msg *tgbotapi.Message) (bool, telegram.Handler, error) {
		reply := tgbotapi.NewMessage(msg.Chat.ID, "pong")
		bot.Send(reply)
		return true, nil, nil
	}
}
