package main

import (
	"context"
	"log"
	"sultra-otomotif-api/internal/config"
	"sultra-otomotif-api/internal/handler"
	"sultra-otomotif-api/internal/middleware"
	"sultra-otomotif-api/internal/repository"
	"sultra-otomotif-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	cfg := config.LoadConfig()

	// 2. Connect to the Database
	db, err := pgxpool.New(context.Background(), cfg.DBSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()
	log.Println("Database connection successful")
	// 3. Inisialisasi Semua Layer (Dependency Injection)

	// --- REPOSITORIES ---
	// Inisialisasi semua repository yang berinteraksi langsung dengan database.
	userRepository := repository.NewUserRepository(db)
	vehicleRepository := repository.NewVehicleRepository(db)
	imageRepository := repository.NewImageRepository(db)
	bookingRepository := repository.NewBookingRepository(db)
	reviewRepository := repository.NewReviewRepository(db)

	// --- SERVICES ---
	userService := service.NewUserService(userRepository)
	vehicleService := service.NewVehicleService(vehicleRepository, imageRepository)
	bookingService := service.NewBookingService(bookingRepository, vehicleRepository)
	reviewService := service.NewReviewService(reviewRepository, bookingRepository)

	// --- HANDLERS ---
	userHandler := handler.NewUserHandler(userService, cfg.JWTSecretKey)
	vehicleHandler := handler.NewVehicleHandler(vehicleService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	reviewHandler := handler.NewReviewHandler(reviewService)

	// 4. Setup Router Gin
	router := gin.Default()
	apiV1 := router.Group("/api/v1")

	setupAuthRoutes(apiV1, userHandler)
	setupVehicleRoutes(apiV1, vehicleHandler, cfg.JWTSecretKey)
	setupBookingRoutes(apiV1, bookingHandler, cfg.JWTSecretKey)
	setupPaymentRoutes(apiV1, bookingHandler)
	setupReviewRoutes(apiV1, reviewHandler, cfg.JWTSecretKey)

	// 6. Menjalankan Server
	log.Printf("Server starting on port %s", cfg.AppPort)
	err = router.Run(":" + cfg.AppPort)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupAuthRoutes mendaftarkan semua rute yang berhubungan dengan autentikasi.
func setupAuthRoutes(group *gin.RouterGroup, handler *handler.UserHandler) {
	authRoutes := group.Group("/auth")
	{
		authRoutes.POST("/register", handler.Register)
		authRoutes.POST("/login", handler.Login)
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
