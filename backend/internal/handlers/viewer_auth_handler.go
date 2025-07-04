package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"sukimise/internal/models"
	"sukimise/internal/services"
)

type ViewerAuthHandler struct {
	service *services.ViewerAuthService
}

func NewViewerAuthHandler(service *services.ViewerAuthService) *ViewerAuthHandler {
	return &ViewerAuthHandler{service: service}
}

func (h *ViewerAuthHandler) AuthenticateViewer(c *gin.Context) {
	var req models.ViewerAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	response, err := h.service.AuthenticateViewer(req.Password, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Set cookie for session persistence
	c.SetCookie(
		"viewer_session",
		response.Token,
		int(response.ExpiresAt.Sub(time.Now()).Seconds()),
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, response)
}

func (h *ViewerAuthHandler) ValidateViewerSession(c *gin.Context) {
	sessionToken := c.GetHeader("X-Viewer-Token")
	if sessionToken == "" {
		sessionToken, _ = c.Cookie("viewer_session")
	}

	if sessionToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No session token provided"})
		return
	}

	valid, err := h.service.ValidateViewerSession(sessionToken)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
}

func (h *ViewerAuthHandler) GetViewerSettings(c *gin.Context) {
	settings, err := h.service.GetViewerSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get settings"})
		return
	}

	// Don't expose password hash
	response := gin.H{
		"id":                     settings.ID,
		"session_duration_days":  settings.SessionDurationDays,
		"created_at":             settings.CreatedAt,
		"updated_at":             settings.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ViewerAuthHandler) UpdateViewerSettings(c *gin.Context) {
	var req models.ViewerSettingsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateViewerSettings(req.Password, req.SessionDurationDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

func (h *ViewerAuthHandler) GetViewerLoginHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	history, total, err := h.service.GetViewerLoginHistory(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get login history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history": history,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func (h *ViewerAuthHandler) CleanupExpiredSessions(c *gin.Context) {
	err := h.service.CleanupExpiredSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expired sessions cleaned up"})
}