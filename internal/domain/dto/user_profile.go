package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/domain/model"
)

// CreateUserProfileRequest represents the request payload for creating a user profile
type CreateUserProfileRequest struct {
	IsGroom                bool    `json:"is_groom" binding:"required"`
	ProfileCreatedBy       string  `json:"profile_created_by" binding:"required,oneof=Self Brother Sister Parents Friend Relative"`
	Name                   string  `json:"name" binding:"required,min=2,max=100"`
	DateOfBirth            string  `json:"date_of_birth" binding:"required,datetime=2006-01-02"`
	Community              string  `json:"community" binding:"required,oneof=A\ muslim Hanafi Salafi Sunni Thableegh Shia Jamat\ Islami"`
	Nationality            string  `json:"nationality" binding:"required,oneof=India UAE UK USA"`
	Height                 float64 `json:"height" binding:"required,gt=0"`
	Weight                 float64 `json:"weight" binding:"required,gt=0"`
	MaritalStatus          string  `json:"marital_status" binding:"required,oneof=Never\ married Widower Divorced Nikah\ Divorce"`
	IsPhysicallyChallenged bool    `json:"is_physically_challenged"`
	HomeDistrict           string  `json:"home_district" binding:"required,min=2,max=50"`
}

// UserProfileResponse represents the response after creating a user profile
type UserProfileResponse struct {
	ID                     uuid.UUID `json:"id"`
	UserID                 uuid.UUID `json:"user_id"`
	IsGroom                bool      `json:"is_groom"`
	ProfileCreatedBy       string    `json:"profile_created_by"`
	Name                   string    `json:"name"`
	DateOfBirth            string    `json:"date_of_birth"`
	Age                    int       `json:"age"`
	Community              string    `json:"community"`
	Nationality            string    `json:"nationality"`
	Height                 float64   `json:"height"`
	Weight                 float64   `json:"weight"`
	MaritalStatus          string    `json:"marital_status"`
	IsPhysicallyChallenged bool      `json:"is_physically_challenged"`
	HomeDistrict           string    `json:"home_district"`
	CreatedAt              time.Time `json:"created_at"`
}

// ToModel converts the DTO to a model.UserProfile
func (req *CreateUserProfileRequest) ToModel(userID uuid.UUID) (*model.UserProfile, error) {
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, err
	}

	return &model.UserProfile{
		UserID:                 userID,
		IsGroom:                req.IsGroom,
		ProfileCreatedBy:       model.ProfileCreatedBy(req.ProfileCreatedBy),
		Name:                   req.Name,
		DateOfBirth:            dob,
		Community:              model.Community(req.Community),
		Nationality:            model.Nationality(req.Nationality),
		Height:                 req.Height,
		Weight:                 req.Weight,
		MaritalStatus:          model.MaritalStatus(req.MaritalStatus),
		IsPhysicallyChallenged: req.IsPhysicallyChallenged,
		HomeDistrict:           model.HomeDistrict(req.HomeDistrict),
	}, nil
}

// FromModel creates a UserProfileResponse from a model.UserProfile
func FromModel(profile *model.UserProfile) *UserProfileResponse {
	return &UserProfileResponse{
		ID:                     profile.ID,
		UserID:                 profile.UserID,
		IsGroom:                profile.IsGroom,
		ProfileCreatedBy:       string(profile.ProfileCreatedBy),
		Name:                   profile.Name,
		DateOfBirth:            profile.DateOfBirth.Format("2006-01-02"),
		Age:                    profile.Age(),
		Community:              string(profile.Community),
		Nationality:            string(profile.Nationality),
		Height:                 profile.Height,
		Weight:                 profile.Weight,
		MaritalStatus:          string(profile.MaritalStatus),
		IsPhysicallyChallenged: profile.IsPhysicallyChallenged,
		HomeDistrict:           string(profile.HomeDistrict),
		CreatedAt:              profile.CreatedAt,
	}
}
