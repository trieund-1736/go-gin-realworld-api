package services

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/repository"
	"go-gin-realworld-api/internal/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Login logs in a user and returns user with JWT token
func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	// Find user by email
	user, err := s.userRepo.FindUserByEmail(email)
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
