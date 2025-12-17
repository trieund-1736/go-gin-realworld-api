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
	userRepo := repository.NewUserRepository()
	profileRepo := repository.NewProfileRepository()
	followRepo := repository.NewFollowRepository()
	articleRepo := repository.NewArticleRepository()
	commentRepo := repository.NewCommentRepository()
	favoriteRepo := repository.NewFavoriteRepository()
	tagRepo := repository.NewTagRepository()

	// Initialize services
	authService := services.NewAuthService(config.DB, userRepo)
	userService := services.NewUserService(config.DB, userRepo, profileRepo, followRepo)
	profileService := services.NewProfileService(config.DB, userRepo, profileRepo, followRepo)
	articleService := services.NewArticleService(config.DB, articleRepo)
	commentService := services.NewCommentService(config.DB, commentRepo, articleRepo)
	favoriteService := services.NewFavoriteService(config.DB, favoriteRepo, articleRepo)
	tagService := services.NewTagService(config.DB, tagRepo)

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
