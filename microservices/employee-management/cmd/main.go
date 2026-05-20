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
	"context"
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
	"employee-management/tracing"

	_ "employee-management/docs" // <-- Swagger docs

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
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

	// Initialize OpenTelemetry tracing with Zipkin exporter
	tp, err := tracing.InitializeTracer("employees-service")
	if err != nil {
		log.Printf("Warning: could not initialize tracing: %v", err)
	} else {
		defer tracing.Shutdown(context.Background(), tp)
		log.Println("Tracing initialized with Zipkin exporter")
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.SetTrustedProxies([]string{"127.0.0.1"})

	// Tracing middleware: create a span for every request
	router.Use(func(c *gin.Context) {
		tracer := otel.Tracer("employees-service")
		ctx, span := tracer.Start(c.Request.Context(), c.Request.Method+" "+c.FullPath())
		defer span.End()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

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

	// Prometheus /metrics endpoint (register before other routes)
	router.GET("/metrics", gin.WrapF(promhttp.Handler().ServeHTTP))

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
