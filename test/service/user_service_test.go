package service

import (
	"context"
	"errors"
	"testing"

	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/services"
	"go-gin-realworld-api/test/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper function to setup test dependencies for UserService
func setupUserServiceTest(t *testing.T) (context.Context, *services.UserService, *mocks.MockUserRepository, *mocks.MockProfileRepository, *mocks.MockFollowRepository, sqlmock.Sqlmock) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockProfileRepo := new(mocks.MockProfileRepository)
	mockFollowRepo := new(mocks.MockFollowRepository)
	gormDB, sqlMock := CreateMockDB(t)
	userService := services.NewUserService(gormDB, mockUserRepo, mockProfileRepo, mockFollowRepo)
	ctxForTest := context.Background()

	return ctxForTest, userService, mockUserRepo, mockProfileRepo, mockFollowRepo, sqlMock
}

func TestUserService_RegisterUser_Success(t *testing.T) {
	// 1. Setup test dependencies
	ctxForTest, userService, mockUserRepo, mockProfileRepo, _, sqlMock := setupUserServiceTest(t)
	username := "testuser"
	email := "test@example.com"
	password := "password123"
	hashedPassword := HashPassword(password)

	// Setup sqlmock expectations for transaction
	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	// 2. Define mock behavior
	// When CreateUser is called, it should set the ID on the user
	mockUserRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.Username == username && u.Email == email
	})).Run(func(args mock.Arguments) {
		user := args.Get(1).(*models.User)
		user.ID = 1
	}).Return(nil)

	// When CreateProfile is called, it should succeed
	mockProfileRepo.On("CreateProfile", mock.Anything, mock.MatchedBy(func(p *models.Profile) bool {
		return p.UserID == int64(1)
	})).Return(nil)

	// 3. Call the service method under test
	user, err := userService.RegisterUser(ctxForTest, username, email, password)

	// 4. Assert results
	assert.NoError(t, err)                         // Check for no error
	assert.NotNil(t, user)                         // User should not be nil
	assert.Equal(t, username, user.Username)       // Username should match
	assert.Equal(t, email, user.Email)             // Email should match
	assert.Equal(t, hashedPassword, user.Password) // Password should be hashed
	assert.Equal(t, int64(1), user.ID)             // ID should be set

	// 5. Assert that mock expectations were met
	mockUserRepo.AssertExpectations(t)
	mockProfileRepo.AssertExpectations(t)
}

func TestUserService_RegisterUser_CreateUserError(t *testing.T) {
	// 1. Setup test dependencies
	ctxForTest, userService, mockUserRepo, mockProfileRepo, _, sqlMock := setupUserServiceTest(t)
	username := "testuser"
	email := "test@example.com"
	password := "password123"

	// Setup sqlmock expectations for transaction that will rollback
	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	// 2. Define mock behavior: CreateUser fails
	expectedError := errors.New("duplicate email")
	mockUserRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.Username == username && u.Email == email
	})).Return(expectedError)

	// CreateProfile should NOT be called because CreateUser fails
	mockProfileRepo.On("CreateProfile", mock.Anything, mock.Anything).Return(nil)

	// 3. Call the service method under test
	_, err := userService.RegisterUser(ctxForTest, username, email, password)

	// 4. Assert results
	assert.Error(t, err)                // Check that an error occurred
	assert.Equal(t, expectedError, err) // Error should match the expected error
	// Note: user is still populated with data even though there's an error (transaction behavior)

	// 5. Assert that mock expectations were met
	mockUserRepo.AssertExpectations(t)
	// Verify that CreateProfile was not called in case of transaction rollback
	mockProfileRepo.AssertNotCalled(t, "CreateProfile")
}

func TestUserService_RegisterUser_CreateProfileError(t *testing.T) {
	// 1. Setup test dependencies
	ctxForTest, userService, mockUserRepo, mockProfileRepo, _, sqlMock := setupUserServiceTest(t)
	username := "testuser"
	email := "test@example.com"
	password := "password123"

	// Setup sqlmock expectations for transaction that will rollback
	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	// 2. Define mock behavior
	// CreateUser succeeds and sets ID
	mockUserRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
		return u.Username == username && u.Email == email
	})).Run(func(args mock.Arguments) {
		user := args.Get(1).(*models.User)
		user.ID = 1
	}).Return(nil)

	// CreateProfile fails
	expectedError := errors.New("failed to create profile")
	mockProfileRepo.On("CreateProfile", mock.Anything, mock.MatchedBy(func(p *models.Profile) bool {
		return p.UserID == int64(1)
	})).Return(expectedError)

	// 3. Call the service method under test
	_, err := userService.RegisterUser(ctxForTest, username, email, password)

	// 4. Assert results
	assert.Error(t, err)                // Check that an error occurred
	assert.Equal(t, expectedError, err) // Error should match the expected error
	// Note: user is still populated with data even though there's an error (transaction behavior)

	// 5. Assert that mock expectations were met
	mockUserRepo.AssertExpectations(t)
	mockProfileRepo.AssertExpectations(t)
}
