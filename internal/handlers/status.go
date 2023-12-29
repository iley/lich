package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/iley/lich/internal/telegram"
	"github.com/iley/lich/internal/torrents"
)

func MakeStatusHandler(downloader *torrents.Downloader) telegram.Handler {
	return func(bot *telegram.Bot, msg *tgbotapi.Message) (bool, telegram.Handler, error) {
		text := downloader.StatusString()
		reply := tgbotapi.NewMessage(msg.Chat.ID, text)
		bot.Send(reply)
		return true, nil, nil
	}
}
