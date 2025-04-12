package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/service"
)

// StandardResponse defines the standard API response format
type StandardResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success sends a successful response with data
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, StandardResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// Created sends a successful creation response with data
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, StandardResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends a 400 bad request response
func BadRequest(c *gin.Context, message string, err error) {
	var svcErr *service.ServiceError
	if errors.As(err, &svcErr) && errors.Is(svcErr.Unwrap(), service.ErrValidation) {
		// Handle validation errors specifically
		c.JSON(http.StatusBadRequest, StandardResponse{
			Status:  false,
			Message: message,
			Error:   svcErr.GetValidationErrors(),
		})
		return
	}

	c.JSON(http.StatusBadRequest, StandardResponse{
		Status:  false,
		Message: message,
		Error:   err.Error(),
	})
}

// Unauthorized sends a 401 unauthorized response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, StandardResponse{
		Status:  false,
		Message: message,
	})
}

// Forbidden sends a 403 forbidden response
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, StandardResponse{
		Status:  false,
		Message: message,
	})
}

// NotFound sends a 404 not found response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, StandardResponse{
		Status:  false,
		Message: message,
	})
}

// Conflict sends a 409 conflict response
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, StandardResponse{
		Status:  false,
		Message: message,
	})
}

// InternalServerError sends a 500 internal server error response
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, StandardResponse{
		Status:  false,
		Message: message,
	})
}

// HandleServiceError handles common service errors and maps them to appropriate HTTP responses
func HandleServiceError(c *gin.Context, err error, operation string) {
	var svcErr *service.ServiceError
	if errors.As(err, &svcErr) {
		switch {
		case errors.Is(svcErr.Unwrap(), service.ErrValidation):
			BadRequest(c, "Validation failed", svcErr)
		case errors.Is(svcErr.Unwrap(), service.ErrNotFound):
			NotFound(c, "Resource not found")
		case errors.Is(svcErr.Unwrap(), service.ErrDuplicate):
			Conflict(c, "Resource already exists")
		case errors.Is(svcErr.Unwrap(), service.ErrUnauthorized):
			Forbidden(c, "You don't have permission to perform this action")
		default:
			InternalServerError(c, "An unexpected error occurred")
		}
		return
	}

	// Generic error handling
	InternalServerError(c, "An unexpected error occurred")
}
