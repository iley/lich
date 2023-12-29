package handlers

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/iley/lich/internal/telegram"
	"github.com/iley/lich/internal/torrents"
)

func MakeStatusHandler(downloader *torrents.Downloader) telegram.Handler {
	return func(bot *telegram.Bot, msg *tgbotapi.Message) (bool, telegram.Handler, error) {
		list := downloader.List()
		if len(list) == 0 {
			bot.SendReply(msg.Chat.ID, "No active downloads")
			return true, nil, nil
		}

		textEntries := make([]string, len(list))
		for i, entry := range list {
			cancelCommand := fmt.Sprintf("/cancel_%s", entry.TorrentId)
			textEntries[i] = fmt.Sprintf("%d: [%s] %s (%s)", i+1, entry.Category, entry.Name, cancelCommand)
		}

		fullText := fmt.Sprintf("Active downloads:\n%s", strings.Join(textEntries, "\n"))
		reply := tgbotapi.NewMessage(msg.Chat.ID, fullText)
		bot.Send(reply)
		return true, nil, nil
	}
}
