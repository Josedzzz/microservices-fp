// Package config provides func for loading app configuration
package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBUrl       string
	RabbitMQUrl string
	JWTSecret   string
	ServerPort  string
}

func Load() *Config {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSL := os.Getenv("DB_SSLMODE")

	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitPort := os.Getenv("RABBITMQ_PORT")
	rabbitUser := os.Getenv("RABBITMQ_USER")
	rabbitPass := os.Getenv("RABBITMQ_PASSWORD")
	rabbitVHost := os.Getenv("RABBITMQ_VHOST")

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		dbSSL,
	)

	rabbitURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/%s",
		rabbitUser,
		rabbitPass,
		rabbitHost,
		rabbitPort,
		rabbitVHost,
	)

	return &Config{
		DBUrl:       dbURL,
		RabbitMQUrl: rabbitURL,
		JWTSecret:   os.Getenv("JWT_SECRET"),
		ServerPort:  os.Getenv("SERVER_PORT"),
	}
}
