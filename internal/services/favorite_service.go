package services

import (
	"errors"
	"go-gin-realworld-api/internal/dtos"
	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/repository"

	"gorm.io/gorm"
)

type FavoriteService struct {
	favoriteRepo *repository.FavoriteRepository
	articleRepo  *repository.ArticleRepository
}

func NewFavoriteService(favoriteRepo *repository.FavoriteRepository, articleRepo *repository.ArticleRepository) *FavoriteService {
	return &FavoriteService{
		favoriteRepo: favoriteRepo,
		articleRepo:  articleRepo,
	}
}

// FavoriteArticle adds article to user's favorites
func (s *FavoriteService) FavoriteArticle(slug string, userID int64) (*dtos.ArticleDetailResponse, error) {
	// Get article by slug
	article, err := s.articleRepo.FindArticleBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("article not found")
		}
		return nil, err
	}

	// Check if already favorited
	isFavorited, err := s.favoriteRepo.IsFavorited(userID, article.ID)
	if err != nil {
		return nil, err
	}

	if !isFavorited {
		// Add favorite
		if err := s.favoriteRepo.AddFavorite(userID, article.ID); err != nil {
			return nil, err
		}

		// Increment favorites count
		article.FavoritesCount++
		if err := s.articleRepo.UpdateArticle(article); err != nil {
			return nil, err
		}
	}

	// Get updated article with favorites info
	updatedArticle, err := s.favoriteRepo.GetArticleWithFavorites(article.ID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	resp, err := articleToResponse(updatedArticle, &userID)
	if err != nil {
		return nil, err
	}

	return &dtos.ArticleDetailResponse{
		Article: resp,
	}, nil
}

// UnfavoriteArticle removes article from user's favorites
func (s *FavoriteService) UnfavoriteArticle(slug string, userID int64) (*dtos.ArticleDetailResponse, error) {
	// Get article by slug
	article, err := s.articleRepo.FindArticleBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("article not found")
		}
		return nil, err
	}

	// Check if favorited
	isFavorited, err := s.favoriteRepo.IsFavorited(userID, article.ID)
	if err != nil {
		return nil, err
	}

	if isFavorited {
		// Remove favorite
		if err := s.favoriteRepo.RemoveFavorite(userID, article.ID); err != nil {
			return nil, err
		}

		// Decrement favorites count
		if article.FavoritesCount > 0 {
			article.FavoritesCount--
			if err := s.articleRepo.UpdateArticle(article); err != nil {
				return nil, err
			}
		}
	}

	// Get updated article with favorites info
	updatedArticle, err := s.favoriteRepo.GetArticleWithFavorites(article.ID)
	if err != nil {
		return nil, err
	}

	// Convert to response
	resp, err := articleToResponse(updatedArticle, &userID)
	if err != nil {
		return nil, err
	}

	return &dtos.ArticleDetailResponse{
		Article: resp,
	}, nil
}

// articleToResponse converts a model Article to ArticleResponse DTO (reuse from ArticleService)
func articleToResponse(article *models.Article, currentUserID *int64) (dtos.ArticleResponse, error) {
	// Convert tags from preloaded ArticleTags
	tagList := make([]string, 0)
	if article.ArticleTags != nil {
		for _, at := range article.ArticleTags {
			if at.Tag != nil {
				tagList = append(tagList, at.Tag.Name)
			}
		}
	}

	// Check if current user favorited this article
	favorited := false
	if currentUserID != nil {
		for _, fav := range article.Favorites {
			if fav.UserID == *currentUserID {
				favorited = true
				break
			}
		}
	}

	return dtos.ArticleResponse{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        tagList,
		CreatedAt:      article.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      article.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Favorited:      favorited,
		FavoritesCount: article.FavoritesCount,
		Author: dtos.ArticleAuthorResponse{
			Username: article.Author.Username,
		},
	}, nil
}
