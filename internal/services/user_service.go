package services

import (
	"database/sql"
	"go-gin-realworld-api/internal/dtos"
	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/repository"
)

type UserService struct {
	userRepo    *repository.UserRepository
	profileRepo *repository.ProfileRepository
	followRepo  *repository.FollowRepository
}

func NewUserService(userRepo *repository.UserRepository, profileRepo *repository.ProfileRepository, followRepo *repository.FollowRepository) *UserService {
	return &UserService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		followRepo:  followRepo,
	}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(username, email, password string) (*models.User, error) {
	// Hash password
	hashedPassword := hashPassword(password)

	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	// Create user in database
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	// Create associated profile
	profile := &models.Profile{
		UserID: user.ID,
	}
	if err := s.profileRepo.CreateProfile(profile); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID gets user by ID
func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	return s.userRepo.FindUserByID(id)
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(userID int64, req *dtos.UpdateUserRequest) (*models.User, error) {
	// Get user
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Update user fields if provided
	if req.User.Email != "" {
		user.Email = req.User.Email
	}
	if req.User.Username != "" {
		user.Username = req.User.Username
	}
	if req.User.Password != "" {
		user.Password = hashPassword(req.User.Password)
	}

	// Update user
	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	// Get or create profile
	profile, err := s.profileRepo.FindProfileByUserID(userID)
	if err != nil {
		// Create new profile if not exists
		profile = &models.Profile{
			UserID: userID,
		}
	}

	// Update profile fields if provided
	if req.User.Image != "" {
		profile.Image = sql.NullString{String: req.User.Image, Valid: true}
	}
	if req.User.Bio != "" {
		profile.Bio = sql.NullString{String: req.User.Bio, Valid: true}
	}

	// Update profile
	if err := s.profileRepo.UpdateProfile(profile); err != nil {
		return nil, err
	}

	return user, nil
}
