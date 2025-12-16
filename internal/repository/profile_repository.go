package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// CreateProfile creates a new profile in the database
func (r *ProfileRepository) CreateProfile(profile *models.Profile) error {
	if err := r.db.Create(profile).Error; err != nil {
		return err
	}
	return nil
}

// FindProfileByUserID finds a profile by user ID
func (r *ProfileRepository) FindProfileByUserID(userID int64) (*models.Profile, error) {
	var profile *models.Profile
	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

// UpdateProfile updates a profile in the database
func (r *ProfileRepository) UpdateProfile(profile *models.Profile) error {
	if err := r.db.Save(profile).Error; err != nil {
		return err
	}
	return nil
}
