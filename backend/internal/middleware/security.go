package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds common security headers to responses
func SecurityHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Prevent page from being displayed in a frame
		c.Header("X-Frame-Options", "DENY")
		
		// Strict transport security (HTTPS only)
		// Only set this if you're using HTTPS in production
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// Content Security Policy - basic policy
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Feature policy to restrict dangerous features
		c.Header("Permissions-Policy", "geolocation=(self), microphone=(), camera=()")
		
		c.Next()
	})
}

// RateLimitHeaders adds rate limiting headers (for future use)
func RateLimitHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// These headers would be set by an actual rate limiter
		// c.Header(\"X-RateLimit-Limit\", \"100\")
		// c.Header(\"X-RateLimit-Remaining\", \"99\")
		// c.Header(\"X-RateLimit-Reset\", strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10))
		
		c.Next()
	})
}