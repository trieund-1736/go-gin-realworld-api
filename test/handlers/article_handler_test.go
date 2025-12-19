package handlers

import (
	"bytes"
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

type articleHandlerMocks struct {
	articleRepo *mocks.MockArticleRepository
	sqlMock     sqlmock.Sqlmock
}

func setupArticleHandlerTest(t *testing.T) (*gin.Engine, *handlers.ArticleHandler, articleHandlerMocks) {
	m := articleHandlerMocks{
		articleRepo: new(mocks.MockArticleRepository),
	}

	mockDB, sqlMock := CreateMockDB(t)
	m.sqlMock = sqlMock
	articleService := services.NewArticleService(mockDB, m.articleRepo)
	articleHandler := handlers.NewArticleHandler(articleService)

	router := SetupRouter()
	return router, articleHandler, m
}

func TestArticleHandler_ListArticles_Success(t *testing.T) {
	router, articleHandler, m := setupArticleHandlerTest(t)
	router.GET("/api/articles", articleHandler.ListArticles)

	articles := []*models.Article{
		{
			ID:          1,
			Slug:        "test-article",
			Title:       "Test Article",
			Description: "Description",
			Body:        "Body",
			Author:      &models.User{Username: "author1"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	m.articleRepo.On("ListArticles", mock.Anything, "", "", (*bool)(nil), (*int64)(nil), 20, 0).
		Return(articles, int64(1), nil)

	req, _ := http.NewRequest("GET", "/api/articles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dtos.ArticlesListResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.ArticlesCount)
	assert.Equal(t, "test-article", resp.Articles[0].Slug)
	assert.Equal(t, "author1", resp.Articles[0].Author.Username)

	m.articleRepo.AssertExpectations(t)
}

func TestArticleHandler_GetArticle_Success(t *testing.T) {
	router, articleHandler, m := setupArticleHandlerTest(t)
	router.GET("/api/articles/:slug", articleHandler.GetArticle)

	slug := "test-article"
	article := &models.Article{
		ID:          1,
		Slug:        slug,
		Title:       "Test Article",
		Description: "Description",
		Body:        "Body",
		Author:      &models.User{Username: "author1"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.articleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(article, nil)

	req, _ := http.NewRequest("GET", "/api/articles/"+slug, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dtos.ArticleDetailResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, slug, resp.Article.Slug)

	m.articleRepo.AssertExpectations(t)
}

func TestArticleHandler_FeedArticles_Success(t *testing.T) {
	router, articleHandler, m := setupArticleHandlerTest(t)

	// Mock middleware to set user_id
	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.GET("/api/articles/feed", articleHandler.FeedArticles)

	articles := []*models.Article{
		{
			ID:          1,
			Slug:        "feed-article",
			Title:       "Feed Article",
			Description: "Description",
			Body:        "Body",
			Author:      &models.User{Username: "followed_user"},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	m.articleRepo.On("FeedArticles", mock.Anything, int64(1), 20, 0).
		Return(articles, int64(1), nil)

	req, _ := http.NewRequest("GET", "/api/articles/feed", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dtos.ArticlesListResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.ArticlesCount)
	assert.Equal(t, "feed-article", resp.Articles[0].Slug)

	m.articleRepo.AssertExpectations(t)
}

func TestArticleHandler_CreateArticle_Success(t *testing.T) {
	router, articleHandler, m := setupArticleHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.POST("/api/articles", articleHandler.CreateArticle)

	reqBody := dtos.CreateArticleRequest{}
	reqBody.Article.Title = "New Article"
	reqBody.Article.Description = "Description"
	reqBody.Article.Body = "Body"
	reqBody.Article.TagList = []string{"tag1", "tag2"}
	jsonBody, _ := json.Marshal(reqBody)

	m.sqlMock.ExpectBegin()
	m.articleRepo.On("CreateArticle", mock.Anything, mock.AnythingOfType("*models.Article")).Return(nil).Run(func(args mock.Arguments) {
		article := args.Get(1).(*models.Article)
		article.ID = 1
	})
	m.articleRepo.On("AssignTagsToArticle", mock.Anything, int64(1), reqBody.Article.TagList).Return(nil)
	m.sqlMock.ExpectCommit()

	// After creation, the service fetches the article again
	createdArticle := &models.Article{
		ID:          1,
		Slug:        "new-article",
		Title:       "New Article",
		Description: "Description",
		Body:        "Body",
		Author:      &models.User{Username: "testuser"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.articleRepo.On("FindArticleBySlug", mock.Anything, "new-article").Return(createdArticle, nil)

	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dtos.ArticleDetailResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "new-article", resp.Article.Slug)

	m.articleRepo.AssertExpectations(t)
}

func TestArticleHandler_UpdateArticle_Success(t *testing.T) {
	router, articleHandler, m := setupArticleHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.PUT("/api/articles/:slug", articleHandler.UpdateArticle)

	slug := "old-article"
	existingArticle := &models.Article{
		ID:          1,
		Slug:        slug,
		Title:       "Old Title",
		Description: "Old Description",
		Body:        "Old Body",
		AuthorID:    1,
	}

	reqBody := dtos.UpdateArticleRequest{}
	reqBody.Article.Title = "New Title"
	jsonBody, _ := json.Marshal(reqBody)

	m.sqlMock.ExpectBegin()
	m.articleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(existingArticle, nil)
	m.articleRepo.On("UpdateArticle", mock.Anything, mock.AnythingOfType("*models.Article")).Return(nil)
	m.sqlMock.ExpectCommit()

	updatedArticle := &models.Article{
		ID:          1,
		Slug:        "new-title",
		Title:       "New Title",
		Description: "Old Description",
		Body:        "Old Body",
		Author:      &models.User{Username: "testuser"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.articleRepo.On("FindArticleBySlug", mock.Anything, "new-title").Return(updatedArticle, nil)

	req, _ := http.NewRequest("PUT", "/api/articles/"+slug, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dtos.ArticleDetailResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "new-title", resp.Article.Slug)

	m.articleRepo.AssertExpectations(t)
}

func TestArticleHandler_DeleteArticle_Success(t *testing.T) {
	router, articleHandler, m := setupArticleHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.DELETE("/api/articles/:slug", articleHandler.DeleteArticle)

	slug := "test-article"
	article := &models.Article{ID: 1, Slug: slug}

	m.sqlMock.ExpectBegin()
	m.articleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(article, nil)
	m.articleRepo.On("DeleteArticleBySlug", mock.Anything, slug).Return(nil)
	m.sqlMock.ExpectCommit()

	req, _ := http.NewRequest("DELETE", "/api/articles/"+slug, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	m.articleRepo.AssertExpectations(t)
}

func TestArticleHandler_GetArticle_NotFound(t *testing.T) {
	router, articleHandler, m := setupArticleHandlerTest(t)
	router.GET("/api/articles/:slug", articleHandler.GetArticle)

	slug := "non-existent"
	m.articleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(nil, gorm.ErrRecordNotFound)

	req, _ := http.NewRequest("GET", "/api/articles/"+slug, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	AssertAPIError(t, w, http.StatusNotFound, "article not found")
}

func TestArticleHandler_CreateArticle_Unauthorized(t *testing.T) {
	router, articleHandler, _ := setupArticleHandlerTest(t)
	router.POST("/api/articles", articleHandler.CreateArticle)

	req, _ := http.NewRequest("POST", "/api/articles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	AssertAPIError(t, w, http.StatusUnauthorized, "authentication required")
}

func TestArticleHandler_UpdateArticle_NotFound(t *testing.T) {
	router, articleHandler, m := setupArticleHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.PUT("/api/articles/:slug", articleHandler.UpdateArticle)

	slug := "non-existent"
	reqBody := dtos.UpdateArticleRequest{}
	reqBody.Article.Title = "New Title"
	jsonBody, _ := json.Marshal(reqBody)

	m.sqlMock.ExpectBegin()
	m.articleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(nil, gorm.ErrRecordNotFound)
	m.sqlMock.ExpectRollback()

	req, _ := http.NewRequest("PUT", "/api/articles/"+slug, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	AssertAPIError(t, w, http.StatusNotFound, "article not found")
}

func TestArticleHandler_DeleteArticle_NotFound(t *testing.T) {
	router, articleHandler, m := setupArticleHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.DELETE("/api/articles/:slug", articleHandler.DeleteArticle)

	slug := "non-existent"
	m.sqlMock.ExpectBegin()
	m.articleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(nil, gorm.ErrRecordNotFound)
	m.sqlMock.ExpectRollback()

	req, _ := http.NewRequest("DELETE", "/api/articles/"+slug, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	AssertAPIError(t, w, http.StatusNotFound, "article not found")
}
