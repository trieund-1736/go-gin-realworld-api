package routes

import (
	"go-gin-realworld-api/internal/bootstrap"
	"go-gin-realworld-api/internal/handlers"
	"go-gin-realworld-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, appContainer *bootstrap.AppContainer) {
	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// API v1 routes
	api := router.Group("/api")
	{
		// User routes
		users := api.Group("/users")
		{
			users.POST("", appContainer.UserHandler.RegisterUser) // Register
			users.POST("/login", appContainer.AuthHandler.Login)  // Login
		}

		// Current user routes (requires auth middleware)
		user := api.Group("/user")
		user.Use(middleware.JWTAuthMiddleware())
		{
			user.GET("", appContainer.UserHandler.GetCurrentUser) // Get current user
			user.PUT("", appContainer.UserHandler.UpdateUser)     // Update current user
		}
		// Profile routes
		profiles := api.Group("/profiles")
		{
			profiles.GET("/:username", middleware.JWTOptionalAuthMiddleware(), appContainer.ProfileHandler.GetProfile)     // Get profile (optional auth)
			profiles.POST("/:username/follow", middleware.JWTAuthMiddleware(), appContainer.ProfileHandler.FollowUser)     // Follow user (required auth)
			profiles.DELETE("/:username/follow", middleware.JWTAuthMiddleware(), appContainer.ProfileHandler.UnfollowUser) // Unfollow user (required auth)
		} // Article routes
		articles := api.Group("/articles")
		{
			articles.GET("", middleware.JWTOptionalAuthMiddleware(), appContainer.ArticleHandler.ListArticles)     // List articles
			articles.GET("/feed", middleware.JWTAuthMiddleware(), appContainer.ArticleHandler.FeedArticles)        // Get feed (auth required)
			articles.GET("/:slug", middleware.JWTOptionalAuthMiddleware(), appContainer.ArticleHandler.GetArticle) // Get article by slug
			articles.POST("", middleware.JWTAuthMiddleware(), appContainer.ArticleHandler.CreateArticle)           // Create article (auth required)
			articles.PUT("/:slug", middleware.JWTAuthMiddleware(), appContainer.ArticleHandler.UpdateArticle)      // Update article (auth required)
			articles.DELETE("/:slug", middleware.JWTAuthMiddleware(), appContainer.ArticleHandler.DeleteArticle)   // Delete article (auth required)

			// Comments
			articles.POST("/:slug/comments", middleware.JWTAuthMiddleware(), appContainer.CommentHandler.CreateComment)       // Add comment (auth required)
			articles.GET("/:slug/comments", middleware.JWTOptionalAuthMiddleware(), appContainer.CommentHandler.GetComments)  // Get comments (optional auth)
			articles.DELETE("/:slug/comments/:id", middleware.JWTAuthMiddleware(), appContainer.CommentHandler.DeleteComment) // Delete comment (auth required)

			// Favorites
			articles.POST("/:slug/favorite", middleware.JWTAuthMiddleware(), appContainer.FavoriteHandler.FavoriteArticle)     // Favorite (auth required)
			articles.DELETE("/:slug/favorite", middleware.JWTAuthMiddleware(), appContainer.FavoriteHandler.UnfavoriteArticle) // Unfavorite (auth required)
		}

		// Tags (TODO: implement)
		api.GET("/tags", getTags)
	}
}

// Handler stubs (TODO: implement)
func getTags(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
