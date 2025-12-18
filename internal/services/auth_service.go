package services

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/repository"
	"go-gin-realworld-api/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	db       *gorm.DB
	userRepo repository.UserRepository
}

func NewAuthService(db *gorm.DB, userRepo repository.UserRepository) *AuthService {
	return &AuthService{db: db, userRepo: userRepo}
}

// Login logs in a user and returns user with JWT token
func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	db := s.db.WithContext(ctx)
	// Find user by email
	user, err := s.userRepo.FindUserByEmail(db, email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Verify password
	hashedPassword := hashPassword(password)
	if user.Password != hashedPassword {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(user.ID, user.Email)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

// hashPassword hashes the password using SHA256
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}
