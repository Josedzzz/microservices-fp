package main

import (
	"log"
	"net/http"

	"employee-management/internal/config"
	"employee-management/internal/db"
	"employee-management/internal/handlers"
	"employee-management/internal/repository"
	"employee-management/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	dbPool := db.NewPostgresPool(cfg.DBUrl)
	defer dbPool.Close()

	repo := repository.NewEmployeeRepository(dbPool)
	service := service.NewEmployeeService(repo)
	handler := handlers.NewEmployeeHandler(service)

	router := gin.New()
	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Global handlers for unsupported routes/methods (Challenge 1)
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

	// Routes
	router.GET("/health", handlers.HealthCheck)
	router.POST("/employees", handler.CreateEmployee)
	router.GET("/employees/:id", handler.GetEmployeeByID)

	log.Printf("Employee service running on :%s", cfg.ServerPort)
	router.Run(":" + cfg.ServerPort)
}
