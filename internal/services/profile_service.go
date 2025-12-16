package services

import (
	"go-gin-realworld-api/internal/dtos"
	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/repository"
)

type ProfileService struct {
	userRepo    *repository.UserRepository
	profileRepo *repository.ProfileRepository
	followRepo  *repository.FollowRepository
}

func NewProfileService(userRepo *repository.UserRepository, profileRepo *repository.ProfileRepository, followRepo *repository.FollowRepository) *ProfileService {
	return &ProfileService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		followRepo:  followRepo,
	}
}

// GetProfileByUsername gets profile by username with following status
func (s *ProfileService) GetProfileByUsername(username string, currentUserID int64) (*dtos.ProfileResponse, error) {
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
func (s *ProfileService) FollowUser(followerID int64, followeeUsername string) (*dtos.ProfileResponse, error) {
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
func (s *ProfileService) UnfollowUser(followerID int64, followeeUsername string) (*dtos.ProfileResponse, error) {
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
