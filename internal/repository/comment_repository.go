package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type CommentRepository struct {
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{}
}

// CreateComment creates a new comment
func (r *CommentRepository) CreateComment(db *gorm.DB, comment *models.Comment) error {
	if err := db.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

// GetCommentsByArticleID gets all comments for an article
func (r *CommentRepository) GetCommentsByArticleID(db *gorm.DB, articleID int64) ([]*models.Comment, error) {
	var comments []*models.Comment
	if err := db.
		Where("article_id = ?", articleID).
		Preload("Author").
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// GetCommentByID gets a comment by ID with author and article preloaded
func (r *CommentRepository) GetCommentByID(db *gorm.DB, id int64) (*models.Comment, error) {
	var comment *models.Comment
	if err := db.
		Preload("Author").
		Preload("Article").
		Where("id = ?", id).
		First(&comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// DeleteComment deletes a comment by ID
func (r *CommentRepository) DeleteComment(db *gorm.DB, id int64) error {
	if err := db.Delete(&models.Comment{}, id).Error; err != nil {
		return err
	}
	return nil
}
