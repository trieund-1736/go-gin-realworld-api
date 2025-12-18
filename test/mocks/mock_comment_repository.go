package mocks

import (
	"go-gin-realworld-api/internal/models"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockCommentRepository is a mock implementation of CommentRepository
type MockCommentRepository struct {
	mock.Mock
}

// CreateComment mock method
func (m *MockCommentRepository) CreateComment(db *gorm.DB, comment *models.Comment) error {
	args := m.Called(db, comment)
	return args.Error(0)
}

// GetCommentsByArticleID mock method
func (m *MockCommentRepository) GetCommentsByArticleID(db *gorm.DB, articleID int64) ([]*models.Comment, error) {
	args := m.Called(db, articleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Comment), args.Error(1)
}

// GetCommentByID mock method
func (m *MockCommentRepository) GetCommentByID(db *gorm.DB, id int64) (*models.Comment, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

// DeleteComment mock method
func (m *MockCommentRepository) DeleteComment(db *gorm.DB, id int64) error {
	args := m.Called(db, id)
	return args.Error(0)
}
