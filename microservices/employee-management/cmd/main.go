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

	"employee-management/internal/config"
	"employee-management/internal/db"
	"employee-management/internal/handlers"
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

	router := gin.New()

	// Trusted proxies
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Global handlers
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "resource not found",
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "method not allowed",
		})
	})

	// Health
	router.GET("employees-service/api/health", handlers.HealthCheck)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Employee routes
	router.POST("employees-service/api/employees", handler.CreateEmployee)
	router.GET("employees-service/api/employees/:id", handler.GetEmployeeByID)
	router.GET("employees-service/api/employees", handler.GetAllEmployees)
	router.PUT("employees-service/api/employees/:id", handler.UpdateEmployee)
	router.DELETE("employees-service/api/employees/:id", handler.DeleteEmployee)

	log.Printf("Employee service running on :%s", cfg.ServerPort)
	router.Run(":" + cfg.ServerPort)
}
