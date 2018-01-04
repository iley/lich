package main

import (
	"gopkg.in/telegram-bot-api.v1"
	"strings"
)

func MakeHelpHandler(commands []string) TelegramHandler {
	text := "Available commands: " + strings.Join(commands, ", ")
	return func(bot *TelegramBot, msg *tgbotapi.Message) (bool, TelegramHandler, error) {
		reply := tgbotapi.NewMessage(msg.Chat.ID, text)
		bot.Api.Send(reply)
		return true, nil, nil
	}
}
