package middleware

import (
	"go-gin-realworld-api/internal/middleware"
	"go-gin-realworld-api/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success with valid token", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.JWTAuthMiddleware())
		r.GET("/test", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			email, _ := c.Get("email")
			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"email":   email,
			})
		})

		token, _ := utils.GenerateJWTToken(1, "test@example.com")
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "\"user_id\":1")
		assert.Contains(t, w.Body.String(), "\"email\":\"test@example.com\"")
	})

	t.Run("Fail with missing header", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.JWTAuthMiddleware())
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "missing authorization header")
	})

	t.Run("Fail with invalid format", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.JWTAuthMiddleware())
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "InvalidToken")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid authorization header format")
	})

	t.Run("Fail with invalid token", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.JWTAuthMiddleware())
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid or expired token")
	})
}

func TestJWTOptionalAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success with valid token", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.JWTOptionalAuthMiddleware())
		r.GET("/test", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"exists":  exists,
			})
		})

		token, _ := utils.GenerateJWTToken(1, "test@example.com")
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "\"user_id\":1")
		assert.Contains(t, w.Body.String(), "\"exists\":true")
	})

	t.Run("Success without token", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.JWTOptionalAuthMiddleware())
		r.GET("/test", func(c *gin.Context) {
			_, exists := c.Get("user_id")
			c.JSON(http.StatusOK, gin.H{
				"exists": exists,
			})
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "\"exists\":false")
	})
}
