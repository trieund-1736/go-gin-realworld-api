package mocks

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockTagRepository is a mock implementation of TagRepository
type MockTagRepository struct {
	mock.Mock
}

// GetAllTags mock method
func (m *MockTagRepository) GetAllTags(db *gorm.DB) ([]string, error) {
	args := m.Called(db)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}
