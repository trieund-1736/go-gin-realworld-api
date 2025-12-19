package handlers

import (
	"go-gin-realworld-api/internal/dtos"
	appErrors "go-gin-realworld-api/internal/errors"
	"go-gin-realworld-api/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// CreateComment handles creating a new comment
func (h *CommentHandler) CreateComment(c *gin.Context) {
	slug := c.Param("slug")

	// Get current user ID (required)
	userID, exists := c.Get("user_id")
	if !exists {
		appErrors.RespondError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var req dtos.CreateCommentRequest

	if appErrors.HandleBindError(c, c.ShouldBindJSON(&req)) {
		return
	}

	comment, err := h.commentService.CreateComment(c.Request.Context(), &req, slug, userID.(int64))
	if err != nil {
		appErrors.RespondError(c, http.StatusNotFound, "article not found")
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetComments handles getting all comments for an article
func (h *CommentHandler) GetComments(c *gin.Context) {
	slug := c.Param("slug")

	comments, err := h.commentService.GetCommentsByArticleSlug(c.Request.Context(), slug)
	if err != nil {
		appErrors.RespondError(c, http.StatusNotFound, "article not found")
		return
	}

	c.JSON(http.StatusOK, comments)
}

// DeleteComment handles deleting a comment
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")

	// Get current user ID (required)
	userID, exists := c.Get("user_id")
	if !exists {
		appErrors.RespondError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	// Parse comment ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "invalid comment id")
		return
	}

	if err := h.commentService.DeleteComment(c.Request.Context(), id, userID.(int64)); err != nil {
		if err.Error() == "forbidden" {
			appErrors.RespondError(c, http.StatusForbidden, "you can only delete your own comments")
			return
		}
		appErrors.RespondError(c, http.StatusNotFound, "comment not found")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
