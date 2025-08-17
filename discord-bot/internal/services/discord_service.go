package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"sukimise-discord-bot/internal/models"
	"sukimise-discord-bot/internal/utils"

	"github.com/google/uuid"
)

type DiscordService struct {
	db             *sql.DB
	sukimiseAPIURL string
}

func NewDiscordService(db *sql.DB, sukimiseAPIURL string) *DiscordService {
	return &DiscordService{
		db:             db,
		sukimiseAPIURL: sukimiseAPIURL,
	}
}

func (s *DiscordService) ConnectDiscordUser(discordID, username, password string) (*models.DiscordLink, error) {
	// Check if Discord user is already linked
	existingLink, err := s.GetDiscordLink(discordID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check existing link: %v", err)
	}
	if existingLink != nil {
		return nil, fmt.Errorf("discord user is already linked to username: %s", existingLink.Username)
	}

	// Authenticate with Sukimise API
	authResp, err := s.authenticateWithSukimise(username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Sukimise: %v", err)
	}

	// Check if Sukimise user is already linked to another Discord account
	existingUserLink, err := s.GetDiscordLinkByUserID(authResp.User.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check existing user link: %v", err)
	}
	if existingUserLink != nil {
		return nil, fmt.Errorf("sukimise user is already linked to another Discord account")
	}

	// Create new Discord link with tokens
	link := &models.DiscordLink{
		ID:           uuid.New(),
		DiscordID:    discordID,
		UserID:       authResp.User.ID,
		Username:     username,
		AccessToken:  authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
		TokenExpiry:  time.Now().Add(24 * time.Hour), // Tokens typically expire in 24 hours
		LinkedAt:     time.Now(),
		LastUsedAt:   time.Now(),
	}

	err = s.CreateDiscordLink(link)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord link: %v", err)
	}

	return link, nil
}

func (s *DiscordService) DisconnectDiscordUser(discordID string) error {
	// Check if Discord user is linked
	_, err := s.GetDiscordLink(discordID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("discord user is not linked to any Sukimise account")
	}
	if err != nil {
		return fmt.Errorf("failed to check Discord link: %v", err)
	}

	// Delete Discord link
	err = s.DeleteDiscordLink(discordID)
	if err != nil {
		return fmt.Errorf("failed to delete Discord link: %v", err)
	}

	return nil
}

func (s *DiscordService) AddStoreFromGoogleMaps(discordID, googleMapURL string) (*models.StoreCreateResponse, error) {
	// Get Discord link
	_, err := s.GetDiscordLink(discordID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("discord user is not linked to any Sukimise account. Use /connect first")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get Discord link: %v", err)
	}

	// Update last used time
	err = s.UpdateLastUsedAt(discordID)
	if err != nil {
		return nil, fmt.Errorf("failed to update last used time: %v", err)
	}

	// Validate Google Maps URL
	if !utils.IsValidGoogleMapsURL(googleMapURL) {
		return nil, fmt.Errorf("invalid Google Maps URL. Must start with https://www.google.com/maps/place/")
	}

	// Extract store information from Google Maps URL
	storeInfo, err := utils.ExtractStoreInfoFromURL(googleMapURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract store information: %v", err)
	}

	// Create store via Sukimise API
	storeResp, err := s.createStoreViaSukimise(discordID, storeInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to create store via Sukimise API: %v", err)
	}

	return storeResp, nil
}

