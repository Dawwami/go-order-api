package main

import (
	"github.com/Dawwami/go-order-api/internal/database"
	"github.com/Dawwami/go-order-api/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Ensure DB cleanup on exit
	defer database.CloseDB()

	// Initialize database singleton
	db := database.GetDB()

	// TODO: AutoMigrate models (Phase 2)
	db.AutoMigrate(&model.User{}, &model.Product{}, &model.Order{})

	// Setup Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	r.Run(":8080")
}
