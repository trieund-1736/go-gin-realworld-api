package handlers

import (
	"go-gin-realworld-api/internal/dtos"
	appErrors "go-gin-realworld-api/internal/errors"
	"go-gin-realworld-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleService *services.ArticleService
}

func NewArticleHandler(articleService *services.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

// ListArticles handles listing articles with optional filtering and pagination
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	var query dtos.ListArticlesQuery

	// Bind query parameters
	if appErrors.HandleBindError(c, c.ShouldBindQuery(&query)) {
		return
	}

	// Get current user ID if authenticated
	var currentUserID *int64
	if userID, exists := c.Get("user_id"); exists {
		id := userID.(int64)
		currentUserID = &id
	}

	// Get articles from service
	response, err := h.articleService.ListArticles(c.Request.Context(), &query, currentUserID)
	if err != nil {
		appErrors.RespondError(c, http.StatusInternalServerError, "failed to fetch articles")
		return
	}

	c.JSON(http.StatusOK, response)
}

// FeedArticles handles getting feed of articles from followed users
func (h *ArticleHandler) FeedArticles(c *gin.Context) {
	// Get current user ID (required)
	userID, exists := c.Get("user_id")
	if !exists {
		appErrors.RespondError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	// Parse pagination parameters
	type FeedQuery struct {
		Limit  int `form:"limit,default=20"`
		Offset int `form:"offset,default=0"`
	}

	var query FeedQuery
	if appErrors.HandleBindError(c, c.ShouldBindQuery(&query)) {
		return
	}

	response, err := h.articleService.GetFeedArticles(c.Request.Context(), userID.(int64), query.Limit, query.Offset)
	if err != nil {
		appErrors.RespondError(c, http.StatusInternalServerError, "failed to fetch feed")
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetArticle handles getting a single article by slug
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	slug := c.Param("slug")

	// Get current user ID if authenticated
	var currentUserID *int64
	if userID, exists := c.Get("user_id"); exists {
		id := userID.(int64)
		currentUserID = &id
	}

	article, err := h.articleService.GetArticleBySlug(c.Request.Context(), slug, currentUserID)
	if err != nil {
		appErrors.RespondError(c, http.StatusNotFound, "article not found")
		return
	}

	c.JSON(http.StatusOK, article)
}

// CreateArticle handles creating a new article
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	// Get current user ID (required)
	userID, exists := c.Get("user_id")
	if !exists {
		appErrors.RespondError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var req dtos.CreateArticleRequest

	if appErrors.HandleBindError(c, c.ShouldBindJSON(&req)) {
		return
	}

	article, err := h.articleService.CreateArticle(c.Request.Context(), &req, userID.(int64))
	if err != nil {
		appErrors.RespondError(c, http.StatusInternalServerError, "failed to create article")
		return
	}

	c.JSON(http.StatusCreated, article)
}

// UpdateArticle handles updating an article
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	slug := c.Param("slug")

	// Get current user ID (required)
	userID, exists := c.Get("user_id")
	if !exists {
		appErrors.RespondError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	var req dtos.UpdateArticleRequest

	if appErrors.HandleBindError(c, c.ShouldBindJSON(&req)) {
		return
	}

	article, err := h.articleService.UpdateArticle(c.Request.Context(), slug, &req, userID.(int64))
	if err != nil {
		appErrors.RespondError(c, http.StatusNotFound, "article not found")
		return
	}

	c.JSON(http.StatusOK, article)
}

// DeleteArticle handles deleting an article
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	slug := c.Param("slug")

	// Get current user ID (required)
	_, exists := c.Get("user_id")
	if !exists {
		appErrors.RespondError(c, http.StatusUnauthorized, "authentication required")
		return
	}

	if err := h.articleService.DeleteArticle(c.Request.Context(), slug); err != nil {
		appErrors.RespondError(c, http.StatusNotFound, "article not found")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
