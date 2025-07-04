package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sukimise/internal/services"
)

func ViewerAuthMiddleware(service *services.ViewerAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken := c.GetHeader("X-Viewer-Token")
		if sessionToken == "" {
			sessionToken, _ = c.Cookie("viewer_session")
		}

		if sessionToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No session token provided"})
			c.Abort()
			return
		}

		valid, err := service.ValidateViewerSession(sessionToken)
		if err != nil || !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			c.Abort()
			return
		}

		c.Set("viewer_session", sessionToken)
		c.Next()
	}
}