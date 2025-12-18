package repository

import (
	"gorm.io/gorm"
)

type TagRepository interface {
	GetAllTags(db *gorm.DB) ([]string, error)
}
