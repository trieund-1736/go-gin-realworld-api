package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

// AddFavorite adds an article to user's favorites
func (r *FavoriteRepository) AddFavorite(userID, articleID int64) error {
	favorite := &models.Favorite{
		UserID:    userID,
		ArticleID: articleID,
	}
	return r.db.Create(favorite).Error
}

// RemoveFavorite removes an article from user's favorites
func (r *FavoriteRepository) RemoveFavorite(userID, articleID int64) error {
	return r.db.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&models.Favorite{}).Error
}

// IsFavorited checks if user has favorited an article
func (r *FavoriteRepository) IsFavorited(userID, articleID int64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Favorite{}).Where("user_id = ? AND article_id = ?", userID, articleID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetArticleWithFavorites gets article with favorites count and list (for checking if favorited)
func (r *FavoriteRepository) GetArticleWithFavorites(articleID int64) (*models.Article, error) {
	var article *models.Article
	err := r.db.
		Preload("Author").
		Preload("ArticleTags.Tag").
		Preload("Favorites").
		Where("id = ?", articleID).
		First(&article).Error

	if err != nil {
		return nil, err
	}
	return article, nil
}
