package handlers

import (
	appErrors "go-gin-realworld-api/internal/errors"
	"go-gin-realworld-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetProfile handles getting a user profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	username := c.Param("username")

	// Get current user ID if authenticated (optional)
	currentUserID := int64(0)
	if userID, exists := c.Get("user_id"); exists {
		currentUserID = userID.(int64)
	}

	// Get profile
	profile, err := h.profileService.GetProfileByUsername(c.Request.Context(), username, currentUserID)
	if err != nil {
		appErrors.RespondError(c, http.StatusNotFound, "profile not found")
		return
	}

	c.JSON(http.StatusOK, profile)
}

// FollowUser handles following a user
func (h *ProfileHandler) FollowUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		appErrors.RespondError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	username := c.Param("username")

	// Follow user
	profile, err := h.profileService.FollowUser(c.Request.Context(), userID.(int64), username)
	if err != nil {
		appErrors.RespondError(c, http.StatusNotFound, "user not found")
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UnfollowUser handles unfollowing a user
func (h *ProfileHandler) UnfollowUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		appErrors.RespondError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	username := c.Param("username")

	// Unfollow user
	profile, err := h.profileService.UnfollowUser(c.Request.Context(), userID.(int64), username)
	if err != nil {
		appErrors.RespondError(c, http.StatusNotFound, "user not found")
		return
	}

	c.JSON(http.StatusOK, profile)
}
