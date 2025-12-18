package services

import (
	"context"
	"database/sql"
	"go-gin-realworld-api/internal/dtos"
	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/repository"

	"gorm.io/gorm"
)

type UserService struct {
	db          *gorm.DB
	userRepo    repository.UserRepository
	profileRepo repository.ProfileRepository
	followRepo  repository.FollowRepository
}

func NewUserService(db *gorm.DB, userRepo repository.UserRepository, profileRepo repository.ProfileRepository, followRepo repository.FollowRepository) *UserService {
	return &UserService{
		db:          db,
		userRepo:    userRepo,
		profileRepo: profileRepo,
		followRepo:  followRepo,
	}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(c context.Context, username, email, password string) (*models.User, error) {

	var user *models.User
	err := s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// Hash password
		hashedPassword := hashPassword(password)

		user = &models.User{
			Username: username,
			Email:    email,
			Password: hashedPassword,
		}

		// Create user in database
		if err := s.userRepo.CreateUser(tx, user); err != nil {
			return err
		}

		// Create associated profile
		profile := &models.Profile{
			UserID: user.ID,
		}
		if err := s.profileRepo.CreateProfile(tx, profile); err != nil {
			return err
		}
		return nil
	})

	return user, err
}

// GetUserByID gets user by ID
func (s *UserService) GetUserByID(c context.Context, id int64) (*models.User, error) {
	return s.userRepo.FindUserByID(s.db.WithContext(c), id)
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(c context.Context, userID int64, req *dtos.UpdateUserRequest) (*models.User, error) {
	var user *models.User

	err := s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// Get user
		user, err := s.userRepo.FindUserByID(tx, userID)
		if err != nil {
			return err
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
		if err := s.userRepo.UpdateUser(tx, user); err != nil {
			return err
		}

		// Get or create profile
		profile, err := s.profileRepo.FindProfileByUserID(tx, userID)
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
		if err := s.profileRepo.UpdateProfile(tx, profile); err != nil {
			return err
		}

		return nil
	})

	return user, err
}
