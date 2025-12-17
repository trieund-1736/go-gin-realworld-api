package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type TagRepository struct {
}

func NewTagRepository() *TagRepository {
	return &TagRepository{}
}

// GetAllTags retrieves all unique tags from the database
func (r *TagRepository) GetAllTags(db *gorm.DB) ([]string, error) {
	var tags []string
	if err := db.
		Model(&models.Tag{}).
		Distinct().
		Pluck("name", &tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
