package services

import (
	"context"
	"errors"
	"go-gin-realworld-api/internal/dtos"
	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/repository"
	"time"

	"gorm.io/gorm"
)

var ErrForbidden = errors.New("forbidden")

type CommentService struct {
	db          *gorm.DB
	commentRepo repository.CommentRepository
	articleRepo repository.ArticleRepository
}

func NewCommentService(db *gorm.DB, commentRepo repository.CommentRepository, articleRepo repository.ArticleRepository) *CommentService {
	return &CommentService{
		db:          db,
		commentRepo: commentRepo,
		articleRepo: articleRepo,
	}
}

// CreateComment creates a new comment
func (s *CommentService) CreateComment(ctx context.Context, req *dtos.CreateCommentRequest, slug string, authorID int64) (*dtos.CommentDetailResponse, error) {
	db := s.db.WithContext(ctx)

	var createdComment *models.Comment
	if err := db.Transaction(func(tx *gorm.DB) error {
		article, err := s.articleRepo.FindArticleBySlug(tx, slug)
		if err != nil {
			return err
		}

		comment := &models.Comment{
			Body:      req.Comment.Body,
			ArticleID: article.ID,
			AuthorID:  authorID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.commentRepo.CreateComment(tx, comment); err != nil {
			return err
		}

		createdComment, err = s.commentRepo.GetCommentByID(tx, comment.ID)
		return err
	}); err != nil {
		return nil, err
	}

	resp, err := s.commentToResponse(createdComment)
	if err != nil {
		return nil, err
	}

	return &dtos.CommentDetailResponse{
		Comment: resp,
	}, nil
}

// GetCommentsByArticleSlug gets all comments for an article by slug
func (s *CommentService) GetCommentsByArticleSlug(ctx context.Context, slug string) (*dtos.CommentsListResponse, error) {
	db := s.db.WithContext(ctx)
	// Get article by slug to get article ID
	article, err := s.articleRepo.FindArticleBySlug(db, slug)
	if err != nil {
		return nil, err
	}

	comments, err := s.commentRepo.GetCommentsByArticleID(db, article.ID)
	if err != nil {
		return nil, err
	}

	commentResponses := make([]dtos.CommentResponse, 0)
	for _, comment := range comments {
		resp, err := s.commentToResponse(comment)
		if err != nil {
			return nil, err
		}
		commentResponses = append(commentResponses, resp)
	}

	return &dtos.CommentsListResponse{
		Comments: commentResponses,
	}, nil
}

// DeleteComment deletes a comment by ID. Only comment author or article author can delete.
func (s *CommentService) DeleteComment(ctx context.Context, id int64, currentUserID int64) error {
	db := s.db.WithContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		comment, err := s.commentRepo.GetCommentByID(tx, id)
		if err != nil {
			return err
		}

		// Check if current user is comment author or article author
		if comment.AuthorID != currentUserID && comment.Article.AuthorID != currentUserID {
			return ErrForbidden
		}

		return s.commentRepo.DeleteComment(tx, id)
	})
}

// commentToResponse converts a model Comment to CommentResponse DTO
func (s *CommentService) commentToResponse(comment *models.Comment) (dtos.CommentResponse, error) {
	return dtos.CommentResponse{
		ID:        comment.ID,
		CreatedAt: comment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: comment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Body:      comment.Body,
		Author: dtos.CommentAuthorResponse{
			Username: comment.Author.Username,
		},
	}, nil
}
