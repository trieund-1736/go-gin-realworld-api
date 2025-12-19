package handlers

import (
	"go-gin-realworld-api/internal/dtos"
	appErrors "go-gin-realworld-api/internal/errors"
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
	if appErrors.HandleBindError(c, c.ShouldBindJSON(&req)) {
		return
	}

	// Login user
	user, token, err := h.authService.Login(c.Request.Context(), req.User.Email, req.User.Password)
	if err != nil {
		appErrors.RespondError(c, http.StatusUnauthorized, "Invalid email or password")
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
