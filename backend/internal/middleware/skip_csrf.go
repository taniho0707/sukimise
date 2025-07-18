package middleware

import "github.com/gin-gonic/gin"

// SkipCSRF returns a middleware that marks the request to skip CSRF protection
func SkipCSRF() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// CSRFスキップフラグを設定
		c.Set("skip_csrf", true)
		c.Next()
	})
}