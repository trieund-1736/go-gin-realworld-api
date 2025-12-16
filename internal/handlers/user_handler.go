package handlers

import (
	"go-gin-realworld-api/internal/dtos"
	"go-gin-realworld-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterUser handles user registration
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req dtos.RegisterUserRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Register user
	user, err := h.userService.RegisterUser(req.User.Username, req.User.Email, req.User.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to register user"})
		return
	}

	// Return user response
	resp := dtos.UserResponse{}
	resp.User.ID = user.ID
	resp.User.Username = user.Username
	resp.User.Email = user.Email

	c.JSON(http.StatusCreated, resp)
}

// GetCurrentUser handles getting current user information
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Get user from database
	user, err := h.userService.GetUserByID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Return user response
	resp := dtos.UserResponse{}
	resp.User.ID = user.ID
	resp.User.Username = user.Username
	resp.User.Email = user.Email

	c.JSON(http.StatusOK, resp)
}
