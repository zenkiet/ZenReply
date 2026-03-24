package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type Meta struct {
	Total     int `json:"total"`
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	TotalPage int `json:"total_page"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// 200 OK
func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// 200 OK with Meta (Pagination, etc...)
func OKWithMeta(c *gin.Context, message string, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// 201 Created
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// 204 No Content
func NoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, Response{
		Success: true,
	})
}

// 400 Bad Request
func BadRequest(c *gin.Context, code, message, details string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// 401 Unauthorized
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "UNAUTHORIZED",
			Message: message,
		},
	})
}

// 403 Forbidden
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "FORBIDDEN",
			Message: message,
		},
	})
}

// 404 Not Found
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "NOT_FOUND",
			Message: message,
		},
	})
}

// 405 Method Not Allowed
func MethodNotAllowed(c *gin.Context, message string) {
	c.JSON(http.StatusMethodNotAllowed, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "METHOD_NOT_ALLOWED",
			Message: message,
		},
	})
}

// 409 Conflict
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "CONFLICT",
			Message: message,
		},
	})
}

// 429 Too Many Requests
func TooManyRequests(c *gin.Context, message string) {
	c.JSON(http.StatusTooManyRequests, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "TOO_MANY_REQUESTS",
			Message: message,
		},
	})
}

// 500 Internal Server Error
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: message,
		},
	})
}

// 503 Service Unavailable
func ServiceUnavailable(c *gin.Context, message string) {
	c.JSON(http.StatusServiceUnavailable, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "SERVICE_UNAVAILABLE",
			Message: message,
		},
	})
}
