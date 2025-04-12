package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/domain/dto"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/repository"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/utils/logger"
	"go.uber.org/zap"
)

const (
	serviceName = "UserProfileService"
	minAge      = 18
	maxAge      = 80
)

// userProfileService implements UserProfileService
type userProfileService struct {
	repo   repository.UserProfileRepository
	logger *logger.Logger
}

// NewUserProfileService creates a new user profile service
func NewUserProfileService(repo repository.UserProfileRepository, logger *logger.Logger) UserProfileService {
	return &userProfileService{
		repo:   repo,
		logger: logger,
	}
}

// CreateProfile creates a new user profile
func (s *userProfileService) CreateProfile(
	ctx context.Context,
	userID uuid.UUID,
	req *dto.CreateUserProfileRequest,
) (*dto.UserProfileResponse, error) {
	const op = "CreateProfile"

	// Validate user ID
	if userID == uuid.Nil {
		return nil, NewError(ErrValidation, op, serviceName, "user ID is required")
	}

	// Validate the request
	validationErrors := s.validateProfileRequest(req)
	if len(validationErrors) > 0 {
		return nil, NewValidationError(op, serviceName, validationErrors)
	}

	// Check if profile already exists for this user
	existingProfile, err := s.repo.GetByUserID(ctx, userID)
	if err == nil && existingProfile != nil {
		return nil, NewError(ErrDuplicate, op, serviceName, "profile already exists for this user")
	} else if err != nil && !errors.Is(errors.Unwrap(err), repository.ErrNotFound) {
		s.logger.Error("Failed to check for existing profile",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return nil, NewError(ErrInternal, op, serviceName, "failed to check for existing profile")
	}

	// Convert request to model
	profile, err := req.ToModel(userID)
	if err != nil {
		s.logger.Error("Failed to convert request to model",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return nil, NewError(ErrValidation, op, serviceName, err.Error())
	}

	// Create the profile
	err = s.repo.Create(ctx, profile)
	if err != nil {
		s.logger.Error("Failed to create profile",
			zap.String("user_id", userID.String()),
			zap.Error(err))

		var repoErr *repository.RepositoryError
		if errors.As(err, &repoErr) && errors.Is(repoErr.Unwrap(), repository.ErrDuplicateKey) {
			return nil, NewError(ErrDuplicate, op, serviceName, "profile already exists for this user")
		}

		return nil, NewError(ErrInternal, op, serviceName, "failed to create profile")
	}

	// Log success
	s.logger.UserProfileEvent(ctx, "profile_created", userID.String(), profile.ID.String(),
		zap.String("name", profile.Name),
		zap.Bool("is_groom", profile.IsGroom))

	// Return response
	return dto.FromModel(profile), nil
}

// GetProfileByID retrieves a profile by ID
func (s *userProfileService) GetProfileByID(
	ctx context.Context,
	profileID uuid.UUID,
	requestingUserID uuid.UUID,
) (*dto.UserProfileResponse, error) {
	const op = "GetProfileByID"

	// Get the profile
	profile, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		s.logger.Error("Failed to get profile by ID",
			zap.String("profile_id", profileID.String()),
			zap.Error(err))

		var repoErr *repository.RepositoryError
		if errors.As(err, &repoErr) && errors.Is(repoErr.Unwrap(), repository.ErrNotFound) {
			return nil, NewError(ErrNotFound, op, serviceName, fmt.Sprintf("profile with ID %s not found", profileID))
		}

		return nil, NewError(ErrInternal, op, serviceName, "failed to retrieve profile")
	}

	// Check if user is authorized to view this profile
	// In this case, we're allowing any authenticated user to view profiles
	// but logging the access for auditing purposes
	if profile.UserID != requestingUserID {
		s.logger.UserProfileEvent(ctx, "profile_accessed", requestingUserID.String(), profileID.String(),
			zap.String("owner_id", profile.UserID.String()))
	} else {
		s.logger.UserProfileEvent(ctx, "profile_accessed_by_owner", profile.UserID.String(), profileID.String())
	}

	return dto.FromModel(profile), nil
}

// GetProfileByUserID retrieves a profile by user ID
func (s *userProfileService) GetProfileByUserID(
	ctx context.Context,
	userID uuid.UUID,
	requestingUserID uuid.UUID,
) (*dto.UserProfileResponse, error) {
	const op = "GetProfileByUserID"

	// Get the profile
	profile, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get profile by user ID",
			zap.String("user_id", userID.String()),
			zap.Error(err))

		var repoErr *repository.RepositoryError
		if errors.As(err, &repoErr) && errors.Is(repoErr.Unwrap(), repository.ErrNotFound) {
			return nil, NewError(ErrNotFound, op, serviceName, fmt.Sprintf("profile for user %s not found", userID))
		}

		return nil, NewError(ErrInternal, op, serviceName, "failed to retrieve profile")
	}

	// Check if user is authorized to view this profile
	// Similar to GetProfileByID, allowing any authenticated user to view
	if userID != requestingUserID {
		s.logger.UserProfileEvent(ctx, "profile_accessed", requestingUserID.String(), profile.ID.String(),
			zap.String("owner_id", userID.String()))
	} else {
		s.logger.UserProfileEvent(ctx, "profile_accessed_by_owner", userID.String(), profile.ID.String())
	}

	return dto.FromModel(profile), nil
}

