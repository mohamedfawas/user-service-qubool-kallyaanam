package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/domain/dto"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/repository"
)

// UserProfileService defines operations available for user profiles
type UserProfileService interface {
	// CreateProfile creates a new user profile
	CreateProfile(ctx context.Context, userID uuid.UUID, req *dto.CreateUserProfileRequest) (*dto.UserProfileResponse, error)

	// GetProfileByID retrieves a profile by ID
	GetProfileByID(ctx context.Context, profileID uuid.UUID, requestingUserID uuid.UUID) (*dto.UserProfileResponse, error)

	// GetProfileByUserID retrieves a profile by user ID
	GetProfileByUserID(ctx context.Context, userID uuid.UUID, requestingUserID uuid.UUID) (*dto.UserProfileResponse, error)

	// UpdateProfile updates an existing profile
	UpdateProfile(ctx context.Context, userID uuid.UUID, profileID uuid.UUID, req *dto.CreateUserProfileRequest) (*dto.UserProfileResponse, error)

	// DeleteProfile deletes a profile
	DeleteProfile(ctx context.Context, userID uuid.UUID, profileID uuid.UUID) error

	// SearchProfiles searches for profiles based on criteria
	SearchProfiles(ctx context.Context, filter repository.ProfileFilter, page, limit int) ([]*dto.UserProfileResponse, int64, error)
}
