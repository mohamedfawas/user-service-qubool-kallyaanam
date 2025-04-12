package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/utils/logger"
	"go.uber.org/zap"
)

const (
	// AuthUserIDKey is the context key for the authenticated user ID
	AuthUserIDKey = "auth_user_id"
)

// AuthError represents authentication related errors
type AuthError struct {
	Code    int
	Message string
}

// Authentication middleware extracts the user ID from the JWT token
// It expects the user ID to be forwarded from the API gateway in the X-User-ID header
func Authentication(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			logger.Warn("Authentication failed: missing X-User-ID header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Authentication required",
				"error":   "Missing authentication information",
			})
			return
		}

		// Validate UUID format
		uid, err := uuid.Parse(userID)
		if err != nil {
			logger.Warn("Authentication failed: invalid user ID format",
				zap.String("user_id", userID),
				zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Authentication failed",
				"error":   "Invalid authentication information",
			})
			return
		}

		// Store user ID in context for handlers to use
		c.Set(AuthUserIDKey, uid)

		// Add user ID to request context for logging
		reqID := c.GetHeader("X-Request-ID")
		logger.Debug("Authenticated request",
			zap.String("user_id", uid.String()),
			zap.String("request_id", reqID),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method))

		c.Next()
	}
}

// GetAuthenticatedUserID retrieves the authenticated user ID from the context
func GetAuthenticatedUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get(AuthUserIDKey)
	if !exists {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user ID has invalid type")
	}

	return uid, nil
}
