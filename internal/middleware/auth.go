package middleware

import (
	appErrors "go-gin-realworld-api/internal/errors"
	"go-gin-realworld-api/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// extractAndValidateToken extracts and validates JWT token from Authorization header
// Returns (userID, email, error) with error constants from errors package
func extractAndValidateToken(c *gin.Context) (int64, string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, "", appErrors.ErrMissingAuthHeader
	}

	// Extract token from "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, "", appErrors.ErrInvalidAuthHeader
	}

	tokenString := parts[1]
	claims, err := utils.ParseJWTToken(tokenString)
	if err != nil {
		return 0, "", appErrors.ErrInvalidToken
	}

	return claims.UserID, claims.Email, nil
}

// JWTAuthMiddleware checks JWT token from Authorization header (required)
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, email, err := extractAndValidateToken(c)
		if err != nil {
			switch err {
			case appErrors.ErrMissingAuthHeader, appErrors.ErrInvalidAuthHeader, appErrors.ErrInvalidToken:
				appErrors.RespondError(c, http.StatusUnauthorized, err.Error())
			default:
				appErrors.RespondError(c, http.StatusUnauthorized, "authentication failed")
			}
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("email", email)
		c.Next()
	}
}

// JWTOptionalAuthMiddleware checks JWT token if provided, but doesn't require it
func JWTOptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, email, _ := extractAndValidateToken(c)
		if userID > 0 {
			c.Set("user_id", userID)
			c.Set("email", email)
		}
		c.Next()
	}
}
