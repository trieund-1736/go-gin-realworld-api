package services

import (
	"errors"
	"go-gin-realworld-api/internal/dtos"
	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/repository"
	"time"
)

var ErrForbidden = errors.New("forbidden")

type CommentService struct {
	commentRepo *repository.CommentRepository
	articleRepo *repository.ArticleRepository
}

func NewCommentService(commentRepo *repository.CommentRepository, articleRepo *repository.ArticleRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
	}
}

// CreateComment creates a new comment
func (s *CommentService) CreateComment(req *dtos.CreateCommentRequest, slug string, authorID int64) (*dtos.CommentDetailResponse, error) {
	// Get article by slug to get article ID
	article, err := s.articleRepo.FindArticleBySlug(slug)
	if err != nil {
		return nil, err
	}

	comment := &models.Comment{
		Body:      req.Comment.Body,
		ArticleID: article.ID,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.commentRepo.CreateComment(comment); err != nil {
		return nil, err
	}

	// Fetch the created comment with preloaded data
	createdComment, err := s.commentRepo.GetCommentByID(comment.ID)
	if err != nil {
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
func (s *CommentService) GetCommentsByArticleSlug(slug string) (*dtos.CommentsListResponse, error) {
	// Get article by slug to get article ID
	article, err := s.articleRepo.FindArticleBySlug(slug)
	if err != nil {
		return nil, err
	}

	comments, err := s.commentRepo.GetCommentsByArticleID(article.ID)
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
func (s *CommentService) DeleteComment(id int64, currentUserID int64) error {
	// Get comment to check ownership
	comment, err := s.commentRepo.GetCommentByID(id)
	if err != nil {
		return err
	}

	// Check if current user is comment author or article author
	if comment.AuthorID != currentUserID && comment.Article.AuthorID != currentUserID {
		return ErrForbidden
	}

	return s.commentRepo.DeleteComment(id)
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