func (s *DiscordService) CreateDiscordLink(link *models.DiscordLink) error {
	query := `
		INSERT INTO discord_links (id, discord_id, user_id, username, access_token, refresh_token, token_expiry, linked_at, last_used_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := s.db.Exec(query, link.ID, link.DiscordID, link.UserID, link.Username, link.AccessToken, link.RefreshToken, link.TokenExpiry, link.LinkedAt, link.LastUsedAt)
	return err
}

func (s *DiscordService) GetDiscordLink(discordID string) (*models.DiscordLink, error) {
	query := `
		SELECT id, discord_id, user_id, username, access_token, refresh_token, token_expiry, linked_at, last_used_at
		FROM discord_links
		WHERE discord_id = $1
	`
	var link models.DiscordLink
	err := s.db.QueryRow(query, discordID).Scan(
		&link.ID, &link.DiscordID, &link.UserID, &link.Username, &link.AccessToken, &link.RefreshToken, &link.TokenExpiry, &link.LinkedAt, &link.LastUsedAt,
	)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (s *DiscordService) GetDiscordLinkByUserID(userID uuid.UUID) (*models.DiscordLink, error) {
	query := `
		SELECT id, discord_id, user_id, username, access_token, refresh_token, token_expiry, linked_at, last_used_at
		FROM discord_links
		WHERE user_id = $1
	`
	var link models.DiscordLink
	err := s.db.QueryRow(query, userID).Scan(
		&link.ID, &link.DiscordID, &link.UserID, &link.Username, &link.AccessToken, &link.RefreshToken, &link.TokenExpiry, &link.LinkedAt, &link.LastUsedAt,
	)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (s *DiscordService) DeleteDiscordLink(discordID string) error {
	query := `DELETE FROM discord_links WHERE discord_id = $1`
	_, err := s.db.Exec(query, discordID)
	return err
}

func (s *DiscordService) UpdateLastUsedAt(discordID string) error {
	query := `UPDATE discord_links SET last_used_at = NOW() WHERE discord_id = $1`
	_, err := s.db.Exec(query, discordID)
	return err
}

func (s *DiscordService) UpdateTokens(discordID, accessToken, refreshToken string, expiry time.Time) error {
	query := `
		UPDATE discord_links 
		SET access_token = $2, refresh_token = $3, token_expiry = $4, last_used_at = NOW() 
		WHERE discord_id = $1
	`
	_, err := s.db.Exec(query, discordID, accessToken, refreshToken, expiry)
	return err
}

func (s *DiscordService) authenticateWithSukimise(username, password string) (*models.AuthResponse, error) {
	loginReq := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal login request: %v", err)
	}

	resp, err := http.Post(s.sukimiseAPIURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make login request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var authResp models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, fmt.Errorf("failed to decode authentication response: %v", err)
	}

	return &authResp, nil
}

func (s *DiscordService) createStoreViaSukimise(discordID string, storeInfo *models.StoreCreateRequest) (*models.StoreCreateResponse, error) {
	// Ensure we have a valid access token (refresh if necessary)
	accessToken, err := s.ensureValidToken(discordID)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(storeInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal store creation request: %v", err)
	}

	req, err := http.NewRequest("POST", s.sukimiseAPIURL+"/api/v1/stores", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create store request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	
	// Debug: Print request headers
	fmt.Printf("DEBUG: Request headers:\n")
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make store creation request: %v", err)
	}
	defer resp.Body.Close()

	// Handle authentication errors - token might have been invalidated after our check
	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Printf("DEBUG: Got 401 Unauthorized, attempting token refresh\n")
		
		// Try to refresh token one more time
		newAccessToken, refreshErr := s.ensureValidToken(discordID)
		if refreshErr != nil {
			return nil, refreshErr
		}
		
		// Create a new request with the refreshed token (need to recreate body)
		newReq, err := http.NewRequest("POST", s.sukimiseAPIURL+"/api/v1/stores", bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("failed to create retry request: %v", err)
		}
		newReq.Header.Set("Content-Type", "application/json")
		newReq.Header.Set("Authorization", "Bearer "+newAccessToken)
		
		resp, err = client.Do(newReq)
		if err != nil {
			return nil, fmt.Errorf("failed to retry store creation request: %v", err)
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		
		// Handle duplicate store error (HTTP 409 Conflict)
		if resp.StatusCode == http.StatusConflict {
			// Parse error response to get details
			var errorResp struct {
				Success bool `json:"success"`
				Error   struct {
					Code    string `json:"code"`
					Message string `json:"message"`
					Details string `json:"details"`
				} `json:"error"`
			}
			
			if err := json.Unmarshal(bodyBytes, &errorResp); err == nil {
				return nil, fmt.Errorf("重複する店舗が見つかりました: %s", errorResp.Error.Details)
			}
			
			return nil, fmt.Errorf("重複する店舗が見つかりました: この店舗は既に登録されています")
		}
		
		return nil, fmt.Errorf("store creation failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read response body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	
	fmt.Printf("DEBUG: API Response: %s\n", string(bodyBytes))

	var apiResp models.APIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode store creation response: %v", err)
	}
	
	fmt.Printf("DEBUG: Parsed response - Name: '%s', Address: '%s'\n", apiResp.Data.Name, apiResp.Data.Address)

	return &apiResp.Data, nil
}

func (s *DiscordService) refreshToken(refreshToken string) (*models.AuthResponse, error) {
	refreshReq := map[string]string{
		"refresh_token": refreshToken,
	}

	jsonData, err := json.Marshal(refreshReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal refresh request: %v", err)
	}

	resp, err := http.Post(s.sukimiseAPIURL+"/api/v1/auth/refresh", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make refresh request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("refresh failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var refreshResp models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err != nil {
		return nil, fmt.Errorf("failed to decode refresh response: %v", err)
	}

	return &refreshResp, nil
}

// ensureValidToken ensures that the user has a valid access token, refreshing if necessary,
// or prompting for re-authentication if the refresh token is also invalid
func (s *DiscordService) ensureValidToken(discordID string) (string, error) {
	// Get Discord link with stored tokens
	link, err := s.GetDiscordLink(discordID)
	if err != nil {
		return "", fmt.Errorf("failed to get Discord link: %v", err)
	}

	// If token is still valid, return it
	if time.Now().Before(link.TokenExpiry) {
		return link.AccessToken, nil
	}

	fmt.Printf("DEBUG: Token expired for user %s, attempting refresh\n", link.Username)

	// Token is expired, try to refresh
	authResp, err := s.refreshToken(link.RefreshToken)
	if err != nil {
		fmt.Printf("DEBUG: Token refresh failed for user %s: %v\n", link.Username, err)
		
		// Check if it's a 401 error (invalid refresh token)
		if isRefreshTokenInvalid(err) {
			// Remove the invalid Discord link
			if deleteErr := s.DeleteDiscordLink(discordID); deleteErr != nil {
				fmt.Printf("DEBUG: Failed to delete invalid Discord link: %v\n", deleteErr)
			}
			return "", fmt.Errorf("認証が期限切れです。/connect コマンドを使用してアカウントを再接続してください")
		}
		
		return "", fmt.Errorf("トークンの更新に失敗しました: %v", err)
	}

	fmt.Printf("DEBUG: Token refreshed successfully for user %s\n", link.Username)

	// Update stored tokens
	newExpiry := time.Now().Add(24 * time.Hour) // Tokens typically expire in 24 hours
	if err := s.UpdateTokens(discordID, authResp.AccessToken, authResp.RefreshToken, newExpiry); err != nil {
		fmt.Printf("DEBUG: Failed to update tokens in database: %v\n", err)
		// Continue with the new token even if DB update failed
	}

	return authResp.AccessToken, nil
}

// isRefreshTokenInvalid checks if the error indicates an invalid refresh token (401 status)
func isRefreshTokenInvalid(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "refresh failed with status 401")
}