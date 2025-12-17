package services

import (
	"go-gin-realworld-api/internal/dtos"
	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/repository"
	"go-gin-realworld-api/internal/utils"
	"time"
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

// GetFeedArticles gets articles from followed users
func (s *ArticleService) GetFeedArticles(userID int64, limit, offset int) (*dtos.ArticlesListResponse, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	articles, total, err := s.articleRepo.FeedArticles(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	articleResponses := make([]dtos.ArticleResponse, 0)
	for _, article := range articles {
		resp, err := s.articleToResponse(article, &userID)
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

// GetArticleBySlug gets article by slug
func (s *ArticleService) GetArticleBySlug(slug string, currentUserID *int64) (*dtos.ArticleDetailResponse, error) {
	article, err := s.articleRepo.FindArticleBySlug(slug)
	if err != nil {
		return nil, err
	}

	resp, err := s.articleToResponse(article, currentUserID)
	if err != nil {
		return nil, err
	}

	return &dtos.ArticleDetailResponse{
		Article: resp,
	}, nil
}

// CreateArticle creates a new article
func (s *ArticleService) CreateArticle(req *dtos.CreateArticleRequest, authorID int64) (*dtos.ArticleDetailResponse, error) {
	slug := utils.GenerateSlug(req.Article.Title)

	article := &models.Article{
		Slug:        slug,
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		AuthorID:    authorID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.articleRepo.CreateArticle(article); err != nil {
		return nil, err
	}

	// Add tags if provided
	if len(req.Article.TagList) > 0 {
		if err := s.articleRepo.AssignTagsToArticle(article.ID, req.Article.TagList); err != nil {
			return nil, err
		}
	}

	// Fetch the created article with preloaded data
	createdArticle, err := s.articleRepo.FindArticleBySlug(slug)
	if err != nil {
		return nil, err
	}

	resp, err := s.articleToResponse(createdArticle, &authorID)
	if err != nil {
		return nil, err
	}

	return &dtos.ArticleDetailResponse{
		Article: resp,
	}, nil
}

// UpdateArticle updates an article
func (s *ArticleService) UpdateArticle(slug string, req *dtos.UpdateArticleRequest, authorID int64) (*dtos.ArticleDetailResponse, error) {
	article, err := s.articleRepo.FindArticleBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Article.Title != "" {
		article.Title = req.Article.Title
		article.Slug = utils.GenerateSlug(req.Article.Title) // Regenerate slug from new title
	}
	if req.Article.Description != "" {
		article.Description = req.Article.Description
	}
	if req.Article.Body != "" {
		article.Body = req.Article.Body
	}

	article.UpdatedAt = time.Now()

	if err := s.articleRepo.UpdateArticle(article); err != nil {
		return nil, err
	}

	// Update tags if provided
	if len(req.Article.TagList) > 0 {
		if err := s.articleRepo.AssignTagsToArticle(article.ID, req.Article.TagList); err != nil {
			return nil, err
		}
	}

	// Fetch updated article
	updatedArticle, err := s.articleRepo.FindArticleBySlug(article.Slug)
	if err != nil {
		return nil, err
	}

	resp, err := s.articleToResponse(updatedArticle, &authorID)
	if err != nil {
		return nil, err
	}

	return &dtos.ArticleDetailResponse{
		Article: resp,
	}, nil
}

// DeleteArticle deletes an article
func (s *ArticleService) DeleteArticle(slug string) error {
	// Check if article exists first
	_, err := s.articleRepo.FindArticleBySlug(slug)
	if err != nil {
		return err
	}

	return s.articleRepo.DeleteArticleBySlug(slug)
}
