package telegram

import (
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/proxy"
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"github.com/iley/lich/internal/config"
)

type Handler func(*Bot, *tgbotapi.Message) (done bool, nextHandler Handler, err error)

type Bot struct {
	config          *config.Config         // Effectively immutable.
	api             *tgbotapi.BotAPI       // Effectively immutable.
	commandHandlers map[string]Handler     // Effectively immutable.
	globalHandlers  []Handler              // Effectively immutable.
	userWhiltelist  map[string]struct{}    // Effectively immutable.
	chatSessions    map[int64]*chatSession // Protected by mutex.
	mutex           sync.Mutex
}

func NewBot(
	cfg *config.Config,
	commandHandlers map[string]Handler,
	globalHandlers []Handler) (*Bot, error) {
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
		api, err = tgbotapi.NewBotAPIWithClient(cfg.Token, httpClient)
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
	if len(cfg.Whitelist) == 0 {
		log.Println("Warning! No Telegram user whlitelist enforced.")
	} else {
		for _, username := range cfg.Whitelist {
			log.Println("Whitelisting user ", username)
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
	updatesChan, err := bot.api.GetUpdatesChan(updateConfig)
	if err != nil {
		return nil
	}
	go bot.RunGCLoop()
	for update := range updatesChan {
		if bot.UserAllowed(update.Message.From.UserName) {
			err = bot.EnqueueMessage(update.Message)
			if err != nil {
				log.Printf("Error enqueueing a message for chat %d: %s", update.Message.Chat.ID, err.Error())
			}
		} else {
			log.Printf("Unauthorized access attempt from user %s", update.Message.From.UserName)
			reply := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know you! Go away!")
			reply.ReplyToMessageID = update.Message.MessageID
			_, err = bot.api.Send(reply)
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
