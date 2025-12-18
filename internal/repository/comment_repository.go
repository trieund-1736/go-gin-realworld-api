package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(db *gorm.DB, comment *models.Comment) error
	GetCommentsByArticleID(db *gorm.DB, articleID int64) ([]*models.Comment, error)
	GetCommentByID(db *gorm.DB, id int64) (*models.Comment, error)
	DeleteComment(db *gorm.DB, id int64) error
}
