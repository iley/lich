package handlers

import (
	"io"
	"fmt"
	"net/http"
	"log"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/iley/lich/internal/config"
	"github.com/iley/lich/internal/telegram"
	"github.com/iley/lich/internal/torrents"
)

func MakeTorrentFileHandler() telegram.Handler {
	return func(bot *telegram.Bot, msg *tgbotapi.Message) (bool, telegram.Handler, error) {
		if msg.Document == nil || msg.Document.FileID == "" || !isTorrent(msg.Document) {
			return false, nil, nil
		}
		text := fmt.Sprintf("Received torrent file %s. Torrent file downloading not implemented.", msg.Document.FileName)
		bot.SendReply(msg.Chat.ID, text)
		return true, nil, nil
	}
}

func MakeMagnetLinkHandler(cfg *config.Config, down *torrents.Downloader) telegram.Handler {
	return func(bot *telegram.Bot, msg *tgbotapi.Message) (bool, telegram.Handler, error) {
		magnetLink := extractMagnetLink(msg.Text)
		if magnetLink == "" {
			return false, nil, nil
		}
		categories := cfg.Categories()
		text := fmt.Sprintf("What category does this torrent belong to? (%s)", strings.Join(categories, ", "))
		reply := tgbotapi.NewMessage(msg.Chat.ID, text)
		reply.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{Keyboard: makeKeyboard(categories), OneTimeKeyboard: true}
		bot.Send(reply)
		return true, makeCategoryHandler(cfg, down, magnetLink), nil
	}
}

func makeCategoryHandler(cfg *config.Config, down *torrents.Downloader, torrentUrl string) telegram.Handler {
	return func(bot *telegram.Bot, msg *tgbotapi.Message) (bool, telegram.Handler, error) {
		category := msg.Text
		_, found := cfg.TargetDirs[category]
		if found {
			request := torrents.DownloadRequest{MagnetLink: torrentUrl, Category: category, ChatId: msg.Chat.ID}
			err := down.Add(&request)
			return true, nil, err
		}

		text := fmt.Sprintf("Unknown category %s. Pick one of %s", category, strings.Join(cfg.Categories(), ", "))
		bot.SendReply(msg.Chat.ID, text)
		return true, makeCategoryHandler(cfg, down, torrentUrl), nil
	}
}

func isTorrent(document *tgbotapi.Document) bool {
	return strings.HasSuffix(document.FileName, ".torrent")
}

func makeKeyboard(keys []string) [][]tgbotapi.KeyboardButton {
	keysPerRow := 3
	keyboard := [][]tgbotapi.KeyboardButton{}
	for i := 0; i < len(keys); i += keysPerRow {
		j := i + keysPerRow
		if j > len(keys) {
			j = len(keys)
		}
		row := make([]tgbotapi.KeyboardButton, j-i)
		for k := 0; k < j-i; k++ {
			row[k] = tgbotapi.NewKeyboardButton(keys[i+k])
		}
		keyboard = append(keyboard, row)
	}
	return keyboard
}

var magnetLinkRegex = regexp.MustCompile(`magnet:\S+`)

func extractMagnetLink(text string) string {
	if isURL(text) {
		var err error
		text, err = downloadWebPage(text)
		if err != nil {
			log.Printf("Could not download web page: %s", err.Error())
			return ""
		}
	}
	return magnetLinkRegex.FindString(text)
}

func isURL(text string) bool {
	return strings.HasPrefix(text, "http://") || strings.HasPrefix(text, "https://")
}

func downloadWebPage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func MakeCancelHandler(down *torrents.Downloader) telegram.Handler {
	return func(bot *telegram.Bot, msg *tgbotapi.Message) (bool, telegram.Handler, error) {
		torrentId := strings.TrimPrefix(msg.Text, "/cancel_")
		torrentId = strings.TrimSpace(torrentId)
		err := down.Cancel(torrentId)
		if err != nil {
			text := fmt.Sprintf("Could not cancel torrent %s: %s", torrentId, err.Error())
			bot.SendReply(msg.Chat.ID, text)
			return true, nil, nil
		}
		bot.SendReply(msg.Chat.ID, fmt.Sprintf("Successfully cancelled torrent %s", torrentId))
		return true, nil, nil
	}
}
