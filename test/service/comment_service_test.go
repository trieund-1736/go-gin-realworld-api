package service

import (
	"context"
	"errors"
	"testing"

	"go-gin-realworld-api/internal/dtos"
	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/services"
	"go-gin-realworld-api/test/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupCommentServiceTest(t *testing.T) (context.Context, *services.CommentService, *mocks.MockCommentRepository, *mocks.MockArticleRepository, sqlmock.Sqlmock) {
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockArticleRepo := new(mocks.MockArticleRepository)
	gormDB, sqlMock := CreateMockDB(t)
	commentService := services.NewCommentService(gormDB, mockCommentRepo, mockArticleRepo)
	ctxForTest := context.Background()

	return ctxForTest, commentService, mockCommentRepo, mockArticleRepo, sqlMock
}

func TestCommentService_CreateComment_Success(t *testing.T) {
	ctxForTest, commentService, mockCommentRepo, mockArticleRepo, sqlMock := setupCommentServiceTest(t)
	slug := "test-article"
	authorID := int64(1)
	req := &dtos.CreateCommentRequest{
		Comment: struct {
			Body string `json:"body" binding:"required"`
		}{
			Body: "Test comment body",
		},
	}

	article := &models.Article{
		ID:   10,
		Slug: slug,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockArticleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(article, nil)
	mockCommentRepo.On("CreateComment", mock.Anything, mock.MatchedBy(func(c *models.Comment) bool {
		return c.Body == req.Comment.Body && c.ArticleID == article.ID && c.AuthorID == authorID
	})).Run(func(args mock.Arguments) {
		comment := args.Get(1).(*models.Comment)
		comment.ID = 100
	}).Return(nil)

	createdComment := &models.Comment{
		ID:        100,
		Body:      req.Comment.Body,
		ArticleID: article.ID,
		AuthorID:  authorID,
		Author: &models.User{
			Username: "commenter",
		},
	}
	mockCommentRepo.On("GetCommentByID", mock.Anything, int64(100)).Return(createdComment, nil)

	resp, err := commentService.CreateComment(ctxForTest, req, slug, authorID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(100), resp.Comment.ID)
	assert.Equal(t, req.Comment.Body, resp.Comment.Body)
	mockArticleRepo.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentService_CreateComment_ArticleNotFound(t *testing.T) {
	ctxForTest, commentService, _, mockArticleRepo, sqlMock := setupCommentServiceTest(t)
	slug := "non-existent"
	authorID := int64(1)
	req := &dtos.CreateCommentRequest{}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	expectedError := errors.New("article not found")
	mockArticleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(nil, expectedError)

	resp, err := commentService.CreateComment(ctxForTest, req, slug, authorID)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedError, err)
	mockArticleRepo.AssertExpectations(t)
}

func TestCommentService_GetCommentsByArticleSlug_Success(t *testing.T) {
	ctxForTest, commentService, mockCommentRepo, mockArticleRepo, _ := setupCommentServiceTest(t)
	slug := "test-article"
	articleID := int64(10)

	article := &models.Article{
		ID:   articleID,
		Slug: slug,
	}

	comments := []*models.Comment{
		{
			ID:   1,
			Body: "Comment 1",
			Author: &models.User{
				Username: "user1",
			},
		},
		{
			ID:   2,
			Body: "Comment 2",
			Author: &models.User{
				Username: "user2",
			},
		},
	}

	mockArticleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(article, nil)
	mockCommentRepo.On("GetCommentsByArticleID", mock.Anything, articleID).Return(comments, nil)

	resp, err := commentService.GetCommentsByArticleSlug(ctxForTest, slug)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Comments, 2)
	assert.Equal(t, "Comment 1", resp.Comments[0].Body)
	mockArticleRepo.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentService_DeleteComment_Success(t *testing.T) {
	ctxForTest, commentService, mockCommentRepo, _, sqlMock := setupCommentServiceTest(t)
	commentID := int64(100)
	currentUserID := int64(1)

	comment := &models.Comment{
		ID:       commentID,
		AuthorID: currentUserID, // Current user is the author
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockCommentRepo.On("GetCommentByID", mock.Anything, commentID).Return(comment, nil)
	mockCommentRepo.On("DeleteComment", mock.Anything, commentID).Return(nil)

	err := commentService.DeleteComment(ctxForTest, commentID, currentUserID)

	assert.NoError(t, err)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentService_DeleteComment_Forbidden(t *testing.T) {
	ctxForTest, commentService, mockCommentRepo, _, sqlMock := setupCommentServiceTest(t)
	commentID := int64(100)
	currentUserID := int64(1)

	comment := &models.Comment{
		ID:       commentID,
		AuthorID: 999, // Someone else's comment
		Article: &models.Article{
			AuthorID: 888, // Someone else's article
		},
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockCommentRepo.On("GetCommentByID", mock.Anything, commentID).Return(comment, nil)

	err := commentService.DeleteComment(ctxForTest, commentID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, services.ErrForbidden, err)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentService_DeleteComment_NotFound(t *testing.T) {
	ctxForTest, commentService, mockCommentRepo, _, sqlMock := setupCommentServiceTest(t)
	commentID := int64(100)
	currentUserID := int64(1)

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	expectedError := errors.New("comment not found")
	mockCommentRepo.On("GetCommentByID", mock.Anything, commentID).Return(nil, expectedError)

	err := commentService.DeleteComment(ctxForTest, commentID, currentUserID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockCommentRepo.AssertExpectations(t)
}
