package services

import (
	"go-gin-realworld-api/internal/dtos"
	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/repository"
)

type ArticleService struct {
	articleRepo *repository.ArticleRepository
}

func NewArticleService(articleRepo *repository.ArticleRepository) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
	}
}

// ListArticles lists articles with optional filtering and pagination
func (s *ArticleService) ListArticles(query *dtos.ListArticlesQuery, currentUserID *int64) (*dtos.ArticlesListResponse, error) {
	// Set defaults for limit and offset
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	// Get articles from repository
	articles, total, err := s.articleRepo.ListArticles(query.Tag, query.Author, query.Favorited, query.Limit, query.Offset)
	if err != nil {
		return nil, err
	}

	// Convert articles to response DTOs
	articleResponses := make([]dtos.ArticleResponse, 0)
	for _, article := range articles {
		resp, err := s.articleToResponse(article, currentUserID)
		if err != nil {
			return nil, err
		}
		articleResponses = append(articleResponses, resp)
	}

	return &dtos.ArticlesListResponse{
		Articles:      articleResponses,
		ArticlesCount: int(total),
	}, nil
}

// articleToResponse converts a model Article to ArticleResponse DTO
func (s *ArticleService) articleToResponse(article *models.Article, currentUserID *int64) (dtos.ArticleResponse, error) {
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
