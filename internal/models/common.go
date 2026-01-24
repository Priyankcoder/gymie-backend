package models

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool                `json:"success"`
	Data       interface{}         `json:"data"`
	Pagination *PaginationResponse `json:"pagination"`
}

// ListQuery represents common list query parameters
type ListQuery struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	SortBy   string `form:"sort_by"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// DateRangeQuery represents date range query parameters
type DateRangeQuery struct {
	StartDate string `form:"start_date"` // YYYY-MM-DD format
	EndDate   string `form:"end_date"`   // YYYY-MM-DD format
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(error, message string, details interface{}) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Error:   error,
		Message: message,
		Details: details,
	}
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, total int64, page, pageSize int) *PaginatedResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: &PaginationResponse{
			Total:       total,
			Page:        page,
			PageSize:    pageSize,
			TotalPages:  totalPages,
			HasNext:     page < totalPages,
			HasPrevious: page > 1,
		},
	}
}

// SetDefaults sets default values for ListQuery
func (q *ListQuery) SetDefaults() {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 20
	}
	if q.Order == "" {
		q.Order = "desc"
	}
}

// GetOffset calculates the offset for pagination
func (q *ListQuery) GetOffset() int {
	return (q.Page - 1) * q.PageSize
}
