package handlers

import (
	appErrors "go-gin-realworld-api/internal/errors"
	"go-gin-realworld-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	favoriteService *services.FavoriteService
}

func NewFavoriteHandler(favoriteService *services.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
	}
}

// FavoriteArticle adds article to user's favorites
// POST /api/articles/:slug/favorite
func (h *FavoriteHandler) FavoriteArticle(c *gin.Context) {
	slug := c.Param("slug")

	// Get current user ID from context (set by JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		appErrors.RespondError(c, http.StatusUnauthorized, "missing authorization")
		return
	}

	currentUserID := userID.(int64)

	// Call service to favorite article
	result, err := h.favoriteService.FavoriteArticle(c.Request.Context(), slug, currentUserID)
	if err != nil {
		if err.Error() == "article not found" {
			appErrors.RespondError(c, http.StatusNotFound, "article not found")
			return
		}
		appErrors.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// UnfavoriteArticle removes article from user's favorites
// DELETE /api/articles/:slug/favorite
func (h *FavoriteHandler) UnfavoriteArticle(c *gin.Context) {
	slug := c.Param("slug")

	// Get current user ID from context (set by JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		appErrors.RespondError(c, http.StatusUnauthorized, "missing authorization")
		return
	}

	currentUserID := userID.(int64)

	// Call service to unfavorite article
	result, err := h.favoriteService.UnfavoriteArticle(c.Request.Context(), slug, currentUserID)
	if err != nil {
		if err.Error() == "article not found" {
			appErrors.RespondError(c, http.StatusNotFound, "article not found")
			return
		}
		appErrors.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
