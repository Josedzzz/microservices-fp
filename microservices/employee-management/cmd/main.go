package main

//	@title			Employee Management API
//	@version		1.0
//	@description	API for managing employees
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.email	josed.amayar@uqvirtual.edu.co

//	@host		localhost:8081
//	@BasePath	/api

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer " followed by a valid JWT token.

import (
	"log"
	"net/http"

	"employee-management/internal/api"
	"employee-management/internal/config"
	"employee-management/internal/db"
	"employee-management/internal/handlers"
	"employee-management/internal/messaging"
	"employee-management/internal/middleware"
	"employee-management/internal/repository"
	"employee-management/internal/service"

	_ "employee-management/docs" // <-- Swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()

	dbPool := db.NewPostgresPool(cfg.DatabaseURL())
	defer dbPool.Close()

	publisher, err := messaging.NewPublisher(cfg.RabbitMQURL())
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer publisher.Close()

	repo := repository.NewEmployeeRepository(dbPool)
	service := service.NewEmployeeService(repo, publisher)
	handler := handlers.NewEmployeeHandler(service)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.Use(middleware.Recovery())
	router.Use(middleware.ErrorHandler())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.NoRoute(func(c *gin.Context) {
		api.NotFound(c, "Resource not found")
	})

	router.NoMethod(func(c *gin.Context) {
		api.Error(c, http.StatusMethodNotAllowed, "Method not allowed")
	})

	// Swagger (Prefix-blind)
	// Using a relative path allows the browser to find the YAML regardless of the Gateway prefix
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("../api/swagger.yaml")))
	router.StaticFile("/api/swagger.yaml", "./docs/swagger.yaml")

	// API Routes (Prefix-blind)
	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/health", handlers.HealthCheck)

		// Apply AuthMiddleware to all protected routes
		protected := apiGroup.Group("/", middleware.AuthMiddleware())
		{
			employees := protected.Group("/employees")
			{
				employees.POST("/", handler.CreateEmployee)
				employees.GET("/:id", handler.GetEmployeeByID)
				employees.GET("/", handler.GetAllEmployees)
				employees.PUT("/:id", handler.UpdateEmployee)
				employees.DELETE("/:id", handler.DeleteEmployee)
			}
		}
	}

	log.Printf("Employee service running on :%s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
