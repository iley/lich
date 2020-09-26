package handlers

import (
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/iley/lich/internal/telegram"
)

func MakeHelpHandler(commands []string) telegram.Handler {
	text := "Available commands: " + strings.Join(commands, ", ")
	return func(bot *telegram.Bot, msg *tgbotapi.Message) (bool, telegram.Handler, error) {
		reply := tgbotapi.NewMessage(msg.Chat.ID, text)
		bot.Send(reply)
		return true, nil, nil
	}
}
