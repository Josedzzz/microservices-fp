// Package middleware contains error handler for the middlewares
package middleware

import (
	"log"
	"net/http"

	"employee-management/internal/api"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Verify unhandled errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("unhandled error %v", err)

			api.Error(c, http.StatusInternalServerError, "Internal server error")

			c.Abort()
			return
		}
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				api.Error(c, http.StatusInternalServerError, "Internal server error")
				c.Abort()
			}
		}()

		c.Next()
	}
}
