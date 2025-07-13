package types

// APIResponse represents the standard response format for all API endpoints
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

// ErrorInfo represents error information in API responses
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// MetaInfo represents metadata in API responses
type MetaInfo struct {
	Total      int  `json:"total,omitempty"`
	Page       *int `json:"page,omitempty"`
	Limit      int  `json:"limit,omitempty"`
	Offset     int  `json:"offset,omitempty"`
	TotalPages *int `json:"total_pages,omitempty"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`
}

// Validate validates pagination parameters
func (p *PaginationRequest) Validate() error {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	return nil
}

// GetOffset calculates offset for database queries
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowOrigins []string `json:"allow_origins"`
	AllowMethods []string `json:"allow_methods"`
	AllowHeaders []string `json:"allow_headers"`
}