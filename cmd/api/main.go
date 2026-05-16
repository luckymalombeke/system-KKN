package main

import (
	"fmt"
	"kkn-system/config"
	"kkn-system/database"
	"kkn-system/handlers"
	"kkn-system/models/entity"
	"kkn-system/repositories"
	"kkn-system/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Connect Database
	database.ConnectDB(cfg)

	// 3. Auto Migration
	err := database.DB.AutoMigrate(&entity.Peserta{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	fmt.Println("Database migration completed")

	// 4. Setup Layers (Dependency Injection)
	pesertaRepo := repositories.NewPesertaRepository(database.DB)
	pesertaService := services.NewPesertaService(pesertaRepo)
	pesertaHandler := handlers.NewPesertaHandler(pesertaService)

	// 5. Setup Router
	r := gin.Default()

	// Health Check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// API Routes
	api := r.Group("/api/v1")
	{
		peserta := api.Group("/peserta")
		{
			peserta.POST("/register", pesertaHandler.Register)
			peserta.GET("/", pesertaHandler.GetAll)
			peserta.GET("/:id", pesertaHandler.GetByID)
		}
	}

	// 6. Start Server
	fmt.Printf("Server starting on port %s...\n", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
