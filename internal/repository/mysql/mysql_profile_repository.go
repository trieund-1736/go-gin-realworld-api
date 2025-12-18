package mysql

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type MySqlProfileRepository struct {
}

func NewMySqlProfileRepository() *MySqlProfileRepository {
	return &MySqlProfileRepository{}
}

// CreateProfile creates a new profile in the database
func (r *MySqlProfileRepository) CreateProfile(db *gorm.DB, profile *models.Profile) error {
	if err := db.Create(profile).Error; err != nil {
		return err
	}
	return nil
}

// FindProfileByUserID finds a profile by user ID
func (r *MySqlProfileRepository) FindProfileByUserID(db *gorm.DB, userID int64) (*models.Profile, error) {
	var profile *models.Profile
	if err := db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

// UpdateProfile updates a profile in the database
func (r *MySqlProfileRepository) UpdateProfile(db *gorm.DB, profile *models.Profile) error {
	if err := db.Save(profile).Error; err != nil {
		return err
	}
	return nil
}
