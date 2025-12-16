package routes

import (
	"github.com/gin-gonic/gin"
	"go-gin-realworld-api/internal/handlers"
)

func SetupRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// API v1 routes
	api := router.Group("/api")
	{
		// User routes (TODO: implement)
		users := api.Group("/users")
		{
			users.POST("", registerUser)    // Register
			users.POST("/login", loginUser) // Login
		}

		// Get current user (TODO: implement with auth middleware)
		api.GET("/user", getCurrentUser)

		// Article routes (TODO: implement)
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

func registerUser(c *gin.Context)      { c.JSON(501, gin.H{"error": "not implemented"}) }
func loginUser(c *gin.Context)         { c.JSON(501, gin.H{"error": "not implemented"}) }
func getCurrentUser(c *gin.Context)    { c.JSON(501, gin.H{"error": "not implemented"}) }
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
