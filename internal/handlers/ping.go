package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/iley/lich/internal/telegram"
)

func MakePingHandler() telegram.Handler {
	return func(bot *telegram.Bot, msg *tgbotapi.Message) (bool, telegram.Handler, error) {
		reply := tgbotapi.NewMessage(msg.Chat.ID, "pong")
		bot.Send(reply)
		return true, nil, nil
	}
}
