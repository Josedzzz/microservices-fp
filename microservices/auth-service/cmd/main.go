package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"

	_ "github.com/lib/pq"

	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/messaging"
	"auth-service/internal/repository"
	"auth-service/internal/security"
	"auth-service/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// -----------------------------
	// Database connection
	// -----------------------------
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("database not reachable:", err)
	}

	log.Println("Connected to database")

	// -----------------------------
	// Ensure schema
	// -----------------------------
	err = repository.EnsureSchema(db)
	if err != nil {
		log.Fatal("database initialization failed:", err)
	}

	// -----------------------------
	// RabbitMQ connection
	// -----------------------------
	conn, err := amqp091.Dial(cfg.RabbitMQUrl)
	if err != nil {
		log.Fatal("failed to connect to RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to RabbitMQ")

	// -----------------------------
	// Dependency injection
	// -----------------------------
	repo := repository.NewUserRepository(db)

	jwtManager := security.NewJWTManager(cfg.JWTSecret)

	authService := service.NewAuthService(repo, jwtManager)

	authHandler := handler.NewAuthHandler(authService)

	// -----------------------------
	// RabbitMQ consumer
	// -----------------------------
	consumer := messaging.NewConsumer(ch, authService)

	err = consumer.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("RabbitMQ consumer started")

	// -----------------------------
	// HTTP server
	// -----------------------------
	router := gin.Default()

	auth := router.Group("/auth-service/api")
	{
		auth.POST("/login", authHandler.Login)
	}

	log.Println("Auth service running on port 8083")

	err = router.Run(":8083")
	if err != nil {
		log.Fatal(err)
	}
}
