package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs incoming requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if query != "" {
			path = path + "?" + query
		}

		log.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s %s",
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			errorMessage,
		)
	}
}

// ErrorRecovery recovers from panics and returns 500 error
func ErrorRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("Panic recovered: %s", err)
			c.JSON(500, gin.H{
				"success": false,
				"error": map[string]interface{}{
					"code":    "INTERNAL_ERROR",
					"message": "Internal server error",
					"details": map[string]string{
						"recovered": "panic recovered",
					},
				},
			})
		} else if err, ok := recovered.(error); ok {
			log.Printf("Panic recovered: %v", err)
			c.JSON(500, gin.H{
				"success": false,
				"error": map[string]interface{}{
					"code":    "INTERNAL_ERROR",
					"message": "Internal server error",
					"details": map[string]string{
						"recovered": "panic recovered",
					},
				},
			})
		} else {
			log.Printf("Panic recovered: %v", recovered)
			c.JSON(500, gin.H{
				"success": false,
				"error": map[string]interface{}{
					"code":    "INTERNAL_ERROR",
					"message": "Internal server error",
				},
			})
		}
		c.Abort()
	})
}

// CORS handles Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// JSONContentType sets Content-Type to application/json
func JSONContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Next()
	}
}

// SecurityHeaders adds basic security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}
