package service

import (
	"context"
	"errors"
	"testing"

	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/services"
	"go-gin-realworld-api/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper function to setup test dependencies
func setupAuthServiceTest(t *testing.T) (context.Context, *services.AuthService, *mocks.MockUserRepository) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockDB, _ := CreateMockDB(t)
	authService := services.NewAuthService(mockDB, mockUserRepo)
	ctxForTest := context.Background()

	return ctxForTest, authService, mockUserRepo
}

func TestAuthService_Login_Success(t *testing.T) {
	// 1. Setup test dependencies
	ctxForTest, authService, mockUserRepo := setupAuthServiceTest(t)
	email := "test@example.com"
	password := "password123"
	hashedPassword := HashPassword(password)

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
	hashedPassword := HashPassword(correctPassword)

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
