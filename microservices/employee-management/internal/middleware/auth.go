package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"employee-management/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates the JWT token and handles RBAC
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			api.Error(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			api.Error(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "supersecretkey" // Fallback for dev, but should be in env
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			api.Error(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			api.Error(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		role, _ := claims["role"].(string)
		
		// RBAC: Rol ADMIN has total access. Rol USER has read-only access.
		if c.Request.Method != http.MethodGet {
			if role != "ADMIN" {
				api.Error(c, http.StatusForbidden, "Forbidden: ADMIN role required for write operations")
				c.Abort()
				return
			}
		}

		// Store claims and raw token in context
		c.Set("user_email", claims["sub"])
		c.Set("user_role", role)
		c.Set("auth_token", tokenString)

		c.Next()
	}
}
