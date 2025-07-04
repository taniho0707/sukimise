package utils

import (
	"sukimise/internal/constants"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectError bool
	}{
		{
			name:        "valid email",
			email:       "test@example.com",
			expectError: false,
		},
		{
			name:        "valid email with subdomain",
			email:       "user@mail.example.com",
			expectError: false,
		},
		{
			name:        "valid email with numbers",
			email:       "user123@example.co.jp",
			expectError: false,
		},
		{
			name:        "empty email",
			email:       "",
			expectError: true,
		},
		{
			name:        "invalid email - no @",
			email:       "invalidEmail",
			expectError: true,
		},
		{
			name:        "invalid email - no domain",
			email:       "user@",
			expectError: true,
		},
		{
			name:        "invalid email - no TLD",
			email:       "user@example",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expectError bool
	}{
		{
			name:        "valid HTTP URL",
			url:         "http://example.com",
			expectError: false,
		},
		{
			name:        "valid HTTPS URL",
			url:         "https://example.com",
			expectError: false,
		},
		{
			name:        "valid URL with path",
			url:         "https://example.com/path/to/page",
			expectError: false,
		},
		{
			name:        "valid URL with query",
			url:         "https://example.com?param=value",
			expectError: false,
		},
		{
			name:        "empty URL",
			url:         "",
			expectError: false, // URL is optional
		},
		{
			name:        "invalid URL - no protocol",
			url:         "example.com",
			expectError: true,
		},
		{
			name:        "invalid URL - invalid protocol",
			url:         "ftp://example.com",
			expectError: true,
		},
		{
			name:        "invalid URL - malformed",
			url:         "http://",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCoordinates(t *testing.T) {
	tests := []struct {
		name        string
		lat         float64
		lng         float64
		expectError bool
	}{
		{
			name:        "valid coordinates - Tokyo",
			lat:         35.6762,
			lng:         139.6503,
			expectError: false,
		},
		{
			name:        "valid coordinates - minimum values",
			lat:         -90.0,
			lng:         -180.0,
			expectError: false,
		},
		{
			name:        "valid coordinates - maximum values",
			lat:         90.0,
			lng:         180.0,
			expectError: false,
		},
		{
			name:        "invalid latitude - too low",
			lat:         -91.0,
			lng:         0.0,
			expectError: true,
		},
		{
			name:        "invalid latitude - too high",
			lat:         91.0,
			lng:         0.0,
			expectError: true,
		},
		{
			name:        "invalid longitude - too low",
			lat:         0.0,
			lng:         -181.0,
			expectError: true,
		},
		{
			name:        "invalid longitude - too high",
			lat:         0.0,
			lng:         181.0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCoordinates(tt.lat, tt.lng)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUserRole(t *testing.T) {
	tests := []struct {
		name        string
		role        string
		expectError bool
	}{
		{
			name:        "valid role - admin",
			role:        constants.RoleAdmin,
			expectError: false,
		},
		{
			name:        "valid role - editor",
			role:        constants.RoleEditor,
			expectError: false,
		},
		{
			name:        "valid role - viewer",
			role:        constants.RoleViewer,
			expectError: false,
		},
		{
			name:        "invalid role",
			role:        "invalid",
			expectError: true,
		},
		{
			name:        "empty role",
			role:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserRole(tt.role)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateBusinessDay(t *testing.T) {
	tests := []struct {
		name        string
		day         string
		expectError bool
	}{
		{
			name:        "valid day - monday",
			day:         constants.BusinessDayMonday,
			expectError: false,
		},
		{
			name:        "valid day - sunday",
			day:         constants.BusinessDaySunday,
			expectError: false,
		},
		{
			name:        "empty day",
			day:         "",
			expectError: false, // Business day is optional
		},
		{
			name:        "invalid day",
			day:         "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBusinessDay(tt.day)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateStringArray(t *testing.T) {
	tests := []struct {
		name        string
		array       []string
		fieldName   string
		maxLength   int
		expectError bool
	}{
		{
			name:        "valid array",
			array:       []string{"item1", "item2", "item3"},
			fieldName:   "categories",
			maxLength:   5,
			expectError: false,
		},
		{
			name:        "empty array",
			array:       []string{},
			fieldName:   "categories",
			maxLength:   5,
			expectError: false,
		},
		{
			name:        "array too long",
			array:       []string{"item1", "item2", "item3", "item4", "item5", "item6"},
			fieldName:   "categories",
			maxLength:   5,
			expectError: true,
		},
		{
			name:        "array with empty item",
			array:       []string{"item1", "", "item3"},
			fieldName:   "categories",
			maxLength:   5,
			expectError: true,
		},
		{
			name:        "array with whitespace-only item",
			array:       []string{"item1", "   ", "item3"},
			fieldName:   "categories",
			maxLength:   5,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStringArray(tt.array, tt.fieldName, tt.maxLength)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateRating(t *testing.T) {
	tests := []struct {
		name        string
		rating      int
		expectError bool
	}{
		{
			name:        "valid rating - 1",
			rating:      1,
			expectError: false,
		},
		{
			name:        "valid rating - 3",
			rating:      3,
			expectError: false,
		},
		{
			name:        "valid rating - 5",
			rating:      5,
			expectError: false,
		},
		{
			name:        "invalid rating - 0",
			rating:      0,
			expectError: true,
		},
		{
			name:        "invalid rating - 6",
			rating:      6,
			expectError: true,
		},
		{
			name:        "invalid rating - negative",
			rating:      -1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRating(tt.rating)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePaymentAmount(t *testing.T) {
	tests := []struct {
		name        string
		amount      *int
		expectError bool
	}{
		{
			name:        "nil amount",
			amount:      nil,
			expectError: false,
		},
		{
			name:        "valid amount - zero",
			amount:      intPtr(0),
			expectError: false,
		},
		{
			name:        "valid amount - positive",
			amount:      intPtr(1000),
			expectError: false,
		},
		{
			name:        "invalid amount - negative",
			amount:      intPtr(-100),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePaymentAmount(tt.amount)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}