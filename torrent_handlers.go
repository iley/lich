package main

import (
	"errors"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"strings"
)

func (d *TorrentDownloader) HandleTorrentFile(bot *TelegramBot, msg *tgbotapi.Message) (bool, TelegramHandler, error) {
	if msg.Document == nil || msg.Document.FileID == "" || !IsTorrent(msg.Document) {
		return false, nil, nil
	}
	text := fmt.Sprintf("Received torrent file %s. Torrent file downloading not implemented.", msg.Document.FileName)
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	bot.Api.Send(reply)
	return true, nil, nil
}

func (d *TorrentDownloader) HandleMagnetLink(bot *TelegramBot, msg *tgbotapi.Message) (bool, TelegramHandler, error) {
	if !strings.HasPrefix(msg.Text, "magnet:") {
		return false, nil, nil
	}
	categories := d.Categories()
	text := fmt.Sprintf("What category does this torrent belong to? (%s)", strings.Join(categories, ", "))
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{Keyboard: MakeKeyboard(categories), OneTimeKeyboard: true}
	bot.Api.Send(reply)
	return true, d.MakeCategoryHandler(msg.Text), nil
}

func IsTorrent(document *tgbotapi.Document) bool {
	return strings.HasSuffix(document.FileName, ".torrent")
}

func MakeKeyboard(keys []string) [][]tgbotapi.KeyboardButton {
	keysPerRow := 3
	keyboard := [][]tgbotapi.KeyboardButton{}
	for i := 0; i < len(keys); i += keysPerRow {
		j := i + keysPerRow
		if j > len(keys) {
			j = len(keys)
		}
		row := make([]tgbotapi.KeyboardButton, j-i)
		for k := i; k < j; k++ {
			row[k] = tgbotapi.NewKeyboardButton(keys[k])
		}
		keyboard = append(keyboard, row)
	}
	return keyboard
}

func (d *TorrentDownloader) MakeCategoryHandler(torrentUrl string) TelegramHandler {
	return func(bot *TelegramBot, msg *tgbotapi.Message) (bool, TelegramHandler, error) {
		category := msg.Text
		_, found := d.config.TargetDirs[category]
		if found {
			replyFunc := func(text string) {
				reply := tgbotapi.NewMessage(msg.Chat.ID, text)
				bot.Api.Send(reply)
			}
			request := DownloadRequest{MagnetLink: torrentUrl, Category: category, Reply: replyFunc}
			select {
			case d.requests <- &request:
			default:
				return true, nil, errors.New("Download queue full")
			}
			return true, nil, nil
		} else {
			text := fmt.Sprintf("Unknown category %s. Pick one of %s", category, strings.Join(d.Categories(), ", "))
			reply := tgbotapi.NewMessage(msg.Chat.ID, text)
			bot.Api.Send(reply)
			return true, d.MakeCategoryHandler(torrentUrl), nil
		}
	}
}
