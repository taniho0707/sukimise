package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	csrfTokenLength = 32
	csrfHeader      = "X-CSRF-Token"
	csrfFormField   = "_csrf_token"
	csrfCookie      = "csrf_token"
)

// CSRFConfig holds CSRF protection configuration
type CSRFConfig struct {
	Secret    string
	TokenName string
	Header    string
	Cookie    string
	Secure    bool
	SameSite  http.SameSite
}

// CSRF returns a CSRF protection middleware
func CSRF(config CSRFConfig) gin.HandlerFunc {
	if config.TokenName == "" {
		config.TokenName = csrfFormField
	}
	if config.Header == "" {
		config.Header = csrfHeader
	}
	if config.Cookie == "" {
		config.Cookie = csrfCookie
	}

	return gin.HandlerFunc(func(c *gin.Context) {
		// GETリクエストとOPTIONSリクエストはCSRF保護をスキップ
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			// CSRFトークンを生成してクッキーに設定
			token := generateCSRFToken()
			setCSRFCookie(c, token, config)
			c.Header(config.Header, token) // レスポンスヘッダーにも設定
			c.Next()
			return
		}

		// POST, PUT, DELETE, PATCHリクエストはCSRF保護を適用
		if !validateCSRFToken(c, config) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "CSRF token validation failed",
				"message": "Invalid or missing CSRF token",
			})
			c.Abort()
			return
		}

		c.Next()
	})
}

// generateCSRFToken generates a cryptographically secure CSRF token
func generateCSRFToken() string {
	bytes := make([]byte, csrfTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		// フォールバック: 時間ベースのトークン（本来は避けるべき）
		return hex.EncodeToString([]byte("fallback_token"))
	}
	return hex.EncodeToString(bytes)
}

// setCSRFCookie sets the CSRF token in a cookie
func setCSRFCookie(c *gin.Context, token string, config CSRFConfig) {
	c.SetSameSite(config.SameSite)
	c.SetCookie(
		config.Cookie,
		token,
		3600, // 1時間
		"/",
		"",
		config.Secure,
		false, // HttpOnly=false (JavaScriptからアクセス可能にする必要がある)
	)
}

// validateCSRFToken validates the CSRF token from header or form
func validateCSRFToken(c *gin.Context, config CSRFConfig) bool {
	// クッキーからトークンを取得
	cookieToken, err := c.Cookie(config.Cookie)
	if err != nil || cookieToken == "" {
		return false
	}

	// ヘッダーまたはフォームフィールドからトークンを取得
	var requestToken string
	
	// まずヘッダーを確認
	requestToken = c.GetHeader(config.Header)
	
	// ヘッダーになければフォームフィールドを確認
	if requestToken == "" {
		requestToken = c.PostForm(config.TokenName)
	}

	// JSONリクエストの場合、リクエストボディからも確認
	if requestToken == "" {
		if jsonToken, exists := c.Get("csrf_token"); exists {
			if tokenStr, ok := jsonToken.(string); ok {
				requestToken = tokenStr
			}
		}
	}

	// トークンが空でないかチェック
	if requestToken == "" {
		return false
	}

	// 定数時間比較でトークンを検証
	return secureCompare(cookieToken, requestToken)
}

// secureCompare performs constant-time comparison of two strings
func secureCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	
	return result == 0
}

// GetCSRFToken extracts CSRF token from various sources for API endpoints
func GetCSRFToken(c *gin.Context) string {
	// ヘッダーから取得を試行
	if token := c.GetHeader(csrfHeader); token != "" {
		return token
	}
	
	// Authorizationヘッダーの代替として
	if auth := c.GetHeader("Authorization"); auth != "" {
		if strings.HasPrefix(auth, "CSRF ") {
			return strings.TrimPrefix(auth, "CSRF ")
		}
	}
	
	// クッキーから取得
	if token, err := c.Cookie(csrfCookie); err == nil {
		return token
	}
	
	return ""
}