package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Check if request ID is already provided in headers
		requestID := c.GetHeader("X-Request-ID")
		
		// Generate a new UUID if not provided
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// Set the request ID in the context and response header
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		
		c.Next()
	})
}