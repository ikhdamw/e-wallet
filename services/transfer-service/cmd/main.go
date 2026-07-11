package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ikhdamw/e-wallet/transfer-service/internal/handler"
	"github.com/ikhdamw/e-wallet/transfer-service/internal/repository"
	"github.com/ikhdamw/e-wallet/transfer-service/internal/service"
	"github.com/ikhdamw/e-wallet/transfer-service/pkg/config"
	"github.com/ikhdamw/e-wallet/transfer-service/pkg/database"
	"github.com/ikhdamw/e-wallet/transfer-service/pkg/middleware"

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

	// Initialize RabbitMQ
	rabbitMQ, err := database.NewRabbitMQConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Initialize repositories
	transferRepo := repository.NewTransferRepository(db)

	// Initialize services
	transferService := service.NewTransferService(transferRepo, rabbitMQ, cfg)

	// Initialize handlers
	transferHandler := handler.NewTransferHandler(transferService)

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "transfer-service",
		})
	})

	// API routes
	api := router.Group("/api/transfer")
	{
		api.POST("/internal", middleware.AuthMiddleware(cfg), transferHandler.InternalTransfer)
		api.POST("/external", middleware.AuthMiddleware(cfg), transferHandler.ExternalTransfer)
		api.GET("/status/:id", middleware.AuthMiddleware(cfg), transferHandler.GetStatus)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	fmt.Printf("🚀 Transfer Service running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
