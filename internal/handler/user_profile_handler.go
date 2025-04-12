package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/domain/dto"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/middleware"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/service"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/utils/logger"
	"go.uber.org/zap"
)

// UserProfileHandler handles HTTP requests for user profiles
type UserProfileHandler struct {
	profileService service.UserProfileService
	logger         *logger.Logger
}

// NewUserProfileHandler creates a new user profile handler
func NewUserProfileHandler(profileService service.UserProfileService, logger *logger.Logger) *UserProfileHandler {
	return &UserProfileHandler{
		profileService: profileService,
		logger:         logger,
	}
}

// RegisterRoutes registers the user profile routes
func (h *UserProfileHandler) RegisterRoutes(router *gin.RouterGroup) {
	profileRoutes := router.Group("/profile")
	{
		// POST /user/profile - Create a user profile
		profileRoutes.POST("", h.CreateProfile)

	}
}

func (h *UserProfileHandler) CreateProfile(c *gin.Context) {
	// Get the authenticated user ID from context
	userID, err := middleware.GetAuthenticatedUserID(c)
	if err != nil {
		h.logger.Warn("Failed to get authenticated user ID", zap.Error(err))
		Unauthorized(c, "Authentication required")
		return
	}

	// Parse request body
	var req dto.CreateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		BadRequest(c, "Invalid request body", err)
		return
	}

	// Call service to create profile
	profile, err := h.profileService.CreateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		h.logger.Error("Failed to create profile",
			zap.String("user_id", userID.String()),
			zap.Error(err))

		var svcErr *service.ServiceError
		if err != nil {
			if errors.As(err, &svcErr) && errors.Is(svcErr.Unwrap(), service.ErrDuplicate) {
				Conflict(c, "A profile already exists for this user")
				return
			}
		}

		HandleServiceError(c, err, "CreateProfile")
		return
	}

	// Return successful response
	h.logger.Info("Profile created",
		zap.String("user_id", userID.String()),
		zap.String("profile_id", profile.ID.String()))
	Created(c, "Profile created successfully", profile)
}
