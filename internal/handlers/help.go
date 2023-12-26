package handlers

import (
	"fmt"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"

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
