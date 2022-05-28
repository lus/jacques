package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config represents the application configuration
type Config struct {
	BotToken string `split_words:"true"`
}

// LoadFromEnv loads the configuration using environment variables or/and a .env file
func LoadFromEnv() (*Config, error) {
	_ = godotenv.Overload()

	config := new(Config)
	if err := envconfig.Process("jcq", config); err != nil {
		return nil, err
	}
	return config, nil
}
