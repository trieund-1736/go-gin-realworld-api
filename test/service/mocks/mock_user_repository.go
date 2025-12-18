package mocks

import (
	"go-gin-realworld-api/internal/models"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

// CreateUser mock method
func (m *MockUserRepository) CreateUser(db *gorm.DB, user *models.User) error {
	args := m.Called(db, user)
	return args.Error(0)
}

// FindUserByEmail mock method
func (m *MockUserRepository) FindUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	args := m.Called(db, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// FindUserByID mock method
func (m *MockUserRepository) FindUserByID(db *gorm.DB, id int64) (*models.User, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// FindUserByUsername mock method
func (m *MockUserRepository) FindUserByUsername(db *gorm.DB, username string, withProfile ...bool) (*models.User, error) {
	args := m.Called(db, username, withProfile)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// UpdateUser mock method
func (m *MockUserRepository) UpdateUser(db *gorm.DB, user *models.User) error {
	args := m.Called(db, user)
	return args.Error(0)
}
