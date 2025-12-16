package main

import (
	"fmt"
	"log"

	"go-gin-realworld-api/internal/bootstrap"
	"go-gin-realworld-api/internal/config"
	"go-gin-realworld-api/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize database
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize app container (contains all repositories, services and handlers)
	appContainer := bootstrap.NewAppContainer()

	// Create Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, appContainer)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
