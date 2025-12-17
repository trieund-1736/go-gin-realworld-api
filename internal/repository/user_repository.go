package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(db *gorm.DB, user *models.User) error {
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// FindUserByEmail finds a user by email
func (r *UserRepository) FindUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user *models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserByID finds a user by ID
func (r *UserRepository) FindUserByID(db *gorm.DB, id int64) (*models.User, error) {
	var user *models.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserByUsername finds a user by username with optional profile preload
func (r *UserRepository) FindUserByUsername(db *gorm.DB, username string, withProfile ...bool) (*models.User, error) {
	var user *models.User
	query := db

	// Check if preload profile flag is set (default: false)
	if len(withProfile) > 0 && withProfile[0] {
		query = query.Preload("Profile")
	}

	if err := query.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates a user in the database
func (r *UserRepository) UpdateUser(db *gorm.DB, user *models.User) error {
	if err := db.Save(user).Error; err != nil {
		return err
	}
	return nil
}
