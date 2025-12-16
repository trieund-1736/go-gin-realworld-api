package handlers

import (
	"go-gin-realworld-api/internal/dtos"
	"go-gin-realworld-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dtos.LoginRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Login user
	user, token, err := h.authService.Login(req.User.Email, req.User.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return login response with token
	resp := dtos.LoginResponse{}
	resp.User.ID = user.ID
	resp.User.Username = user.Username
	resp.User.Email = user.Email
	resp.User.Token = token

	c.JSON(http.StatusOK, resp)
}
