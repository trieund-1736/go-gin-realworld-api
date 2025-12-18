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

func setupArticleServiceTest(t *testing.T) (context.Context, *services.ArticleService, *mocks.MockArticleRepository, sqlmock.Sqlmock) {
	mockArticleRepo := new(mocks.MockArticleRepository)
	gormDB, sqlMock := CreateMockDB(t)
	articleService := services.NewArticleService(gormDB, mockArticleRepo)
	ctxForTest := context.Background()

	return ctxForTest, articleService, mockArticleRepo, sqlMock
}

func TestArticleService_ListArticles_Success(t *testing.T) {
	ctxForTest, articleService, mockArticleRepo, _ := setupArticleServiceTest(t)
	query := &dtos.ListArticlesQuery{
		Limit:  20,
		Offset: 0,
	}
	currentUserID := int64(1)

	articles := []*models.Article{
		{
			ID:          1,
			Slug:        "test-article",
			Title:       "Test Article",
			Description: "Description",
			Body:        "Body",
			Author: &models.User{
				Username: "author1",
			},
		},
	}
	total := int64(1)

	mockArticleRepo.On("ListArticles", mock.Anything, "", "", (*bool)(nil), &currentUserID, 20, 0).Return(articles, total, nil)

	resp, err := articleService.ListArticles(ctxForTest, query, &currentUserID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.ArticlesCount)
	assert.Equal(t, "test-article", resp.Articles[0].Slug)
	mockArticleRepo.AssertExpectations(t)
}

func TestArticleService_GetArticleBySlug_Success(t *testing.T) {
	ctxForTest, articleService, mockArticleRepo, _ := setupArticleServiceTest(t)
	slug := "test-article"
	currentUserID := int64(1)

	article := &models.Article{
		ID:          1,
		Slug:        slug,
		Title:       "Test Article",
		Description: "Description",
		Body:        "Body",
		Author: &models.User{
			Username: "author1",
		},
	}

	mockArticleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(article, nil)

	resp, err := articleService.GetArticleBySlug(ctxForTest, slug, &currentUserID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, slug, resp.Article.Slug)
	mockArticleRepo.AssertExpectations(t)
}

func TestArticleService_GetArticleBySlug_Error(t *testing.T) {
	ctxForTest, articleService, mockArticleRepo, _ := setupArticleServiceTest(t)
	slug := "non-existent"
	expectedError := errors.New("article not found")

	mockArticleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(nil, expectedError)

	resp, err := articleService.GetArticleBySlug(ctxForTest, slug, nil)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedError, err)
	mockArticleRepo.AssertExpectations(t)
}

func TestArticleService_GetFeedArticles_Success(t *testing.T) {
	ctxForTest, articleService, mockArticleRepo, _ := setupArticleServiceTest(t)
	userID := int64(1)
	limit := 20
	offset := 0

	articles := []*models.Article{
		{
			ID:    1,
			Slug:  "feed-article",
			Title: "Feed Article",
			Author: &models.User{
				Username: "followed_user",
			},
		},
	}
	total := int64(1)

	mockArticleRepo.On("FeedArticles", mock.Anything, userID, limit, offset).Return(articles, total, nil)

	resp, err := articleService.GetFeedArticles(ctxForTest, userID, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.ArticlesCount)
	assert.Equal(t, "feed-article", resp.Articles[0].Slug)
	mockArticleRepo.AssertExpectations(t)
}

