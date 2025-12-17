package handlers

import (
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve tags"})
		return
	}
	c.JSON(200, gin.H{"tags": tags})
}