// UpdateProfile updates an existing profile
func (s *userProfileService) UpdateProfile(
	ctx context.Context,
	userID uuid.UUID,
	profileID uuid.UUID,
	req *dto.CreateUserProfileRequest,
) (*dto.UserProfileResponse, error) {
	const op = "UpdateProfile"

	// Validate request
	validationErrors := s.validateProfileRequest(req)
	if len(validationErrors) > 0 {
		return nil, NewValidationError(op, serviceName, validationErrors)
	}

	// Get existing profile
	existingProfile, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		s.logger.Error("Failed to get profile for update",
			zap.String("profile_id", profileID.String()),
			zap.Error(err))

		var repoErr *repository.RepositoryError
		if errors.As(err, &repoErr) && errors.Is(repoErr.Unwrap(), repository.ErrNotFound) {
			return nil, NewError(ErrNotFound, op, serviceName, fmt.Sprintf("profile with ID %s not found", profileID))
		}

		return nil, NewError(ErrInternal, op, serviceName, "failed to retrieve profile for update")
	}

	// Check authorization - only the owner can update their profile
	if existingProfile.UserID != userID {
		s.logger.Warn("Unauthorized profile update attempt",
			zap.String("requester_id", userID.String()),
			zap.String("profile_id", profileID.String()),
			zap.String("owner_id", existingProfile.UserID.String()))
		return nil, NewError(ErrUnauthorized, op, serviceName, "you can only update your own profile")
	}

	// Convert request to model
	updatedProfile, err := req.ToModel(userID)
	if err != nil {
		return nil, NewError(ErrValidation, op, serviceName, err.Error())
	}

	// Preserve ID and creation time
	updatedProfile.ID = profileID
	updatedProfile.CreatedAt = existingProfile.CreatedAt
	updatedProfile.UpdatedAt = time.Now()

	// Update profile
	err = s.repo.Update(ctx, updatedProfile)
	if err != nil {
		s.logger.Error("Failed to update profile",
			zap.String("profile_id", profileID.String()),
			zap.Error(err))
		return nil, NewError(ErrInternal, op, serviceName, "failed to update profile")
	}

	// Log the update
	s.logger.UserProfileEvent(ctx, "profile_updated", userID.String(), profileID.String(),
		zap.String("name", updatedProfile.Name))

	return dto.FromModel(updatedProfile), nil
}

