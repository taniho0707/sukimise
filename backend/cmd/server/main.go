package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	// Validate user environment variables
	if err := validateUserEnvironmentVariables(); err != nil {
		log.Fatal("User environment validation failed:", err)
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
	viewerAuthRepo := repositories.NewViewerAuthRepository(db)
	categoryCustomizationRepo := repositories.NewCategoryCustomizationRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	storeService := services.NewStoreService(storeRepo)
	reviewService := services.NewReviewService(reviewRepo)
	viewerAuthService := services.NewViewerAuthService(viewerAuthRepo)
	categoryCustomizationService := services.NewCategoryCustomizationService(categoryCustomizationRepo)

	// Initialize handlers
	handler := handlers.NewHandler(userService, storeService, reviewService)
	viewerAuthHandler := handlers.NewViewerAuthHandler(viewerAuthService)
	categoryCustomizationHandler := handlers.NewCategoryCustomizationHandler(categoryCustomizationService, storeService)

	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	r := gin.New()

	// Add middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.NewCORSMiddleware(cfg.CORS))
	r.Use(middleware.RequestID())
	r.Use(middleware.SecurityHeaders())
	
	// CSRF protection
	csrfConfig := middleware.CSRFConfig{
		Secret:   cfg.JWT.Secret, // JWT秘密鍵を再利用
		Secure:   cfg.IsProduction(),
		SameSite: http.SameSiteStrictMode,
	}
	r.Use(middleware.CSRF(csrfConfig))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		// Check database connection
		dbStatus := "ok"
		if err := db.Ping(); err != nil {
			dbStatus = "error"
			c.JSON(503, types.APIResponse{
				Success: false,
				Error: &types.APIError{
					Code:    "HEALTH_CHECK_FAILED",
					Message: "Database connection failed",
					Details: err.Error(),
				},
			})
			return
		}

		c.JSON(200, types.APIResponse{
			Success: true,
			Data: map[string]interface{}{
				"status":     "ok",
				"message":    "Sukimise API is running",
				"version":    "1.0.0",
				"timestamp":  time.Now().Unix(),
				"environment": cfg.Server.Environment,
				"database":   dbStatus,
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

		// Viewer authentication routes
		viewer := api.Group("/viewer")
		{
			viewer.POST("/auth", viewerAuthHandler.AuthenticateViewer)
			viewer.GET("/validate", viewerAuthHandler.ValidateViewerSession)
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

		categoryCustomizations := api.Group("/category-customizations")
		{
			categoryCustomizations.GET("", categoryCustomizationHandler.GetCategoryCustomizations)
			categoryCustomizations.GET("/:categoryName", categoryCustomizationHandler.GetCategoryCustomization)
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

			// Admin-only routes for viewer settings
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/viewer-settings", viewerAuthHandler.GetViewerSettings)
				admin.PUT("/viewer-settings", viewerAuthHandler.UpdateViewerSettings)
				admin.GET("/viewer-history", viewerAuthHandler.GetViewerLoginHistory)
				admin.POST("/viewer-cleanup", viewerAuthHandler.CleanupExpiredSessions)

				// Category customization management (admin only)
				admin.POST("/category-customizations", categoryCustomizationHandler.CreateCategoryCustomization)
				admin.PUT("/category-customizations/:categoryName", categoryCustomizationHandler.UpdateCategoryCustomization)
				admin.DELETE("/category-customizations/:categoryName", categoryCustomizationHandler.DeleteCategoryCustomization)
				admin.POST("/category-customizations/sync", categoryCustomizationHandler.SyncCategoriesWithStores)
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

// validateUserEnvironmentVariables validates that required user environment variables are set
func validateUserEnvironmentVariables() error {
	adminUsers := os.Getenv("ADMIN_USERS")
	editorUsers := os.Getenv("EDITOR_USERS")

	if adminUsers == "" {
		return fmt.Errorf("ADMIN_USERS environment variable is required. Please set admin users in format: username1:bcrypt_hash1;username2:bcrypt_hash2")
	}

	if editorUsers == "" {
		return fmt.Errorf("EDITOR_USERS environment variable is required. Please set editor users in format: username1:bcrypt_hash1;username2:bcrypt_hash2")
	}

	// Validate admin users format
	if err := validateUserFormatting("ADMIN_USERS", adminUsers); err != nil {
		return err
	}

	// Validate editor users format
	if err := validateUserFormatting("EDITOR_USERS", editorUsers); err != nil {
		return err
	}

	log.Printf("User environment variables validated successfully")
	return nil
}

// validateUserFormatting validates the format of user entries
func validateUserFormatting(envVar, envValue string) error {
	userEntries := strings.Split(envValue, ";")
	
	if len(userEntries) == 0 {
		return fmt.Errorf("%s must contain at least one user entry", envVar)
	}

	for i, entry := range userEntries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue // Skip empty entries
		}

		parts := strings.SplitN(entry, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("%s entry %d has invalid format: %s (expected username:bcrypt_hash)", envVar, i+1, entry)
		}

		username := strings.TrimSpace(parts[0])
		password := strings.TrimSpace(parts[1])

		if username == "" {
			return fmt.Errorf("%s entry %d has empty username", envVar, i+1)
		}

		if password == "" {
			return fmt.Errorf("%s entry %d has empty password hash", envVar, i+1)
		}

		// Basic bcrypt hash validation (should start with $2a$, $2b$, or $2y$ and be around 60 chars)
		if !strings.HasPrefix(password, "$2a$") && !strings.HasPrefix(password, "$2b$") && !strings.HasPrefix(password, "$2y$") {
			return fmt.Errorf("%s entry %d has invalid bcrypt hash format: %s (should start with $2a$, $2b$, or $2y$)", envVar, i+1, username)
		}

		if len(password) < 50 || len(password) > 70 {
			return fmt.Errorf("%s entry %d has invalid bcrypt hash length for user: %s (should be around 60 characters)", envVar, i+1, username)
		}
	}

	return nil
}