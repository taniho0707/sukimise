package models

import (
	"time"

	"github.com/google/uuid"
)

type ViewerSettings struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	PasswordHash         string    `json:"-" db:"password_hash"`
	SessionDurationDays  int       `json:"session_duration_days" db:"session_duration_days"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type ViewerLoginHistory struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	IPAddress    string     `json:"ip_address" db:"ip_address"`
	UserAgent    string     `json:"user_agent" db:"user_agent"`
	LoginTime    time.Time  `json:"login_time" db:"login_time"`
	SessionToken string     `json:"session_token" db:"session_token"`
	ExpiresAt    time.Time  `json:"expires_at" db:"expires_at"`
}

type ViewerAuthRequest struct {
	Password string `json:"password" binding:"required"`
}

type ViewerAuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ViewerSettingsUpdateRequest struct {
	Password            string `json:"password" binding:"required"`
	SessionDurationDays int    `json:"session_duration_days" binding:"required,min=1,max=365"`
}