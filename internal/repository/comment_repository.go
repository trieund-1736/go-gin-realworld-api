package repository

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// CreateComment creates a new comment
func (r *CommentRepository) CreateComment(comment *models.Comment) error {
	if err := r.db.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

// GetCommentsByArticleID gets all comments for an article
func (r *CommentRepository) GetCommentsByArticleID(articleID int64) ([]*models.Comment, error) {
	var comments []*models.Comment
	if err := r.db.
		Where("article_id = ?", articleID).
		Preload("Author").
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// GetCommentByID gets a comment by ID with author and article preloaded
func (r *CommentRepository) GetCommentByID(id int64) (*models.Comment, error) {
	var comment *models.Comment
	if err := r.db.
		Preload("Author").
		Preload("Article").
		Where("id = ?", id).
		First(&comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// DeleteComment deletes a comment by ID
func (r *CommentRepository) DeleteComment(id int64) error {
	if err := r.db.Delete(&models.Comment{}, id).Error; err != nil {
		return err
	}
	return nil
}
