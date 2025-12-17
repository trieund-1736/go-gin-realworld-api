package bootstrap

import (
	"go-gin-realworld-api/internal/config"
	"go-gin-realworld-api/internal/handlers"
	"go-gin-realworld-api/internal/repository"
	"go-gin-realworld-api/internal/services"
)

type AppContainer struct {
	// Handlers
	UserHandler     *handlers.UserHandler
	AuthHandler     *handlers.AuthHandler
	ProfileHandler  *handlers.ProfileHandler
	ArticleHandler  *handlers.ArticleHandler
	CommentHandler  *handlers.CommentHandler
	FavoriteHandler *handlers.FavoriteHandler
	TagHandler      *handlers.TagHandler
}

func NewAppContainer() *AppContainer {
	// Initialize repositories
	userRepo := repository.NewUserRepository(config.DB)
	profileRepo := repository.NewProfileRepository(config.DB)
	followRepo := repository.NewFollowRepository(config.DB)
	articleRepo := repository.NewArticleRepository(config.DB)
	commentRepo := repository.NewCommentRepository(config.DB)
	favoriteRepo := repository.NewFavoriteRepository(config.DB)
	tagRepo := repository.NewTagRepository(config.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo, profileRepo, followRepo)
	profileService := services.NewProfileService(userRepo, profileRepo, followRepo)
	articleService := services.NewArticleService(articleRepo)
	commentService := services.NewCommentService(commentRepo, articleRepo)
	favoriteService := services.NewFavoriteService(favoriteRepo, articleRepo)
	tagService := services.NewTagService(tagRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	profileHandler := handlers.NewProfileHandler(profileService)
	articleHandler := handlers.NewArticleHandler(articleService)
	commentHandler := handlers.NewCommentHandler(commentService)
	favoriteHandler := handlers.NewFavoriteHandler(favoriteService)
	tagHandler := handlers.NewTagHandler(tagService)

	return &AppContainer{
		UserHandler:     userHandler,
		AuthHandler:     authHandler,
		ProfileHandler:  profileHandler,
		ArticleHandler:  articleHandler,
		CommentHandler:  commentHandler,
		FavoriteHandler: favoriteHandler,
		TagHandler:      tagHandler,
	}
}
