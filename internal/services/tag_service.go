package services

import "go-gin-realworld-api/internal/repository"

type TagService struct {
	tagRepo *repository.TagRepository
}

func NewTagService(tagRepo *repository.TagRepository) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

// GetAllTags retrieves all unique tags
func (s *TagService) GetAllTags() ([]string, error) {
	return s.tagRepo.GetAllTags()
}
