// Package config provides configuration management from enviroment variables
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds configuration loaded from env
type Config struct {
	ServerPort string
	DBUrl      string
}

// Load gets the config from env variables
// Exits if DATABASE_URL is not set
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", "8081"),
		DBUrl:      getEnv("DATABASE_URL", ""),
	}

	if cfg.DBUrl == "" {
		log.Fatal("DATABASE_URL is required")
	}

	return cfg
}

// getEnv returns env variable value or default if not set
func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
