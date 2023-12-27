package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/iley/lich/internal/config"
	"github.com/iley/lich/internal/handlers"
	"github.com/iley/lich/internal/telegram"
	"github.com/iley/lich/internal/torrents"
)

var (
	//go:embed version.txt
	versionFileContents string
)

func versionString() string {
	return strings.TrimSuffix(versionFileContents, "\n")
}

func main() {
	log.Printf("Lich v%s\n", versionString())

	ctx, cancel := context.WithCancel(context.Background())

	// handle SIGINT and SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Signal received, cancelling context.")
		cancel()

		time.Sleep(5 * time.Second)
		log.Println("Exiting after 5 seconds.")
		os.Exit(0)
	}()

	configPath := flag.String("config", "config.json", "Configuration file path")
	flag.Parse()
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Could not load config from %s: %s\n", *configPath, err)
	}

	down, err := torrents.NewDownloader(ctx, cfg)
	if err != nil {
		log.Fatalf("Could not create the torrent downloader: %s", err)
		os.Exit(1)
	}

	commandHandlers := map[string]telegram.Handler{
		"ping":   handlers.MakePingHandler(),
		"status": handlers.MakeStatusHandler(down),
	}
	commands := make([]string, 1, len(commandHandlers)+1)
	commands[0] = "/help"
	for command := range commandHandlers {
		commands = append(commands, "/"+command)
	}
	commandHandlers["help"] = handlers.MakeHelpHandler(commands, versionString())

	globalHandlers := []telegram.Handler{
		handlers.MakeTorrentFileHandler(),
		handlers.MakeMagnetLinkHandler(cfg, down),
	}

	bot, err := telegram.NewBot(cfg, commandHandlers, globalHandlers)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not create the Telegram bot:", err)
		os.Exit(1)
	}

	err = bot.RunLoop()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while running the Telegram bot:", err)
		os.Exit(1)
	}
}
