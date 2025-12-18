package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gin-realworld-api/internal/handlers"
	"go-gin-realworld-api/internal/services"
	"go-gin-realworld-api/test/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type tagHandlerMocks struct {
	tagRepo *mocks.MockTagRepository
	sqlMock sqlmock.Sqlmock
}

func setupTagHandlerTest(t *testing.T) (*gin.Engine, *handlers.TagHandler, tagHandlerMocks) {
	m := tagHandlerMocks{
		tagRepo: new(mocks.MockTagRepository),
	}

	mockDB, sqlMock := CreateMockDB(t)
	m.sqlMock = sqlMock
	tagService := services.NewTagService(mockDB, m.tagRepo)
	tagHandler := handlers.NewTagHandler(tagService)

	router := SetupRouter()
	return router, tagHandler, m
}

func TestTagHandler_GetTags_Success(t *testing.T) {
	router, tagHandler, m := setupTagHandlerTest(t)
	router.GET("/api/tags", tagHandler.GetTags)

	expectedTags := []string{"reactjs", "angularjs", "dragons"}
	m.tagRepo.On("GetAllTags", mock.Anything).Return(expectedTags, nil)

	req, _ := http.NewRequest("GET", "/api/tags", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string][]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedTags, resp["tags"])

	m.tagRepo.AssertExpectations(t)
}

func TestTagHandler_GetTags_Error(t *testing.T) {
	router, tagHandler, m := setupTagHandlerTest(t)
	router.GET("/api/tags", tagHandler.GetTags)

	m.tagRepo.On("GetAllTags", mock.Anything).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/tags", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to retrieve tags")

	m.tagRepo.AssertExpectations(t)
}
