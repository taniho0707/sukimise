package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	// Get Discord link with stored token
	link, err := s.GetDiscordLink(discordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Discord link: %v", err)
	}

	// Check if token is expired and refresh if needed
	accessToken := link.AccessToken
	if time.Now().After(link.TokenExpiry) {
		// Token is expired, try to refresh
		newToken, err := s.refreshToken(link.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token, please reconnect: %v", err)
		}
		accessToken = newToken
		// Update stored token (simplified, should also update refresh token and expiry)
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

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make store creation request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
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

func (s *DiscordService) refreshToken(refreshToken string) (string, error) {
	refreshReq := map[string]string{
		"refresh_token": refreshToken,
	}

	jsonData, err := json.Marshal(refreshReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal refresh request: %v", err)
	}

	resp, err := http.Post(s.sukimiseAPIURL+"/api/v1/auth/refresh", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to make refresh request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("refresh failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var refreshResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err != nil {
		return "", fmt.Errorf("failed to decode refresh response: %v", err)
	}

	return refreshResp.AccessToken, nil
}