package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/domain/model"
)

// UserProfileRepository defines operations for working with user profiles
type UserProfileRepository interface {
	// Create adds a new user profile
	Create(ctx context.Context, profile *model.UserProfile) error

	// GetByID retrieves a profile by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*model.UserProfile, error)

	// GetByUserID retrieves a profile by the associated user ID
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.UserProfile, error)

	// Update updates an existing profile
	Update(ctx context.Context, profile *model.UserProfile) error

	// Delete soft-deletes a profile
	Delete(ctx context.Context, id uuid.UUID) error

	// SearchProfiles searches for profiles with pagination based on filter criteria
	SearchProfiles(ctx context.Context, filter ProfileFilter, page, limit int) ([]*model.UserProfile, int64, error)

	// WithTransaction executes operations within a database transaction
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

// ProfileFilter defines criteria for filtering profiles
type ProfileFilter struct {
	IsGroom                *bool
	Community              []model.Community
	Nationality            []model.Nationality
	MaritalStatus          []model.MaritalStatus
	HomeDistrict           []model.HomeDistrict
	MinAge                 *int
	MaxAge                 *int
	MinHeight              *float64
	MaxHeight              *float64
	IsPhysicallyChallenged *bool
	CreatedAfter           *time.Time
	CreatedBefore          *time.Time
}
