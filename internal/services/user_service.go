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

// GetProfileByUsername gets profile by username with following status
func (s *UserService) GetProfileByUsername(username string, currentUserID int64) (*dtos.ProfileResponse, error) {
	// Find user by username with profile preloaded
	user, err := s.userRepo.FindUserByUsername(username, true)
	if err != nil {
		return nil, err
	}

	// Use preloaded profile or create default empty profile
	profile := user.Profile
	if profile == nil {
		profile = &models.Profile{
			UserID: user.ID,
		}
	}

	// Check if current user follows this user
	following := false
	if currentUserID > 0 {
		following, err = s.followRepo.IsFollowing(currentUserID, user.ID)
		if err != nil {
			// If error occurs, set following to false
			following = false
		}
	}

	resp := &dtos.ProfileResponse{
		Profile: dtos.ProfileUserResponse{
			Username:  user.Username,
			Bio:       profile.Bio.String,
			Image:     profile.Image.String,
			Following: following,
		},
	}

	return resp, nil
}

// FollowUser creates a follow relationship
func (s *UserService) FollowUser(followerID int64, followeeUsername string) (*dtos.ProfileResponse, error) {
	// Find followee by username
	followee, err := s.userRepo.FindUserByUsername(followeeUsername)
	if err != nil {
		return nil, err
	}

	// Check if already following
	isFollowing, err := s.followRepo.IsFollowing(followerID, followee.ID)
	if err != nil {
		return nil, err
	}

	if !isFollowing {
		follow := &models.Follow{
			FollowerID: followerID,
			FolloweeID: followee.ID,
		}
		if err := s.followRepo.CreateFollow(follow); err != nil {
			return nil, err
		}
	}

	// Return updated profile
	return s.GetProfileByUsername(followeeUsername, followerID)
}

// UnfollowUser deletes a follow relationship
func (s *UserService) UnfollowUser(followerID int64, followeeUsername string) (*dtos.ProfileResponse, error) {
	// Find followee by username
	followee, err := s.userRepo.FindUserByUsername(followeeUsername)
	if err != nil {
		return nil, err
	}

	// Delete follow relationship
	if err := s.followRepo.DeleteFollow(followerID, followee.ID); err != nil {
		return nil, err
	}

	// Return updated profile
	return s.GetProfileByUsername(followeeUsername, followerID)
}
