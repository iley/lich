package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	configPath := flag.String("config", "config.json", "Configuration file path")
	flag.Parse()
	config, err := LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load config from %s: %s\n", *configPath, err)
		os.Exit(1)
	}

	commandHandlers := map[string]TelegramHandler{
		"ping":    HandlePing,
		"version": HandleVersion,
	}
	commands := make([]string, 1, len(commandHandlers)+1)
	commands[0] = "/help"
	for command, _ := range commandHandlers {
		commands = append(commands, "/"+command)
	}
	commandHandlers["help"] = MakeHelpHandler(commands)

	downloader, err := NewTorrentDownloader(&config.Torrent)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not create the torrent downloader:", err)
		os.Exit(1)
	}
	globalHandlers := []TelegramHandler{
		downloader.HandleMagnetLink,
		downloader.HandleTorrentFile,
	}

	bot, err := NewTelegramBot(&config.Telegram, commandHandlers, globalHandlers)
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
