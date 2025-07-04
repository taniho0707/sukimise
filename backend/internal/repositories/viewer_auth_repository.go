package repositories

import (
	"database/sql"
	"time"

	"sukimise/internal/models"
)

type ViewerAuthRepository struct {
	db *sql.DB
}

func NewViewerAuthRepository(db *sql.DB) *ViewerAuthRepository {
	return &ViewerAuthRepository{db: db}
}

func (r *ViewerAuthRepository) GetViewerSettings() (*models.ViewerSettings, error) {
	var settings models.ViewerSettings
	query := `
		SELECT id, password_hash, session_duration_days, created_at, updated_at
		FROM viewer_settings
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := r.db.QueryRow(query).Scan(
		&settings.ID,
		&settings.PasswordHash,
		&settings.SessionDurationDays,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *ViewerAuthRepository) UpdateViewerSettings(passwordHash string, sessionDurationDays int) error {
	query := `
		UPDATE viewer_settings 
		SET password_hash = $1, session_duration_days = $2, updated_at = NOW()
		WHERE id = (SELECT id FROM viewer_settings ORDER BY created_at DESC LIMIT 1)
	`
	_, err := r.db.Exec(query, passwordHash, sessionDurationDays)
	return err
}

func (r *ViewerAuthRepository) CreateViewerLoginHistory(ipAddress, userAgent, sessionToken string, expiresAt time.Time) error {
	query := `
		INSERT INTO viewer_login_history (ip_address, user_agent, session_token, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, ipAddress, userAgent, sessionToken, expiresAt)
	return err
}

func (r *ViewerAuthRepository) GetViewerLoginHistory(limit, offset int) ([]models.ViewerLoginHistory, error) {
	var history []models.ViewerLoginHistory
	query := `
		SELECT id, ip_address, user_agent, login_time, session_token, expires_at
		FROM viewer_login_history
		ORDER BY login_time DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.ViewerLoginHistory
		err := rows.Scan(
			&item.ID,
			&item.IPAddress,
			&item.UserAgent,
			&item.LoginTime,
			&item.SessionToken,
			&item.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, item)
	}
	return history, nil
}

func (r *ViewerAuthRepository) GetViewerLoginHistoryCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM viewer_login_history`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

func (r *ViewerAuthRepository) ValidateViewerSession(sessionToken string) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM viewer_login_history
		WHERE session_token = $1 AND expires_at > NOW()
	`
	err := r.db.QueryRow(query, sessionToken).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ViewerAuthRepository) CleanupExpiredSessions() error {
	query := `DELETE FROM viewer_login_history WHERE expires_at < NOW()`
	_, err := r.db.Exec(query)
	return err
}