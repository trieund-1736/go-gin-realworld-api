package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gin-realworld-api/internal/dtos"
	"go-gin-realworld-api/internal/handlers"
	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/services"
	"go-gin-realworld-api/test/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type profileHandlerMocks struct {
	userRepo    *mocks.MockUserRepository
	profileRepo *mocks.MockProfileRepository
	followRepo  *mocks.MockFollowRepository
	sqlMock     sqlmock.Sqlmock
}

func setupProfileHandlerTest(t *testing.T) (*gin.Engine, *handlers.ProfileHandler, profileHandlerMocks) {
	m := profileHandlerMocks{
		userRepo:    new(mocks.MockUserRepository),
		profileRepo: new(mocks.MockProfileRepository),
		followRepo:  new(mocks.MockFollowRepository),
	}

	mockDB, sqlMock := CreateMockDB(t)
	m.sqlMock = sqlMock
	profileService := services.NewProfileService(mockDB, m.userRepo, m.profileRepo, m.followRepo)
	profileHandler := handlers.NewProfileHandler(profileService)

	router := SetupRouter()
	return router, profileHandler, m
}

func TestProfileHandler_GetProfile_Success(t *testing.T) {
	router, profileHandler, m := setupProfileHandlerTest(t)
	router.GET("/api/profiles/:username", profileHandler.GetProfile)

	username := "testuser"
	user := &models.User{
		ID:       1,
		Username: username,
		Profile: &models.Profile{
			UserID: 1,
			Bio:    sql.NullString{String: "test bio", Valid: true},
			Image:  sql.NullString{String: "test image", Valid: true},
		},
	}

	m.userRepo.On("FindUserByUsername", mock.Anything, username, []bool{true}).Return(user, nil)

	req, _ := http.NewRequest("GET", "/api/profiles/"+username, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dtos.ProfileResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, username, resp.Profile.Username)
	assert.Equal(t, "test bio", resp.Profile.Bio)
	assert.Equal(t, "test image", resp.Profile.Image)
	assert.False(t, resp.Profile.Following)

	m.userRepo.AssertExpectations(t)
}

func TestProfileHandler_GetProfile_NotFound(t *testing.T) {
	router, profileHandler, m := setupProfileHandlerTest(t)
	router.GET("/api/profiles/:username", profileHandler.GetProfile)

	username := "nonexistent"
	m.userRepo.On("FindUserByUsername", mock.Anything, username, []bool{true}).Return(nil, gorm.ErrRecordNotFound)

	req, _ := http.NewRequest("GET", "/api/profiles/"+username, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "profile not found")
}

func TestProfileHandler_GetProfile_Authenticated_Following(t *testing.T) {
	router, profileHandler, m := setupProfileHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.GET("/api/profiles/:username", profileHandler.GetProfile)

	username := "testuser"
	user := &models.User{
		ID:       2,
		Username: username,
		Profile: &models.Profile{
			UserID: 2,
		},
	}

	m.userRepo.On("FindUserByUsername", mock.Anything, username, []bool{true}).Return(user, nil)
	m.followRepo.On("IsFollowing", mock.Anything, int64(1), int64(2)).Return(true, nil)

	req, _ := http.NewRequest("GET", "/api/profiles/"+username, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dtos.ProfileResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Profile.Following)

	m.userRepo.AssertExpectations(t)
	m.followRepo.AssertExpectations(t)
}

func TestProfileHandler_FollowUser_Success(t *testing.T) {
	router, profileHandler, m := setupProfileHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.POST("/api/profiles/:username/follow", profileHandler.FollowUser)

	username := "followee"
	followee := &models.User{ID: 2, Username: username}

	m.sqlMock.ExpectBegin()
	m.userRepo.On("FindUserByUsername", mock.Anything, username, []bool(nil)).Return(followee, nil)
	m.followRepo.On("IsFollowing", mock.Anything, int64(1), int64(2)).Return(false, nil).Once()
	m.followRepo.On("CreateFollow", mock.Anything, mock.AnythingOfType("*models.Follow")).Return(nil)
	m.sqlMock.ExpectCommit()

	// After follow, it calls GetProfileByUsername
	userWithProfile := &models.User{
		ID:       2,
		Username: username,
		Profile:  &models.Profile{UserID: 2},
	}
	m.userRepo.On("FindUserByUsername", mock.Anything, username, []bool{true}).Return(userWithProfile, nil)
	m.followRepo.On("IsFollowing", mock.Anything, int64(1), int64(2)).Return(true, nil).Once()

	req, _ := http.NewRequest("POST", "/api/profiles/"+username+"/follow", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dtos.ProfileResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Profile.Following)

	m.userRepo.AssertExpectations(t)
	m.followRepo.AssertExpectations(t)
}

func TestProfileHandler_UnfollowUser_Success(t *testing.T) {
	router, profileHandler, m := setupProfileHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.DELETE("/api/profiles/:username/follow", profileHandler.UnfollowUser)

	username := "followee"
	followee := &models.User{ID: 2, Username: username}

	m.sqlMock.ExpectBegin()
	m.userRepo.On("FindUserByUsername", mock.Anything, username, []bool(nil)).Return(followee, nil)
	m.followRepo.On("DeleteFollow", mock.Anything, int64(1), int64(2)).Return(nil)
	m.sqlMock.ExpectCommit()

	// After unfollow, it calls GetProfileByUsername
	userWithProfile := &models.User{
		ID:       2,
		Username: username,
		Profile:  &models.Profile{UserID: 2},
	}
	m.userRepo.On("FindUserByUsername", mock.Anything, username, []bool{true}).Return(userWithProfile, nil)
	m.followRepo.On("IsFollowing", mock.Anything, int64(1), int64(2)).Return(false, nil)

	req, _ := http.NewRequest("DELETE", "/api/profiles/"+username+"/follow", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dtos.ProfileResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Profile.Following)

	m.userRepo.AssertExpectations(t)
	m.followRepo.AssertExpectations(t)
}
