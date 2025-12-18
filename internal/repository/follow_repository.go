package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type FollowRepository interface {
	CreateFollow(db *gorm.DB, follow *models.Follow) error
	DeleteFollow(db *gorm.DB, followerID, followeeID int64) error
	IsFollowing(db *gorm.DB, followerID, followeeID int64) (bool, error)
}
