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
