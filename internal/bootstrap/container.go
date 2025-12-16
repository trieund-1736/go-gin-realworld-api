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
	userService := services.NewUserService(userRepo, profileRepo, followRepo)
	authService := services.NewAuthService(userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	profileHandler := handlers.NewProfileHandler(userService)

	return &AppContainer{
		UserHandler:    userHandler,
		AuthHandler:    authHandler,
		ProfileHandler: profileHandler,
	}
}
