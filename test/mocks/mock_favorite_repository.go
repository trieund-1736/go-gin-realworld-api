package mocks

import (
	"go-gin-realworld-api/internal/models"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockFavoriteRepository is a mock implementation of FavoriteRepository
type MockFavoriteRepository struct {
	mock.Mock
}

// AddFavorite mock method
func (m *MockFavoriteRepository) AddFavorite(db *gorm.DB, userID, articleID int64) error {
	args := m.Called(db, userID, articleID)
	return args.Error(0)
}

// RemoveFavorite mock method
func (m *MockFavoriteRepository) RemoveFavorite(db *gorm.DB, userID, articleID int64) error {
	args := m.Called(db, userID, articleID)
	return args.Error(0)
}

// IsFavorited mock method
func (m *MockFavoriteRepository) IsFavorited(db *gorm.DB, userID, articleID int64) (bool, error) {
	args := m.Called(db, userID, articleID)
	return args.Bool(0), args.Error(1)
}

// GetArticleWithFavorites mock method
func (m *MockFavoriteRepository) GetArticleWithFavorites(db *gorm.DB, articleID int64) (*models.Article, error) {
	args := m.Called(db, articleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Article), args.Error(1)
}
