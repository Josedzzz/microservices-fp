package main

// @title Employee Management API
// @version 1.0
// @description API for managing employees
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@company.com

// @host localhost:8081
// @BasePath /

import (
	"log"
	"net/http"

	"employee-management/internal/api"
	"employee-management/internal/config"
	"employee-management/internal/db"
	"employee-management/internal/handlers"
	"employee-management/internal/middleware"
	"employee-management/internal/repository"
	"employee-management/internal/service"

	_ "employee-management/docs" // <-- Swagger docs (IMPORTANT)

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()

	dbPool := db.NewPostgresPool(cfg.DatabaseURL())
	defer dbPool.Close()

	repo := repository.NewEmployeeRepository(dbPool)
	service := service.NewEmployeeService(repo)
	handler := handlers.NewEmployeeHandler(service)

	// Gin config
	gin.SetMode(gin.ReleaseMode) // Change mode for development
	router := gin.New()

	// Trusted proxies
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// Middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.ErrorHandler())
	router.Use(gin.Logger())
	router.Use(gin.Recovery()) // Recovery fallback

	// Global handlers
	router.NoRoute(func(c *gin.Context) {
		api.NotFound(c, "Resource not found")
	})

	router.NoMethod(func(c *gin.Context) {
		api.Error(c, http.StatusMethodNotAllowed, "Method not allowed")
	})

	apiGroup := router.Group("/employees-service/api")
	{
		// Health
		apiGroup.GET("/health", handlers.HealthCheck)

		// Swagger
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// Employee routes
		employees := apiGroup.Group("/employees")
		{
			employees.POST("/", handler.CreateEmployee)
			employees.GET("/:id", handler.GetEmployeeByID)
			employees.GET("/", handler.GetAllEmployees)
			employees.PUT("/:id", handler.UpdateEmployee)
			employees.DELETE("/:id", handler.DeleteEmployee)
		}
	}

	log.Printf("Employee service running on :%s", cfg.ServerPort)
	log.Printf("Swagger UI available at http://localhost:%s/swagger/index.html", cfg.ServerPort)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
