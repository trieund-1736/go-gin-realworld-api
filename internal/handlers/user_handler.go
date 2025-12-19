package handlers

import (
	"go-gin-realworld-api/internal/dtos"
	customErr "go-gin-realworld-api/internal/errors"
	"go-gin-realworld-api/internal/services"
	"log"
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
	if customErr.HandleBindError(c, c.ShouldBindJSON(&req)) {
		return
	}

	// Register user
	user, err := h.userService.RegisterUser(c.Request.Context(), req.User.Username, req.User.Email, req.User.Password)
	if err != nil {
		log.Println("Error registering user:", err)
		switch err {
		case customErr.ErrUserAlreadyExists:
			customErr.RespondError(c, http.StatusConflict, "User already exists")
		default:
			customErr.RespondError(c, http.StatusInternalServerError, "Failed to register user")
		}
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
		customErr.RespondError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	// Get user from database
	user, err := h.userService.GetUserByID(c.Request.Context(), userID.(int64))
	if err != nil {
		customErr.RespondError(c, http.StatusNotFound, "user not found")
		return
	}

	// Return user response
	resp := dtos.UserResponse{}
	resp.User.ID = user.ID
	resp.User.Username = user.Username
	resp.User.Email = user.Email

	c.JSON(http.StatusOK, resp)
}

// UpdateUser handles updating user information
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		customErr.RespondError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req dtos.UpdateUserRequest

	// Validate request body
	if customErr.HandleBindError(c, c.ShouldBindJSON(&req)) {
		return
	}

	// Update user
	user, err := h.userService.UpdateUser(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		customErr.RespondError(c, http.StatusInternalServerError, "failed to update user")
		return
	}

	// Return update response
	resp := dtos.UpdateUserResponse{}
	resp.User.ID = user.ID
	resp.User.Username = user.Username
	resp.User.Email = user.Email
	resp.User.Image = req.User.Image
	resp.User.Bio = req.User.Bio

	c.JSON(http.StatusOK, resp)
}
