package telegram

import (
	"log"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/proxy"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/iley/lich/internal/config"
)

// Different types of handlers supported.
const (
	// Global handlers are called for every message.
	HANDLER_GLOBAL = iota
	// Command handlers are called for messages that start with an exact command.
	HANDLER_COMMAND = iota
	// Wildcard handlers are the same as command handlers except that the command can include an arbitrary suffix.
	HANDLER_WILDCARD_COMMAND = iota
)

type Handler func(*Bot, *tgbotapi.Message) (done bool, nextHandler Handler, err error)

type HandlerDesc struct {
	Handler Handler
	Command string
	Scope   int
}

type Bot struct {
	config          *config.Config         // Effectively immutable.
	api             *tgbotapi.BotAPI       // Effectively immutable.
	commandHandlers map[string]Handler     // Effectively immutable.
	globalHandlers  []Handler              // Effectively immutable.
	userWhiltelist  map[string]struct{}    // Effectively immutable.
	chatSessions    map[int64]*chatSession // Protected by mutex.
	mutex           sync.Mutex
}

func NewBot(cfg *config.Config, handlers []HandlerDesc) (*Bot, error) {
	globalHandlers := make([]Handler, 0)
	commandHandlers := make(map[string]Handler)
	for _, handlerDesc := range handlers {
		switch handlerDesc.Scope {
		case HANDLER_GLOBAL:
			globalHandlers = append(globalHandlers, handlerDesc.Handler)
		case HANDLER_COMMAND:
			if handlerDesc.Command == "" {
				return nil, fmt.Errorf("empty command for command handler")
			}
			commandHandlers[handlerDesc.Command] = handlerDesc.Handler
		case HANDLER_WILDCARD_COMMAND:
			// TODO: Implement wildcard command handlers.
			commandHandlers[handlerDesc.Command] = handlerDesc.Handler
		default:
			return nil, fmt.Errorf("invalid handler scope %d", handlerDesc.Scope)
		}
	}

	var api *tgbotapi.BotAPI
	var err error
	if cfg.Proxy == nil {
		api, err = tgbotapi.NewBotAPI(cfg.Token)
		if err != nil {
			return nil, err
		}
	} else {
		httpClient, err := proxyHTTPClient(cfg.Proxy.Address, cfg.Proxy.Username, cfg.Proxy.Password)
		if err != nil {
			return nil, err
		}
		api, err = tgbotapi.NewBotAPIWithClient(cfg.Token, tgbotapi.APIEndpoint, httpClient)
		if err != nil {
			return nil, err
		}
	}

	bot := Bot{
		config:          cfg,
		api:             api,
		commandHandlers: commandHandlers,
		globalHandlers:  globalHandlers,
		userWhiltelist:  make(map[string]struct{}),
		chatSessions:    make(map[int64]*chatSession),
	}
	if len(cfg.UsersAllowlist) == 0 {
		log.Println("Warning! No Telegram user whlitelist enforced.")
	} else {
		for _, username := range cfg.UsersAllowlist {
			log.Println("Allowlisting user ", username)
			bot.userWhiltelist[username] = struct{}{}
		}
	}
	return &bot, nil
}

func (bot *Bot) UserAllowed(username string) bool {
	if len(bot.userWhiltelist) == 0 {
		return true
	}
	_, ok := bot.userWhiltelist[username]
	return ok
}

func (bot *Bot) RunLoop() error {
	log.Println("Running the Telegram bot")
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 120
	updatesChan := bot.api.GetUpdatesChan(updateConfig)
	go bot.RunGCLoop()
	for update := range updatesChan {
		if bot.UserAllowed(update.Message.From.UserName) {
			err := bot.EnqueueMessage(update.Message)
			if err != nil {
				log.Printf("Error enqueueing a message for chat %d: %s", update.Message.Chat.ID, err.Error())
			}
		} else {
			log.Printf("Unauthorized access attempt from user %s", update.Message.From.UserName)
			reply := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know you! Go away!")
			reply.ReplyToMessageID = update.Message.MessageID
			_, err := bot.api.Send(reply)
			if err != nil {
				log.Printf("error sending message: %v", err)
			}
		}
	}
	return nil
}

func (bot *Bot) Send(c tgbotapi.Chattable) {
	_, err := bot.api.Send(c)
	if err != nil {
		log.Printf("error sending message: %v", err)
	}
}

func (bot *Bot) RunGCLoop() {
	for {
		// TODO: Make the GC timeout configurable.
		time.Sleep(time.Minute * 5)
		bot.mutex.Lock()
		for chatID, session := range bot.chatSessions {
			if session.IsStale() {
				log.Printf("Session for chat %d is stale. Closing it", chatID)
				delete(bot.chatSessions, chatID)
				session.Close()
			}
		}
		bot.mutex.Unlock()
	}
}

func (bot *Bot) EnqueueMessage(message *tgbotapi.Message) error {
	bot.mutex.Lock()
	defer bot.mutex.Unlock()
	session, ok := bot.chatSessions[message.Chat.ID]
	if ok {
		log.Printf("Reusing session for chat %d", message.Chat.ID)
	} else {
		log.Printf("Creating new session for chat %d", message.Chat.ID)
		session = newChatSession(bot)
		bot.chatSessions[message.Chat.ID] = session
	}
	return session.EnqueueMessage(message)
}

func proxyHTTPClient(addr, username, password string) (*http.Client, error) {
	var auth *proxy.Auth = nil
	if username != "" || password != "" {
		auth = &proxy.Auth{
			User:     username,
			Password: password,
		}
	}
	dialer, err := proxy.SOCKS5("tcp", addr, auth, proxy.Direct)
	if err != nil {
		return nil, err
	}
	httpTransport := &http.Transport{Dial: dialer.Dial}
	httpClient := &http.Client{Transport: httpTransport}
	return httpClient, nil
}
