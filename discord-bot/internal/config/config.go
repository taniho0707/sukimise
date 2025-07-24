package config

import (
	"errors"
	"os"
)

type Config struct {
	DiscordToken       string
	DatabaseURL        string
	SukimiseAPIURL     string
	SukimiseFrontendURL string
	BotPort            string
}

func Load() (*Config, error) {
	config := &Config{
		DiscordToken:        os.Getenv("DISCORD_TOKEN"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		SukimiseAPIURL:      os.Getenv("SUKIMISE_API_URL"),
		SukimiseFrontendURL: os.Getenv("VITE_API_BASE_URL"),
		BotPort:             os.Getenv("BOT_PORT"),
	}

	// Set default values
	if config.SukimiseAPIURL == "" {
		config.SukimiseAPIURL = "http://backend:8081"
	}
	if config.SukimiseFrontendURL == "" {
		config.SukimiseFrontendURL = "http://localhost"
	}
	if config.BotPort == "" {
		config.BotPort = "8082"
	}

	// Validate required fields
	if config.DiscordToken == "" {
		return nil, errors.New("DISCORD_TOKEN is required")
	}
	if config.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	return config, nil
}