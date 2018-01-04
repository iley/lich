package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Telegram TelegramConfig `json:"telegram"`
	Torrent  TorrentConfig  `json:"torrent"`
}

type TelegramConfig struct {
	Token     string   `json:"token"`
	Whitelist []string `json:"whitelist"`
}

const UnsortedCategory = "unsorted"

type TorrentConfig struct {
	WorkDir    string            `json:"work_dir"`
	TargetDirs map[string]string `json:"target_dirs"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
