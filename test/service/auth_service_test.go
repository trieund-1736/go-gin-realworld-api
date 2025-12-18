package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"testing"

	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/services"
	"go-gin-realworld-api/test/service/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Helper function to hash password (same as in auth_service.go)
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// Helper function to create a mock DB
func createMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	return gormDB, mock
}

// Helper function to setup test dependencies
func setupAuthServiceTest(t *testing.T) (context.Context, *services.AuthService, *mocks.MockUserRepository) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockDB, _ := createMockDB(t)
	authService := services.NewAuthService(mockDB, mockUserRepo)
	ctxForTest := context.Background()

	return ctxForTest, authService, mockUserRepo
}

func TestAuthService_Login_Success(t *testing.T) {
	// 1. Setup test dependencies
	ctxForTest, authService, mockUserRepo := setupAuthServiceTest(t)
	email := "test@example.com"
	password := "password123"
	hashedPassword := hashPassword(password)

	expectedUser := &models.User{
		ID:       1,
		Username: "testuser",
		Email:    email,
		Password: hashedPassword,
	}

	// 2. Define mock behavior: When FindUserByEmail is called,
	// it will return expectedUser and no error
	mockUserRepo.On("FindUserByEmail", mock.Anything, email).Return(expectedUser, nil)

	// 3. Call the service method under test
	user, token, err := authService.Login(ctxForTest, email, password)

	// 4. Assert results
	assert.NoError(t, err) // Check for no error
	assert.NotNil(t, user) // User should not be nil
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.NotEmpty(t, token) // Token should be generated

	// 5. Assert that mock expectations were met
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	// 1. Setup test dependencies
	ctxForTest, authService, mockRepo := setupAuthServiceTest(t)
	email := "notfound@example.com"
	password := "password123"

	// 2. Define mock behavior: FindUserByEmail returns an error
	expectedError := errors.New("user not found")
	mockRepo.On("FindUserByEmail", mock.Anything, email).Return(nil, expectedError)

	// 3. Call the service method under test
	user, token, err := authService.Login(ctxForTest, email, password)

	// 4. Assert results
	assert.Error(t, err)                                // Check that an error occurred
	assert.Equal(t, "invalid credentials", err.Error()) // Should return generic error message
	assert.Nil(t, user)                                 // User should be nil
	assert.Empty(t, token)                              // Token should be empty

	// 5. Assert that mock expectations were met
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	// 1. Setup test dependencies
	ctxForTest, authService, mockRepo := setupAuthServiceTest(t)
	email := "test@example.com"
	correctPassword := "password123"
	wrongPassword := "wrongpassword"
	hashedPassword := hashPassword(correctPassword)

	existingUser := &models.User{
		ID:       1,
		Username: "testuser",
		Email:    email,
		Password: hashedPassword,
	}

	// 2. Define mock behavior: User is found but password is wrong
	mockRepo.On("FindUserByEmail", mock.Anything, email).Return(existingUser, nil)

	// 3. Call the service method with wrong password
	user, token, err := authService.Login(ctxForTest, email, wrongPassword)

	// 4. Assert results
	assert.Error(t, err)                                // Check that an error occurred
	assert.Equal(t, "invalid credentials", err.Error()) // Should return generic error message
	assert.Nil(t, user)                                 // User should be nil
	assert.Empty(t, token)                              // Token should be empty

	// 5. Assert that mock expectations were met
	mockRepo.AssertExpectations(t)
}
