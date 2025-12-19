package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

type favoriteHandlerMocks struct {
	favoriteRepo *mocks.MockFavoriteRepository
	articleRepo  *mocks.MockArticleRepository
	sqlMock      sqlmock.Sqlmock
}

func setupFavoriteHandlerTest(t *testing.T) (*gin.Engine, *handlers.FavoriteHandler, favoriteHandlerMocks) {
	m := favoriteHandlerMocks{
		favoriteRepo: new(mocks.MockFavoriteRepository),
		articleRepo:  new(mocks.MockArticleRepository),
	}

	mockDB, sqlMock := CreateMockDB(t)
	m.sqlMock = sqlMock
	favoriteService := services.NewFavoriteService(mockDB, m.favoriteRepo, m.articleRepo)
	favoriteHandler := handlers.NewFavoriteHandler(favoriteService)

	router := SetupRouter()
	return router, favoriteHandler, m
}

func TestFavoriteHandler_FavoriteArticle_Success(t *testing.T) {
	router, favoriteHandler, m := setupFavoriteHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.POST("/api/articles/:slug/favorite", favoriteHandler.FavoriteArticle)

	slug := "test-article"
	article := &models.Article{
		ID:             1,
		Slug:           slug,
		FavoritesCount: 0,
	}

	m.sqlMock.ExpectBegin()
	m.articleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(article, nil)
	m.favoriteRepo.On("IsFavorited", mock.Anything, int64(1), int64(1)).Return(false, nil)
	m.favoriteRepo.On("AddFavorite", mock.Anything, int64(1), int64(1)).Return(nil)
	m.articleRepo.On("UpdateArticle", mock.Anything, mock.AnythingOfType("*models.Article")).Return(nil)
	m.sqlMock.ExpectCommit()

	updatedArticle := &models.Article{
		ID:             1,
		Slug:           slug,
		Title:          "Test Article",
		Description:    "Description",
		Body:           "Body",
		FavoritesCount: 1,
		Author:         &models.User{Username: "author1"},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Favorites:      []*models.Favorite{{UserID: 1, ArticleID: 1}},
	}
	m.favoriteRepo.On("GetArticleWithFavorites", mock.Anything, int64(1)).Return(updatedArticle, nil)

	req, _ := http.NewRequest("POST", "/api/articles/"+slug+"/favorite", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dtos.ArticleDetailResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Article.Favorited)
	assert.Equal(t, 1, resp.Article.FavoritesCount)

	m.articleRepo.AssertExpectations(t)
	m.favoriteRepo.AssertExpectations(t)
}

func TestFavoriteHandler_FavoriteArticle_Unauthorized(t *testing.T) {
	router, favoriteHandler, _ := setupFavoriteHandlerTest(t)
	router.POST("/api/articles/:slug/favorite", favoriteHandler.FavoriteArticle)

	req, _ := http.NewRequest("POST", "/api/articles/test-article/favorite", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	AssertAPIError(t, w, http.StatusUnauthorized, "missing authorization")
}

func TestFavoriteHandler_FavoriteArticle_NotFound(t *testing.T) {
	router, favoriteHandler, m := setupFavoriteHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.POST("/api/articles/:slug/favorite", favoriteHandler.FavoriteArticle)

	slug := "non-existent"
	m.sqlMock.ExpectBegin()
	m.articleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(nil, gorm.ErrRecordNotFound)
	m.sqlMock.ExpectRollback()

	req, _ := http.NewRequest("POST", "/api/articles/"+slug+"/favorite", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	AssertAPIError(t, w, http.StatusNotFound, "article not found")
}

func TestFavoriteHandler_UnfavoriteArticle_Success(t *testing.T) {
	router, favoriteHandler, m := setupFavoriteHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.DELETE("/api/articles/:slug/favorite", favoriteHandler.UnfavoriteArticle)

	slug := "test-article"
	article := &models.Article{
		ID:             1,
		Slug:           slug,
		FavoritesCount: 1,
	}

	m.sqlMock.ExpectBegin()
	m.articleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(article, nil)
	m.favoriteRepo.On("IsFavorited", mock.Anything, int64(1), int64(1)).Return(true, nil)
	m.favoriteRepo.On("RemoveFavorite", mock.Anything, int64(1), int64(1)).Return(nil)
	m.articleRepo.On("UpdateArticle", mock.Anything, mock.AnythingOfType("*models.Article")).Return(nil)
	m.sqlMock.ExpectCommit()

	updatedArticle := &models.Article{
		ID:             1,
		Slug:           slug,
		Title:          "Test Article",
		Description:    "Description",
		Body:           "Body",
		FavoritesCount: 0,
		Author:         &models.User{Username: "author1"},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Favorites:      []*models.Favorite{},
	}
	m.favoriteRepo.On("GetArticleWithFavorites", mock.Anything, int64(1)).Return(updatedArticle, nil)

	req, _ := http.NewRequest("DELETE", "/api/articles/"+slug+"/favorite", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dtos.ArticleDetailResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Article.Favorited)
	assert.Equal(t, 0, resp.Article.FavoritesCount)

	m.articleRepo.AssertExpectations(t)
	m.favoriteRepo.AssertExpectations(t)
}

func TestFavoriteHandler_UnfavoriteArticle_NotFound(t *testing.T) {
	router, favoriteHandler, m := setupFavoriteHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.DELETE("/api/articles/:slug/favorite", favoriteHandler.UnfavoriteArticle)

	slug := "non-existent"
	m.sqlMock.ExpectBegin()
	m.articleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(nil, gorm.ErrRecordNotFound)
	m.sqlMock.ExpectRollback()

	req, _ := http.NewRequest("DELETE", "/api/articles/"+slug+"/favorite", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	AssertAPIError(t, w, http.StatusNotFound, "article not found")
}
