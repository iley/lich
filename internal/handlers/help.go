package handlers

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/iley/lich/internal/telegram"
)

func MakeHelpHandler(commands []string, versionString string) telegram.Handler {
	commandList := strings.Join(commands, ", ")
	text := fmt.Sprintf("Lich v%s\nAvailable commands: %s", versionString, commandList)
	return func(bot *telegram.Bot, msg *tgbotapi.Message) (bool, telegram.Handler, error) {
		reply := tgbotapi.NewMessage(msg.Chat.ID, text)
		bot.Send(reply)
		return true, nil, nil
	}
}
