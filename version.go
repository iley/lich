package main

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"os"
)

var version string

func HandleVersion(bot *TelegramBot, msg *tgbotapi.Message) (bool, TelegramHandler, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown machine"
	}
	ver := version
	if ver == "" {
		ver = "unknown"
	}
	text := fmt.Sprintf("Lich %s running on %s", ver, hostname)
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	bot.Api.Send(reply)
	return true, nil, nil
}
