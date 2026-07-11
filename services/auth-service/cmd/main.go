package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ikhdamw/e-wallet/auth-service/internal/handler"
	"github.com/ikhdamw/e-wallet/auth-service/internal/repository"
	"github.com/ikhdamw/e-wallet/auth-service/internal/service"
	"github.com/ikhdamw/e-wallet/auth-service/pkg/config"
	"github.com/ikhdamw/e-wallet/auth-service/pkg/database"
	"github.com/ikhdamw/e-wallet/auth-service/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient := database.NewRedisConnection(cfg)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, redisClient, cfg)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "auth-service",
		})
	})

	// API routes
	api := router.Group("/api/auth")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/refresh", authHandler.RefreshToken)
		api.GET("/me", middleware.AuthMiddleware(cfg), authHandler.GetMe)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("🚀 Auth Service running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
