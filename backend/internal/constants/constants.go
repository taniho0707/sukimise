package constants

// HTTP Status Codes
const (
	StatusSuccess = "success"
	StatusError   = "error"
)

// User Roles
const (
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleViewer = "viewer"
)

// Pagination
const (
	DefaultLimit = 20
	MaxLimit     = 100
	DefaultPage  = 1
)

// File Upload
const (
	MaxFileSize    = 10 << 20 // 10MB
	AllowedImageTypes = "image/jpeg,image/png,image/gif,image/webp"
)

// Database
const (
	DefaultTimeout = 30 // seconds
)

// JWT
const (
	AccessTokenDuration  = 24   // hours
	RefreshTokenDuration = 168  // hours (7 days)
)

// Business Hours
const (
	BusinessDayMonday    = "monday"
	BusinessDayTuesday   = "tuesday"
	BusinessDayWednesday = "wednesday"
	BusinessDayThursday  = "thursday"
	BusinessDayFriday    = "friday"
	BusinessDaySaturday  = "saturday"
	BusinessDaySunday    = "sunday"
)

// Error Codes
const (
	ErrorCodeValidation     = "VALIDATION_ERROR"
	ErrorCodeNotFound       = "NOT_FOUND"
	ErrorCodeUnauthorized   = "UNAUTHORIZED"
	ErrorCodeForbidden      = "FORBIDDEN"
	ErrorCodeInternalError  = "INTERNAL_ERROR"
	ErrorCodeDatabaseError  = "DATABASE_ERROR"
	ErrorCodeFileUploadError = "FILE_UPLOAD_ERROR"
)