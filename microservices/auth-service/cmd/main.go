package main

//	@title			Auth Service API
//	@version		1.0
//	@description	API for authentication and user management

//	@contact.name	API Support
//	@contact.email	josed.amayar@uqvirtual.edu.co

//	@host		localhost:8083
//	@BasePath	/api

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/lib/pq"

	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/messaging"
	"auth-service/internal/repository"
	"auth-service/internal/security"
	"auth-service/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("database not reachable:", err)
	}

	log.Println("Connected to database")

	err = repository.EnsureSchema(db)
	if err != nil {
		log.Fatal("database initialization failed:", err)
	}

	conn, err := amqp091.Dial(cfg.RabbitMQUrl)
	if err != nil {
		log.Fatal("failed to connect to RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to RabbitMQ")

	repo := repository.NewUserRepository(db)
	jwtManager := security.NewJWTManager(cfg.JWTSecret)
	producer := messaging.NewProducer(ch)
	authService := service.NewAuthService(repo, jwtManager, producer)
	authHandler := handler.NewAuthHandler(authService)
	consumer := messaging.NewConsumer(ch, authService)

	err = consumer.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("RabbitMQ consumer started")

	router := gin.Default()

	// Prometheus /metrics endpoint (register before other routes)
	router.GET("/metrics", gin.WrapF(promhttp.Handler().ServeHTTP))

	// Swagger documentation
	// Using a relative path makes it portable and avoids Gateway prefix issues
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("../api/swagger.yaml")))
	router.StaticFile("/api/swagger.yaml", "./docs/swagger.yaml")

	// API Routes (Prefix-blind)
	auth := router.Group("/api")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/recover-password", authHandler.RecoverPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	log.Println("Auth service running on port 8083")

	err = router.Run(":8083")
	if err != nil {
		log.Fatal(err)
	}
}
