package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an API error
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, code, message, details string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: "Request failed",
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request data", err.Error())
}

// NotFoundResponse sends a not found error response
func NotFoundResponse(c *gin.Context, resource string) {
	ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", resource+" not found", "")
}

// InternalErrorResponse sends an internal server error response
func InternalErrorResponse(c *gin.Context, err error) {
	GetLogger().Error("Internal server error: ", err)
	ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", "")
}

// UnauthorizedResponse sends an unauthorized error response
func UnauthorizedResponse(c *gin.Context) {
	ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", "")
}

// ForbiddenResponse sends a forbidden error response
func ForbiddenResponse(c *gin.Context) {
	ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "Access denied", "")
}
