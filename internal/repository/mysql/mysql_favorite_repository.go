package mysql

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type MySqlFavoriteRepository struct {
}

func NewMySqlFavoriteRepository() *MySqlFavoriteRepository {
	return &MySqlFavoriteRepository{}
}

// AddFavorite adds an article to user's favorites
func (r *MySqlFavoriteRepository) AddFavorite(db *gorm.DB, userID, articleID int64) error {
	favorite := &models.Favorite{
		UserID:    userID,
		ArticleID: articleID,
	}
	return db.Create(favorite).Error
}

// RemoveFavorite removes an article from user's favorites
func (r *MySqlFavoriteRepository) RemoveFavorite(db *gorm.DB, userID, articleID int64) error {
	return db.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&models.Favorite{}).Error
}

// IsFavorited checks if user has favorited an article
func (r *MySqlFavoriteRepository) IsFavorited(db *gorm.DB, userID, articleID int64) (bool, error) {
	var count int64
	err := db.Model(&models.Favorite{}).Where("user_id = ? AND article_id = ?", userID, articleID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetArticleWithFavorites gets article with favorites count and list (for checking if favorited)
func (r *MySqlFavoriteRepository) GetArticleWithFavorites(db *gorm.DB, articleID int64) (*models.Article, error) {
	var article *models.Article
	err := db.
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
