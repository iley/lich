package telegram

import (
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type chatSession struct {
	input      chan *tgbotapi.Message
	lastAccess time.Time // Protected by mutex.
	mutex      sync.Mutex
}

func newChatSession(bot *Bot) *chatSession {
	session := &chatSession{
		input: make(chan *tgbotapi.Message, 32),
	}
	go session.RunLoop(bot)
	return session
}

func (session *chatSession) EnqueueMessage(message *tgbotapi.Message) error {
	select {
	case session.input <- message:
		return nil
	default:
		return errors.New("Queue is full")
	}
}

func (session *chatSession) Close() {
	close(session.input)
}

func (session *chatSession) RunLoop(bot *Bot) {
	var nextHandler Handler = nil
	for message := range session.input {
		session.mutex.Lock()
		session.lastAccess = time.Now()
		session.mutex.Unlock()
		done := false
		var err error
		if nextHandler != nil {
			done, nextHandler, err = nextHandler(bot, message)
		}
		if !done && message.IsCommand() {
			command := strings.TrimLeft(message.Text, "/")
			if spaceIndex := strings.IndexRune(command, ' '); spaceIndex != -1 {
				command = command[:spaceIndex]
			}
			handler, found := bot.commandHandlers[command]
			if found {
				done, nextHandler, err = handler(bot, message)
			}
		}
		for _, globalHandler := range bot.globalHandlers {
			if done || err != nil {
				break
			}
			done, nextHandler, err = globalHandler(bot, message)
		}
		if !done && err == nil {
			err = errors.New("I don't understand you")
		}
		if err != nil {
			log.Printf("Error while handling a user message: %s", err.Error())
			reply := tgbotapi.NewMessage(message.Chat.ID, err.Error())
			reply.ReplyToMessageID = message.MessageID
			bot.Send(reply)
		}
	}
}

func (session *chatSession) IsStale() bool {
	session.mutex.Lock()
	defer session.mutex.Unlock()
	// TODO: Make the session lifetime configurable.
	return time.Now().Sub(session.lastAccess) > time.Minute*20
}
