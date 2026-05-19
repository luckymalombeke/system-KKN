package main

import (
	"fmt"
	"kkn-system/config"
	"kkn-system/database"
	"kkn-system/handlers"
	"kkn-system/middleware"
	"kkn-system/models/entity"
	"kkn-system/repositories"
	"kkn-system/services"
	"kkn-system/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	if err := utils.InitJWTSecret(cfg.JWTSecret); err != nil {
		log.Fatal("JWT config error:", err)
	}

	// 2. Connect Database
	database.ConnectDB(cfg)

	// 3. Auto Migration
	err := database.DB.AutoMigrate(&entity.Peserta{}, &entity.Pembayaran{}, &entity.Notifikasi{}, &entity.Lokasi{}, &entity.Admin{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	fmt.Println("Database migration completed")

	// 4. Setup Layers (Dependency Injection)
	pesertaRepo := repositories.NewPesertaRepository(database.DB)
	pembayaranRepo := repositories.NewPembayaranRepository(database.DB)
	notifikasiRepo := repositories.NewNotifikasiRepository(database.DB)
	lokasiRepo := repositories.NewLokasiRepository(database.DB)
	adminRepo := repositories.NewAdminRepository(database.DB)

	paymentService := services.NewPaymentService(cfg)

	pesertaService := services.NewPesertaService(pesertaRepo)
	pembayaranService := services.NewPembayaranService(pembayaranRepo, pesertaRepo, paymentService, cfg.PaymentExpiryHours)
	notifikasiService := services.NewNotifikasiService(notifikasiRepo)
	lokasiService := services.NewLokasiService(lokasiRepo)
	emailService := services.NewEmailService(cfg.SMTP)
	authService := services.NewAuthService(pesertaRepo, adminRepo, emailService, cfg.OTPExpiryMinutes)

	pesertaHandler := handlers.NewPesertaHandler(pesertaService)
	pembayaranHandler := handlers.NewPembayaranHandler(pembayaranService)
	notifikasiHandler := handlers.NewNotifikasiHandler(notifikasiService)
	lokasiHandler := handlers.NewLokasiHandler(lokasiService)
	authHandler := handlers.NewAuthHandler(authService)

	// 5. Setup Router
	r := gin.Default()

	// 5.5 Setup CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health Check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// API Routes
	api := r.Group("/api/v1")
	{
		// 1. PUBLIC ROUTES
		auth := api.Group("/auth")
		{
			auth.POST("/request-otp", authHandler.RequestOTP)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
			auth.POST("/admin-login", authHandler.AdminLogin) // Admin Login URL
		}

		pesertaPublic := api.Group("/peserta")
		{
			pesertaPublic.POST("/register", pesertaHandler.Register) // Mendaftar tidak perlu login
		}

		// Webhook Pembayaran Midtrans (Harus public)
		api.POST("/pembayaran/notification", pembayaranHandler.HandleNotification)

		// 2. ADMIN PROTECTED ROUTES
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin", "superadmin"))
		{
			// Kelola Peserta
			admin.GET("/peserta", pesertaHandler.GetAll)
			admin.GET("/peserta/:id", pesertaHandler.GetByID)
			admin.PATCH("/peserta/:id/status", pesertaHandler.UpdateStatus)
			admin.PATCH("/peserta/:id/lokasi", pesertaHandler.AssignLocation)
			
			// Kelola Lokasi
			admin.POST("/lokasi", lokasiHandler.Create)
			admin.GET("/lokasi", lokasiHandler.GetAll) // Bisa juga public jika mhs butuh lihat, tapi untuk saat ini admin

			// Kelola Pembayaran
			admin.GET("/pembayaran/:id", pembayaranHandler.GetByID) // Admin bisa lihat spesifik
		}

		// 3. STUDENT PROTECTED ROUTES (user_id dari JWT, tanpa ID di URL)
		mhs := api.Group("/mahasiswa")
		mhs.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("mahasiswa"))
		{
			mhs.GET("/profile", pesertaHandler.GetMyProfile)

			mhs.POST("/pembayaran/invoice", pembayaranHandler.CreateMyInvoice)
			mhs.GET("/pembayaran", pembayaranHandler.GetMyPembayaran)

			mhs.GET("/notifikasi", notifikasiHandler.GetMyNotifikasi)
			mhs.PATCH("/notifikasi/:id/read", notifikasiHandler.MarkMyNotifikasiAsRead)
		}
	}

	// 6. Start Server
	fmt.Printf("Server starting on port %s...\n", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
