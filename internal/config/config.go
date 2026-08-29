package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Query struct {
	Query      string `yaml:"query"`
	MaxResults int64  `yaml:"max_results,omitempty"`
	Language   string `yaml:"-"` 
}

type SearchConfig struct {
	DaysLookback       int      `yaml:"days_lookback"`
	MaxResultsPerQuery int64    `yaml:"max_results_per_query"`
	EnabledLanguages   []string `yaml:"enabled_languages"`
}

type Config struct {
	Search           SearchConfig
	Queries          []Query
	YouTubeAPIKey    string
	TelegramBotToken string
	TelegramChatID   string
	DBPath           string
}

type rawConfig struct {
	Search SearchConfig `yaml:"search"`
}

func Load(configDir string) (*Config, error) {

	_ = godotenv.Load()

	if configDir == "" {
		configDir = "config"
	}

	cfg := &Config{
		DBPath:           "data/videos.db",
		YouTubeAPIKey:    os.Getenv("YOUTUBE_API_KEY"),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
	}

	mainConfigPath, err := findConfigFile(configDir, "config")
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(mainConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read main config: %w", err)
	}

	var raw rawConfig
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse main config: %w", err)
	}
	cfg.Search = raw.Search

	if cfg.Search.DaysLookback <= 0 {
		cfg.Search.DaysLookback = 14
	}
	if cfg.Search.MaxResultsPerQuery <= 0 {
		cfg.Search.MaxResultsPerQuery = 15
	}

	queriesDir := filepath.Join(configDir, "queries")
	entries, err := os.ReadDir(queriesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read queries directory: %w", err)
	}

	enabled := make(map[string]bool)
	for _, lang := range cfg.Search.EnabledLanguages {
		enabled[strings.ToLower(lang)] = true
	}
	loadAll := len(enabled) == 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}

		lang := strings.TrimSuffix(name, ".yaml")
		lang = strings.TrimSuffix(lang, ".yml")
		lang = strings.ToLower(lang)

		if !loadAll && !enabled[lang] {
			continue
		}

		filePath := filepath.Join(queriesDir, name)
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read query file %s: %w", name, err)
		}

		var fileQueries struct {
			Queries []Query `yaml:"queries"`
		}
		if err := yaml.Unmarshal(fileData, &fileQueries); err != nil {
			return nil, fmt.Errorf("failed to parse query file %s: %w", name, err)
		}

		for _, q := range fileQueries.Queries {
			if strings.TrimSpace(q.Query) == "" {
				continue
			}
			q.Language = lang
			if q.MaxResults <= 0 {
				q.MaxResults = cfg.Search.MaxResultsPerQuery
			}
			cfg.Queries = append(cfg.Queries, q)
		}
	}

	if cfg.YouTubeAPIKey == "" {
		return nil, fmt.Errorf("YOUTUBE_API_KEY is required")
	}
	if len(cfg.Queries) == 0 {
		return nil, fmt.Errorf("no search queries loaded")
	}

	return cfg, nil
}

func findConfigFile(dir, name string) (string, error) {
	candidates := []string{
		filepath.Join(dir, name+".yaml"),
		filepath.Join(dir, name+".yml"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("config file not found (tried .yaml and .yml)")
}