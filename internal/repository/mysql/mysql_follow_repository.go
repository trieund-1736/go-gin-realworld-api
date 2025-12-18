package mysql

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type MySqlFollowRepository struct {
}

func NewMySqlFollowRepository() *MySqlFollowRepository {
	return &MySqlFollowRepository{}
}

// CreateFollow creates a follow relationship
func (r *MySqlFollowRepository) CreateFollow(db *gorm.DB, follow *models.Follow) error {
	if err := db.Create(follow).Error; err != nil {
		return err
	}
	return nil
}

// DeleteFollow deletes a follow relationship
func (r *MySqlFollowRepository) DeleteFollow(db *gorm.DB, followerID, followeeID int64) error {
	if err := db.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&models.Follow{}).Error; err != nil {
		return err
	}
	return nil
}

// IsFollowing checks if a user follows another user
func (r *MySqlFollowRepository) IsFollowing(db *gorm.DB, followerID, followeeID int64) (bool, error) {
	var count int64
	if err := db.Model(&models.Follow{}).Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
