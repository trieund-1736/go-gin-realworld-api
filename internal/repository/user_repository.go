package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// FindUserByEmail finds a user by email
func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	var user *models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserByID finds a user by ID
func (r *UserRepository) FindUserByID(id int64) (*models.User, error) {
	var user *models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserByUsername finds a user by username with optional profile preload
func (r *UserRepository) FindUserByUsername(username string, withProfile ...bool) (*models.User, error) {
	var user *models.User
	query := r.db

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
func (r *UserRepository) UpdateUser(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}
