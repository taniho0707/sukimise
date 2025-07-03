package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sukimise/internal/constants"
	"sukimise/internal/errors"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	urlRegex   = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
)

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	if email == "" {
		return errors.NewValidationError("Email is required", "")
	}
	if !emailRegex.MatchString(email) {
		return errors.NewValidationError("Invalid email format", "")
	}
	return nil
}

// ValidateURL validates URL format
func ValidateURL(urlStr string) error {
	if urlStr == "" {
		return nil // URL is optional
	}
	
	if !urlRegex.MatchString(urlStr) {
		return errors.NewValidationError("Invalid URL format", fmt.Sprintf("URL: %s", urlStr))
	}
	
	// Additional validation using net/url
	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return errors.NewValidationError("Invalid URL format", err.Error())
	}
	
	return nil
}

// ValidateCoordinates validates latitude and longitude
func ValidateCoordinates(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return errors.NewValidationError("Invalid latitude", "Latitude must be between -90 and 90")
	}
	if lng < -180 || lng > 180 {
		return errors.NewValidationError("Invalid longitude", "Longitude must be between -180 and 180")
	}
	return nil
}

// ValidateUserRole validates user role
func ValidateUserRole(role string) error {
	validRoles := []string{constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer}
	for _, validRole := range validRoles {
		if role == validRole {
			return nil
		}
	}
	return errors.NewValidationError("Invalid user role", fmt.Sprintf("Role must be one of: %s", strings.Join(validRoles, ", ")))
}

// ValidateBusinessDay validates business day
func ValidateBusinessDay(day string) error {
	if day == "" {
		return nil // Business day is optional
	}
	
	validDays := []string{
		constants.BusinessDayMonday,
		constants.BusinessDayTuesday,
		constants.BusinessDayWednesday,
		constants.BusinessDayThursday,
		constants.BusinessDayFriday,
		constants.BusinessDaySaturday,
		constants.BusinessDaySunday,
	}
	
	for _, validDay := range validDays {
		if day == validDay {
			return nil
		}
	}
	
	return errors.NewValidationError("Invalid business day", fmt.Sprintf("Day must be one of: %s", strings.Join(validDays, ", ")))
}

// ValidateStringArray validates string array fields
func ValidateStringArray(arr []string, fieldName string, maxLength int) error {
	if len(arr) > maxLength {
		return errors.NewValidationError(
			fmt.Sprintf("%s array too long", fieldName),
			fmt.Sprintf("Maximum %d items allowed", maxLength),
		)
	}
	
	for i, item := range arr {
		if strings.TrimSpace(item) == "" {
			return errors.NewValidationError(
				fmt.Sprintf("Empty %s item", fieldName),
				fmt.Sprintf("Item at index %d is empty", i),
			)
		}
	}
	
	return nil
}

// ValidateRating validates rating value
func ValidateRating(rating int) error {
	if rating < 1 || rating > 5 {
		return errors.NewValidationError("Invalid rating", "Rating must be between 1 and 5")
	}
	return nil
}

// ValidatePaymentAmount validates payment amount
func ValidatePaymentAmount(amount *int) error {
	if amount != nil && *amount < 0 {
		return errors.NewValidationError("Invalid payment amount", "Payment amount cannot be negative")
	}
	return nil
}