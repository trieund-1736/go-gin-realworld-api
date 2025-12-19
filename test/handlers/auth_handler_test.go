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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupAuthHandlerTest sets up the dependencies for testing AuthHandler
func setupAuthHandlerTest(t *testing.T) (*gin.Engine, *mocks.MockUserRepository) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockDB, _ := CreateMockDB(t)
	authService := services.NewAuthService(mockDB, mockUserRepo)
	authHandler := handlers.NewAuthHandler(authService)

	router := SetupRouter()
	router.POST("/api/users/login", authHandler.Login)

	return router, mockUserRepo
}

func TestAuthHandler_Login_Success(t *testing.T) {
	// Setup
	router, mockUserRepo := setupAuthHandlerTest(t)

	email := "test@example.com"
	password := "password123"
	hashedPassword := HashPassword(password)

	expectedUser := &models.User{
		ID:       1,
		Username: "testuser",
		Email:    email,
		Password: hashedPassword,
	}

	// Mock behavior
	mockUserRepo.On("FindUserByEmail", mock.Anything, email).Return(expectedUser, nil)

	// Request body
	loginReq := dtos.LoginRequest{}
	loginReq.User.Email = email
	loginReq.User.Password = password
	jsonBody, _ := json.Marshal(loginReq)

	// Perform request
	req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var resp dtos.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Username, resp.User.Username)
	assert.Equal(t, expectedUser.Email, resp.User.Email)
	assert.NotEmpty(t, resp.User.Token)

	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Login_Validation(t *testing.T) {
	// Setup
	router, _ := setupAuthHandlerTest(t)

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedError  string
		expectedDetail map[string]string
	}{
		{
			name:           "Missing user object",
			payload:        `{}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
			expectedDetail: map[string]string{"Email": "is required", "Password": "is required"},
		},
		{
			name:           "Missing email",
			payload:        `{"user": {"password": "password123"}}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
			expectedDetail: map[string]string{"Email": "is required"},
		},
		{
			name:           "Invalid email format",
			payload:        `{"user": {"email": "invalid-email", "password": "password123"}}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
			expectedDetail: map[string]string{"Email": "must be a valid email"},
		},
		{
			name:           "Missing password",
			payload:        `{"user": {"email": "test@example.com"}}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
			expectedDetail: map[string]string{"Password": "is required"},
		},
		{
			name:           "Empty email",
			payload:        `{"user": {"email": "", "password": "password123"}}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
			expectedDetail: map[string]string{"Email": "is required"},
		},
		{
			name:           "Empty password",
			payload:        `{"user": {"email": "test@example.com", "password": ""}}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
			expectedDetail: map[string]string{"Password": "is required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			AssertAPIError(t, w, tt.expectedStatus, tt.expectedError, tt.expectedDetail)
		})
	}
}

func TestAuthHandler_Login_Unauthorized(t *testing.T) {
	// Setup
	router, mockUserRepo := setupAuthHandlerTest(t)

	email := "test@example.com"
	password := "wrongpassword"

	// Mock behavior: User not found or wrong password returns error in service
	// The service returns "invalid credentials" for both cases
	mockUserRepo.On("FindUserByEmail", mock.Anything, email).Return(nil, fmt.Errorf("user not found"))

	// Request body
	loginReq := dtos.LoginRequest{}
	loginReq.User.Email = email
	loginReq.User.Password = password
	jsonBody, _ := json.Marshal(loginReq)

	// Perform request
	req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	AssertAPIError(t, w, http.StatusUnauthorized, "Invalid email or password")

	mockUserRepo.AssertExpectations(t)
}
