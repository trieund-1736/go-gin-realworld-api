package bootstrap

import (
	"go-gin-realworld-api/internal/config"
	"go-gin-realworld-api/internal/handlers"
	"go-gin-realworld-api/internal/repository"
	"go-gin-realworld-api/internal/services"
)

type AppContainer struct {
	// Handlers
	UserHandler    *handlers.UserHandler
	AuthHandler    *handlers.AuthHandler
	ProfileHandler *handlers.ProfileHandler
}

func NewAppContainer() *AppContainer {
	// Initialize repositories
	userRepo := repository.NewUserRepository(config.DB)
	profileRepo := repository.NewProfileRepository(config.DB)
	followRepo := repository.NewFollowRepository(config.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo, profileRepo, followRepo)
	profileService := services.NewProfileService(userRepo, profileRepo, followRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	profileHandler := handlers.NewProfileHandler(profileService)

	return &AppContainer{
		UserHandler:    userHandler,
		AuthHandler:    authHandler,
		ProfileHandler: profileHandler,
	}
}
