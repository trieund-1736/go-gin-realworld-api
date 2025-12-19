package service

import (
	"context"
	"testing"

	appErrors "go-gin-realworld-api/internal/errors"
	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/services"
	"go-gin-realworld-api/test/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func setupFavoriteServiceTest(t *testing.T) (context.Context, *services.FavoriteService, *mocks.MockFavoriteRepository, *mocks.MockArticleRepository, sqlmock.Sqlmock) {
	mockFavoriteRepo := new(mocks.MockFavoriteRepository)
	mockArticleRepo := new(mocks.MockArticleRepository)
	gormDB, sqlMock := CreateMockDB(t)
	favoriteService := services.NewFavoriteService(gormDB, mockFavoriteRepo, mockArticleRepo)
	ctxForTest := context.Background()

	return ctxForTest, favoriteService, mockFavoriteRepo, mockArticleRepo, sqlMock
}

func TestFavoriteService_FavoriteArticle_Success_New(t *testing.T) {
	ctxForTest, favoriteService, mockFavoriteRepo, mockArticleRepo, sqlMock := setupFavoriteServiceTest(t)
	slug := "test-article"
	userID := int64(1)
	articleID := int64(10)

	article := &models.Article{
		ID:             articleID,
		Slug:           slug,
		FavoritesCount: 0,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockArticleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(article, nil)
	mockFavoriteRepo.On("IsFavorited", mock.Anything, userID, articleID).Return(false, nil)
	mockFavoriteRepo.On("AddFavorite", mock.Anything, userID, articleID).Return(nil)
	mockArticleRepo.On("UpdateArticle", mock.Anything, mock.MatchedBy(func(a *models.Article) bool {
		return a.ID == articleID && a.FavoritesCount == 1
	})).Return(nil)

	updatedArticle := &models.Article{
		ID:             articleID,
		Slug:           slug,
		FavoritesCount: 1,
		Author: &models.User{
			Username: "author1",
		},
		Favorites: []*models.Favorite{
			{UserID: userID, ArticleID: articleID},
		},
	}
	mockFavoriteRepo.On("GetArticleWithFavorites", mock.Anything, articleID).Return(updatedArticle, nil)

	resp, err := favoriteService.FavoriteArticle(ctxForTest, slug, userID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Article.Favorited)
	assert.Equal(t, 1, resp.Article.FavoritesCount)
	mockArticleRepo.AssertExpectations(t)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestFavoriteService_FavoriteArticle_Success_AlreadyFavorited(t *testing.T) {
	ctxForTest, favoriteService, mockFavoriteRepo, mockArticleRepo, sqlMock := setupFavoriteServiceTest(t)
	slug := "test-article"
	userID := int64(1)
	articleID := int64(10)

	article := &models.Article{
		ID:             articleID,
		Slug:           slug,
		FavoritesCount: 1,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockArticleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(article, nil)
	mockFavoriteRepo.On("IsFavorited", mock.Anything, userID, articleID).Return(true, nil)

	updatedArticle := &models.Article{
		ID:             articleID,
		Slug:           slug,
		FavoritesCount: 1,
		Author: &models.User{
			Username: "author1",
		},
		Favorites: []*models.Favorite{
			{UserID: userID, ArticleID: articleID},
		},
	}
	mockFavoriteRepo.On("GetArticleWithFavorites", mock.Anything, articleID).Return(updatedArticle, nil)

	resp, err := favoriteService.FavoriteArticle(ctxForTest, slug, userID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Article.Favorited)
	mockArticleRepo.AssertExpectations(t)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestFavoriteService_FavoriteArticle_NotFound(t *testing.T) {
	ctxForTest, favoriteService, _, mockArticleRepo, sqlMock := setupFavoriteServiceTest(t)
	slug := "non-existent"
	userID := int64(1)

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockArticleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(nil, gorm.ErrRecordNotFound)

	resp, err := favoriteService.FavoriteArticle(ctxForTest, slug, userID)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, appErrors.ErrNotFound, err)
	mockArticleRepo.AssertExpectations(t)
}

func TestFavoriteService_UnfavoriteArticle_Success(t *testing.T) {
	ctxForTest, favoriteService, mockFavoriteRepo, mockArticleRepo, sqlMock := setupFavoriteServiceTest(t)
	slug := "test-article"
	userID := int64(1)
	articleID := int64(10)

	article := &models.Article{
		ID:             articleID,
		Slug:           slug,
		FavoritesCount: 1,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockArticleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(article, nil)
	mockFavoriteRepo.On("IsFavorited", mock.Anything, userID, articleID).Return(true, nil)
	mockFavoriteRepo.On("RemoveFavorite", mock.Anything, userID, articleID).Return(nil)
	mockArticleRepo.On("UpdateArticle", mock.Anything, mock.MatchedBy(func(a *models.Article) bool {
		return a.ID == articleID && a.FavoritesCount == 0
	})).Return(nil)

	updatedArticle := &models.Article{
		ID:             articleID,
		Slug:           slug,
		FavoritesCount: 0,
		Author: &models.User{
			Username: "author1",
		},
		Favorites: []*models.Favorite{},
	}
	mockFavoriteRepo.On("GetArticleWithFavorites", mock.Anything, articleID).Return(updatedArticle, nil)

	resp, err := favoriteService.UnfavoriteArticle(ctxForTest, slug, userID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Article.Favorited)
	assert.Equal(t, 0, resp.Article.FavoritesCount)
	mockArticleRepo.AssertExpectations(t)
	mockFavoriteRepo.AssertExpectations(t)
}

func TestFavoriteService_UnfavoriteArticle_NotFound(t *testing.T) {
	ctxForTest, favoriteService, _, mockArticleRepo, sqlMock := setupFavoriteServiceTest(t)
	slug := "non-existent"
	userID := int64(1)

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockArticleRepo.On("FindArticleBySlug", mock.Anything, slug).Return(nil, gorm.ErrRecordNotFound)

	resp, err := favoriteService.UnfavoriteArticle(ctxForTest, slug, userID)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, appErrors.ErrNotFound, err)
	mockArticleRepo.AssertExpectations(t)
}
