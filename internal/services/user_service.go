package services

import (
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

// GetUserByID gets user by ID
func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	return s.userRepo.FindUserByID(id)
}
