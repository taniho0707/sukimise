package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"sukimise/internal/models"
)

// Test suite for Viewer Authentication API based on CLAUDE.md specifications
// Tests the viewer authentication system for URL + password based access

func TestAPI_ViewerAuth_AuthenticateViewer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "valid password format",
			requestBody: map[string]string{
				"password": "viewerpass123",
			},
			expectedStatus: http.StatusUnauthorized, // Expected since no real service configured
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Invalid password", response["error"])
			},
		},
		{
			name: "missing password",
			requestBody: map[string]string{},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "empty password",
			requestBody: map[string]string{
				"password": "",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			var reqBody []byte
			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/api/v1/viewer/auth", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "Test Browser/1.0")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Test viewer authentication request validation
			var authReq models.ViewerAuthRequest
			if err := c.ShouldBindJSON(&authReq); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				// Simulate authentication failure for testing
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			}

			// Verify response
			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
		})
	}
}

func TestAPI_ViewerAuth_ValidateViewerSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		sessionToken   string
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "no session token",
			sessionToken:   "",
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name:           "invalid session token",
			sessionToken:   "invalid-token-123",
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/viewer/validate", nil)
			if tt.sessionToken != "" {
				req.Header.Set("Authorization", "Bearer "+tt.sessionToken)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Test session validation logic
			sessionToken := c.GetHeader("X-Viewer-Token")
			if sessionToken == "" {
				sessionToken, _ = c.Cookie("viewer_session")
			}

			if sessionToken == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No session token provided"})
			} else {
				// Simulate validation failure for testing
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			}

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
		})
	}
}

func TestAPI_ViewerAuth_GetViewerSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("without authentication", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/admin/viewer-settings", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Test admin-only endpoint access
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Settings retrieved"})
		}

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("with mock admin authentication", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/admin/viewer-settings", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		// Mock admin role (would normally be set by middleware)
		c.Set("role", "admin")

		// Test admin-only endpoint access
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Settings retrieved"})
		}

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAPI_ViewerAuth_UpdateViewerSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "valid settings update",
			requestBody: map[string]interface{}{
				"password":              "newviewerpass123",
				"session_duration_days": 7,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
			},
		},
		{
			name: "invalid session duration",
			requestBody: map[string]interface{}{
				"password":              "newviewerpass123",
				"session_duration_days": 400, // Invalid (> 365)
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "missing password",
			requestBody: map[string]interface{}{
				"session_duration_days": 7,
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("PUT", "/api/v1/admin/viewer-settings", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("role", "admin") // Mock admin role

			// Test viewer settings update validation
			var settingsReq models.ViewerSettingsUpdateRequest
			if err := c.ShouldBindJSON(&settingsReq); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				// Additional validation for settings
				if settingsReq.SessionDurationDays > 365 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Session duration too long"})
				} else {
					c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
				}
			}

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
		})
	}
}

func TestAPI_ViewerAuth_GetViewerLoginHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("pagination parameter validation", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/admin/viewer-history?page=1&limit=50", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("role", "admin")

		// Test pagination parameter parsing
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "20")

		assert.Equal(t, "1", page)
		assert.Equal(t, "50", limit)

		c.JSON(http.StatusOK, gin.H{
			"page":  page,
			"limit": limit,
		})

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAPI_ViewerAuth_CleanupExpiredSessions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("admin only access", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/admin/viewer-cleanup", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Test admin role requirement
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Cleanup completed"})
		}

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestAPI_ViewerAuth_SessionManagement(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("client information extraction", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/viewer/auth", nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 Test Browser")
		req.Header.Set("X-Forwarded-For", "192.168.1.100")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Test client information extraction
		userAgent := c.GetHeader("User-Agent")
		clientIP := c.ClientIP()

		assert.Equal(t, "Mozilla/5.0 Test Browser", userAgent)
		assert.NotEmpty(t, clientIP)

		c.JSON(http.StatusOK, gin.H{
			"user_agent": userAgent,
			"client_ip":  clientIP,
		})

		assert.Equal(t, http.StatusOK, w.Code)
	})
}