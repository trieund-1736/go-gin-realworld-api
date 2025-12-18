package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(db *gorm.DB, user *models.User) error
	FindUserByEmail(db *gorm.DB, email string) (*models.User, error)
	FindUserByID(db *gorm.DB, id int64) (*models.User, error)
	FindUserByUsername(db *gorm.DB, username string, withProfile ...bool) (*models.User, error)
	UpdateUser(db *gorm.DB, user *models.User) error
}
