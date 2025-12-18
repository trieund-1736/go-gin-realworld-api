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

type commentHandlerMocks struct {
	commentRepo *mocks.MockCommentRepository
	articleRepo *mocks.MockArticleRepository
	sqlMock     sqlmock.Sqlmock
}

func setupCommentHandlerTest(t *testing.T) (*gin.Engine, *handlers.CommentHandler, commentHandlerMocks) {
	m := commentHandlerMocks{
		commentRepo: new(mocks.MockCommentRepository),
		articleRepo: new(mocks.MockArticleRepository),
	}

	mockDB, sqlMock := CreateMockDB(t)
	m.sqlMock = sqlMock
	commentService := services.NewCommentService(mockDB, m.commentRepo, m.articleRepo)
	commentHandler := handlers.NewCommentHandler(commentService)

	router := SetupRouter()
	return router, commentHandler, m
}

func TestCommentHandler_CreateComment_Success(t *testing.T) {
	router, commentHandler, m := setupCommentHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.POST("/api/articles/:slug/comments", commentHandler.CreateComment)

	slug := "test-article"
	reqBody := dtos.CreateCommentRequest{}
	reqBody.Comment.Body = "Test comment body"
	jsonBody, _ := json.Marshal(reqBody)

	article := &models.Article{ID: 1, Slug: slug}
	comment := &models.Comment{
		ID:        1,
		Body:      "Test comment body",
		ArticleID: 1,
		AuthorID:  1,
		Author:    &models.User{Username: "commenter"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	m.sqlMock.ExpectBegin()
	m.articleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(article, nil)
	m.commentRepo.On("CreateComment", mock.Anything, mock.AnythingOfType("*models.Comment")).Return(nil).Run(func(args mock.Arguments) {
		c := args.Get(1).(*models.Comment)
		c.ID = 1
	})
	m.commentRepo.On("GetCommentByID", mock.Anything, int64(1)).Return(comment, nil)
	m.sqlMock.ExpectCommit()

	req, _ := http.NewRequest("POST", "/api/articles/"+slug+"/comments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dtos.CommentDetailResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.Comment.ID)
	assert.Equal(t, "Test comment body", resp.Comment.Body)

	m.articleRepo.AssertExpectations(t)
	m.commentRepo.AssertExpectations(t)
}

func TestCommentHandler_CreateComment_Unauthorized(t *testing.T) {
	router, commentHandler, _ := setupCommentHandlerTest(t)
	router.POST("/api/articles/:slug/comments", commentHandler.CreateComment)

	req, _ := http.NewRequest("POST", "/api/articles/test-article/comments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "authentication required")
}

func TestCommentHandler_GetComments_Success(t *testing.T) {
	router, commentHandler, m := setupCommentHandlerTest(t)
	router.GET("/api/articles/:slug/comments", commentHandler.GetComments)

	slug := "test-article"
	article := &models.Article{ID: 1, Slug: slug}
	comments := []*models.Comment{
		{
			ID:        1,
			Body:      "Comment 1",
			Author:    &models.User{Username: "user1"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	m.articleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(article, nil)
	m.commentRepo.On("GetCommentsByArticleID", mock.Anything, int64(1)).Return(comments, nil)

	req, _ := http.NewRequest("GET", "/api/articles/"+slug+"/comments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dtos.CommentsListResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp.Comments, 1)
	assert.Equal(t, "Comment 1", resp.Comments[0].Body)

	m.articleRepo.AssertExpectations(t)
	m.commentRepo.AssertExpectations(t)
}

func TestCommentHandler_DeleteComment_Success(t *testing.T) {
	router, commentHandler, m := setupCommentHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.DELETE("/api/articles/:slug/comments/:id", commentHandler.DeleteComment)

	commentID := int64(1)
	comment := &models.Comment{
		ID:       commentID,
		AuthorID: 1, // Current user is author
		Article:  &models.Article{AuthorID: 2},
	}

	m.sqlMock.ExpectBegin()
	m.commentRepo.On("GetCommentByID", mock.Anything, commentID).Return(comment, nil)
	m.commentRepo.On("DeleteComment", mock.Anything, commentID).Return(nil)
	m.sqlMock.ExpectCommit()

	req, _ := http.NewRequest("DELETE", "/api/articles/test-article/comments/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	m.commentRepo.AssertExpectations(t)
}

func TestCommentHandler_DeleteComment_Forbidden(t *testing.T) {
	router, commentHandler, m := setupCommentHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.DELETE("/api/articles/:slug/comments/:id", commentHandler.DeleteComment)

	commentID := int64(1)
	comment := &models.Comment{
		ID:       commentID,
		AuthorID: 2,                            // Not current user
		Article:  &models.Article{AuthorID: 3}, // Not article author either
	}

	m.sqlMock.ExpectBegin()
	m.commentRepo.On("GetCommentByID", mock.Anything, commentID).Return(comment, nil)
	m.sqlMock.ExpectRollback()

	req, _ := http.NewRequest("DELETE", "/api/articles/test-article/comments/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "you can only delete your own comments")

	m.commentRepo.AssertExpectations(t)
}

func TestCommentHandler_DeleteComment_NotFound(t *testing.T) {
	router, commentHandler, m := setupCommentHandlerTest(t)

	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.DELETE("/api/articles/:slug/comments/:id", commentHandler.DeleteComment)

	commentID := int64(1)
	m.sqlMock.ExpectBegin()
	m.commentRepo.On("GetCommentByID", mock.Anything, commentID).Return(nil, gorm.ErrRecordNotFound)
	m.sqlMock.ExpectRollback()

	req, _ := http.NewRequest("DELETE", "/api/articles/test-article/comments/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "comment not found")

	m.commentRepo.AssertExpectations(t)
}
