package mysql

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type MySqlTagRepository struct {
}

func NewMySqlTagRepository() *MySqlTagRepository {
	return &MySqlTagRepository{}
}

// GetAllTags retrieves all unique tags from the database
func (r *MySqlTagRepository) GetAllTags(db *gorm.DB) ([]string, error) {
	var tags []string
	if err := db.
		Model(&models.Tag{}).
		Distinct().
		Pluck("name", &tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
