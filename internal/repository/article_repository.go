package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	ListArticles(db *gorm.DB, tag, author string, favorited *bool, currentUserID *int64, limit, offset int) ([]*models.Article, int64, error)
	FeedArticles(db *gorm.DB, userID int64, limit, offset int) ([]*models.Article, int64, error)
	FindArticleBySlug(db *gorm.DB, slug string) (*models.Article, error)
	CreateArticle(db *gorm.DB, article *models.Article) error
	UpdateArticle(db *gorm.DB, article *models.Article) error
	DeleteArticleBySlug(db *gorm.DB, slug string) error
	AssignTagsToArticle(db *gorm.DB, articleID int64, tagNames []string) error
}
