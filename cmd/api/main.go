package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sultra-otomotif-api/internal/config"
	"sultra-otomotif-api/internal/handler"
	"sultra-otomotif-api/internal/middleware"
	"sultra-otomotif-api/internal/repository"
	"sultra-otomotif-api/internal/service"
	"sultra-otomotif-api/internal/websocket"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Memuat Konfigurasi Aplikasi
	cfg := config.LoadConfig()
	if cfg.JWTSecretKey == "" || cfg.DBSource == "" || cfg.CloudinaryURL == "" || cfg.FrontendURL == "" {
		log.Fatal("FATAL: Required environment variables are not set.")
	}

	// 2. Menghubungkan ke Database
	db, err := pgxpool.New(context.Background(), cfg.DBSource)
	if err != nil {
		log.Fatalf("FATAL: Unable to connect to database: %v\n", err)
	}
	defer db.Close()
	log.Println("Database connection successful")

	// 3. Inisialisasi Semua Layer (Dependency Injection)
	userRepository := repository.NewUserRepository(db)
	vehicleRepository := repository.NewVehicleRepository(db)
	imageRepository := repository.NewImageRepository(db)
	bookingRepository := repository.NewBookingRepository(db)
	reviewRepository := repository.NewReviewRepository(db)
	salesRepository := repository.NewSalesRepository(db)
	chatRepository := repository.NewChatRepository(db)

	userService := service.NewUserService(userRepository)
	vehicleService := service.NewVehicleService(vehicleRepository, imageRepository, userRepository)
	bookingService := service.NewBookingService(bookingRepository, vehicleRepository)
	reviewService := service.NewReviewService(reviewRepository, bookingRepository)
	adminService := service.NewAdminService(userRepository, vehicleRepository)
	salesService := service.NewSalesService(salesRepository, vehicleRepository)
	chatService := service.NewChatService(chatRepository, vehicleRepository)

	userHandler := handler.NewUserHandler(userService, cfg.JWTSecretKey)
	vehicleHandler := handler.NewVehicleHandler(vehicleService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	reviewHandler := handler.NewReviewHandler(reviewService)
	adminHandler := handler.NewAdminHandler(adminService)
	salesHandler := handler.NewSalesHandler(salesService)
	chatHandler := handler.NewChatHandler(chatService)

	hub := websocket.NewHub(chatService)
	go hub.Run()

	// 4. Setup Router Gin
	// Set GIN_MODE dari environment variable, default ke "debug"
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" || ginMode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Konfigurasi CORS yang fleksibel
	corsConfig := cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL, "https://penjualan-dan-penyewaan-kendaraan-f.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	apiV1 := router.Group("/api/v1")

	// Health Check Endpoint
	apiV1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// 5. Mendaftarkan Semua Rute API
	setupAuthRoutes(apiV1, userHandler, cfg.JWTSecretKey)
	setupVehicleRoutes(apiV1, vehicleHandler, cfg.JWTSecretKey)
	setupBookingRoutes(apiV1, bookingHandler, cfg.JWTSecretKey)
	setupPaymentRoutes(apiV1, bookingHandler)
	setupReviewRoutes(apiV1, reviewHandler, cfg.JWTSecretKey)
	setupAdminRoutes(apiV1, adminHandler, cfg.JWTSecretKey)
	setupSalesRoutes(apiV1, salesHandler, cfg.JWTSecretKey)
	setupChatRoutes(apiV1, chatHandler, cfg.JWTSecretKey)

	// Daftarkan Rute WebSocket
	apiV1.GET("/ws", middleware.AuthMiddleware(cfg.JWTSecretKey), func(c *gin.Context) {
		handler.ServeWs(hub, c)
	})

	// 6. Menjalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.AppPort // Default ke 8080 untuk lokal
	}
	log.Printf("Server starting on port %s in %s mode", port, gin.Mode())
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("FATAL: Failed to start server: %v", err)
	}
}

// setupAuthRoutes mendaftarkan semua rute yang berhubungan dengan autentikasi.
func setupAuthRoutes(group *gin.RouterGroup, handler *handler.UserHandler, jwtSecret string) {
	authRoutes := group.Group("/auth")
	{
		authRoutes.POST("/register", handler.Register)
		authRoutes.POST("/login", handler.Login)
		authRoutes.GET("/me", middleware.AuthMiddleware(jwtSecret), handler.GetMe)
	}
}

// setupVehicleRoutes mendaftarkan semua rute yang berhubungan dengan kendaraan.
func setupVehicleRoutes(group *gin.RouterGroup, handler *handler.VehicleHandler, jwtSecret string) {
	vehicleRoutes := group.Group("/vehicles")
	{
		// Rute Publik
		vehicleRoutes.GET("/", handler.GetAllVehicles)
		vehicleRoutes.GET("/:id", handler.GetVehicleByID)

		// Rute yang dilindungi (hanya untuk Vendor)
		protectedVendorRoutes := vehicleRoutes.Use(middleware.AuthMiddleware(jwtSecret), middleware.RoleMiddleware("vendor"))
		{
			protectedVendorRoutes.POST("/", handler.CreateVehicle)
			protectedVendorRoutes.PUT("/:id", handler.UpdateVehicle)
			protectedVendorRoutes.DELETE("/:id", handler.DeleteVehicle)
			protectedVendorRoutes.POST("/:id/images", handler.UploadVehicleImage)
			protectedVendorRoutes.GET("/my-listings", handler.GetMyListings)
		}
	}
}

