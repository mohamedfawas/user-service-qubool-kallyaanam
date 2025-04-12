package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProfileCreatedBy represents who created the profile
type ProfileCreatedBy string

// Community represents the religious community of the profile
type Community string

// Nationality represents the nationality of the profile
type Nationality string

// MaritalStatus represents the marital status of the profile
type MaritalStatus string

// HomeDistrict represents the home district of the profile in Kerala
type HomeDistrict string

// Enum values for ProfileCreatedBy
const (
	ProfileCreatedBySelf     ProfileCreatedBy = "Self"
	ProfileCreatedByBrother  ProfileCreatedBy = "Brother"
	ProfileCreatedBySister   ProfileCreatedBy = "Sister"
	ProfileCreatedByParents  ProfileCreatedBy = "Parents"
	ProfileCreatedByFriend   ProfileCreatedBy = "Friend"
	ProfileCreatedByRelative ProfileCreatedBy = "Relative"
)

// Enum values for Community
const (
	CommunityAMuslim     Community = "A muslim"
	CommunityHanafi      Community = "Hanafi"
	CommunitySalafi      Community = "Salafi"
	CommunitySunni       Community = "Sunni"
	CommunityThableegh   Community = "Thableegh"
	CommunityShia        Community = "Shia"
	CommunityJamatIslami Community = "Jamat Islami"
)

// Enum values for Nationality
const (
	NationalityIndia Nationality = "India"
	NationalityUAE   Nationality = "UAE"
	NationalityUK    Nationality = "UK"
	NationalityUSA   Nationality = "USA"
)

// Enum values for MaritalStatus
const (
	MaritalStatusNeverMarried MaritalStatus = "Never married"
	MaritalStatusWidower      MaritalStatus = "Widower"
	MaritalStatusDivorced     MaritalStatus = "Divorced"
	MaritalStatusNikahDivorce MaritalStatus = "Nikah Divorce"
)

// Enum values for HomeDistrict (the 14 districts of Kerala)
const (
	HomeDistrictThiruvananthapuram HomeDistrict = "Thiruvananthapuram"
	HomeDistrictKollam             HomeDistrict = "Kollam"
	HomeDistrictPathanamthitta     HomeDistrict = "Pathanamthitta"
	HomeDistrictAlappuzha          HomeDistrict = "Alappuzha"
	HomeDistrictKottayam           HomeDistrict = "Kottayam"
	HomeDistrictIdukki             HomeDistrict = "Idukki"
	HomeDistrictErnakulam          HomeDistrict = "Ernakulam"
	HomeDistrictThrissur           HomeDistrict = "Thrissur"
	HomeDistrictPalakkad           HomeDistrict = "Palakkad"
	HomeDistrictMalappuram         HomeDistrict = "Malappuram"
	HomeDistrictKozhikode          HomeDistrict = "Kozhikode"
	HomeDistrictWayanad            HomeDistrict = "Wayanad"
	HomeDistrictKannur             HomeDistrict = "Kannur"
	HomeDistrictKasaragod          HomeDistrict = "Kasaragod"
)

// UserProfile represents the profile information of a user in the matrimony platform
type UserProfile struct {
	ID                     uuid.UUID        `gorm:"type:uuid;primary_key" json:"id"`
	UserID                 uuid.UUID        `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	IsGroom                bool             `gorm:"not null" json:"is_groom"`
	ProfileCreatedBy       ProfileCreatedBy `gorm:"type:profile_created_by;not null" json:"profile_created_by"`
	Name                   string           `gorm:"type:varchar(100);not null" json:"name"`
	DateOfBirth            time.Time        `gorm:"type:date;not null" json:"date_of_birth"`
	Community              Community        `gorm:"type:community_type;not null" json:"community"`
	Nationality            Nationality      `gorm:"type:nationality_type;not null" json:"nationality"`
	Height                 float64          `gorm:"type:decimal(5,2);not null" json:"height"`
	Weight                 float64          `gorm:"type:decimal(5,2);not null" json:"weight"`
	MaritalStatus          MaritalStatus    `gorm:"type:marital_status_type;not null" json:"marital_status"`
	IsPhysicallyChallenged bool             `gorm:"not null;default:false" json:"is_physically_challenged"`
	HomeDistrict           HomeDistrict     `gorm:"type:home_district_type;not null" json:"home_district"`
	CreatedAt              time.Time        `gorm:"not null" json:"created_at"`
	UpdatedAt              time.Time        `gorm:"not null" json:"updated_at"`
	DeletedAt              gorm.DeletedAt   `gorm:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (up *UserProfile) BeforeCreate(tx *gorm.DB) error {
	if up.ID == uuid.Nil {
		up.ID = uuid.New()
	}
	return nil
}

// Age calculates the user's age based on the date of birth
func (up *UserProfile) Age() int {
	now := time.Now()
	years := now.Year() - up.DateOfBirth.Year()

	// Adjust age if birthday hasn't occurred yet this year
	birthMonth, birthDay := up.DateOfBirth.Month(), up.DateOfBirth.Day()
	currentMonth, currentDay := now.Month(), now.Day()

	if currentMonth < birthMonth || (currentMonth == birthMonth && currentDay < birthDay) {
		years--
	}

	return years
}

// TableName specifies the table name for UserProfile model
func (UserProfile) TableName() string {
	return "user_profiles"
}