// DeleteProfile deletes a profile
func (s *userProfileService) DeleteProfile(
	ctx context.Context,
	userID uuid.UUID,
	profileID uuid.UUID,
) error {
	const op = "DeleteProfile"

	// Get existing profile
	existingProfile, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		s.logger.Error("Failed to get profile for deletion",
			zap.String("profile_id", profileID.String()),
			zap.Error(err))

		var repoErr *repository.RepositoryError
		if errors.As(err, &repoErr) && errors.Is(repoErr.Unwrap(), repository.ErrNotFound) {
			return NewError(ErrNotFound, op, serviceName, fmt.Sprintf("profile with ID %s not found", profileID))
		}

		return NewError(ErrInternal, op, serviceName, "failed to retrieve profile for deletion")
	}

	// Check authorization - only the owner can delete their profile
	if existingProfile.UserID != userID {
		s.logger.Warn("Unauthorized profile deletion attempt",
			zap.String("requester_id", userID.String()),
			zap.String("profile_id", profileID.String()),
			zap.String("owner_id", existingProfile.UserID.String()))
		return NewError(ErrUnauthorized, op, serviceName, "you can only delete your own profile")
	}

	// Delete the profile
	err = s.repo.Delete(ctx, profileID)
	if err != nil {
		s.logger.Error("Failed to delete profile",
			zap.String("profile_id", profileID.String()),
			zap.Error(err))
		return NewError(ErrInternal, op, serviceName, "failed to delete profile")
	}

	// Log the deletion
	s.logger.UserProfileEvent(ctx, "profile_deleted", userID.String(), profileID.String())

	return nil
}

// SearchProfiles searches for profiles based on criteria
func (s *userProfileService) SearchProfiles(
	ctx context.Context,
	filter repository.ProfileFilter,
	page, limit int,
) ([]*dto.UserProfileResponse, int64, error) {
	const op = "SearchProfiles"

	// Search profiles
	profiles, total, err := s.repo.SearchProfiles(ctx, filter, page, limit)
	if err != nil {
		s.logger.Error("Failed to search profiles", zap.Error(err))
		return nil, 0, NewError(ErrInternal, op, serviceName, "failed to search profiles")
	}

	// Convert to DTOs
	results := make([]*dto.UserProfileResponse, len(profiles))
	for i, profile := range profiles {
		results[i] = dto.FromModel(profile)
	}

	return results, total, nil
}

// validateProfileRequest validates the profile request
func (s *userProfileService) validateProfileRequest(req *dto.CreateUserProfileRequest) []ValidationError {
	var errors []ValidationError

	// Validate name
	if len(req.Name) < 2 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name must be at least 2 characters long",
		})
	}

	// Validate date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		errors = append(errors, ValidationError{
			Field:   "date_of_birth",
			Message: "Invalid date format, expected YYYY-MM-DD",
		})
	} else {
		// Check age range
		age := calculateAge(dob)
		if age < minAge {
			errors = append(errors, ValidationError{
				Field:   "date_of_birth",
				Message: fmt.Sprintf("Age must be at least %d years", minAge),
			})
		} else if age > maxAge {
			errors = append(errors, ValidationError{
				Field:   "date_of_birth",
				Message: fmt.Sprintf("Age must not exceed %d years", maxAge),
			})
		}
	}

	// Validate height
	if req.Height <= 0 {
		errors = append(errors, ValidationError{
			Field:   "height",
			Message: "Height must be greater than 0",
		})
	} else if req.Height < 100 || req.Height > 250 {
		errors = append(errors, ValidationError{
			Field:   "height",
			Message: "Height must be between 100 and 250 cm",
		})
	}

	// Validate weight
	if req.Weight <= 0 {
		errors = append(errors, ValidationError{
			Field:   "weight",
			Message: "Weight must be greater than 0",
		})
	} else if req.Weight < 30 || req.Weight > 200 {
		errors = append(errors, ValidationError{
			Field:   "weight",
			Message: "Weight must be between 30 and 200 kg",
		})
	}

	return errors
}

// calculateAge calculates age from date of birth
func calculateAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()

	// Adjust age if birthday hasn't occurred yet this year
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		years--
	}

	return years
}
