package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
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
)

type userHandlerMocks struct {
	userRepo    *mocks.MockUserRepository
	profileRepo *mocks.MockProfileRepository
	followRepo  *mocks.MockFollowRepository
	sqlMock     sqlmock.Sqlmock
}

func setupUserHandlerTest(t *testing.T) (*gin.Engine, *handlers.UserHandler, userHandlerMocks) {
	m := userHandlerMocks{
		userRepo:    new(mocks.MockUserRepository),
		profileRepo: new(mocks.MockProfileRepository),
		followRepo:  new(mocks.MockFollowRepository),
	}

	mockDB, sqlMock := CreateMockDB(t)
	m.sqlMock = sqlMock
	userService := services.NewUserService(mockDB, m.userRepo, m.profileRepo, m.followRepo)
	userHandler := handlers.NewUserHandler(userService)

	router := SetupRouter()
	return router, userHandler, m
}

func TestUserHandler_RegisterUser_Success(t *testing.T) {
	router, userHandler, m := setupUserHandlerTest(t)
	router.POST("/api/users", userHandler.RegisterUser)

	username := "testuser"
	email := "test@example.com"
	password := "password123"

	// Mock behavior
	m.sqlMock.ExpectBegin()
	m.sqlMock.ExpectCommit()

	m.userRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(1).(*models.User)
		user.ID = 1
	})
	m.profileRepo.On("CreateProfile", mock.Anything, mock.AnythingOfType("*models.Profile")).Return(nil)

	// Request body
	reqBody := dtos.RegisterUserRequest{}
	reqBody.User.Username = username
	reqBody.User.Email = email
	reqBody.User.Password = password
	jsonBody, _ := json.Marshal(reqBody)

	// Perform request
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dtos.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, username, resp.User.Username)
	assert.Equal(t, email, resp.User.Email)
	assert.Equal(t, int64(1), resp.User.ID)

	m.userRepo.AssertExpectations(t)
	m.profileRepo.AssertExpectations(t)
}

func TestUserHandler_RegisterUser_Validation(t *testing.T) {
	router, userHandler, _ := setupUserHandlerTest(t)
	router.POST("/api/users", userHandler.RegisterUser)

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Missing user object",
			payload:        `{}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid request body",
		},
		{
			name:           "Missing username",
			payload:        `{"user": {"email": "test@example.com", "password": "password123"}}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid request body",
		},
		{
			name:           "Invalid email format",
			payload:        `{"user": {"username": "testuser", "email": "invalid-email", "password": "password123"}}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedError)
		})
	}
}

func TestUserHandler_GetCurrentUser_Success(t *testing.T) {
	router, userHandler, m := setupUserHandlerTest(t)

	// Mock middleware to set user_id
	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.GET("/api/user", userHandler.GetCurrentUser)

	expectedUser := &models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Mock behavior
	m.userRepo.On("FindUserByID", mock.Anything, int64(1)).Return(expectedUser, nil)

	// Perform request
	req, _ := http.NewRequest("GET", "/api/user", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var resp dtos.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Username, resp.User.Username)
	assert.Equal(t, expectedUser.Email, resp.User.Email)

	m.userRepo.AssertExpectations(t)
}

func TestUserHandler_GetCurrentUser_Unauthorized(t *testing.T) {
	router, userHandler, _ := setupUserHandlerTest(t)
	// No middleware to set user_id
	router.GET("/api/user", userHandler.GetCurrentUser)

	// Perform request
	req, _ := http.NewRequest("GET", "/api/user", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "user not authenticated")
}

func TestUserHandler_GetCurrentUser_NotFound(t *testing.T) {
	router, userHandler, m := setupUserHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.GET("/api/user", userHandler.GetCurrentUser)

	// Mock behavior
	m.userRepo.On("FindUserByID", mock.Anything, int64(1)).Return(nil, fmt.Errorf("user not found"))

	// Perform request
	req, _ := http.NewRequest("GET", "/api/user", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "user not found")

	m.userRepo.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_Success(t *testing.T) {
	router, userHandler, m := setupUserHandlerTest(t)

	// Mock middleware to set user_id
	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.PUT("/api/user", userHandler.UpdateUser)

	existingUser := &models.User{
		ID:       1,
		Username: "olduser",
		Email:    "old@example.com",
	}

	// Mock behavior
	m.sqlMock.ExpectBegin()
	m.sqlMock.ExpectCommit()

	m.userRepo.On("FindUserByID", mock.Anything, int64(1)).Return(existingUser, nil)
	m.userRepo.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
	m.profileRepo.On("FindProfileByUserID", mock.Anything, int64(1)).Return(&models.Profile{UserID: 1}, nil)
	m.profileRepo.On("UpdateProfile", mock.Anything, mock.AnythingOfType("*models.Profile")).Return(nil)

	// Request body
	updateReq := dtos.UpdateUserRequest{}
	updateReq.User.Username = "newuser"
	updateReq.User.Bio = "new bio"
	jsonBody, _ := json.Marshal(updateReq)

	// Perform request
	req, _ := http.NewRequest("PUT", "/api/user", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var resp dtos.UpdateUserResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "newuser", resp.User.Username)
	assert.Equal(t, "new bio", resp.User.Bio)

	m.userRepo.AssertExpectations(t)
	m.profileRepo.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_Validation(t *testing.T) {
	router, userHandler, _ := setupUserHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.PUT("/api/user", userHandler.UpdateUser)

	// Malformed JSON
	payload := `{"user": { "email": "invalid"`
	req, _ := http.NewRequest("PUT", "/api/user", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request body")
}
