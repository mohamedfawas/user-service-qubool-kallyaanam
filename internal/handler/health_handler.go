package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/utils/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB, logger *logger.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		logger: logger,
	}
}

// RegisterRoutes registers the health check routes
func (h *HealthHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.HealthCheck)
	router.GET("/readiness", h.ReadinessCheck)
}

// HealthCheck handles basic health check requests
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "UP",
		"service":   "user-service",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// ReadinessCheck performs a more thorough health check including DB connection
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		h.logger.Error("Database connection error during readiness check", zap.Error(err))
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "DOWN",
			"service":   "user-service",
			"database":  "DOWN",
			"error":     "Database connection error",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	// Ping the database
	err = sqlDB.Ping()
	if err != nil {
		h.logger.Error("Database ping failed during readiness check", zap.Error(err))
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "DOWN",
			"service":   "user-service",
			"database":  "DOWN",
			"error":     "Database ping failed",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "UP",
		"service":   "user-service",
		"database":  "UP",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
