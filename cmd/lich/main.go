package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/iley/lich/internal/config"
	"github.com/iley/lich/internal/handlers"
	"github.com/iley/lich/internal/telegram"
	"github.com/iley/lich/internal/torrents"
)

var (
	//go:embed VERSION
	versionFileContents string
)

func versionString() string {
	return strings.TrimSuffix(versionFileContents, "\n")
}

func main() {
	fmt.Printf("Lich v%s\n", versionString())

	configPath := flag.String("config", "config.json", "Configuration file path")
	flag.Parse()
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load config from %s: %s\n", *configPath, err)
		os.Exit(1)
	}

	// TODO: Move standard handlers into the telegram package.
	commandHandlers := map[string]telegram.Handler{
		"ping": handlers.MakePingHandler(),
	}
	commands := make([]string, 1, len(commandHandlers)+1)
	commands[0] = "/help"
	for command := range commandHandlers {
		commands = append(commands, "/"+command)
	}
	commandHandlers["help"] = handlers.MakeHelpHandler(commands)

	down, err := torrents.NewDownloader(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not create the torrent downloader:", err)
		os.Exit(1)
	}

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
