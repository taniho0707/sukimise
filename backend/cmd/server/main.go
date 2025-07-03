package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sukimise/internal/config"
	"sukimise/internal/database"
	"sukimise/internal/handlers"
	"sukimise/internal/middleware"
	"sukimise/internal/repositories"
	"sukimise/internal/services"
	"sukimise/internal/types"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatal("Configuration validation failed:", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Configure database connection pool
	db.SetMaxOpenConns(cfg.Database.MaxConnections)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	storeRepo := repositories.NewStoreRepository(db)
	reviewRepo := repositories.NewReviewRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	storeService := services.NewStoreService(storeRepo)
	reviewService := services.NewReviewService(reviewRepo)

	// Initialize handlers
	handler := handlers.NewHandler(userService, storeService, reviewService)

	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	r := gin.New()

	// Add middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.SecurityHeaders())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, types.APIResponse{
			Success: true,
			Data: map[string]interface{}{
				"status":     "ok",
				"message":    "Sukimise API is running",
				"version":    "1.0.0",
				"timestamp":  time.Now().Unix(),
				"environment": cfg.Server.Environment,
			},
		})
	})

	// 静的ファイル配信
	r.GET("/uploads/:filename", handler.ServeUpload)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", handler.Login)
			auth.POST("/refresh", handler.RefreshToken)
		}

		stores := api.Group("/stores")
		{
			stores.GET("", handler.GetStores)
			stores.GET("/export/csv", handler.ExportStoresCSV)
			stores.GET("/categories", handler.GetCategories)
			stores.GET("/tags", handler.GetTags)
			stores.GET("/:id", handler.GetStore)
			stores.GET("/:id/reviews", handler.GetReviewsByStore)
		}

		protected := api.Group("")
		protected.Use(middleware.Auth())
		{
			protectedStores := protected.Group("/stores")
			{
				protectedStores.POST("", handler.CreateStore)
				protectedStores.PUT("/:id", handler.UpdateStore)
				protectedStores.DELETE("/:id", handler.DeleteStore)
			}

			reviews := protected.Group("/reviews")
			{
				reviews.POST("", handler.CreateReview)
				reviews.PUT("/:id", handler.UpdateReview)
				reviews.DELETE("/:id", handler.DeleteReview)
			}

			users := protected.Group("/users")
			{
				users.GET("/me", handler.GetCurrentUser)
				users.PUT("/me", handler.UpdateCurrentUser)
				users.GET("/me/reviews", handler.GetMyReviews)
			}

			upload := protected.Group("/upload")
			{
				upload.POST("/image", handler.UploadImage)
				upload.DELETE("/:filename", handler.DeleteUpload)
			}
		}
	}

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s (environment: %s)", cfg.GetServerAddress(), cfg.Server.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}