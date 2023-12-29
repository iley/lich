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

const shutdownDelay = 5 * time.Second

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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Printf("Termination signal received, shutting down in %s...", shutdownDelay)
		cancel()

		time.Sleep(shutdownDelay)
		log.Println("Exiting")
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

	handlers := []telegram.HandlerDesc{
		{
			Scope:   telegram.HANDLER_GLOBAL,
			Handler: handlers.MakeTorrentFileHandler(),
		},
		{
			Scope:   telegram.HANDLER_GLOBAL,
			Handler: handlers.MakeMagnetLinkHandler(cfg, down),
		},
		{
			Scope:   telegram.HANDLER_COMMAND,
			Command: "ping",
			Handler: handlers.MakePingHandler(),
		},
		{
			Scope:   telegram.HANDLER_COMMAND,
			Command: "status",
			Handler: handlers.MakeStatusHandler(down),
		},
		{
			Scope:   telegram.HANDLER_COMMAND,
			Command: "help",
			Handler: handlers.MakeHelpHandler([]string{"/ping", "/status"}, versionString()),
		},
	}

	bot, err := telegram.NewBot(cfg, handlers)
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
