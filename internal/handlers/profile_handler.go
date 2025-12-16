package handlers

import (
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
	profile, err := h.profileService.GetProfileByUsername(username, currentUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// FollowUser handles following a user
func (h *ProfileHandler) FollowUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	username := c.Param("username")

	// Follow user
	profile, err := h.profileService.FollowUser(userID.(int64), username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UnfollowUser handles unfollowing a user
func (h *ProfileHandler) UnfollowUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	username := c.Param("username")

	// Unfollow user
	profile, err := h.profileService.UnfollowUser(userID.(int64), username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}
