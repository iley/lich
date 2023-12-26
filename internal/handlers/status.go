package handlers

import (
	"fmt"

	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/iley/lich/internal/telegram"
	"github.com/iley/lich/internal/torrents"
)

func MakeStatusHandler(downloader *torrents.Downloader) telegram.Handler {
	return func(bot *telegram.Bot, msg *tgbotapi.Message) (bool, telegram.Handler, error) {
		text := ""
		if downloader.GetInProgressCount() > 0 {
			pluralSuffix := ""
			if downloader.GetInProgressCount() > 1 {
				pluralSuffix = "s"
			}
			text = fmt.Sprintf("Downloading %d torrent%s", downloader.GetInProgressCount(), pluralSuffix)
		} else {
			text = "No torrents are being downloaded"
		}
		reply := tgbotapi.NewMessage(msg.Chat.ID, text)
		bot.Send(reply)
		return true, nil, nil
	}
}
