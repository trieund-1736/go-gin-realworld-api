package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/services"
	"go-gin-realworld-api/test/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupProfileServiceTest(t *testing.T) (context.Context, *services.ProfileService, *mocks.MockUserRepository, *mocks.MockProfileRepository, *mocks.MockFollowRepository, sqlmock.Sqlmock) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockProfileRepo := new(mocks.MockProfileRepository)
	mockFollowRepo := new(mocks.MockFollowRepository)
	gormDB, sqlMock := CreateMockDB(t)
	profileService := services.NewProfileService(gormDB, mockUserRepo, mockProfileRepo, mockFollowRepo)
	ctxForTest := context.Background()

	return ctxForTest, profileService, mockUserRepo, mockProfileRepo, mockFollowRepo, sqlMock
}

func TestProfileService_GetProfileByUsername_Success(t *testing.T) {
	ctxForTest, profileService, mockUserRepo, _, mockFollowRepo, _ := setupProfileServiceTest(t)
	username := "testuser"
	currentUserID := int64(1)
	targetUserID := int64(2)

	user := &models.User{
		ID:       targetUserID,
		Username: username,
		Profile: &models.Profile{
			UserID: targetUserID,
			Bio:    sql.NullString{String: "test bio", Valid: true},
			Image:  sql.NullString{String: "test image", Valid: true},
		},
	}

	mockUserRepo.On("FindUserByUsername", mock.Anything, username, []bool{true}).Return(user, nil)
	mockFollowRepo.On("IsFollowing", mock.Anything, currentUserID, targetUserID).Return(true, nil)

	resp, err := profileService.GetProfileByUsername(ctxForTest, username, currentUserID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, username, resp.Profile.Username)
	assert.Equal(t, "test bio", resp.Profile.Bio)
	assert.True(t, resp.Profile.Following)
	mockUserRepo.AssertExpectations(t)
	mockFollowRepo.AssertExpectations(t)
}

func TestProfileService_GetProfileByUsername_NotFound(t *testing.T) {
	ctxForTest, profileService, mockUserRepo, _, _, _ := setupProfileServiceTest(t)
	username := "nonexistent"
	expectedError := errors.New("user not found")

	mockUserRepo.On("FindUserByUsername", mock.Anything, username, []bool{true}).Return(nil, expectedError)

	resp, err := profileService.GetProfileByUsername(ctxForTest, username, 0)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedError, err)
	mockUserRepo.AssertExpectations(t)
}

func TestProfileService_FollowUser_Success(t *testing.T) {
	ctxForTest, profileService, mockUserRepo, _, mockFollowRepo, sqlMock := setupProfileServiceTest(t)
	followerID := int64(1)
	followeeUsername := "followee"
	followeeID := int64(2)

	followee := &models.User{
		ID:       followeeID,
		Username: followeeUsername,
		Profile:  &models.Profile{UserID: followeeID},
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	// Transaction calls
	mockUserRepo.On("FindUserByUsername", mock.Anything, followeeUsername, []bool(nil)).Return(followee, nil).Once()
	mockFollowRepo.On("IsFollowing", mock.Anything, followerID, followeeID).Return(false, nil).Once()
	mockFollowRepo.On("CreateFollow", mock.Anything, mock.MatchedBy(func(f *models.Follow) bool {
		return f.FollowerID == followerID && f.FolloweeID == followeeID
	})).Return(nil).Once()

	// GetProfileByUsername calls (after transaction)
	mockUserRepo.On("FindUserByUsername", mock.Anything, followeeUsername, []bool{true}).Return(followee, nil).Once()
	mockFollowRepo.On("IsFollowing", mock.Anything, followerID, followeeID).Return(true, nil).Once()

	resp, err := profileService.FollowUser(ctxForTest, followerID, followeeUsername)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Profile.Following)
	mockUserRepo.AssertExpectations(t)
	mockFollowRepo.AssertExpectations(t)
}

func TestProfileService_UnfollowUser_Success(t *testing.T) {
	ctxForTest, profileService, mockUserRepo, _, mockFollowRepo, sqlMock := setupProfileServiceTest(t)
	followerID := int64(1)
	followeeUsername := "followee"
	followeeID := int64(2)

	followee := &models.User{
		ID:       followeeID,
		Username: followeeUsername,
		Profile:  &models.Profile{UserID: followeeID},
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	// Transaction calls
	mockUserRepo.On("FindUserByUsername", mock.Anything, followeeUsername, []bool(nil)).Return(followee, nil)
	mockFollowRepo.On("DeleteFollow", mock.Anything, followerID, followeeID).Return(nil)

	// GetProfileByUsername calls (after transaction)
	mockUserRepo.On("FindUserByUsername", mock.Anything, followeeUsername, []bool{true}).Return(followee, nil)
	mockFollowRepo.On("IsFollowing", mock.Anything, followerID, followeeID).Return(false, nil)

	resp, err := profileService.UnfollowUser(ctxForTest, followerID, followeeUsername)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Profile.Following)
	mockUserRepo.AssertExpectations(t)
	mockFollowRepo.AssertExpectations(t)
}
