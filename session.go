package main

import (
	"errors"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strings"
	"sync"
	"time"
)

type ChatSession struct {
	input      chan *tgbotapi.Message
	lastAccess time.Time // Protected by mutex.
	mutex      sync.Mutex
}

func NewChatSession(bot *TelegramBot) *ChatSession {
	session := &ChatSession{
		input: make(chan *tgbotapi.Message, 32),
	}
	go session.RunLoop(bot)
	return session
}

func (session *ChatSession) EnqueueMessage(message *tgbotapi.Message) error {
	select {
	case session.input <- message:
		return nil
	default:
		return errors.New("Queue is full")
	}
}

func (session *ChatSession) Close() {
	close(session.input)
}

func (session *ChatSession) RunLoop(bot *TelegramBot) {
	var nextHandler TelegramHandler = nil
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
			handler, found := bot.CommandHandlers[command]
			if found {
				done, nextHandler, err = handler(bot, message)
			}
		}
		for _, globalHandler := range bot.GlobalHandlers {
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
			bot.Api.Send(reply)
		}
	}
}

func (session *ChatSession) IsStale() bool {
	session.mutex.Lock()
	defer session.mutex.Unlock()
	// TODO: Make the session lifetime configurable.
	return time.Now().Sub(session.lastAccess) > time.Minute*20
}