// setupBookingRoutes mendaftarkan semua rute yang berhubungan dengan booking.
func setupBookingRoutes(group *gin.RouterGroup, handler *handler.BookingHandler, jwtSecret string) {
	bookingRoutes := group.Group("/bookings")
	bookingRoutes.Use(middleware.AuthMiddleware(jwtSecret)) // Semua rute booking butuh login
	{
		// Rute khusus Customer
		bookingRoutes.POST("/", middleware.RoleMiddleware("customer"), handler.CreateBooking)
		bookingRoutes.GET("/my-bookings", middleware.RoleMiddleware("customer"), handler.GetMyBookings)

		// Rute khusus Vendor
		bookingRoutes.GET("/vendor", middleware.RoleMiddleware("vendor"), handler.GetVendorBookings)
		bookingRoutes.PATCH("/:id/status", middleware.RoleMiddleware("vendor"), handler.UpdateBookingStatus)

		// Rute yang bisa diakses oleh Customer atau Vendor yang bersangkutan
		bookingRoutes.GET("/:id", handler.GetBookingByID)
	}
}

// setupPaymentRoutes mendaftarkan rute untuk callback dari payment gateway.
func setupPaymentRoutes(group *gin.RouterGroup, handler *handler.BookingHandler) {
	paymentRoutes := group.Group("/payments")
	{
		// Rute ini publik karena akan diakses oleh server payment gateway
		paymentRoutes.POST("/callback", handler.PaymentCallback)
	}
}

// setupReviewRoutes mendaftarkan semua rute yang berhubungan dengan review.
func setupReviewRoutes(group *gin.RouterGroup, handler *handler.ReviewHandler, jwtSecret string) {
	// Endpoint publik untuk melihat review sebuah kendaraan
	group.GET("/vehicles/:id/reviews", handler.GetVehicleReviews)

	// Endpoint dilindungi untuk membuat review
	reviewCreationRoutes := group.Group("/bookings/:booking_id/reviews")
	reviewCreationRoutes.Use(middleware.AuthMiddleware(jwtSecret), middleware.RoleMiddleware("customer"))
	{
		reviewCreationRoutes.POST("/", handler.CreateReview)
	}
}

func setupAdminRoutes(group *gin.RouterGroup, handler *handler.AdminHandler, jwtSecret string) {
	adminRoutes := group.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(jwtSecret), middleware.RoleMiddleware("admin"))
	{
		// Rute Manajemen Vendor
		adminRoutes.GET("/vendors", handler.GetVendors)
		adminRoutes.PATCH("/vendors/:id/verify", handler.VerifyVendor)

		// Rute Manajemen User
		adminRoutes.GET("/users", handler.GetAllUsers)
		adminRoutes.DELETE("/users/:id", handler.DeleteUser)

		// Rute Manajemen Listing
		adminRoutes.GET("/vehicles", handler.GetAllVehicles)
		adminRoutes.DELETE("/vehicles/:id", handler.DeleteVehicle)
	}
}

func setupSalesRoutes(group *gin.RouterGroup, handler *handler.SalesHandler, jwtSecret string) {
	salesRoutes := group.Group("/sales")
	salesRoutes.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// Rute untuk melihat riwayat pembelian (customer) & penjualan (vendor)
		salesRoutes.GET("/purchases", middleware.RoleMiddleware("customer"), handler.GetMyPurchases)
		salesRoutes.GET("/sales", middleware.RoleMiddleware("vendor"), handler.GetMySales)
	}

	// Rute untuk customer memulai pembelian sebuah mobil
	purchaseRoutes := group.Group("/vehicles/:id/purchase")
	purchaseRoutes.Use(middleware.AuthMiddleware(jwtSecret), middleware.RoleMiddleware("customer"))
	{
		purchaseRoutes.POST("/", handler.InitiatePurchase)
	}

	// Rute callback pembayaran (mirip booking)
	group.POST("/sales/callback", handler.PaymentCallback)
}

func setupChatRoutes(group *gin.RouterGroup, handler *handler.ChatHandler, jwtSecret string) {
	chatRoutes := group.Group("/conversations")
	chatRoutes.Use(middleware.AuthMiddleware(jwtSecret))
	{
		chatRoutes.GET("/", handler.ListConversations)
		chatRoutes.GET("/:id/messages", handler.GetMessages)
	}

	// Rute untuk memulai percakapan
	startChatRoutes := group.Group("/vehicles/:id/conversations")
	startChatRoutes.Use(middleware.AuthMiddleware(jwtSecret), middleware.RoleMiddleware("customer"))
	{
		startChatRoutes.POST("/", handler.StartConversation)
	}
}
