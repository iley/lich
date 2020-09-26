package config

import "sort"

func (cfg *Config) Categories() []string {
	categories := make([]string, 0, len(cfg.TargetDirs))
	for category := range cfg.TargetDirs {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}
