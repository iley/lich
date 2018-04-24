package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"sync"
	"time"
)

type TelegramBot struct {
	Config          *TelegramConfig            // Effectively immutable.
	Api             *tgbotapi.BotAPI           // Effectively immutable.
	CommandHandlers map[string]TelegramHandler // Effectively immutable.
	GlobalHandlers  []TelegramHandler          // Effectively immutable.
	userWhiltelist  map[string]struct{}        // Effectively immutable.
	chatSessions    map[int64]*ChatSession     // Protected by mutex.
	mutex           sync.Mutex
}

type TelegramHandler func(*TelegramBot, *tgbotapi.Message) (done bool, nextHandler TelegramHandler, err error)

func NewTelegramBot(
	config *TelegramConfig,
	commandHandlers map[string]TelegramHandler,
	globalHandlers []TelegramHandler) (*TelegramBot, error) {
	botApi, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}
	bot := TelegramBot{
		Config:          config,
		Api:             botApi,
		CommandHandlers: commandHandlers,
		GlobalHandlers:  globalHandlers,
		userWhiltelist:  make(map[string]struct{}),
		chatSessions:    make(map[int64]*ChatSession),
	}
	if len(config.Whitelist) == 0 {
		log.Println("Warning! No Telegram user whlitelist enforced.")
	} else {
		for _, username := range config.Whitelist {
			log.Println("Whitelisting user ", username)
			bot.userWhiltelist[username] = struct{}{}
		}
	}
	return &bot, nil
}

func (bot *TelegramBot) UserAllowed(username string) bool {
	if len(bot.userWhiltelist) == 0 {
		return true
	}
	_, ok := bot.userWhiltelist[username]
	return ok
}

func (bot *TelegramBot) RunLoop() error {
	log.Println("Running the Telegram bot")
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 120
	updatesChan, err := bot.Api.GetUpdatesChan(updateConfig)
	if err != nil {
		return nil
	}
	go bot.RunGCLoop()
	for update := range updatesChan {
		if bot.UserAllowed(update.Message.From.UserName) {
			bot.EnqueueMessage(update.Message)
			if err != nil {
				log.Printf("Error enqueueing a message for chat %d: %s", update.Message.Chat.ID, err.Error())
			}
		} else {
			log.Printf("Unauthorized access attempt from user %s", update.Message.From.UserName)
			reply := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know you! Go away!")
			reply.ReplyToMessageID = update.Message.MessageID
			bot.Api.Send(reply)
		}
	}
	return nil
}

func (bot *TelegramBot) RunGCLoop() {
	for {
		// TODO: Make the GC timeout configurable.
		time.Sleep(time.Minute * 5)
		bot.mutex.Lock()
		for chatId, session := range bot.chatSessions {
			if session.IsStale() {
				log.Printf("Session for chat %d is stale. Closing it", chatId)
				delete(bot.chatSessions, chatId)
				session.Close()
			}
		}
		bot.mutex.Unlock()
	}
}

func (bot *TelegramBot) EnqueueMessage(message *tgbotapi.Message) error {
	bot.mutex.Lock()
	defer bot.mutex.Unlock()
	session, ok := bot.chatSessions[message.Chat.ID]
	if ok {
		log.Printf("Reusing session for chat %d", message.Chat.ID)
	} else {
		log.Printf("Creating new session for chat %d", message.Chat.ID)
		session = NewChatSession(bot)
		bot.chatSessions[message.Chat.ID] = session
	}
	return session.EnqueueMessage(message)
}
