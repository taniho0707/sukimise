package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log"
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
		// APIエンドポイント（ログイン、リフレッシュ）はCSRF保護をスキップ
		path := c.Request.URL.Path
		method := c.Request.Method
		userAgent := c.GetHeader("User-Agent")
		
		log.Printf("DEBUG CSRF: Processing %s %s (User-Agent: %s)", method, path, userAgent)
		
		if path == "/api/v1/auth/login" || path == "/api/v1/auth/refresh" {
			// ログイン/リフレッシュエンドポイントはCSRF保護をスキップ
			log.Printf("DEBUG CSRF: Skipping CSRF protection for auth endpoint: %s", path)
			c.Next()
			return
		}

		// JWT認証されたAPIリクエストはCSRF保護をスキップ
		if isJWTAuthenticated(c) {
			log.Printf("DEBUG CSRF: Skipping CSRF protection for JWT authenticated request to %s", path)
			c.Next()
			return
		}
		
		// CSRFスキップフラグがある場合はスキップ
		if skipCSRF, exists := c.Get("skip_csrf"); exists && skipCSRF.(bool) {
			c.Next()
			return
		}
		
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
		log.Printf("DEBUG CSRF: Applying CSRF protection to %s %s", c.Request.Method, path)
		if !validateCSRFToken(c, config) {
			log.Printf("DEBUG CSRF: CSRF token validation failed for %s %s", c.Request.Method, path)
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

// isJWTAuthenticated checks if the request has a valid JWT Authorization header
func isJWTAuthenticated(c *gin.Context) bool {
	authHeader := c.GetHeader("Authorization")
	log.Printf("DEBUG JWT: Authorization header: '%s'", authHeader)
	
	if authHeader == "" {
		log.Printf("DEBUG JWT: No Authorization header found")
		return false
	}

	// Check if it's a Bearer token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		log.Printf("DEBUG JWT: Authorization header is not Bearer token (prefix check failed)")
		return false
	}

	// Extract the token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		log.Printf("DEBUG JWT: Empty token string after removing Bearer prefix")
		return false
	}

	// Basic validation - if it's a JWT format (has dots)
	// More thorough validation would require parsing the JWT, but for CSRF skipping,
	// just checking the presence of a Bearer token is sufficient since the actual
	// JWT validation happens in the Auth middleware later
	parts := strings.Split(tokenString, ".")
	isValid := len(parts) == 3 // JWT has 3 parts separated by dots
	log.Printf("DEBUG JWT: Token validation - token length: %d, parts: %d, isValid: %v", len(tokenString), len(parts), isValid)
	
	if isValid {
		log.Printf("DEBUG JWT: JWT token format is valid, skipping CSRF protection")
	} else {
		log.Printf("DEBUG JWT: JWT token format is invalid")
	}
	
	return isValid
}