package services

import (
	"context"
	"go-gin-realworld-api/internal/repository"

	"gorm.io/gorm"
)

type TagService struct {
	db      *gorm.DB
	tagRepo repository.TagRepository
}

func NewTagService(db *gorm.DB, tagRepo repository.TagRepository) *TagService {
	return &TagService{
		db:      db,
		tagRepo: tagRepo,
	}
}

// GetAllTags retrieves all unique tags
func (s *TagService) GetAllTags(ctx context.Context) ([]string, error) {
	db := s.db.WithContext(ctx)
	return s.tagRepo.GetAllTags(db)
}
