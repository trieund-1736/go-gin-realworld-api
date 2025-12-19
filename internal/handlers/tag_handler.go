package handlers

import (
	appErrors "go-gin-realworld-api/internal/errors"
	"go-gin-realworld-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagService *services.TagService
}

func NewTagHandler(tagService *services.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// GetTags handles the request to get all unique tags
func (h *TagHandler) GetTags(c *gin.Context) {
	tags, err := h.tagService.GetAllTags(c.Request.Context())
	if err != nil {
		appErrors.RespondError(c, http.StatusInternalServerError, "failed to retrieve tags")
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}
