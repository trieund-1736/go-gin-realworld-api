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
	response, err := h.articleService.ListArticles(&query, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch articles"})
		return
	}

	c.JSON(http.StatusOK, response)
}
