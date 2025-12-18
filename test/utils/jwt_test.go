package utils

import (
	"go-gin-realworld-api/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTToken(t *testing.T) {
	userID := int64(123)
	email := "test@example.com"

	t.Run("Generate and Parse Success", func(t *testing.T) {
		token, err := utils.GenerateJWTToken(userID, email)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := utils.ParseJWTToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, email, claims.Email)
	})

	t.Run("Parse Invalid Token", func(t *testing.T) {
		claims, err := utils.ParseJWTToken("invalid.token.string")
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Parse Empty Token", func(t *testing.T) {
		claims, err := utils.ParseJWTToken("")
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}
