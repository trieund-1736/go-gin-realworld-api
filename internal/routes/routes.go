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
			profiles.GET("/:username", appContainer.ProfileHandler.GetProfile)                                             // Get profile (optional auth)
			profiles.POST("/:username/follow", middleware.JWTAuthMiddleware(), appContainer.ProfileHandler.FollowUser)     // Follow user
			profiles.DELETE("/:username/follow", middleware.JWTAuthMiddleware(), appContainer.ProfileHandler.UnfollowUser) // Unfollow user
		} // Article routes (TODO: implement)
		articles := api.Group("/articles")
		{
			articles.GET("", listArticles)           // List articles
			articles.GET("/feed", feedArticles)      // Get feed (auth required)
			articles.GET("/:slug", getArticle)       // Get article by slug
			articles.POST("", createArticle)         // Create article (auth required)
			articles.PUT("/:slug", updateArticle)    // Update article (auth required)
			articles.DELETE("/:slug", deleteArticle) // Delete article (auth required)

			// Comments (TODO: implement)
			articles.POST("/:slug/comments", addComment)          // Add comment
			articles.GET("/:slug/comments", getComments)          // Get comments
			articles.DELETE("/:slug/comments/:id", deleteComment) // Delete comment

			// Favorites (TODO: implement)
			articles.POST("/:slug/favorite", favoriteArticle)     // Favorite
			articles.DELETE("/:slug/favorite", unfavoriteArticle) // Unfavorite
		}

		// Tags (TODO: implement)
		api.GET("/tags", getTags)
	}
}

// Handler stubs (TODO: implement)
func listArticles(c *gin.Context)      { c.JSON(501, gin.H{"error": "not implemented"}) }
func feedArticles(c *gin.Context)      { c.JSON(501, gin.H{"error": "not implemented"}) }
func getArticle(c *gin.Context)        { c.JSON(501, gin.H{"error": "not implemented"}) }
func createArticle(c *gin.Context)     { c.JSON(501, gin.H{"error": "not implemented"}) }
func updateArticle(c *gin.Context)     { c.JSON(501, gin.H{"error": "not implemented"}) }
func deleteArticle(c *gin.Context)     { c.JSON(501, gin.H{"error": "not implemented"}) }
func addComment(c *gin.Context)        { c.JSON(501, gin.H{"error": "not implemented"}) }
func getComments(c *gin.Context)       { c.JSON(501, gin.H{"error": "not implemented"}) }
func deleteComment(c *gin.Context)     { c.JSON(501, gin.H{"error": "not implemented"}) }
func favoriteArticle(c *gin.Context)   { c.JSON(501, gin.H{"error": "not implemented"}) }
func unfavoriteArticle(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) }
func getTags(c *gin.Context)           { c.JSON(501, gin.H{"error": "not implemented"}) }
