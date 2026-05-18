package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	jwtSecret   = getEnv("JWT_SECRET", "supersecretkey")
	gatewayPort = getEnv("GATEWAY_PORT", "8000")
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// ProxyConfig holds the configuration for a proxied service
type ProxyConfig struct {
	TargetURL   string
	Prefix      string
	PublicPaths []string
	Proxy       *httputil.ReverseProxy
}

func NewProxyConfig(targetURL, prefix string, publicPaths []string) *ProxyConfig {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Failed to parse target URL %s: %v", targetURL, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	
	// Customize the director to handle prefix stripping and headers
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Forwarded-Prefix", prefix)
		req.Host = target.Host
	}

	return &ProxyConfig{
		TargetURL:   targetURL,
		Prefix:      prefix,
		PublicPaths: publicPaths,
		Proxy:       proxy,
	}
}

func main() {
	router := gin.Default()
	router.RedirectTrailingSlash = false

	// Load HTML templates from the templates directory
	router.LoadHTMLGlob("templates/*")

	// Prometheus /metrics endpoint (register before other routes)
	router.GET("/metrics", gin.WrapF(promhttp.Handler().ServeHTTP))

	// Standard CORS Middleware
	router.Use(corsMiddleware())

	// Initialize Service Proxies
	services := []*ProxyConfig{
		NewProxyConfig(getEnv("AUTH_SERVICE_URL", "http://auth-service:8083"), "/auth-service", []string{"/api/login", "/api/recover-password", "/api/reset-password", "/swagger", "/api/swagger.yaml"}),
		NewProxyConfig(getEnv("EMPLOYEES_SERVICE_URL", "http://employees-service:8081"), "/employees-service", []string{"/swagger", "/api/swagger.yaml", "/api/health"}),
		NewProxyConfig(getEnv("DEPARTMENTS_SERVICE_URL", "http://departments-service:8082"), "/departments-service", []string{"/docs", "/openapi.json", "/api/"}),
		NewProxyConfig(getEnv("NOTIFICATIONS_SERVICE_URL", "http://notifications-service:8084"), "/notifications-service", []string{"/swagger-ui", "/v3/api-docs", "/swagger-ui.html"}),
		NewProxyConfig(getEnv("PROFILES_SERVICE_URL", "http://profiles-service:8085"), "/profiles-service", []string{"/swagger", "/swagger/", "/health"}),
	}

	// Register Routes
	for _, svc := range services {
		svc := svc // capture loop var
		router.Any(svc.Prefix+"/*path", authMiddleware(svc), rbacMiddleware(), func(c *gin.Context) {
			svc.Proxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	// Presentation Layer: Gateway Landing Page
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	log.Printf("API Gateway running on port %s", gatewayPort)
	router.Run(":" + gatewayPort)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-User-Email, X-User-Role")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func authMiddleware(svc *ProxyConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Check if path is public
		for _, p := range svc.PublicPaths {
			if strings.HasPrefix(path, svc.Prefix+p) {
				c.Set("isPublic", true)
				c.Next()
				return
			}
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth header format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("email", claims["sub"])
		c.Set("role", claims["role"])
		
		// Set headers for downstream services
		c.Request.Header.Set("X-User-Email", claims["sub"].(string))
		c.Request.Header.Set("X-User-Role", claims["role"].(string))
		
		c.Next()
	}
}

func rbacMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If we already aborted in authMiddleware, stop here
		if c.IsAborted() {
			return
		}

		// Skip RBAC for public paths
		if isPublic, _ := c.Get("isPublic"); isPublic == true {
			c.Next()
			return
		}

		role, _ := c.Get("role")

		// PDF Requirement (Reto 4): 
		// ADMIN: Total access.
		// USER: Read-only (GET).
		if c.Request.Method != http.MethodGet && role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "ADMIN role required for non-read operations"})
			c.Abort()
			return
		}

		c.Next()
	}
}
