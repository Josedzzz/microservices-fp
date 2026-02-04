package main

import (
	"log"

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

	router := gin.Default()

	router.GET("/health", handlers.HealthCheck)
	router.POST("/employees", handler.CreateEmployee)

	log.Printf("Employee service running on :%s", cfg.ServerPort)
	router.Run(":" + cfg.ServerPort)
}
