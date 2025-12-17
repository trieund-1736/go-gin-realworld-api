package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// ListArticles lists articles with optional filtering and pagination
func (r *ArticleRepository) ListArticles(tag, author, favorited string, limit, offset int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	query := r.db

	// Filter by tag
	if tag != "" {
		query = query.
			Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Joins("JOIN tags ON tags.id = article_tags.tag_id").
			Where("tags.name = ?", tag)
	}

	// Filter by author
	if author != "" {
		query = query.
			Joins("JOIN users ON users.id = articles.author_id").
			Where("users.username = ?", author)
	}

	// Filter by favorited (user who favorited the article)
	if favorited != "" {
		query = query.
			Joins("JOIN favorites ON favorites.article_id = articles.id").
			Joins("JOIN users ON users.id = favorites.user_id").
			Where("users.username = ?", favorited)
	}

	// Get total count before pagination
	if err := query.Model(&models.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting and pagination
	if err := query.
		Preload("Author").
		Preload("ArticleTags.Tag").
		Preload("Favorites").
		Order("articles.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}
