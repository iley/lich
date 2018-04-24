package main

import (
	"gopkg.in/telegram-bot-api.v4"
)

func HandlePing(bot *TelegramBot, msg *tgbotapi.Message) (bool, TelegramHandler, error) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, "pong")
	reply.ReplyToMessageID = msg.MessageID
	bot.Api.Send(reply)
	return true, nil, nil
}
