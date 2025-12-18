package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type FavoriteRepository interface {
	AddFavorite(db *gorm.DB, userID, articleID int64) error
	RemoveFavorite(db *gorm.DB, userID, articleID int64) error
	IsFavorited(db *gorm.DB, userID, articleID int64) (bool, error)
	GetArticleWithFavorites(db *gorm.DB, articleID int64) (*models.Article, error)
}
