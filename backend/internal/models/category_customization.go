package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CategoryCustomization represents a category with custom icon and color
type CategoryCustomization struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	CategoryName string     `json:"category_name" db:"category_name"`
	Icon         *string    `json:"icon" db:"icon"`       // nullable for categories without custom icons
	Color        *string    `json:"color" db:"color"`     // nullable for categories without custom colors
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// CategoryCustomizationRequest represents the request body for category customization operations
type CategoryCustomizationRequest struct {
	CategoryName string  `json:"category_name" binding:"required"`
	Icon         *string `json:"icon"`
	Color        *string `json:"color"`
}

// ValidateIcon validates that the icon is a single character (emoji or text)
func (r *CategoryCustomizationRequest) ValidateIcon() error {
	if r.Icon != nil && len([]rune(*r.Icon)) != 1 {
		return fmt.Errorf("icon must be a single character or emoji")
	}
	return nil
}

// ValidateColor validates that the color is a valid hex color code
func (r *CategoryCustomizationRequest) ValidateColor() error {
	if r.Color != nil {
		color := *r.Color
		if len(color) != 7 || color[0] != '#' {
			return fmt.Errorf("color must be a valid hex color code (e.g., #FF5733)")
		}
		// Check if the remaining characters are valid hex
		for _, char := range color[1:] {
			if !((char >= '0' && char <= '9') || (char >= 'A' && char <= 'F') || (char >= 'a' && char <= 'f')) {
				return fmt.Errorf("color must be a valid hex color code (e.g., #FF5733)")
			}
		}
	}
	return nil
}

// Validate validates the entire request
func (r *CategoryCustomizationRequest) Validate() error {
	if err := r.ValidateIcon(); err != nil {
		return err
	}
	if err := r.ValidateColor(); err != nil {
		return err
	}
	return nil
}