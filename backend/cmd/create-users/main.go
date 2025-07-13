package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"sukimise/internal/config"
	"sukimise/internal/models"
	"sukimise/internal/repositories"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// UserInfo represents user information from environment variables
type UserInfo struct {
	Username string
	Password string // bcrypt hash
	Role     string
}

func main() {
	log.Println("Starting user creation from environment variables...")

	// Load configuration
	cfg := config.LoadConfig()
	if cfg.Database.URL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Parse admin and editor users from environment variables
	adminUsers := parseUsersFromEnv("ADMIN_USERS", "admin")
	editorUsers := parseUsersFromEnv("EDITOR_USERS", "editor")

	if len(adminUsers) == 0 {
		log.Fatal("ADMIN_USERS environment variable is required and must contain at least one admin user")
	}

	if len(editorUsers) == 0 {
		log.Fatal("EDITOR_USERS environment variable is required and must contain at least one editor user")
	}

	log.Printf("Found %d admin users and %d editor users in environment variables", len(adminUsers), len(editorUsers))

	// Connect to database
	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create user repository
	userRepo := repositories.NewUserRepository(db)

	// Create all users
	allUsers := append(adminUsers, editorUsers...)
	for _, userInfo := range allUsers {
		err := createUser(userRepo, userInfo)
		if err != nil {
			log.Printf("Failed to create user '%s': %v", userInfo.Username, err)
			continue
		}
		log.Printf("Successfully created %s user: %s", userInfo.Role, userInfo.Username)
	}

	log.Println("User creation completed successfully")
}

// parseUsersFromEnv parses users from environment variable in format: username1:hash1;username2:hash2
func parseUsersFromEnv(envVar, role string) []UserInfo {
	envValue := os.Getenv(envVar)
	if envValue == "" {
		return []UserInfo{}
	}

	var users []UserInfo
	userEntries := strings.Split(envValue, ";")

	for _, entry := range userEntries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		parts := strings.SplitN(entry, ":", 2)
		if len(parts) != 2 {
			log.Printf("Invalid user entry format in %s: %s (expected username:hash)", envVar, entry)
			continue
		}

		username := strings.TrimSpace(parts[0])
		password := strings.TrimSpace(parts[1])

		if username == "" || password == "" {
			log.Printf("Empty username or password in %s: %s", envVar, entry)
			continue
		}

		users = append(users, UserInfo{
			Username: username,
			Password: password,
			Role:     role,
		})
	}

	return users
}

// createUser creates a user in the database
func createUser(userRepo repositories.UserRepositoryInterface, userInfo UserInfo) error {
	// Check if user already exists
	existingUser, err := userRepo.GetByUsername(userInfo.Username)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check if user exists: %v", err)
	}

	if existingUser != nil {
		log.Printf("User '%s' already exists, skipping", userInfo.Username)
		return nil
	}

	// Create new user
	user := &models.User{
		ID:        uuid.New(),
		Username:  userInfo.Username,
		Email:     "", // Email is not used in the new system
		Password:  userInfo.Password, // Already bcrypt hashed
		Role:      userInfo.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return userRepo.Create(user)
}