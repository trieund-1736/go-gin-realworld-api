package mocks

import (
	"go-gin-realworld-api/internal/models"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockProfileRepository is a mock implementation of ProfileRepository
type MockProfileRepository struct {
	mock.Mock
}

// CreateProfile mock method
func (m *MockProfileRepository) CreateProfile(db *gorm.DB, profile *models.Profile) error {
	args := m.Called(db, profile)
	return args.Error(0)
}

// FindProfileByUserID mock method
func (m *MockProfileRepository) FindProfileByUserID(db *gorm.DB, userID int64) (*models.Profile, error) {
	args := m.Called(db, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Profile), args.Error(1)
}

// UpdateProfile mock method
func (m *MockProfileRepository) UpdateProfile(db *gorm.DB, profile *models.Profile) error {
	args := m.Called(db, profile)
	return args.Error(0)
}
