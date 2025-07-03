package errors

import (
	"fmt"
	"net/http"
	"sukimise/internal/constants"
	"sukimise/internal/types"

	"github.com/gin-gonic/gin"
)

// AppError represents an application error with HTTP status code and error code
type AppError struct {
	Message    string
	StatusCode int
	Code       string
	Details    string
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(message string, statusCode int, code string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
		Code:       code,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string, details string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Code:       constants.ErrorCodeValidation,
		Details:    details,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
		Code:       constants.ErrorCodeNotFound,
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
		Code:       constants.ErrorCodeUnauthorized,
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusForbidden,
		Code:       constants.ErrorCodeForbidden,
	}
}

// NewInternalError creates an internal server error
func NewInternalError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Code:       constants.ErrorCodeInternalError,
	}
}

// NewDatabaseError creates a database error
func NewDatabaseError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Code:       constants.ErrorCodeDatabaseError,
	}
}

// HandleError handles application errors and sends appropriate HTTP response
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		c.JSON(appErr.StatusCode, types.APIResponse{
			Success: false,
			Error: &types.ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Details,
			},
		})
		return
	}

	// Handle unknown errors
	c.JSON(http.StatusInternalServerError, types.APIResponse{
		Success: false,
		Error: &types.ErrorInfo{
			Code:    constants.ErrorCodeInternalError,
			Message: "An unexpected error occurred",
		},
	})
}

// SendSuccess sends a successful response
func SendSuccess(c *gin.Context, data interface{}, meta ...*types.MetaInfo) {
	response := types.APIResponse{
		Success: true,
		Data:    data,
	}

	if len(meta) > 0 && meta[0] != nil {
		response.Meta = meta[0]
	}

	c.JSON(http.StatusOK, response)
}

// SendCreated sends a created response
func SendCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, types.APIResponse{
		Success: true,
		Data:    data,
	})
}

// SendNoContent sends a no content response
func SendNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, types.APIResponse{
		Success: true,
	})
}