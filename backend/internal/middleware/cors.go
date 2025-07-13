package middleware

import (
	"fmt"
	"strings"
	"sukimise/internal/config"

	"github.com/gin-gonic/gin"
)

// CORS returns a CORS middleware with secure default configuration
func CORS() gin.HandlerFunc {
	// 開発環境のデフォルト設定
	defaultConfig := config.CORSConfig{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token", "X-Viewer-Token"},
		AllowCredentials: true,
		MaxAge:           86400,
	}
	return NewCORSMiddleware(defaultConfig)
}

// NewCORSMiddleware creates a CORS middleware with the provided configuration
func NewCORSMiddleware(corsConfig config.CORSConfig) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// セキュリティ: 許可されたオリジンのみを許可
		var allowedOrigin string
		if isOriginAllowed(origin, corsConfig.AllowedOrigins) {
			allowedOrigin = origin
		} else if len(corsConfig.AllowedOrigins) > 0 {
			// 明示的に許可されたオリジンがある場合、最初のものを使用
			allowedOrigin = corsConfig.AllowedOrigins[0]
		} else {
			// フォールバック（本番環境では避けるべき）
			allowedOrigin = "http://localhost:3000"
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", boolToString(corsConfig.AllowCredentials))
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
		
		if corsConfig.MaxAge > 0 {
			c.Writer.Header().Set("Access-Control-Max-Age", intToString(corsConfig.MaxAge))
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// isOriginAllowed checks if the origin is in the allowed origins list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}
	
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}

// Helper functions
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func intToString(i int) string {
	return fmt.Sprintf("%d", i)
}