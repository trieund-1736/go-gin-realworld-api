package services

import (
	"context"
	"go-gin-realworld-api/internal/dtos"
	"go-gin-realworld-api/internal/models"
	"go-gin-realworld-api/internal/repository"

	"gorm.io/gorm"
)

type ProfileService struct {
	db          *gorm.DB
	userRepo    *repository.UserRepository
	profileRepo *repository.ProfileRepository
	followRepo  *repository.FollowRepository
}

func NewProfileService(db *gorm.DB, userRepo *repository.UserRepository, profileRepo *repository.ProfileRepository, followRepo *repository.FollowRepository) *ProfileService {
	return &ProfileService{
		db:          db,
		userRepo:    userRepo,
		profileRepo: profileRepo,
		followRepo:  followRepo,
	}
}

// GetProfileByUsername gets profile by username with following status
func (s *ProfileService) GetProfileByUsername(ctx context.Context, username string, currentUserID int64) (*dtos.ProfileResponse, error) {
	db := s.db.WithContext(ctx)
	// Find user by username with profile preloaded
	user, err := s.userRepo.FindUserByUsername(db, username, true)
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
		following, err = s.followRepo.IsFollowing(db, currentUserID, user.ID)
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
func (s *ProfileService) FollowUser(ctx context.Context, followerID int64, followeeUsername string) (*dtos.ProfileResponse, error) {
	db := s.db.WithContext(ctx)
	if err := db.Transaction(func(tx *gorm.DB) error {
		followee, err := s.userRepo.FindUserByUsername(tx, followeeUsername)
		if err != nil {
			return err
		}

		// Check if already following
		isFollowing, err := s.followRepo.IsFollowing(tx, followerID, followee.ID)
		if err != nil {
			return err
		}

		if !isFollowing {
			follow := &models.Follow{
				FollowerID: followerID,
				FolloweeID: followee.ID,
			}
			if err := s.followRepo.CreateFollow(tx, follow); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// Return updated profile
	return s.GetProfileByUsername(ctx, followeeUsername, followerID)
}

// UnfollowUser deletes a follow relationship
func (s *ProfileService) UnfollowUser(ctx context.Context, followerID int64, followeeUsername string) (*dtos.ProfileResponse, error) {
	db := s.db.WithContext(ctx)
	if err := db.Transaction(func(tx *gorm.DB) error {
		followee, err := s.userRepo.FindUserByUsername(tx, followeeUsername)
		if err != nil {
			return err
		}

		// Delete follow relationship
		if err := s.followRepo.DeleteFollow(tx, followerID, followee.ID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// Return updated profile
	return s.GetProfileByUsername(ctx, followeeUsername, followerID)
}
