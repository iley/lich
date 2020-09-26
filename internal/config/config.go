package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Token      string            `json:"token"`
	Whitelist  []string          `json:"whitelist"`
	Proxy      *ProxyConfig      `json:"proxy,omitempty"`
	WorkDir    string            `json:"work_dir"`
	TargetDirs map[string]string `json:"target_dirs"`
}

type ProxyConfig struct {
	Address  string `json:"address"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

const UnsortedCategory = "unsorted"

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var cfg Config
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	err = validateConfig(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	if cfg.WorkDir == "" {
		return errors.New("Missing required option 'work_directory'")
	}
	err := os.MkdirAll(cfg.WorkDir, 0755)
	if err != nil {
		msg := fmt.Sprintf("Could not create work directory %s: %s",
			cfg.WorkDir, err.Error())
		return errors.New(msg)
	}
	hasUnsortedCategory := false
	for category, targetDir := range cfg.TargetDirs {
		if category == UnsortedCategory {
			hasUnsortedCategory = true
		}
		err = os.MkdirAll(targetDir, 0755)
		if err != nil {
			msg := fmt.Sprintf("Could not create target directory %s for category %s: %s",
				targetDir, category, err.Error())
			return errors.New(msg)
		}
	}
	if !hasUnsortedCategory {
		return fmt.Errorf("Required category '%s' not found", UnsortedCategory)
	}
	return nil
}
