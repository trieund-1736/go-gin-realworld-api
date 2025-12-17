package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type ProfileRepository struct {
}

func NewProfileRepository() *ProfileRepository {
	return &ProfileRepository{}
}

// CreateProfile creates a new profile in the database
func (r *ProfileRepository) CreateProfile(db *gorm.DB, profile *models.Profile) error {
	if err := db.Create(profile).Error; err != nil {
		return err
	}
	return nil
}

// FindProfileByUserID finds a profile by user ID
func (r *ProfileRepository) FindProfileByUserID(db *gorm.DB, userID int64) (*models.Profile, error) {
	var profile *models.Profile
	if err := db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

// UpdateProfile updates a profile in the database
func (r *ProfileRepository) UpdateProfile(db *gorm.DB, profile *models.Profile) error {
	if err := db.Save(profile).Error; err != nil {
		return err
	}
	return nil
}
