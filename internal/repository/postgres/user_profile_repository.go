package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/domain/model"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/repository"
	"gorm.io/gorm"
)

const (
	entityUserProfile = "UserProfile"
)

// UserProfileRepository implements repository.UserProfileRepository for PostgreSQL
type UserProfileRepository struct {
	db *gorm.DB
}

// NewUserProfileRepository creates a new UserProfileRepository
func NewUserProfileRepository(db *gorm.DB) repository.UserProfileRepository {
	return &UserProfileRepository{
		db: db,
	}
}

// Create adds a new user profile to the database
func (r *UserProfileRepository) Create(ctx context.Context, profile *model.UserProfile) error {
	const op = "Create"

	if profile.UserID == uuid.Nil {
		return repository.NewError(repository.ErrInvalidOperation, op, entityUserProfile, "user_id is required")
	}

	err := r.db.WithContext(ctx).Create(profile).Error
	if err != nil {
		// Check for duplicate key violation
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "unique_user_profile" {
				return repository.NewError(repository.ErrDuplicateKey, op, entityUserProfile, "profile already exists for this user")
			}
		}
		return repository.NewError(err, op, entityUserProfile, "")
	}

	return nil
}

// GetByID retrieves a user profile by ID
func (r *UserProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.UserProfile, error) {
	const op = "GetByID"

	var profile model.UserProfile
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&profile).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewError(repository.ErrNotFound, op, entityUserProfile, fmt.Sprintf("id: %s", id))
		}
		return nil, repository.NewError(err, op, entityUserProfile, "")
	}

	return &profile, nil
}

// GetByUserID retrieves a user profile by user ID
func (r *UserProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.UserProfile, error) {
	const op = "GetByUserID"

	var profile model.UserProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewError(repository.ErrNotFound, op, entityUserProfile, fmt.Sprintf("user_id: %s", userID))
		}
		return nil, repository.NewError(err, op, entityUserProfile, "")
	}

	return &profile, nil
}

// Update updates an existing user profile
func (r *UserProfileRepository) Update(ctx context.Context, profile *model.UserProfile) error {
	const op = "Update"

	if profile.ID == uuid.Nil {
		return repository.NewError(repository.ErrInvalidOperation, op, entityUserProfile, "id is required")
	}

	// Set updated_at to current time
	profile.UpdatedAt = time.Now()

	err := r.db.WithContext(ctx).Save(profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.NewError(repository.ErrNotFound, op, entityUserProfile, fmt.Sprintf("id: %s", profile.ID))
		}
		return repository.NewError(err, op, entityUserProfile, "")
	}

	return nil
}

// Delete soft-deletes a user profile
func (r *UserProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "Delete"

	result := r.db.WithContext(ctx).Delete(&model.UserProfile{}, id)
	if result.Error != nil {
		return repository.NewError(result.Error, op, entityUserProfile, "")
	}

	if result.RowsAffected == 0 {
		return repository.NewError(repository.ErrNotFound, op, entityUserProfile, fmt.Sprintf("id: %s", id))
	}

	return nil
}

// SearchProfiles searches for profiles with pagination based on filter criteria
func (r *UserProfileRepository) SearchProfiles(
	ctx context.Context,
	filter repository.ProfileFilter,
	page,
	limit int,
) ([]*model.UserProfile, int64, error) {
	const op = "SearchProfiles"

	var profiles []*model.UserProfile
	var total int64

	query := r.db.WithContext(ctx).Model(&model.UserProfile{})

	// Apply filters
	if filter.IsGroom != nil {
		query = query.Where("is_groom = ?", *filter.IsGroom)
	}

	if len(filter.Community) > 0 {
		query = query.Where("community IN ?", filter.Community)
	}

	if len(filter.Nationality) > 0 {
		query = query.Where("nationality IN ?", filter.Nationality)
	}

	if len(filter.MaritalStatus) > 0 {
		query = query.Where("marital_status IN ?", filter.MaritalStatus)
	}

	if len(filter.HomeDistrict) > 0 {
		query = query.Where("home_district IN ?", filter.HomeDistrict)
	}

	if filter.MinAge != nil || filter.MaxAge != nil {
		now := time.Now()
		if filter.MinAge != nil {
			maxDOB := now.AddDate(-*filter.MinAge, 0, 0)
			query = query.Where("date_of_birth <= ?", maxDOB)
		}
		if filter.MaxAge != nil {
			minDOB := now.AddDate(-*filter.MaxAge, 0, 0)
			query = query.Where("date_of_birth >= ?", minDOB)
		}
	}

	if filter.MinHeight != nil {
		query = query.Where("height >= ?", *filter.MinHeight)
	}

	if filter.MaxHeight != nil {
		query = query.Where("height <= ?", *filter.MaxHeight)
	}

	if filter.IsPhysicallyChallenged != nil {
		query = query.Where("is_physically_challenged = ?", *filter.IsPhysicallyChallenged)
	}

	if filter.CreatedAfter != nil {
		query = query.Where("created_at >= ?", *filter.CreatedAfter)
	}

	if filter.CreatedBefore != nil {
		query = query.Where("created_at <= ?", *filter.CreatedBefore)
	}

	// Count total matches before applying pagination
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, repository.NewError(err, op, entityUserProfile, "count failed")
	}

	// Apply pagination
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Execute query with pagination
	err = query.Offset(offset).Limit(limit).Find(&profiles).Error
	if err != nil {
		return nil, 0, repository.NewError(err, op, entityUserProfile, "")
	}

	return profiles, total, nil
}

// WithTransaction executes operations within a database transaction
func (r *UserProfileRepository) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	const op = "WithTransaction"

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new context with the transaction
		txCtx := context.WithValue(ctx, "tx", tx)

		// Execute the function with the transaction context
		err := fn(txCtx)
		if err != nil {
			return repository.NewError(err, op, entityUserProfile, "transaction failed")
		}

		return nil
	})
}
