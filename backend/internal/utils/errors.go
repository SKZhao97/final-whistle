package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorCode represents a standardized error code
type ErrorCode string

const (
	// ValidationError indicates invalid input
	ValidationError ErrorCode = "VALIDATION_ERROR"
	// Unauthorized indicates authentication required
	Unauthorized ErrorCode = "UNAUTHORIZED"
	// Forbidden indicates insufficient permissions
	Forbidden ErrorCode = "FORBIDDEN"
	// NotFound indicates resource not found
	NotFound ErrorCode = "NOT_FOUND"
	// Conflict indicates resource conflict (e.g., duplicate)
	Conflict ErrorCode = "CONFLICT"
	// InternalError indicates server-side error
	InternalError ErrorCode = "INTERNAL_ERROR"
	// ServiceUnavailable indicates service temporarily unavailable
	ServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
)

// APIError represents a standardized API error response
type APIError struct {
	Success bool                `json:"success"`
	Error   APIErrorDetails     `json:"error"`
}

// APIErrorDetails contains error details
type APIErrorDetails struct {
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// NewAPIError creates a new APIError
func NewAPIError(code ErrorCode, message string, details map[string]interface{}) APIError {
	return APIError{
		Success: false,
		Error: APIErrorDetails{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// ErrorResponse sends a standardized error response
func ErrorResponse(c *gin.Context, status int, code ErrorCode, message string, details map[string]interface{}) {
	c.JSON(status, NewAPIError(code, message, details))
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, message string, details map[string]interface{}) {
	if details == nil {
		details = make(map[string]interface{})
	}
	ErrorResponse(c, http.StatusBadRequest, ValidationError, message, details)
}

// UnauthorizedResponse sends an unauthorized error response
func UnauthorizedResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Authentication required"
	}
	ErrorResponse(c, http.StatusUnauthorized, Unauthorized, message, nil)
}

// ForbiddenResponse sends a forbidden error response
func ForbiddenResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Insufficient permissions"
	}
	ErrorResponse(c, http.StatusForbidden, Forbidden, message, nil)
}

// NotFoundResponse sends a not found error response
func NotFoundResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	ErrorResponse(c, http.StatusNotFound, NotFound, message, nil)
}

// ConflictResponse sends a conflict error response
func ConflictResponse(c *gin.Context, message string, details map[string]interface{}) {
	if message == "" {
		message = "Resource conflict"
	}
	ErrorResponse(c, http.StatusConflict, Conflict, message, details)
}

// InternalErrorResponse sends an internal error response
func InternalErrorResponse(c *gin.Context, message string, details map[string]interface{}) {
	if message == "" {
		message = "Internal server error"
	}
	ErrorResponse(c, http.StatusInternalServerError, InternalError, message, details)
}

// ServiceUnavailableResponse sends a service unavailable response
func ServiceUnavailableResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Service temporarily unavailable"
	}
	ErrorResponse(c, http.StatusServiceUnavailable, ServiceUnavailable, message, nil)
}

// SuccessResponse sends a standardized success response
func SuccessResponse(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{
		"success": true,
		"data":    data,
	})
}

// OKResponse sends a 200 OK success response
func OKResponse(c *gin.Context, data interface{}) {
	SuccessResponse(c, http.StatusOK, data)
}

// CreatedResponse sends a 201 Created success response
func CreatedResponse(c *gin.Context, data interface{}) {
	SuccessResponse(c, http.StatusCreated, data)
}

// NoContentResponse sends a 204 No Content response
func NoContentResponse(c *gin.Context) {
	c.Status(http.StatusNoContent)
}