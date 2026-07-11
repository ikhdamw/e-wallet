package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ikhdamw/e-wallet/wallet-service/internal/handler"
	"github.com/ikhdamw/e-wallet/wallet-service/internal/repository"
	"github.com/ikhdamw/e-wallet/wallet-service/internal/service"
	"github.com/ikhdamw/e-wallet/wallet-service/pkg/config"
	"github.com/ikhdamw/e-wallet/wallet-service/pkg/database"
	"github.com/ikhdamw/e-wallet/wallet-service/pkg/middleware"

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
	walletRepo := repository.NewWalletRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize services
	walletService := service.NewWalletService(walletRepo, transactionRepo, redisClient, cfg)

	// Initialize handlers
	walletHandler := handler.NewWalletHandler(walletService)

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "wallet-service",
		})
	})

	// API routes
	api := router.Group("/api/wallet")
	{
		api.GET("/balance", middleware.AuthMiddleware(cfg), walletHandler.GetBalance)
		api.POST("/topup", middleware.AuthMiddleware(cfg), walletHandler.TopUp)
		api.GET("/history", middleware.AuthMiddleware(cfg), walletHandler.GetHistory)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	fmt.Printf("🚀 Wallet Service running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
