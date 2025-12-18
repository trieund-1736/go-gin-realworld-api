package mocks

import (
	"go-gin-realworld-api/internal/models"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockArticleRepository is a mock implementation of ArticleRepository
type MockArticleRepository struct {
	mock.Mock
}

// ListArticles mock method
func (m *MockArticleRepository) ListArticles(db *gorm.DB, tag, author string, favorited *bool, currentUserID *int64, limit, offset int) ([]*models.Article, int64, error) {
	args := m.Called(db, tag, author, favorited, currentUserID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.Article), args.Get(1).(int64), args.Error(2)
}

// FeedArticles mock method
func (m *MockArticleRepository) FeedArticles(db *gorm.DB, userID int64, limit, offset int) ([]*models.Article, int64, error) {
	args := m.Called(db, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.Article), args.Get(1).(int64), args.Error(2)
}

// FindArticleBySlug mock method
func (m *MockArticleRepository) FindArticleBySlug(db *gorm.DB, slug string) (*models.Article, error) {
	args := m.Called(db, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Article), args.Error(1)
}

// CreateArticle mock method
func (m *MockArticleRepository) CreateArticle(db *gorm.DB, article *models.Article) error {
	args := m.Called(db, article)
	return args.Error(0)
}

// UpdateArticle mock method
func (m *MockArticleRepository) UpdateArticle(db *gorm.DB, article *models.Article) error {
	args := m.Called(db, article)
	return args.Error(0)
}

// DeleteArticleBySlug mock method
func (m *MockArticleRepository) DeleteArticleBySlug(db *gorm.DB, slug string) error {
	args := m.Called(db, slug)
	return args.Error(0)
}

// AssignTagsToArticle mock method
func (m *MockArticleRepository) AssignTagsToArticle(db *gorm.DB, articleID int64, tagNames []string) error {
	args := m.Called(db, articleID, tagNames)
	return args.Error(0)
}
