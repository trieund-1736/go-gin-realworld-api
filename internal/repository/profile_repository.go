package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type ProfileRepository interface {
	CreateProfile(db *gorm.DB, profile *models.Profile) error
	FindProfileByUserID(db *gorm.DB, userID int64) (*models.Profile, error)
	UpdateProfile(db *gorm.DB, profile *models.Profile) error
}