func TestArticleService_CreateArticle_Success(t *testing.T) {
	ctxForTest, articleService, mockArticleRepo, sqlMock := setupArticleServiceTest(t)
	authorID := int64(1)
	req := &dtos.CreateArticleRequest{
		Article: struct {
			Title       string   `json:"title" binding:"required"`
			Description string   `json:"description" binding:"required"`
			Body        string   `json:"body" binding:"required"`
			TagList     []string `json:"tagList"`
		}{
			Title:       "New Article",
			Description: "Description",
			Body:        "Body",
			TagList:     []string{"tag1", "tag2"},
		},
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockArticleRepo.On("CreateArticle", mock.Anything, mock.MatchedBy(func(a *models.Article) bool {
		return a.Title == req.Article.Title && a.AuthorID == authorID
	})).Run(func(args mock.Arguments) {
		article := args.Get(1).(*models.Article)
		article.ID = 1
	}).Return(nil)

	mockArticleRepo.On("AssignTagsToArticle", mock.Anything, int64(1), req.Article.TagList).Return(nil)

	createdArticle := &models.Article{
		ID:          1,
		Slug:        "new-article",
		Title:       "New Article",
		Description: "Description",
		Body:        "Body",
		AuthorID:    authorID,
		Author: &models.User{
			Username: "author1",
		},
	}
	mockArticleRepo.On("FindArticleBySlug", mock.Anything, "new-article").Return(createdArticle, nil)

	resp, err := articleService.CreateArticle(ctxForTest, req, authorID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "new-article", resp.Article.Slug)
	mockArticleRepo.AssertExpectations(t)
}

func TestArticleService_UpdateArticle_Success(t *testing.T) {
	ctxForTest, articleService, mockArticleRepo, sqlMock := setupArticleServiceTest(t)
	authorID := int64(1)
	slug := "old-article"
	req := &dtos.UpdateArticleRequest{
		Article: struct {
			Title       string   `json:"title"`
			Description string   `json:"description"`
			Body        string   `json:"body"`
			TagList     []string `json:"tagList"`
		}{
			Title: "Updated Title",
		},
	}

	existingArticle := &models.Article{
		ID:       1,
		Slug:     slug,
		Title:    "Old Title",
		AuthorID: authorID,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockArticleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(existingArticle, nil)
	mockArticleRepo.On("UpdateArticle", mock.Anything, mock.MatchedBy(func(a *models.Article) bool {
		return a.Title == "Updated Title" && a.Slug == "updated-title"
	})).Return(nil)

	updatedArticle := &models.Article{
		ID:       1,
		Slug:     "updated-title",
		Title:    "Updated Title",
		AuthorID: authorID,
		Author: &models.User{
			Username: "author1",
		},
	}
	mockArticleRepo.On("FindArticleBySlug", mock.Anything, "updated-title").Return(updatedArticle, nil)

	resp, err := articleService.UpdateArticle(ctxForTest, slug, req, authorID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "updated-title", resp.Article.Slug)
	mockArticleRepo.AssertExpectations(t)
}

func TestArticleService_DeleteArticle_Success(t *testing.T) {
	ctxForTest, articleService, mockArticleRepo, sqlMock := setupArticleServiceTest(t)
	slug := "to-delete"

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockArticleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(&models.Article{ID: 1, Slug: slug}, nil)
	mockArticleRepo.On("DeleteArticleBySlug", mock.Anything, slug).Return(nil)

	err := articleService.DeleteArticle(ctxForTest, slug)

	assert.NoError(t, err)
	mockArticleRepo.AssertExpectations(t)
}

func TestArticleService_CreateArticle_Error(t *testing.T) {
	ctxForTest, articleService, mockArticleRepo, sqlMock := setupArticleServiceTest(t)
	authorID := int64(1)
	req := &dtos.CreateArticleRequest{
		Article: struct {
			Title       string   `json:"title" binding:"required"`
			Description string   `json:"description" binding:"required"`
			Body        string   `json:"body" binding:"required"`
			TagList     []string `json:"tagList"`
		}{
			Title: "New Article",
		},
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	expectedError := errors.New("db error")
	mockArticleRepo.On("CreateArticle", mock.Anything, mock.Anything).Return(expectedError)

	resp, err := articleService.CreateArticle(ctxForTest, req, authorID)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedError, err)
	mockArticleRepo.AssertExpectations(t)
}

func TestArticleService_UpdateArticle_NotFound(t *testing.T) {
	ctxForTest, articleService, mockArticleRepo, sqlMock := setupArticleServiceTest(t)
	authorID := int64(1)
	slug := "non-existent"
	req := &dtos.UpdateArticleRequest{}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	expectedError := errors.New("article not found")
	mockArticleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(nil, expectedError)

	resp, err := articleService.UpdateArticle(ctxForTest, slug, req, authorID)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedError, err)
	mockArticleRepo.AssertExpectations(t)
}

func TestArticleService_DeleteArticle_NotFound(t *testing.T) {
	ctxForTest, articleService, mockArticleRepo, sqlMock := setupArticleServiceTest(t)
	slug := "non-existent"

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	expectedError := errors.New("article not found")
	mockArticleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(nil, expectedError)

	err := articleService.DeleteArticle(ctxForTest, slug)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockArticleRepo.AssertExpectations(t)
}
