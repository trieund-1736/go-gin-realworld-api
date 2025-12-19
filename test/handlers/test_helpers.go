package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	appErrors "go-gin-realworld-api/internal/errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// HashPassword hashes a password using SHA256 (same as in services)
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// AssertAPIError asserts that the response is a valid API error response
func AssertAPIError(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedMessage string, expectedDetails ...map[string]string) {
	assert.Equal(t, expectedStatus, w.Code)

	var resp appErrors.APIErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, resp.Code)
	assert.Equal(t, expectedMessage, resp.Message)

	if len(expectedDetails) > 0 && expectedDetails[0] != nil {
		details, ok := resp.Details.(map[string]interface{})
		assert.True(t, ok, "Details should be a map")
		for k, v := range expectedDetails[0] {
			assert.Equal(t, v, details[k])
		}
	}
}

// CreateMockDB creates a mock database for testing
func CreateMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, sqlMock, err := sqlmock.New()
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

	return gormDB, sqlMock
}

// SetupRouter sets up a gin router in test mode
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}
