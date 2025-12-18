package mocks

import (
	"go-gin-realworld-api/internal/models"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockFollowRepository is a mock implementation of FollowRepository
type MockFollowRepository struct {
	mock.Mock
}

// CreateFollow mock method
func (m *MockFollowRepository) CreateFollow(db *gorm.DB, follow *models.Follow) error {
	args := m.Called(db, follow)
	return args.Error(0)
}

// DeleteFollow mock method
func (m *MockFollowRepository) DeleteFollow(db *gorm.DB, followerID, followeeID int64) error {
	args := m.Called(db, followerID, followeeID)
	return args.Error(0)
}

// IsFollowing mock method
func (m *MockFollowRepository) IsFollowing(db *gorm.DB, followerID, followeeID int64) (bool, error) {
	args := m.Called(db, followerID, followeeID)
	return args.Bool(0), args.Error(1)
}
