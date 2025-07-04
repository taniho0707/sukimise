package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"sukimise/internal/models"
	"sukimise/internal/repositories"
)

type ViewerAuthService struct {
	repo *repositories.ViewerAuthRepository
}

func NewViewerAuthService(repo *repositories.ViewerAuthRepository) *ViewerAuthService {
	return &ViewerAuthService{repo: repo}
}

func (s *ViewerAuthService) AuthenticateViewer(password, ipAddress, userAgent string) (*models.ViewerAuthResponse, error) {
	settings, err := s.repo.GetViewerSettings()
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(settings.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	// Generate session token
	sessionToken, err := s.generateSessionToken()
	if err != nil {
		return nil, err
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(settings.SessionDurationDays) * 24 * time.Hour)

	// Save login history
	err = s.repo.CreateViewerLoginHistory(ipAddress, userAgent, sessionToken, expiresAt)
	if err != nil {
		return nil, err
	}

	return &models.ViewerAuthResponse{
		Token:     sessionToken,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *ViewerAuthService) ValidateViewerSession(sessionToken string) (bool, error) {
	return s.repo.ValidateViewerSession(sessionToken)
}

func (s *ViewerAuthService) GetViewerSettings() (*models.ViewerSettings, error) {
	return s.repo.GetViewerSettings()
}

func (s *ViewerAuthService) UpdateViewerSettings(password string, sessionDurationDays int) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdateViewerSettings(string(hashedPassword), sessionDurationDays)
}

func (s *ViewerAuthService) GetViewerLoginHistory(page, limit int) ([]models.ViewerLoginHistory, int, error) {
	offset := (page - 1) * limit
	history, err := s.repo.GetViewerLoginHistory(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repo.GetViewerLoginHistoryCount()
	if err != nil {
		return nil, 0, err
	}

	return history, count, nil
}

func (s *ViewerAuthService) CleanupExpiredSessions() error {
	return s.repo.CleanupExpiredSessions()
}

func (s *ViewerAuthService) generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}