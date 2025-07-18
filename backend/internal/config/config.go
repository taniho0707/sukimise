package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Upload   UploadConfig   `yaml:"upload"`
	CORS     CORSConfig     `yaml:"cors"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string        `yaml:"port"`
	Host         string        `yaml:"host"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	Environment  string        `yaml:"environment"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL             string        `yaml:"url"`
	MaxConnections  int           `yaml:"max_connections"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	Timeout         time.Duration `yaml:"timeout"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret                string        `yaml:"secret"`
	AccessTokenDuration   time.Duration `yaml:"access_token_duration"`
	RefreshTokenDuration  time.Duration `yaml:"refresh_token_duration"`
	Issuer                string        `yaml:"issuer"`
}

// UploadConfig holds file upload configuration
type UploadConfig struct {
	MaxFileSize     int64    `yaml:"max_file_size"`
	AllowedTypes    []string `yaml:"allowed_types"`
	UploadDir       string   `yaml:"upload_dir"`
	BaseURL         string   `yaml:"base_url"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8081"),
			Host:         getEnv("HOST", "localhost"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
			Environment:  getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "postgres://sukimise:password@localhost:5432/sukimise_db?sslmode=disable"),
			MaxConnections:  getIntEnv("DB_MAX_CONNECTIONS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			Timeout:         getDurationEnv("DB_TIMEOUT", 30*time.Second),
		},
		JWT: JWTConfig{
			Secret:                getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			AccessTokenDuration:   getDurationEnv("JWT_ACCESS_DURATION", 24*time.Hour),
			RefreshTokenDuration:  getDurationEnv("JWT_REFRESH_DURATION", 7*24*time.Hour),
			Issuer:                getEnv("JWT_ISSUER", "sukimise"),
		},
		Upload: UploadConfig{
			MaxFileSize:  getInt64Env("UPLOAD_MAX_FILE_SIZE", 10*1024*1024), // 10MB
			AllowedTypes: getStringSliceEnv("UPLOAD_ALLOWED_TYPES", []string{"image/jpeg", "image/png", "image/gif", "image/webp"}),
			UploadDir:    getEnv("UPLOAD_DIR", "./uploads"),
			BaseURL:      getEnv("UPLOAD_BASE_URL", "http://localhost:8081"),
		},
		CORS: CORSConfig{
			AllowedOrigins:   getStringSliceEnv("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:5173"}),
			AllowedMethods:   getStringSliceEnv("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders:   getStringSliceEnv("CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization"}),
			AllowCredentials: getBoolEnv("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           getIntEnv("CORS_MAX_AGE", 86400), // 24 hours
		},
	}
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// GetServerAddress returns the server address
func (c *Config) GetServerAddress() string {
	if c.Server.Host == "" {
		return ":" + c.Server.Port
	}
	return c.Server.Host + ":" + c.Server.Port
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.JWT.Secret == "your-secret-key-change-in-production" && c.IsProduction() {
		log.Fatal("JWT secret must be changed in production")
	}
	
	if c.Database.URL == "" {
		log.Fatal("Database URL is required")
	}
	
	return nil
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if int64Value, err := strconv.ParseInt(value, 10, 64); err == nil {
			return int64Value
		}
		log.Printf("Warning: Invalid int64 value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		log.Printf("Warning: Invalid boolean value for %s: %s, using default: %t", key, value, defaultValue)
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("Warning: Invalid duration value for %s: %s, using default: %v", key, value, defaultValue)
	}
	return defaultValue
}

func getStringSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma-separated parsing
		// For more complex parsing, consider using a proper configuration library
		var result []string
		for _, item := range parseCommaSeparated(value) {
			if trimmed := trimSpace(item); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

// Helper function to parse comma-separated values
func parseCommaSeparated(value string) []string {
	var result []string
	current := ""
	
	for _, char := range value {
		if char == ',' {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	
	if current != "" {
		result = append(result, current)
	}
	
	return result
}

// Helper function to trim whitespace
func trimSpace(s string) string {
	start := 0
	end := len(s)
	
	// Trim leading whitespace
	for start < end && isSpace(s[start]) {
		start++
	}
	
	// Trim trailing whitespace
	for start < end && isSpace(s[end-1]) {
		end--
	}
	
	return s[start:end]
}

// Helper function to check if character is whitespace
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}