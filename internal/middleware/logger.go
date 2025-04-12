package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/utils/logger"
	"go.uber.org/zap"
)

// RequestLogger middleware logs HTTP requests
func RequestLogger(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Get or generate request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Header("X-Request-ID", requestID)
		}

		// Add request ID to context for later use in logging
		c.Set("request_id", requestID)

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// Log request details
		logFields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("request_id", requestID),
			zap.String("client_ip", c.ClientIP()),
		}

		switch {
		case statusCode >= 500:
			logger.Error("Server error", logFields...)
		case statusCode >= 400:
			logger.Warn("Client error", logFields...)
		default:
			logger.Info("Request completed", logFields...)
		}
	}
}
