package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type FollowRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

// CreateFollow creates a follow relationship
func (r *FollowRepository) CreateFollow(follow *models.Follow) error {
	if err := r.db.Create(follow).Error; err != nil {
		return err
	}
	return nil
}

// DeleteFollow deletes a follow relationship
func (r *FollowRepository) DeleteFollow(followerID, followeeID int64) error {
	if err := r.db.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&models.Follow{}).Error; err != nil {
		return err
	}
	return nil
}

// IsFollowing checks if a user follows another user
func (r *FollowRepository) IsFollowing(followerID, followeeID int64) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Follow{}).Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
