package services

import (
	"crypto/sha256"
	"fmt"
	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/repository"
	"time"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(username, email, password string) (*models.User, error) {
	// Hash password
	hashedPassword := hashPassword(password)

	user := &models.User{
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create user in database
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// hashPassword hashes the password using SHA256
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}
