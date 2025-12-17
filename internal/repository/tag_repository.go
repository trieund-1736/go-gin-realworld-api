package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// GetAllTags retrieves all unique tags from the database
func (r *TagRepository) GetAllTags() ([]string, error) {
	var tags []string
	if err := r.db.
		Model(&models.Tag{}).
		Distinct().
		Pluck("name", &tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
