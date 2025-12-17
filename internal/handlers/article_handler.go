package handlers

import (
	"go-gin-realworld-api/internal/dtos"
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
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch articles"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// FeedArticles handles getting feed of articles from followed users
func (h *ArticleHandler) FeedArticles(c *gin.Context) {
	// Get current user ID (required)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	// Parse pagination parameters
	type FeedQuery struct {
		Limit  int `form:"limit,default=20"`
		Offset int `form:"offset,default=0"`
	}

	var query FeedQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	response, err := h.articleService.GetFeedArticles(c.Request.Context(), userID.(int64), query.Limit, query.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch feed"})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	c.JSON(http.StatusOK, article)
}

// CreateArticle handles creating a new article
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	// Get current user ID (required)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req dtos.CreateArticleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	article, err := h.articleService.CreateArticle(c.Request.Context(), &req, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create article"})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req dtos.UpdateArticleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	article, err := h.articleService.UpdateArticle(c.Request.Context(), slug, &req, userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	if err := h.articleService.DeleteArticle(c.Request.Context(), slug); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
